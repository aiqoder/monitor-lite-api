package tools

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-resty/resty/v2"
	"strings"
	"github.com/aiqoder/monitor-lite-api/pkg/common/m3u"
)

type IpTvSourceTree struct {
	Group string                  `json:"group"`
	Tv    map[string][]IpTvSource // 缓存tv 名称 以及URL 用 # 分割
}

type IpTvSource struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Group string `json:"group"`
}

func IptvSourceToTree(source []IpTvSource) []IpTvSourceTree {
	var tree []IpTvSourceTree
	with := slice.GroupWith(source, func(item IpTvSource) string {
		return item.Group
	})

	for group, src := range with {
		withSource := slice.GroupWith(src, func(item IpTvSource) string {
			return item.Name
		})

		tv := map[string][]IpTvSource{}

		for name, tvs := range withSource {
			tv[name] = tvs
		}
		tree = append(tree, IpTvSourceTree{
			Group: group,
			Tv:    tv,
		})
	}

	return tree
}

// 聚合URL 用 # 分割
func IptvGatherUrl(inSource []IpTvSource) []IpTvSource {
	var source []IpTvSource

	// 聚合相同名称的链接
	with := slice.GroupWith(inSource, func(item IpTvSource) string {
		return item.Name
	})

	for name, src := range with {
		var urls []string
		for _, tv := range src {
			urls = append(urls, tv.Url)
		}
		source = append(source, IpTvSource{
			Group: src[0].Group,
			Name:  name,
			Url:   slice.Join(urls, "#"),
		})
	}
	return source
}

// ParserIpTvSource 解析ip源 url 是一个链接或者一个文本
func ParserIpTvSource(url string) []IpTvSource {
	//if strings.HasPrefix(url, "https://raw.githubusercontent.com") {
	//	url = fmt.Sprintf("https://ghp.ci/%s", url)
	//}

	fileType := 0 // 1 标识m3u / m3u8 / 或者链接  2 表示 m3u8 纯文本 3 表示txt的URL 4. txt纯文本

	lines := strings.Split(url, "\n")
	if len(lines) == 0 {
		return nil
	}

	line := lines[0]
	line = strings.TrimSpace(line)
	line = strings.Replace(line, "\r", "", -1)

	if strings.HasPrefix(line, "http") {
		if strings.HasSuffix(url, ".m3u") || strings.HasSuffix(url, ".m3u8") {
			fileType = 1
		} else {
			fileType = 3
		}
	} else {
		tv := strings.Split(line, ",")
		if len(tv) == 2 && (strings.HasPrefix(tv[1], "http") || tv[1] == "#genre#") {
			fileType = 4
			// 如果包含 #EXTINF 表示m3u的纯文本
		} else if strings.HasPrefix(line, "#EXTM3U") || strings.HasPrefix(line, "#EXTINF") {
			fileType = 2
		}
	}

	var list []IpTvSource
	// 有逗号 或者 后缀是m3u 或者 m3u8 或者 包含 #EXTINF
	if fileType == 1 || fileType == 2 {
		playlist, err := m3u.Parse(url)
		if err == nil {
			for i := 0; i < len(playlist.Tracks); i++ {
				tv := playlist.Tracks[i]
				list = append(list, IpTvSource{
					Name:  tv.Name,
					Url:   strings.ReplaceAll(tv.URI, "\r", ""),
					Group: "",
				})
			}

			return list
		}
	} else if fileType == 3 {
		get, err := resty.New().R().Get(url)

		if err != nil {
			list = append(list, IpTvSource{
				Name:  "无法访问您提供的连接",
				Url:   url,
				Group: "未知的分组",
			})
			return list
		}

		split := strings.Split(get.String(), "\n")

		tvGroup := ""
		for i := 0; i < len(split); i++ {
			raw := strings.Split(split[i], ",")
			if len(raw) < 2 {
				continue
			}

			tvName := raw[0]
			tvUrl := raw[1]

			if strings.Contains(tvUrl, "genre") {
				tvGroup = tvName
				continue
			}

			if StrInContains(tvName, []string{"公众号", "下载", "免费", "关注", "免费", "软件", "退款", "受骗", "更新", "群"}) {
				continue
			}

			if len(tvName) > 0 && len(tvUrl) > 0 {
				if strings.Contains(tvUrl, "#") {
					tvUrls := strings.Split(tvUrl, "#")
					for i := 0; i < len(tvUrls); i++ {
						u := tvUrls[i]

						if StrInContains(strings.ToLower(u), []string{"p2p://", "p3p://", "p8p://", "mitv://", "tvbus://"}) {
							continue
						}

						list = append(list, IpTvSource{
							Name:  tvName,
							Url:   strings.ReplaceAll(u, "\r", ""),
							Group: tvGroup,
						})
					}
				} else {
					if StrInContains(strings.ToLower(tvUrl), []string{"p2p://", "p3p://", "p8p://", "mitv://", "tvbus://"}) {
						continue
					}
					list = append(list, IpTvSource{
						Name:  tvName,
						Url:   strings.ReplaceAll(tvUrl, "\r", ""),
						Group: tvGroup,
					})
				}
			}
		}

		return list
	} else if fileType == 4 {
		split := strings.Split(url, "\n")

		tvGroup := ""
		for i := 0; i < len(split); i++ {
			raw := strings.Split(split[i], ",")
			if len(raw) < 2 {
				continue
			}

			tvName := raw[0]
			tvUrl := raw[1]

			if strings.Contains(tvUrl, "genre") {
				tvGroup = tvName
				continue
			}

			if StrInContains(tvName, []string{"公众号", "下载", "免费", "关注", "免费", "软件", "退款", "受骗", "更新", "群"}) {
				continue
			}

			list = append(list, IpTvSource{
				Name:  tvName,
				Url:   strings.ReplaceAll(tvUrl, "\r", ""),
				Group: tvGroup,
			})
		}
	}

	return list
}
