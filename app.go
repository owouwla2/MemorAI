package main

import (
	"context"
	"fmt"

	"memory-analyzer/ai"
	"memory-analyzer/config"
	"memory-analyzer/memory"
	"memory-analyzer/monitor"
	"memory-analyzer/startup"
)

// App 应用主结构
type App struct {
	ctx context.Context
}

// NewApp 创建应用
func NewApp() *App {
	return &App{}
}

// startup Wails 启动回调
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ===== 内存分析 =====

// AnalyzeMemory 采集一次内存快照
func (a *App) AnalyzeMemory() (*memory.AnalysisReport, error) {
	cfg, _ := config.Load()
	topN := 30
	if cfg != nil && cfg.TopN > 0 {
		topN = cfg.TopN
	}
	return memory.Analyze(topN)
}

// ===== 30秒采样监控 =====

// StartMonitoring 启动一段时间采样
// duration: 总秒数（默认30）
// interval: 采样间隔秒数（默认1）
func (a *App) StartMonitoring(duration, interval int) (*monitor.Result, error) {
	if duration <= 0 {
		duration = 30
	}
	if interval <= 0 {
		interval = 1
	}
	return monitor.Start(a.ctx, duration, interval)
}

// IsMonitoringRunning 是否正在采样
func (a *App) IsMonitoringRunning() bool {
	return monitor.IsRunning()
}

// ===== 自启项管理 =====

// ListStartupItems 列出所有开机自启项
func (a *App) ListStartupItems() ([]*startup.StartupItem, error) {
	return startup.EnumerateStartupItems()
}

// ToggleStartupItem 启用/禁用自启项
func (a *App) ToggleStartupItem(id string, enable bool) error {
	return startup.SetEnabled(id, enable)
}

// IsAdmin 当前进程是否管理员
func (a *App) IsAdmin() bool {
	return startup.IsAdmin()
}

// ===== 配置 =====

// GetConfig 获取当前配置
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig 保存配置
func (a *App) SaveConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}
	return config.Save(cfg)
}

// ===== AI =====

// AIAnalyze 让 AI 分析内存（一次性，固定流程）
// 采集内存快照 + 自启项列表 一并喂给 AI
func (a *App) AIAnalyze() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("加载配置失败: %w", err)
	}
	if cfg.AIAPIKey == "" {
		return "", fmt.Errorf("尚未配置 AI API Key，请先在设置中填写")
	}

	report, err := memory.Analyze(cfg.TopN)
	if err != nil {
		return "", fmt.Errorf("采集内存失败: %w", err)
	}

	// 自启项采集失败不算致命，继续无自启项分析
	startupItems, _ := startup.EnumerateStartupItems()

	client := ai.NewClient(cfg.AIBaseURL, cfg.AIAPIKey, cfg.AIModel)
	userMsg := ai.BuildUserPrompt(report, startupItems)
	systemPrompt := buildFinalSystemPrompt(cfg)
	return client.Chat(a.ctx, systemPrompt, userMsg)
}

// AIChat 多轮对话接口
// messages: 完整对话历史（user/assistant 交替）
// 后端会自动 prepend system prompt（默认 + 用户自定义）
func (a *App) AIChat(messages []ai.Message) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("加载配置失败: %w", err)
	}
	if cfg.AIAPIKey == "" {
		return "", fmt.Errorf("尚未配置 AI API Key，请先在设置中填写")
	}
	if len(messages) == 0 {
		return "", fmt.Errorf("消息列表不能为空")
	}

	systemPrompt := buildFinalSystemPrompt(cfg)
	full := make([]ai.Message, 0, len(messages)+1)
	full = append(full, ai.Message{Role: "system", Content: systemPrompt})
	full = append(full, messages...)

	client := ai.NewClient(cfg.AIBaseURL, cfg.AIAPIKey, cfg.AIModel)
	return client.ChatMessages(a.ctx, full)
}

// GetMemorySnapshotPrompt 获取一份当前内存快照的 prompt 文本
// 前端可以把它作为一条 user 消息添加到对话开头，给 AI 提供上下文
func (a *App) GetMemorySnapshotPrompt() (string, error) {
	cfg, _ := config.Load()
	topN := 30
	if cfg != nil && cfg.TopN > 0 {
		topN = cfg.TopN
	}
	report, err := memory.Analyze(topN)
	if err != nil {
		return "", err
	}
	startupItems, _ := startup.EnumerateStartupItems()
	return ai.BuildUserPrompt(report, startupItems), nil
}

// GetDefaultSystemPrompt 暴露默认 system prompt（前端设置页参考）
func (a *App) GetDefaultSystemPrompt() string {
	return ai.SystemPrompt
}

// buildFinalSystemPrompt 拼接默认 system prompt 和用户自定义 extra prompt
func buildFinalSystemPrompt(cfg *config.Config) string {
	if cfg == nil || cfg.AIExtraPrompt == "" {
		return ai.SystemPrompt
	}
	return ai.SystemPrompt + "\n\n## 用户自定义指令\n" + cfg.AIExtraPrompt
}

// TestAIConnection 测试 AI 连接
func (a *App) TestAIConnection() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	if cfg.AIAPIKey == "" {
		return "", fmt.Errorf("API Key 未配置")
	}
	client := ai.NewClient(cfg.AIBaseURL, cfg.AIAPIKey, cfg.AIModel)
	return client.Chat(a.ctx, "你是一个简洁的助手。", "请用一句话回复：连接成功。")
}
