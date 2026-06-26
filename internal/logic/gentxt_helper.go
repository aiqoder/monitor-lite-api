package logic

import (
	"strings"
	"github.com/aiqoder/monitor-lite-api/model"
)

func resolveGroupChannelNames(configured string, tvs []model.Tv) []string {
	if strings.TrimSpace(configured) != "" {
		return strings.Split(configured, "#")
	}
	seen := make(map[string]struct{})
	var names []string
	for _, tv := range tvs {
		name := tv.DisplayName
		if name == "" {
			name = tv.Name
		}
		if _, ok := seen[name]; ok || name == "" {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}
