package monitor

import (
	"context"
	"errors"
	"sync"
	"time"

	"memory-analyzer/memory"
	"memory-analyzer/sysstat"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Sample 单次采样
type Sample struct {
	TimestampMs int64   `json:"timestampMs"` // 自采样开始的相对毫秒
	WallTime    string  `json:"wallTime"`    // 墙上时间 HH:MM:SS
	UsedPercent float64 `json:"usedPercent"`
	UsedGB      float64 `json:"usedGB"`
	AvailableGB float64 `json:"availableGB"`
	CPUPercent  float64 `json:"cpuPercent"` // CPU 占用率
	// 该时刻 Top N 软件分组的内存（MB）
	TopGroups []GroupPoint `json:"topGroups"`
}

// GroupPoint 折线图所需的单组数据点
type GroupPoint struct {
	Name string  `json:"name"`
	MB   float64 `json:"mb"`
}

// Result 30秒采样的完整结果
type Result struct {
	StartedAt       string                  `json:"startedAt"`
	DurationSec     int                     `json:"durationSec"`
	IntervalSec     int                     `json:"intervalSec"`
	Samples         []Sample                `json:"samples"`
	FinalGroups     []*memory.SoftwareGroup `json:"finalGroups"`
	System          *memory.SystemMemoryInfo `json:"system"`
	// 用于折线图的稳定 series（采样期间出现过的 Top 8 软件名）
	TrackedSoftware []string `json:"trackedSoftware"`
}

// 全局状态：避免并发采样
var (
	mu      sync.Mutex
	running bool
)

// IsRunning 是否正在采样
func IsRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return running
}

// Start 开始一次采样
// duration: 总时长（秒），interval: 采样间隔（秒）
// ctx 可为 nil（CLI 调用场景，不会发射进度事件）
func Start(ctx context.Context, duration, interval int) (*Result, error) {
	mu.Lock()
	if running {
		mu.Unlock()
		return nil, errors.New("已经有一个采样在进行中")
	}
	running = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		running = false
		mu.Unlock()
	}()

	// wailsCtx 用于发射事件（nil 时不发射，CLI 模式）
	wailsCtx := ctx
	// runCtx 用于 select；nil 时用 Background 避免 panic
	runCtx := ctx
	if runCtx == nil {
		runCtx = context.Background()
	}

	if duration <= 0 {
		duration = 30
	}
	if interval <= 0 {
		interval = 1
	}
	totalSamples := duration / interval
	if totalSamples < 1 {
		totalSamples = 1
	}

	startedAt := time.Now()
	samples := make([]Sample, 0, totalSamples)

	// 用于跟踪 Top 软件
	softwareSeen := make(map[string]float64)
	const maxTrackedTop = 8

	// 预热 CPU 采样器，让后续 0-interval 调用能拿到正确的差值
	sysstat.WarmupCPU()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	collectAndAppend := func(idx int) error {
		report, err := memory.Analyze(0)
		if err != nil {
			return err
		}

		// 取该时刻 Top 8 软件
		topN := 8
		if len(report.Groups) < topN {
			topN = len(report.Groups)
		}
		gpoints := make([]GroupPoint, 0, topN)
		for i := 0; i < topN; i++ {
			g := report.Groups[i]
			gpoints = append(gpoints, GroupPoint{Name: g.Name, MB: g.TotalMB})
			if v, ok := softwareSeen[g.Name]; !ok || g.TotalMB > v {
				softwareSeen[g.Name] = g.TotalMB
			}
		}

		// CPU 占用率（使用上次到现在的差值，所以是非阻塞的）
		cpuPct := sysstat.GetCPUQuick()

		now := time.Now()
		s := Sample{
			TimestampMs: now.Sub(startedAt).Milliseconds(),
			WallTime:    now.Format("15:04:05"),
			UsedPercent: report.System.UsedPercent,
			UsedGB:      report.System.UsedGB,
			AvailableGB: report.System.AvailableGB,
			CPUPercent:  cpuPct,
			TopGroups:   gpoints,
		}
		samples = append(samples, s)

		if wailsCtx != nil {
			wailsruntime.EventsEmit(wailsCtx, "monitor:progress", map[string]interface{}{
				"index":   idx + 1,
				"total":   totalSamples,
				"sample":  s,
				"percent": float64(idx+1) / float64(totalSamples) * 100,
			})
		}
		return nil
	}

	if err := collectAndAppend(0); err != nil {
		return nil, err
	}

	for i := 1; i < totalSamples; i++ {
		select {
		case <-runCtx.Done():
			return nil, runCtx.Err()
		case <-ticker.C:
			if err := collectAndAppend(i); err != nil {
				return nil, err
			}
		}
	}

	// 整理 tracked software：取出现过的最大值排序，留 Top 8
	type kv struct {
		k string
		v float64
	}
	all := make([]kv, 0, len(softwareSeen))
	for k, v := range softwareSeen {
		all = append(all, kv{k, v})
	}
	// 简单的选择排序取前 maxTrackedTop
	for i := 0; i < len(all) && i < maxTrackedTop; i++ {
		maxIdx := i
		for j := i + 1; j < len(all); j++ {
			if all[j].v > all[maxIdx].v {
				maxIdx = j
			}
		}
		all[i], all[maxIdx] = all[maxIdx], all[i]
	}
	if len(all) > maxTrackedTop {
		all = all[:maxTrackedTop]
	}
	tracked := make([]string, 0, len(all))
	for _, x := range all {
		tracked = append(tracked, x.k)
	}

	// 最终快照（用于饼图）
	finalReport, err := memory.Analyze(0)
	if err != nil {
		return nil, err
	}

	res := &Result{
		StartedAt:       startedAt.Format("2006-01-02 15:04:05"),
		DurationSec:     duration,
		IntervalSec:     interval,
		Samples:         samples,
		FinalGroups:     finalReport.Groups,
		System:          finalReport.System,
		TrackedSoftware: tracked,
	}

	// 完成事件
	if wailsCtx != nil {
		wailsruntime.EventsEmit(wailsCtx, "monitor:done", res)
	}

	return res, nil
}
