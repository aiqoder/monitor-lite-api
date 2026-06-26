package aigroup

import (
	"fmt"
	"sort"
	"strings"
)

// ListModels 调用 OpenAI 兼容接口 GET /models 获取可用模型列表
func ListModels(settings AISettings) ([]string, error) {
	client := NewClient(settings)
	if !client.Enabled() {
		return nil, fmt.Errorf("未配置 API Key")
	}

	type modelsResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	var resp modelsResp
	r, err := client.http.R().
		SetHeader("Authorization", "Bearer "+client.settings.APIKey).
		SetResult(&resp).
		Get(client.settings.BaseURL + "/models")
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		msg := r.String()
		if resp.Error != nil && resp.Error.Message != "" {
			msg = resp.Error.Message
		}
		return nil, fmt.Errorf("获取模型列表失败: %s", msg)
	}

	seen := make(map[string]struct{})
	models := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		models = append(models, id)
	}
	sort.Strings(models)
	return models, nil
}
