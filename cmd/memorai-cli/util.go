package main

import (
	"fmt"
	"os"
	"strings"
)

// 是否启用 ANSI 颜色（CI/管道环境自动关闭）
var colorEnabled = isTerminal()

func isTerminal() bool {
	// Win10/11 默认终端支持 ANSI；老版 cmd 不支持但 PowerShell 可以
	// 简单判断：stdout 是字符设备
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ANSI 颜色码
const (
	cReset  = "\x1b[0m"
	cBold   = "\x1b[1m"
	cDim    = "\x1b[2m"
	cRed    = "\x1b[31m"
	cGreen  = "\x1b[32m"
	cYellow = "\x1b[33m"
	cBlue   = "\x1b[34m"
	cCyan   = "\x1b[36m"
	cGray   = "\x1b[90m"
)

func paint(s, color string) string {
	if !colorEnabled {
		return s
	}
	return color + s + cReset
}

// 状态颜色：根据百分比上色
func paintPct(pct float64) string {
	s := fmt.Sprintf("%.1f%%", pct)
	if !colorEnabled {
		return s
	}
	switch {
	case pct >= 85:
		return cRed + cBold + s + cReset
	case pct >= 65:
		return cYellow + s + cReset
	default:
		return cGreen + s + cReset
	}
}

// 简易进度条
func progressBar(pct float64, width int) string {
	if width <= 0 {
		width = 20
	}
	filled := int(float64(width) * pct / 100)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	if !colorEnabled {
		return bar
	}
	switch {
	case pct >= 85:
		return cRed + bar + cReset
	case pct >= 65:
		return cYellow + bar + cReset
	default:
		return cCyan + bar + cReset
	}
}

// 简单表格行（左对齐，右对齐数字）
func formatRow(cols []string, widths []int) string {
	var sb strings.Builder
	for i, c := range cols {
		w := 0
		if i < len(widths) {
			w = widths[i]
		}
		if w == 0 {
			sb.WriteString(c)
		} else {
			sb.WriteString(padRight(c, w))
		}
		if i < len(cols)-1 {
			sb.WriteString("  ")
		}
	}
	return sb.String()
}

func padRight(s string, w int) string {
	// 简单按字节填充（中文显示宽度问题暂不处理）
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

func padLeft(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return strings.Repeat(" ", w-len(s)) + s
}

// 截断字符串到 max 字符（按字节）
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}

// 打印分隔线
func hr(width int) string {
	return strings.Repeat("─", width)
}

// 错误退出
func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, paint("error: ", cRed+cBold)+format+"\n", args...)
	os.Exit(1)
}
