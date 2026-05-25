package main

import (
	"fmt"
	"os"
	"strings"

	"memory-analyzer/config"
)

func runConfig(args []string) {
	if len(args) == 0 {
		printConfigUsage()
		return
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "show", "get":
		runConfigShow()
	case "set":
		runConfigSet(rest)
	case "path":
		runConfigPath()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n\n", sub)
		printConfigUsage()
	}
}

func runConfigShow() {
	cfg, err := config.Load()
	if err != nil {
		fatal("加载配置失败: %v", err)
	}
	fmt.Println()
	fmt.Println(paint("MemorAI Config", cBold+cCyan))
	fmt.Println(paint(hr(50), cGray))

	mask := cfg.AIAPIKey
	if len(mask) > 8 {
		mask = mask[:4] + "****" + mask[len(mask)-4:]
	} else if mask != "" {
		mask = "****"
	}

	fmt.Printf("  ai-base-url     %s\n", cfg.AIBaseURL)
	fmt.Printf("  ai-key          %s\n", mask)
	fmt.Printf("  ai-model        %s\n", cfg.AIModel)
	fmt.Printf("  ai-extra-prompt %s\n", truncOrEmpty(cfg.AIExtraPrompt, 60))
	fmt.Printf("  top-n           %d\n", cfg.TopN)
	fmt.Println()
}

func runConfigSet(args []string) {
	if len(args) < 2 {
		fatal("用法: memorai-cli config set <key> <value>\n可用的 key: ai-base-url, ai-key, ai-model, ai-extra-prompt, top-n")
	}
	key := args[0]
	value := strings.Join(args[1:], " ")

	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	switch key {
	case "ai-base-url", "ai-baseurl", "url":
		cfg.AIBaseURL = value
	case "ai-key", "api-key", "key":
		cfg.AIAPIKey = value
	case "ai-model", "model":
		cfg.AIModel = value
	case "ai-extra-prompt", "extra-prompt", "prompt":
		cfg.AIExtraPrompt = value
	case "top-n", "topn":
		var n int
		_, err := fmt.Sscanf(value, "%d", &n)
		if err != nil || n <= 0 {
			fatal("top-n 必须是正整数")
		}
		cfg.TopN = n
	default:
		fatal("未知的配置项: %s\n可用: ai-base-url, ai-key, ai-model, ai-extra-prompt, top-n", key)
	}

	if err := config.Save(cfg); err != nil {
		fatal("保存失败: %v", err)
	}
	fmt.Printf(paint("✓ %s updated\n", cGreen), key)
}

func runConfigPath() {
	exe, err := os.Executable()
	if err != nil {
		fatal("无法获取 exe 路径: %v", err)
	}
	dir := exe[:strings.LastIndexAny(exe, "/\\")]
	fmt.Println(dir + string(os.PathSeparator) + "config.json")
}

func truncOrEmpty(s string, n int) string {
	if s == "" {
		return paint("(empty)", cDim)
	}
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

func printConfigUsage() {
	fmt.Println(`config - 配置管理

USAGE:
    memorai-cli config show
    memorai-cli config set <key> <value>
    memorai-cli config path

KEYS:
    ai-base-url      OpenAI 兼容接口地址
    ai-key           API Key
    ai-model         模型名（如 deepseek-chat）
    ai-extra-prompt  自定义系统提示词
    top-n            进程列表 Top N

EXAMPLES:
    memorai-cli config show
    memorai-cli config set ai-key sk-xxxxxxxxxxxx
    memorai-cli config set ai-base-url https://api.deepseek.com/v1
    memorai-cli config set ai-model deepseek-chat`)
}
