package ai

import (
	"fmt"
	"path/filepath"
	"strings"

	"memory-analyzer/memory"
	"memory-analyzer/startup"
)

const SystemPrompt = `你是 MemorAI 内置的 Windows 系统优化专家，擅长分析内存与 CPU 占用问题。
你将收到用户电脑的状态快照（系统内存、CPU、按软件分组的进程列表、开机自启项）。

请按下面的格式输出中文分析（用 Markdown）：

## 总体评估
（1-2句话评价当前内存与 CPU 使用情况是否合理）

## 主要资源占用分析
（列出占用最高的几类软件，分析每个是否合理、能否优化；同时关注 CPU 是否异常高）

## 异常或可疑项
（如果发现占用过高、不必要、或可能有泄漏的进程，重点指出）

## 自启项建议
（重点：根据"开机自启项"列表，结合每个软件当前的实际占用，告诉用户：
1. 哪些自启项明显没必要、占内存又大，强烈建议禁用
2. 哪些是必要的（杀毒、显卡驱动、输入法等）不建议禁用
3. 哪些自启项当前没运行但仍设了自启，可以禁用以加快开机
对每条具体建议，明确说出"建议在【自启项】页禁用：XXX"）

## 优化建议
（按优先级给出其他具体可操作的建议）

注意：
- 用户是普通用户，建议要具体可操作（说清在哪里点哪个按钮）
- 杀毒软件如卡巴斯基占用高是正常的，但可以建议调整防护级别
- WSL2/vmmem、Hyper-V、Docker 是常见的内存大户，要明确指出
- CPU 占用瞬时值仅参考，建议观察一段时间再下结论
- 不要捏造数据，只基于提供的快照分析
- 系统服务、显卡驱动、输入法等绝对不要建议禁用
`

// BuildUserPrompt 根据采集到的报告生成用户消息
// startupItems 可以为 nil（当采集失败时）
func BuildUserPrompt(r *memory.AnalysisReport, startupItems []*startup.StartupItem) string {
	var sb strings.Builder

	sb.WriteString("# 内存采集快照\n\n")

	// 系统总览
	sb.WriteString("## 系统信息\n")
	if r.System != nil {
		fmt.Fprintf(&sb, "- 操作系统: %s %s\n", r.System.Platform, r.System.PlatformVersion)
		fmt.Fprintf(&sb, "- 总内存: %.2f GB\n", r.System.TotalGB)
		fmt.Fprintf(&sb, "- 已用内存: %.2f GB (%.1f%%)\n", r.System.UsedGB, r.System.UsedPercent)
		fmt.Fprintf(&sb, "- 可用内存: %.2f GB\n", r.System.AvailableGB)
		fmt.Fprintf(&sb, "- 已开机时长: %.1f 小时\n", float64(r.System.UptimeSeconds)/3600.0)
	}
	if r.CPU != nil {
		sb.WriteString("\n## CPU 信息\n")
		if r.CPU.ModelName != "" {
			fmt.Fprintf(&sb, "- 型号: %s\n", r.CPU.ModelName)
		}
		fmt.Fprintf(&sb, "- 物理核心: %d 个 / 逻辑核心: %d 个\n", r.CPU.Cores, r.CPU.LogicalCores)
		if r.CPU.MHz > 0 {
			fmt.Fprintf(&sb, "- 主频: %.2f GHz\n", r.CPU.MHz/1000.0)
		}
		fmt.Fprintf(&sb, "- 当前总占用率: %.1f%%\n", r.CPU.UsedPercent)
	}

	// 软件分组
	sb.WriteString("\n## 按软件分组的内存占用（前 20 项）\n")
	limit := 20
	if len(r.Groups) < limit {
		limit = len(r.Groups)
	}
	for i := 0; i < limit; i++ {
		g := r.Groups[i]
		fmt.Fprintf(&sb, "- [%s] %s: %.0f MB (%.1f%%) — %d 个进程\n",
			g.Category, g.Name, g.TotalMB, g.TotalPct, g.ProcessCount)
	}

	// 进程详情
	sb.WriteString("\n## 单进程占用 Top 15\n")
	pLimit := 15
	if len(r.TopProcs) < pLimit {
		pLimit = len(r.TopProcs)
	}
	for i := 0; i < pLimit; i++ {
		p := r.TopProcs[i]
		fmt.Fprintf(&sb, "- %s (PID %d): %.0f MB\n", p.Name, p.PID, p.MemoryMB)
	}

	// 自启项 + 与内存关联
	if len(startupItems) > 0 {
		writeStartupSection(&sb, r, startupItems)
	} else {
		sb.WriteString("\n## 开机自启项\n（采集失败或无数据）\n")
	}

	sb.WriteString("\n请基于以上数据进行分析。")
	return sb.String()
}

// writeStartupSection 写入自启项段落，并尝试匹配每个自启项的当前内存占用
func writeStartupSection(sb *strings.Builder, r *memory.AnalysisReport, items []*startup.StartupItem) {
	// 把 TopProcs 索引一下，按进程名小写
	procByName := make(map[string]float64) // 同名进程总和
	for _, p := range r.TopProcs {
		key := strings.ToLower(p.Name)
		procByName[key] += p.MemoryMB
	}

	enabled := make([]*startup.StartupItem, 0)
	disabled := make([]*startup.StartupItem, 0)
	for _, it := range items {
		if it.Enabled {
			enabled = append(enabled, it)
		} else {
			disabled = append(disabled, it)
		}
	}

	fmt.Fprintf(sb, "\n## 开机自启项（共 %d 项，启用 %d / 已禁用 %d）\n",
		len(items), len(enabled), len(disabled))

	if len(enabled) > 0 {
		sb.WriteString("\n### 已启用（开机会自动启动）\n")
		for _, it := range enabled {
			scope := "用户"
			if it.Scope == startup.ScopeSystem {
				scope = "系统"
			}
			memInfo := matchMemory(it, procByName)
			fmt.Fprintf(sb, "- [%s] %s — `%s`%s\n",
				scope, it.Name, truncForPrompt(it.Command, 100), memInfo)
		}
	}

	if len(disabled) > 0 {
		sb.WriteString("\n### 已禁用\n")
		for _, it := range disabled {
			scope := "用户"
			if it.Scope == startup.ScopeSystem {
				scope = "系统"
			}
			fmt.Fprintf(sb, "- [%s] %s\n", scope, it.Name)
		}
	}
}

// matchMemory 尝试把自启项关联到当前进程内存
// 通过 command 中的 exe 文件名匹配 procByName
func matchMemory(it *startup.StartupItem, procByName map[string]float64) string {
	exe := extractExeName(it.Command)
	if exe == "" {
		exe = it.Name + ".exe" // fallback
	}
	exeLower := strings.ToLower(exe)
	if mb, ok := procByName[exeLower]; ok {
		marker := ""
		if mb >= 500 {
			marker = " ⚠️"
		}
		return fmt.Sprintf("，当前运行中占用 **%.0f MB**%s", mb, marker)
	}
	return "，当前未在 Top 进程中（占用较小或未运行）"
}

// extractExeName 从命令行中提取可执行文件名
// 例: `"C:\Program Files\xx\app.exe" --silent` -> "app.exe"
func extractExeName(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}
	// 去引号
	if strings.HasPrefix(cmd, `"`) {
		end := strings.Index(cmd[1:], `"`)
		if end > 0 {
			return filepath.Base(cmd[1 : 1+end])
		}
	}
	// 没引号，按空格切第一个 token
	first := strings.Fields(cmd)[0]
	return filepath.Base(first)
}

func truncForPrompt(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
