package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client OpenAI 兼容协议客户端
type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	HTTP    *http.Client
}

// NewClient 创建客户端
func NewClient(baseURL, apiKey, model string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest OpenAI Chat Completions 请求体
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// chatResponse OpenAI Chat Completions 响应（只解析需要的字段）
type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Chat 调用聊天接口（非流式）
func (c *Client) Chat(ctx context.Context, system, user string) (string, error) {
	return c.ChatMessages(ctx, []Message{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	})
}

// ChatMessages 接收完整对话历史，支持多轮对话
func (c *Client) ChatMessages(ctx context.Context, messages []Message) (string, error) {
	if c.BaseURL == "" {
		return "", errors.New("AI BaseURL 未配置")
	}
	if c.APIKey == "" {
		return "", errors.New("AI API Key 未配置")
	}
	if c.Model == "" {
		return "", errors.New("AI Model 未配置")
	}
	if len(messages) == 0 {
		return "", errors.New("消息列表为空")
	}

	body := chatRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := c.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 500))
	}

	var cr chatResponse
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, 响应: %s", err, truncate(string(respBody), 300))
	}
	if cr.Error != nil {
		return "", fmt.Errorf("AI 返回错误: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return "", errors.New("AI 未返回任何结果")
	}

	return cr.Choices[0].Message.Content, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
