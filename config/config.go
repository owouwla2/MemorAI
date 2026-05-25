package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Config 应用配置
type Config struct {
	// AI 配置（OpenAI 兼容协议）
	AIBaseURL string `json:"aiBaseURL"` // 例: https://api.openai.com/v1, https://api.deepseek.com/v1
	AIAPIKey  string `json:"aiApiKey"`
	AIModel   string `json:"aiModel"`   // 例: gpt-4o-mini, deepseek-chat
	// 自定义提示词：会追加到默认 system prompt 之后
	AIExtraPrompt string `json:"aiExtraPrompt"`
	// 采集参数
	TopN int `json:"topN"` // 单进程列表 Top N，默认 30
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		AIBaseURL:     "https://api.deepseek.com/v1",
		AIAPIKey:      "",
		AIModel:       "deepseek-chat",
		AIExtraPrompt: "",
		TopN:          30,
	}
}

var (
	mu        sync.Mutex
	cachePath string
)

// configPath 返回配置文件路径（与 exe 同目录的 config.json）
func configPath() (string, error) {
	if cachePath != "" {
		return cachePath, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	cachePath = filepath.Join(dir, "config.json")
	return cachePath, nil
}

// Load 加载配置；文件不存在时返回默认配置（不报错）
func Load() (*Config, error) {
	mu.Lock()
	defer mu.Unlock()

	path, err := configPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), err
	}

	// 容错：必填项空时回填默认
	if cfg.TopN <= 0 {
		cfg.TopN = 30
	}
	return cfg, nil
}

// Save 保存配置
func Save(cfg *Config) error {
	mu.Lock()
	defer mu.Unlock()

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
