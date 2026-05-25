package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"memory-analyzer/monitor"
	"memory-analyzer/sysstat"
)

func runMonitor(args []string) {
	fs := flag.NewFlagSet("monitor", flag.ExitOnError)
	seconds := fs.Int("seconds", 30, "采样总时长（秒）")
	interval := fs.Int("interval", 1, "采样间隔（秒）")
	asJSON := fs.Bool("json", false, "输出 JSON")
	out := fs.String("o", "", "结果保存到文件（JSON）")
	fs.Parse(args)

	if *seconds < 1 {
		*seconds = 30
	}
	if *interval < 1 {
		*interval = 1
	}

	fmt.Printf(paint("Sampling %ds (interval %ds)...\n\n", cDim), *seconds, *interval)

	// 用一个简化的采样循环（不依赖 wails ctx）
	sysstat.WarmupCPU()
	startedAt := time.Now()
	totalSamples := *seconds / *interval
	if totalSamples < 1 {
		totalSamples = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*seconds+5)*time.Second)
	defer cancel()

	// 自定义简单采样（不通过 wails EventsEmit），实时打印进度
	result, err := sampleWithProgress(ctx, *seconds, *interval, totalSamples, startedAt)
	if err != nil {
		fatal("采样失败: %v", err)
	}

	if *asJSON || *out != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if *out != "" {
			if err := os.WriteFile(*out, data, 0644); err != nil {
				fatal("写文件失败: %v", err)
			}
			fmt.Printf(paint("Saved to %s\n", cGreen), *out)
		} else {
			os.Stdout.Write(data)
			fmt.Println()
		}
		return
	}

	// 文本摘要
	printMonitorSummary(result)
}

// sampleWithProgress 直接调用 monitor.Start（其内部已经做了采样）
// 但为了在 CLI 显示实时进度，这里手写一个简化版
func sampleWithProgress(ctx context.Context, seconds, interval, totalSamples int, startedAt time.Time) (*monitor.Result, error) {
	// 复用 monitor.Start（它会采样并通过 wails ctx 推事件；CLI 下 ctx 是 nil → 不推事件）
	// 但 monitor.Start 的进度需要从外部观察，简单起见我们循环 totalSamples 次手动采样
	// 这里直接调用 monitor.Start 但在另一个 goroutine 跑进度提示

	doneCh := make(chan struct {
		r   *monitor.Result
		err error
	}, 1)
	go func() {
		// 调用 monitor.Start 时 ctx 传 nil 会避免事件发射
		r, err := monitor.Start(nil, seconds, interval)
		doneCh <- struct {
			r   *monitor.Result
			err error
		}{r, err}
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	doneSamples := 0
	for {
		select {
		case res := <-doneCh:
			fmt.Print("\r" + paint(fmt.Sprintf("[%d/%d] sampling complete   \n", totalSamples, totalSamples), cGreen))
			return res.r, res.err
		case <-ticker.C:
			doneSamples++
			if doneSamples > totalSamples {
				doneSamples = totalSamples
			}
			pct := float64(doneSamples) / float64(totalSamples) * 100
			fmt.Printf("\r%s %s %s",
				paint(fmt.Sprintf("[%d/%d]", doneSamples, totalSamples), cDim),
				progressBar(pct, 25),
				paint(fmt.Sprintf("%.0f%%", pct), cDim),
			)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func printMonitorSummary(r *monitor.Result) {
	if r == nil || len(r.Samples) == 0 {
		fmt.Println("无采样数据")
		return
	}

	// 内存/CPU 统计
	var memMin, memMax, memSum, cpuMin, cpuMax, cpuSum float64
	memMin, cpuMin = 100, 100
	for _, s := range r.Samples {
		if s.UsedPercent < memMin {
			memMin = s.UsedPercent
		}
		if s.UsedPercent > memMax {
			memMax = s.UsedPercent
		}
		memSum += s.UsedPercent
		if s.CPUPercent < cpuMin {
			cpuMin = s.CPUPercent
		}
		if s.CPUPercent > cpuMax {
			cpuMax = s.CPUPercent
		}
		cpuSum += s.CPUPercent
	}
	n := float64(len(r.Samples))
	memAvg := memSum / n
	cpuAvg := cpuSum / n

	fmt.Println()
	fmt.Println(paint("Sampling Result", cBold+cCyan))
	fmt.Println(paint(hr(60), cGray))
	fmt.Printf("  Started:   %s\n", r.StartedAt)
	fmt.Printf("  Duration:  %ds (interval %ds, %d samples)\n", r.DurationSec, r.IntervalSec, len(r.Samples))
	fmt.Println()

	fmt.Printf("  %s   min %5.1f%%   max %5.1f%%   avg %5.1f%%\n",
		paint("Memory", cBold), memMin, memMax, memAvg)
	fmt.Printf("  %s      min %5.1f%%   max %5.1f%%   avg %5.1f%%\n",
		paint("CPU", cBold), cpuMin, cpuMax, cpuAvg)
	fmt.Println()

	// Top 软件占比
	if len(r.FinalGroups) > 0 {
		fmt.Println(paint("Final Top 10 Software:", cBold))
		fmt.Println(paint(hr(60), cGray))
		limit := 10
		if len(r.FinalGroups) < limit {
			limit = len(r.FinalGroups)
		}
		for i := 0; i < limit; i++ {
			g := r.FinalGroups[i]
			fmt.Printf("  %s  %s  %s  %s\n",
				padLeft(fmt.Sprintf("%d.", i+1), 3),
				padRight(truncate(g.Name, 30), 30),
				padLeft(fmt.Sprintf("%.0f MB", g.TotalMB), 9),
				paint(fmt.Sprintf("%.1f%%", g.TotalPct), cDim),
			)
		}
		fmt.Println()
	}

	// ASCII 折线图（简化）
	if len(r.Samples) > 1 {
		fmt.Println(paint("Memory % over time:", cBold))
		printSparkline(r.Samples, "memory")
		fmt.Println()
		fmt.Println(paint("CPU % over time:", cBold))
		printSparkline(r.Samples, "cpu")
		fmt.Println()
	}
}

// printSparkline 用 Unicode 块字符画一条迷你折线
func printSparkline(samples []monitor.Sample, kind string) {
	if len(samples) == 0 {
		return
	}
	// 8级
	bars := []rune("▁▂▃▄▅▆▇█")
	values := make([]float64, len(samples))
	var maxV float64 = 0.0001
	for i, s := range samples {
		v := s.UsedPercent
		if kind == "cpu" {
			v = s.CPUPercent
		}
		values[i] = v
		if v > maxV {
			maxV = v
		}
	}
	// 限制宽度（最多 60 字符）
	width := len(values)
	if width > 60 {
		// 等距抽样
		step := float64(width) / 60
		newVals := make([]float64, 60)
		for i := 0; i < 60; i++ {
			idx := int(float64(i) * step)
			if idx >= width {
				idx = width - 1
			}
			newVals[i] = values[idx]
		}
		values = newVals
		width = 60
	}
	fmt.Print("  ")
	for _, v := range values {
		idx := int(v / maxV * float64(len(bars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		fmt.Printf("%c", bars[idx])
	}
	fmt.Printf("  max=%.1f%%\n", maxV)
}
