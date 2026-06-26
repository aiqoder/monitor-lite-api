package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type PtbtvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPtbtvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PtbtvLogic {
	return &PtbtvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 代理自 https://www.ptbtv.com/live/pt2t/
func (l *PtbtvLogic) Ptbtv(req *types.PtbTvReq) (resp string, err error) {
	vid := map[string]string{
		"pt1": "4",
		"pt2": "5",
		"xy":  "6",
	}
	// 组装请求头
	__Ox1895e := []string{"\x6C\x65\x6E\x67\x74\x68", "\x63\x68\x61\x72\x43\x6F\x64\x65\x41\x74", "", "\x30", "\x73\x75\x62\x73\x74\x72", "\x0A", "\x72\x65\x70\x6C\x61\x63\x65", "\x66\x72\x6F\x6D\x43\x68\x61\x72\x43\x6F\x64\x65", "\x74\x6F\x4C\x6F\x77\x65\x72\x43\x61\x73\x65", "\x65\x78\x74\x65\x6E\x64", "\x66\x33\x33\x62\x61\x31\x35\x65\x66\x66\x61\x35\x63\x31\x30\x65\x38\x37\x33\x62\x66\x33\x38\x34\x32\x61\x66\x62\x34\x36\x61\x36", "\x59\x57\x46\x6b\x5a\x44\x59\x77\x4d\x6a\x4e\x6b\x4e\x7a\x4d\x7a\x4e\x7a\x55\x77\x5a\x57\x4a\x6a\x59\x6a\x45\x34\x4e\x57\x46\x6a\x5a\x6a\x59\x33\x59\x6d\x51\x79\x59\x7a\x45\x3d", "\x31\x2E\x30\x2E\x30", "\x67\x65\x74\x54\x69\x6D\x65", "\x26", "\x6D\x64\x35"}
	_0x7523x47 := __Ox1895e[10]
	_0x7523x48 := __Ox1895e[11]
	_0x7523x49 := __Ox1895e[12]
	_0x7523x4a := time.Now().Unix()
	var _0x7523x4b = _0x7523x47 + __Ox1895e[14] + _0x7523x48 + __Ox1895e[14] + _0x7523x49 + __Ox1895e[14] + strconv.FormatInt(_0x7523x4a, 10)
	var _0x7523x4c = cryptor.Md5String(_0x7523x4b)

	// https://www.ptbtv.com/live/pt2t/
	// 4/5/6 莆田一台/二台/仙游电视
	bstrURL := fmt.Sprintf("https://www.ptbtv.com/m2o/channel/channel_info.php")

	type Raw struct {
		M3u8  string `json:"m3u8"`
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	var result []Raw

	var request = resty.New()

	var headers = map[string]string{
		"X-API-TIMESTAMP": strconv.FormatInt(_0x7523x4a, 10),
		"X-API-KEY":       _0x7523x47,
		"X-AUTH-TYPE":     __Ox1895e[15],
		"X-API-VERSION":   _0x7523x49,
		"X-API-SIGNATURE": _0x7523x4c,
		"User-Agent":      "Apifox/1.0.0 (https://apifox.com)",
		"Host":            "www.ptbtv.com",
	}
	fmt.Println(strconv.FormatInt(_0x7523x4a, 10), _0x7523x4c)
	response, err := request.R().SetResult(&result).ForceContentType("application/json").SetHeaders(headers).SetQueryParams(map[string]string{"channel_id": vid[req.T]}).Get(bstrURL)

	if err != nil {
		err := json.Unmarshal(response.Body(), &result)
		if err != nil {
			return "", err
		}
		return "", err
	}

	if len(result) > 0 {
		return result[0].M3u8, err
	}
	return "", errors.New("未获得任何节目")
}
