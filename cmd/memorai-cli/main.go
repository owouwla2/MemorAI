// Package main 是 MemorAI 的命令行版本入口。
// 使用方式: memorai-cli <command> [flags]
//
// 命令列表:
//   analyze    采集一次内存快照并输出
//   ai         让 AI 分析当前内存（需配置 API Key）
//   monitor    30秒（可调）采样
//   startup    自启项管理 (list / enable / disable)
//   config     配置查看与修改 (show / set)
//   help       查看帮助
//
// 体积/内存目标: <15MB exe, <20MB 运行内存
package main

import (
	"fmt"
	"os"
)

const version = "v0.4.0"

func main() {
	if len(os.Args) < 2 {
		printRoot()
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "analyze", "a":
		runAnalyze(args)
	case "ai":
		runAI(args)
	case "monitor", "mon":
		runMonitor(args)
	case "startup", "su":
		runStartup(args)
	case "config", "cfg":
		runConfig(args)
	case "version", "-v", "--version":
		fmt.Printf("memorai-cli %s\n", version)
	case "help", "-h", "--help":
		printRoot()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printRoot()
		os.Exit(2)
	}
}

func printRoot() {
	fmt.Println(`MemorAI CLI - AI-native memory analyzer for Windows

USAGE:
    memorai-cli <command> [flags]

COMMANDS:
    analyze     采集一次内存快照
    ai          让 AI 分析（需配置 API Key）
    monitor     30秒采样监控
    startup     自启项管理 (list/enable/disable)
    config      配置 (show/set)
    version     显示版本
    help        显示帮助

EXAMPLES:
    memorai-cli analyze              # 文本输出
    memorai-cli analyze --json       # JSON 输出
    memorai-cli analyze --markdown   # Markdown 输出
    memorai-cli ai                   # AI 分析当前快照
    memorai-cli monitor --seconds 30 # 30秒采样
    memorai-cli startup list         # 列出自启项
    memorai-cli startup disable <id> # 禁用某项
    memorai-cli config show          # 查看配置
    memorai-cli config set ai-key sk-xxx

For per-command help, use: memorai-cli <command> --help`)
}
