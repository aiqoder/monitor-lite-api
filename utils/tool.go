package utils

import (
	"strings"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
)

type STv struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Search 远程源池已移除，保留接口供本地搜索逻辑调用。
func Search(key, t string) ([]STv, error) {
	return nil, nil
}

// 中国禁止出现的文字
func ChinaDisabledStr(word string) bool {
	taiWan := []string{"台湾", "台视", "民视", "公视", "中视", "中天", "寰宇", "纬来", "八大", "龙华", "靖天", "靖洋", "东森", "博斯", "阿里郎"}
	hongkong := []string{"香港", "凤凰", "翡翠", "tvb", "亚视", "美亚", "有线直播", "天映", "RHK32", "创世电视台", "奇妙电视"}
	macao := []string{"澳门", "濠江", "莲花"}

	allWords := [][]string{taiWan, hongkong, macao}
	for i := 0; i < len(allWords); i++ {
		aw := allWords[i]
		for j := 0; j < len(aw); j++ {
			a := aw[j]
			if strings.Contains(strings.ToLower(word), a) {
				return true
			}
		}
	}

	if tools.IsJapanString(word) {
		return true
	}

	return false
}
