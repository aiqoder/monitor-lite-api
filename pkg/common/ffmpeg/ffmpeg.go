package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
)

const probeTimeout = 20 * time.Second

type VideoInfo struct {
	Height uint64 `json:"height"`
	Width  uint64 `json:"weight"`
	Fps    uint64 `json:"fps"`   // 码率
	Speed  uint64 `json:"speed"` // 加载毫秒数
	CodeC  string `json:"codeC"` // 解码方式
}

func GetOnlineVideoInfo(url string) (VideoInfo, error) {
	start := time.Now().UnixNano() / 1e6
	info := VideoInfo{}

	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner",
		"-rw_timeout", "15000000",
		"-probesize", "5000000",
		"-analyzeduration", "5000000",
		"-i", url,
	)

	var out bytes.Buffer
	cmd.Stderr = &out

	err := cmd.Run()
	data := out.String()

	if ctx.Err() == context.DeadlineExceeded {
		return info, errors.New("检测超时: " + url)
	}

	if len(data) == 0 {
		if err != nil {
			return info, errors.New("没有获得任何视频信息: " + url)
		}
		return info, errors.New("没有获得任何视频信息: " + url)
	}

	// 没有音频信息判断为无效的视频源
	if !strings.Contains(data, "Audio:") {
		return VideoInfo{}, errors.New("没有音频的源: " + data)
	}

	if strings.Contains(data, "Stream #0") {
		// 提取宽高
		reg := regexp.MustCompile(`, \d+x\d+`)
		size := reg.FindString(data)
		wh := strings.Split(size, "x")

		if size != "" {
			width, _ := convertor.ToInt(strings.TrimPrefix(wh[0], ", "))
			height, _ := convertor.ToInt(wh[1])
			info.Width = uint64(width)
			info.Height = uint64(height)
		}

		// 宽度是0，可能是广播
		if info.Width == 0 && strings.Contains(data, "Audio:") {
			info.Speed = uint64(time.Now().UnixNano()/1e6 - start)
			return VideoInfo{}, nil
		}
		// 忽略过低的分辨率视频
		if info.Width <= 320 {
			return VideoInfo{}, errors.New("视频分辨率太小：" + url)
		}

		// 提取FPS
		reg2 := regexp.MustCompile(`\d+ fps,`)
		fps := reg2.FindString(data)
		fps = strings.TrimSuffix(fps, " fps,")
		fpsNum, _ := convertor.ToInt(fps)
		info.Fps = uint64(fpsNum)

		// 提取Codec
		reg3 := regexp.MustCompile(`Video: [a-z0-9_]+\b`)
		codec := reg3.FindString(data)
		split := strings.Split(codec, " ")
		if len(split) > 1 {
			info.CodeC = split[1]
		}
	}
	info.Speed = uint64(time.Now().UnixNano()/1e6 - start)

	return info, nil
}
