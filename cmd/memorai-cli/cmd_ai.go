package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"memory-analyzer/ai"
	"memory-analyzer/config"
	"memory-analyzer/memory"
	"memory-analyzer/startup"
)

func runAI(args []string) {
	fs := flag.NewFlagSet("ai", flag.ExitOnError)
	timeout := fs.Int("timeout", 120, "AI 请求超时秒数")
	fs.Parse(args)

	cfg, err := config.Load()
	if err != nil {
		fatal("加载配置失败: %v", err)
	}
	if cfg.AIAPIKey == "" {
		fatal("尚未配置 AI API Key。运行: memorai-cli config set ai-key YOUR_KEY")
	}

	fmt.Println(paint("Collecting snapshot...", cDim))
	report, err := memory.Analyze(cfg.TopN)
	if err != nil {
		fatal("采集失败: %v", err)
	}
	startupItems, _ := startup.EnumerateStartupItems()

	systemPrompt := ai.SystemPrompt
	if cfg.AIExtraPrompt != "" {
		systemPrompt += "\n\n## 用户自定义指令\n" + cfg.AIExtraPrompt
	}
	userMsg := ai.BuildUserPrompt(report, startupItems)

	fmt.Printf(paint("Asking %s @ %s ...\n", cDim), cfg.AIModel, cfg.AIBaseURL)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	client := ai.NewClient(cfg.AIBaseURL, cfg.AIAPIKey, cfg.AIModel)
	resp, err := client.Chat(ctx, systemPrompt, userMsg)
	if err != nil {
		fatal("AI 调用失败: %v", err)
	}

	// 终端 markdown 渲染（极简版：标题着色、列表保留）
	printMarkdown(resp)
	fmt.Println()
}

// printMarkdown 极简终端 Markdown 渲染
func printMarkdown(md string) {
	for _, line := range splitLines(md) {
		switch {
		case startsWith(line, "## "):
			fmt.Println(paint(line[3:], cBold+cCyan))
			fmt.Println(paint(hr(len(line)-3), cGray))
		case startsWith(line, "# "):
			fmt.Println(paint(line[2:], cBold+cCyan))
		case startsWith(line, "### "):
			fmt.Println(paint(line[4:], cBold))
		case startsWith(line, "- ") || startsWith(line, "* "):
			fmt.Println(paint("  •", cCyan), line[2:])
		default:
			fmt.Println(line)
		}
	}
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func startsWith(s, p string) bool {
	return len(s) >= len(p) && s[:len(p)] == p
}
