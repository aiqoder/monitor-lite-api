package proxy

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type IQiLuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIQiLuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IQiLuLogic {
	return &IQiLuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

func (l *IQiLuLogic) aesEncrypt(word, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()
	origData := pkcs7Padding([]byte(word), blockSize)

	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted)
}

func (l *IQiLuLogic) aesDecrypt(crypted, key string) (string, error) {
	decodeString, err2 := base64.StdEncoding.DecodeString(crypted)

	if err2 != nil {
		return "", err2
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeString))
	blockMode.CryptBlocks(origData, decodeString)
	origData = pkcs7UnPadding(origData)
	return string(origData), nil
}

func (l *IQiLuLogic) IQiLu(req *types.IQiLuReq) (resp string, err error) {
	// 生成直播链接
	vid := map[string]string{
		"xw": "24602", //  新闻频道 0
		"ql": "24584", //  山东齐鲁频道
		"ty": "24587", //  体育休闲频道
		"sh": "24596", //  生活频道
		"zy": "24593", //  综艺频道
		"nk": "24599", // 农科频道
		"wl": "24590", // 文旅频道
		"sr": "24605", // 少儿频道 0
		"ws": "24581", // 山东卫视 0
	}
	t := strconv.FormatInt(time.Now().UnixMilli(), 10)
	key := "k5x99e1mswelc4vt"
	target := vid[req.T]
	post, err := resty.New().R().
		SetQueryParams(map[string]string{
			"t": t,
			"s": cryptor.Md5String(target + t + "1qkhcc7og9zeftsu"),
		}).
		SetBody(
			l.aesEncrypt(fmt.Sprintf(`{"channelMark":"%s"}`, target), key),
		).
		Post("https://feiying.litenews.cn/api/v1/auth/exchange")

	if err != nil {
		log.Error(err)
		return "", err
	}

	// 解密请求结果
	decrypt, err := l.aesDecrypt(post.String(), key)

	if err != nil {
		log.Error(err)
		return "", err
	}

	// 获得要播放的URL地址
	result := struct {
		Code uint64 `json:"code"`
		Data string `json:"data"`
	}{}

	json.Unmarshal([]byte(decrypt), &result)

	// 向 https://datacenter.live.qcloud.com/ 发请求，激活播放链接
	// 40102 山东卫视
	token, _ := random.UUIdV4()
	devUUID, _ := random.UUIdV4()

	raw := map[string]any{}

	if req.T == "ws" {
		raw = map[string]any{
			"app_id":  0,
			"command": 40102,
			"data": []map[string]any{
				{
					"str_user_id":            devUUID,
					"dev_uuid":               devUUID,
					"str_session_id":         token,
					"bytes_token":            token,
					"str_device_type":        "",
					"str_os_info":            "win",
					"str_package_name":       "",
					"u32_network_type":       "",
					"u32_server_ip":          "",
					"str_stream_url":         result.Data,
					"u64_timestamp":          time.Now().UnixMilli(),
					"u32_link_type":          1,
					"u32_channel_type":       1,
					"str_app_version":        "",
					"platform":               3,
					"uint32_platform":        3,
					"str_browser_version":    "129.0.0.0",
					"str_browser_model":      "chrome",
					"str_user_agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
					"u32_video_drop":         "",
					"u32_drop_usage":         "",
					"float64_rtt":            "",
					"uint64_data_time":       time.Now().Unix(),
					"u32_avg_cpu_usage":      "",
					"u32_avg_memory":         "",
					"u32_result":             14250,
					"u32_speed_cnt":          0,
					"u32_avg_cache_time":     4109.75,
					"u32_max_load":           nil,
					"u32_audio_block_time":   0,
					"u32_avg_load":           0,
					"u32_load_cnt":           0,
					"u32_nodata_cnt":         0,
					"u32_first_i_frame":      2908.5,
					"u32_video_avg_fps":      nil,
					"u32_avg_block_time":     0,
					"u64_block_count":        0,
					"u32_video_block_time":   0,
					"u64_jitter_cache_max":   0,
					"u64_block_duration_max": nil,
					"u64_jitter_cache_avg":   0,
					"u32_ip_count_quic":      "",
					"u32_connect_count_quic": "",
					"u32_connect_count_tcp":  "",
					"u32_is_real_time":       "",
					"u32_first_frame_black":  "",
					"u32_delay_report":       0,
				},
			},
			"module_id": 1005,
		}
	} else {
		raw = map[string]any{
			"app_id":  0,
			"command": 40100,
			"data": []map[string]any{
				{
					"str_user_id":                 devUUID,
					"dev_uuid":                    devUUID,
					"str_session_id":              token,
					"bytes_token":                 token,
					"str_device_type":             "",
					"str_os_info":                 "win",
					"str_package_name":            "",
					"u32_network_type":            "",
					"u32_server_ip":               "",
					"str_stream_url":              result.Data,
					"u32_link_type":               1,
					"u32_channel_type":            1,
					"str_app_version":             "",
					"platform":                    3,
					"uint32_platform":             3,
					"str_browser_version":         "129.0.0.0",
					"str_browser_model":           "chrome",
					"str_user_agent":              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
					"u32_video_drop":              "",
					"u32_drop_usage":              "",
					"float64_rtt":                 "",
					"u32_cpu_usage":               "",
					"u32_app_cpu_usage":           "",
					"u32_avg_memory":              "",
					"u32_avg_cpu_usage":           "",
					"uint64_data_time":            time.Now().Unix(),
					"u32_recv_av_diff_time":       0,
					"u32_play_av_diff_time":       0,
					"u64_playtime":                4059,
					"u32_audio_decode_type":       2,
					"u32_audio_block_count":       0,
					"u32_audio_cache_time":        16277,
					"u32_audio_drop":              "",
					"u32_video_decode_type":       0,
					"u32_video_recv_fps":          126,
					"u32_fps":                     126,
					"u32_video_cache_time":        16277,
					"u32_avg_cache_count":         0,
					"u32_video_block_count":       0,
					"u32_avg_net_speed":           "",
					"u32_video_light_block_count": 0,
					"u32_video_large_block_count": 0,
					"u32_audio_jitter_60ms_count": 0,
					"u32_video_decode_fail":       "",
					"u32_audio_decode_fail":       "",
					"u32_avg_video_bitrate":       0,
					"u32_avg_audio_bitrate":       0,
					"u32_block_usage":             0,
				},
			},
			"module_id": 1005,
		}
	}

	response, err := resty.New().R().
		SetHeaders(map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
			"referer":    "https://v.iqilu.com/",
			"origin":     "https://v.iqilu.com",
		}).
		SetBody(raw).Post("https://datacenter.live.qcloud.com")

	if err != nil {
		return "", err
	}

	fmt.Println(result.Data, response.String())
	time.Sleep(time.Second)
	return result.Data, nil
}
