package rules

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

type NewRule struct {
	Group map[string]string `yaml:"group"`
	Name  map[string]string `yaml:"name"`
}

// 新的RuleMeta结构体
type RuleMeta struct {
	FakeName   string         `json:"fakeName"`
	FakeNameRe *regexp.Regexp `json:"fakeNameRe"` // fakeName 对应的正则表达式
	RightName  string         `json:"rightName"`  // 正确的名称
	GroupName  string         `json:"groupName"`  // 分组名称
}

// rule yaml 规则文件内容
func NewRuleMetas(ruleString string) []RuleMeta {
	var rules []RuleMeta

	rule := NewRule{}
	// 读取规则文件
	err := yaml.Unmarshal([]byte(ruleString), &rule)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	for rightName, fakeName := range rule.Name {
		fakeNames := strings.Split(fakeName, "#")
		for i := 0; i < len(fakeNames); i++ {
			regular := fmt.Sprintf("%s", strings.ToLower(fakeNames[i]))
			re, err := regexp.Compile(regular)
			if err != nil {
				log.Error(err)
				continue
			}
			ruleMeta := RuleMeta{
				FakeName:   fakeNames[i],
				FakeNameRe: re,
				RightName:  rightName,
			}

		GroupLoop:
			for groupName, groupRightName := range rule.Group {
				rightNames := strings.Split(groupRightName, "#")
				for j := 0; j < len(rightNames); j++ {
					if rightNames[j] == rightName {
						ruleMeta.GroupName = groupName
						break GroupLoop
					}
				}
			}

			rules = append(rules, ruleMeta)
		}
	}

	return rules
}
