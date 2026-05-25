package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"memory-analyzer/memory"
)

func runAnalyze(args []string) {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "输出 JSON 格式")
	asMD := fs.Bool("markdown", false, "输出 Markdown 格式")
	topN := fs.Int("top", 30, "进程列表只显示前 N 个")
	fs.Parse(args)

	report, err := memory.Analyze(*topN)
	if err != nil {
		fatal("采集失败: %v", err)
	}

	switch {
	case *asJSON:
		printAnalyzeJSON(report)
	case *asMD:
		printAnalyzeMarkdown(report)
	default:
		printAnalyzeText(report)
	}
}

func printAnalyzeJSON(r *memory.AnalysisReport) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r)
}

func printAnalyzeText(r *memory.AnalysisReport) {
	sys := r.System
	cpu := r.CPU

	// 标题
	fmt.Println()
	fmt.Println(paint("MemorAI · System Snapshot", cBold+cCyan))
	fmt.Println(paint(hr(60), cGray))
	fmt.Println()

	// 内存仪表
	fmt.Printf("  %s  %s  %s\n",
		paint("Memory", cBold),
		progressBar(sys.UsedPercent, 30),
		paintPct(sys.UsedPercent),
	)
	fmt.Printf("  %s used / %s total · %s available\n",
		paint(fmt.Sprintf("%.2f GB", sys.UsedGB), cBold),
		fmt.Sprintf("%.2f GB", sys.TotalGB),
		paint(fmt.Sprintf("%.2f GB", sys.AvailableGB), cGreen),
	)
	fmt.Println()

	// CPU
	if cpu != nil {
		fmt.Printf("  %s     %s  %s\n",
			paint("CPU   ", cBold),
			progressBar(cpu.UsedPercent, 30),
			paintPct(cpu.UsedPercent),
		)
		if cpu.ModelName != "" {
			fmt.Printf("  %s · %d cores / %d threads",
				truncate(cpu.ModelName, 50), cpu.Cores, cpu.LogicalCores)
			if cpu.MHz > 0 {
				fmt.Printf(" @ %.2f GHz", cpu.MHz/1000.0)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// 系统信息
	fmt.Printf("  %s  %s %s · host: %s · uptime: %.1fh\n",
		paint("System", cDim),
		sys.Platform, sys.PlatformVersion, sys.Hostname,
		float64(sys.UptimeSeconds)/3600.0,
	)
	fmt.Println()

	// 软件分组
	fmt.Println(paint("Top software (grouped):", cBold))
	fmt.Println(paint(hr(60), cGray))
	header := formatRow(
		[]string{"#", "SOFTWARE", "MEMORY", "PCT", "PROCS"},
		[]int{3, 32, 10, 7, 5},
	)
	fmt.Println(paint(header, cDim))

	limit := 15
	if len(r.Groups) < limit {
		limit = len(r.Groups)
	}
	for i := 0; i < limit; i++ {
		g := r.Groups[i]
		mem := fmt.Sprintf("%.0f MB", g.TotalMB)
		if g.TotalMB >= 1024 {
			mem = fmt.Sprintf("%.2f GB", g.TotalMB/1024)
		}
		row := formatRow(
			[]string{
				fmt.Sprintf("%d", i+1),
				truncate(g.Name, 30),
				mem,
				fmt.Sprintf("%.1f%%", g.TotalPct),
				fmt.Sprintf("× %d", g.ProcessCount),
			},
			[]int{3, 32, 10, 7, 5},
		)
		// 内存大于 500 MB 标红，大于 200 标黄
		if g.TotalMB >= 500 {
			fmt.Println(paint(row, cYellow))
		} else {
			fmt.Println(row)
		}
	}
	if len(r.Groups) > limit {
		fmt.Printf("  %s\n", paint(fmt.Sprintf("... %d more groups", len(r.Groups)-limit), cDim))
	}
	fmt.Println()
	fmt.Printf("  %s %d processes · %d software groups\n",
		paint("Total:", cDim), r.Total, len(r.Groups))
	fmt.Println()
}

func printAnalyzeMarkdown(r *memory.AnalysisReport) {
	var sb strings.Builder
	sys := r.System
	cpu := r.CPU
	sb.WriteString("# MemorAI Snapshot\n\n")
	sb.WriteString("## System\n")
	fmt.Fprintf(&sb, "- OS: %s %s\n", sys.Platform, sys.PlatformVersion)
	fmt.Fprintf(&sb, "- Memory: %.2f / %.2f GB (%.1f%%)\n", sys.UsedGB, sys.TotalGB, sys.UsedPercent)
	fmt.Fprintf(&sb, "- Available: %.2f GB\n", sys.AvailableGB)
	if cpu != nil {
		fmt.Fprintf(&sb, "- CPU: %.1f%% · %d cores / %d threads", cpu.UsedPercent, cpu.Cores, cpu.LogicalCores)
		if cpu.ModelName != "" {
			fmt.Fprintf(&sb, " · %s", cpu.ModelName)
		}
		sb.WriteString("\n")
	}
	fmt.Fprintf(&sb, "- Uptime: %.1f hours\n\n", float64(sys.UptimeSeconds)/3600.0)

	sb.WriteString("## Top Software (Grouped)\n\n")
	sb.WriteString("| # | Software | Category | Memory | % | Procs |\n")
	sb.WriteString("|---|----------|----------|--------|---|-------|\n")
	limit := 20
	if len(r.Groups) < limit {
		limit = len(r.Groups)
	}
	for i := 0; i < limit; i++ {
		g := r.Groups[i]
		fmt.Fprintf(&sb, "| %d | %s | %s | %.0f MB | %.1f%% | %d |\n",
			i+1, g.Name, g.Category, g.TotalMB, g.TotalPct, g.ProcessCount)
	}
	fmt.Fprintf(&sb, "\n_Total %d processes, %d groups_\n", r.Total, len(r.Groups))

	fmt.Print(sb.String())
}
