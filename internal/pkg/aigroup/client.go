package aigroup

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
)

type AISettings struct {
	BaseURL string
	APIKey  string
	Model   string
}

type ChannelInput struct {
	Index   int    `json:"index"`
	RawName string `json:"rawName"`
}

type ClassifyResult struct {
	Index       int    `json:"index"`
	DisplayName string `json:"displayName"`
	Group       string `json:"group"`
}

type Client struct {
	settings AISettings
	http     *resty.Client
}

func NewClient(settings AISettings) *Client {
	if settings.Model == "" {
		settings.Model = "gpt-4o-mini"
	}
	if settings.BaseURL == "" {
		settings.BaseURL = "https://api.openai.com/v1"
	}
	settings.BaseURL = strings.TrimRight(settings.BaseURL, "/")

	return &Client{
		settings: settings,
		http:     resty.New().SetTimeout(120 * time.Second),
	}
}

func (c *Client) Enabled() bool {
	return c.settings.APIKey != ""
}

func (c *Client) ClassifyBatch(cfg prompt.Config, channels []ChannelInput) ([]ClassifyResult, error) {
	if len(channels) == 0 {
		return nil, nil
	}

	channelsJSON, err := json.Marshal(channels)
	if err != nil {
		return nil, err
	}

	userPrompt := cfg.BuildUserPrompt(string(channelsJSON))

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type reqBody struct {
		Model          string    `json:"model"`
		Messages       []message `json:"messages"`
		Temperature    float64   `json:"temperature"`
		ResponseFormat *struct {
			Type string `json:"type"`
		} `json:"response_format,omitempty"`
	}

	body := reqBody{
		Model: c.settings.Model,
		Messages: []message{
			{Role: "system", Content: cfg.SystemPrompt()},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.1,
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
	}

	type chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	var resp chatResp
	r, err := c.http.R().
		SetHeader("Authorization", "Bearer "+c.settings.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post(c.settings.BaseURL + "/chat/completions")
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		msg := r.String()
		if resp.Error != nil && resp.Error.Message != "" {
			msg = resp.Error.Message
		}
		return nil, fmt.Errorf("AI 接口错误: %s", msg)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI 返回为空")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	results, err := parseResults(content)
	if err != nil {
		log.Error("AI 响应解析失败:", content, err)
		return nil, err
	}
	return results, nil
}

func parseResults(content string) ([]ClassifyResult, error) {
	var wrapper struct {
		Results []ClassifyResult `json:"results"`
	}
	if err := json.Unmarshal([]byte(content), &wrapper); err == nil && len(wrapper.Results) > 0 {
		return wrapper.Results, nil
	}

	var list []ClassifyResult
	if err := json.Unmarshal([]byte(content), &list); err == nil {
		return list, nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &obj); err != nil {
		return nil, fmt.Errorf("无法解析 AI 响应: %w", err)
	}
	for _, v := range obj {
		if err := json.Unmarshal(v, &list); err == nil && len(list) > 0 {
			return list, nil
		}
	}
	return nil, fmt.Errorf("无法解析 AI 响应 JSON")
}

const BatchSize = 25

func ClassifyAll(settings AISettings, cfg prompt.Config, rawNames []string) (map[int]ClassifyResult, error) {
	client := NewClient(settings)
	if !client.Enabled() {
		return nil, fmt.Errorf("未配置 AI API Key")
	}

	out := make(map[int]ClassifyResult, len(rawNames))
	for start := 0; start < len(rawNames); start += BatchSize {
		end := start + BatchSize
		if end > len(rawNames) {
			end = len(rawNames)
		}
		batch := make([]ChannelInput, 0, end-start)
		for i := start; i < end; i++ {
			batch = append(batch, ChannelInput{Index: i, RawName: rawNames[i]})
		}
		results, err := client.ClassifyBatch(cfg, batch)
		if err != nil {
			return out, err
		}
		for _, r := range results {
			out[r.Index] = r
		}
	}
	return out, nil
}
