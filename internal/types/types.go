package types

type PlayReq struct {
	Url string `form:"url"`
}

type TvJsonReq struct {
	TvName string `form:"tvName"`
	Mode   string `form:"mode,options=so|me|re"` // 搜索模式 so 表示从池子搜索 me 表示搜索自己的源 re 表示正则表达式搜素
}

type TvPageReq struct {
	Current     int64  `form:"current"`
	Size        int64  `form:"size"`
	Group       string `form:"group,optional"`
	DisplayName string `form:"displayName,optional"`
	Name        string `form:"name,optional"`
	Url         string `form:"url,optional"`
	Width       string `form:"width,optional"`
	Height      string `form:"height,optional"`
}

type TvListPageResp struct {
	Current int64        `json:"current"`
	Size    int64        `json:"size"`
	Total   int64        `json:"total"`
	Records []TvListResp `json:"records"`
}

type TvListResp struct {
	ID          uint64 `json:"id,optional"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Width       uint64 `json:"width,optional"`
	Height      uint64 `json:"height,optional"`
	Speed       uint64 `json:"speed,optional"`
	FailCount   uint64 `json:"failCount"`
	UpdateTime  string `json:"updateTime"`
	DisplayName string `json:"displayName"`
	Group       string `json:"group"`
	Weight      int64  `json:"weight"` // 用于输出频道排序
}

type Tv struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Group string `json:"group"` // 分组 - 根据规则引擎进行解析
}

type BatchTvDelReq struct {
	Ids []string `json:"ids"`
}

type BatchTvUpdateReq struct {
	Tvs []TvUpdateReq `json:"tvs"`
}

type TvUpdateReq struct {
	ID          uint64 `json:"id,optional"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Width       uint64 `json:"width,optional"`
	Height      uint64 `json:"height,optional"`
	Speed       uint64 `json:"speed,optional"`
	DisplayName string `json:"displayName,optional"`
	Group       string `json:"group,optional"`
	Weight      int64  `json:"weight,optional"`
	Fail        bool   `json:"fail"`
}

type TvCheckReq struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type RuleGroup struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
}

type RuleGetResp struct {
	Groups []RuleGroup `json:"groups"`
}

type RuleUpdateReq struct {
	Groups []RuleGroup `json:"groups"`
}

type IdentifyReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type IdentifyResp struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type SuperTvReq struct {
}

type SuperWatchTvReq struct {
	Path string `uri:"path"`
}

type CheckAllReq struct {
	Type  string `json:"type,options=fail|pix0|all|select"`
	Extra string `json:"extra,optional"`
}

type SyncTVDataReq struct {
	Host              string `form:"host,optional"`
	Port              string `form:"port,optional"`
	User              string `form:"user"`
	Password          string `form:"password"`
	Database          string `form:"database"`
	Table             string `form:"table"`
	InsertBeforeClear bool   `form:"insertBeforeClear"` // 插入数据清空原表
	InsertJson        string `form:"insertJson"`
}

type Setting struct {
	ID    uint64 `json:"id"`
	Key   string `json:"key"`   // 原始名称
	Value string `json:"value"` // 值
	Type  string `json:"type"`
}

type FindReq struct {
}

type FindResp struct {
	Settings []Setting `json:"settings"`
}

type UpdateSettingReq struct {
	Settings []Setting `json:"settings"`
}

type UpdateSettingResp struct {
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResp struct {
}

type AiModelsReq struct {
	BaseURL string `form:"baseUrl,optional"`
	ApiKey  string `form:"apiKey,optional"`
}

type AiModelsResp struct {
	Models []string `json:"models"`
}

type FengShowReq struct {
	T string `form:"t,options=zx|zw|xg"`
}

type PtbTvReq struct {
	T string `form:"t,options=pt1|pt2|xy"`
}

type IQiLuReq struct {
	T string `form:"t,options=xw|ql|ty|sh|nk|wl|sr|ws|zy"`
}

type Subscriber struct {
	ID        uint64 `json:"id,optional"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	Count     uint64 `json:"count,optional"`
	CheckTime string `json:"checkTime,optional"`
}

type SubscriberDeleteReq struct {
	ID uint64 `uri:"id"`
}

type EpgListReq struct {
	Channel string `form:"channel"`
}

type EpgData struct {
	Channel string `json:"channel"` // 此处的channel是EpgChannel当中的ID
	Start   string `json:"start"`
	Stop    string `json:"stop"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Date    string `json:"date"`
}

type EpgSearchReq struct {
	Date string `form:"date,optional"`
	Ch   string `form:"ch,optional"`
}

type EpgSearchReply struct {
	ChannelName string    `json:"channel_name"`
	Date        string    `json:"date"`
	EpgData     []EpgData `json:"epg_data"`
}

type CustomOutReq struct {
	Path string `uri:"path"`
}

type AddSelfoutReq struct {
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Speed  int64  `json:"speed"`
	Key    string `json:"key"`
}

type AddSelfoutResp struct {
}

type DelSelfoutReq struct {
	Id int64 `form:"id"`
}

type DelSelfoutResp struct {
}

type UpdateSelfoutReq struct {
	Id     int64  `json:"id"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Speed  int64  `json:"speed"`
	Key    string `json:"key"`
}

type UpdateSelfoutResp struct {
}

type GetSelfoutByIdReq struct {
	Id int64 `form:"id"`
}

type GetSelfoutByIdResp struct {
	Id     int64  `json:"id"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Speed  int64  `json:"speed"`
	Key    string `json:"key"`
}

type SearchSelfoutReq struct {
	Current int64  `form:"current"`
	Size    int64  `form:"size"`
	Id      int64  `form:"id,optional"`
	Width   int64  `form:"width,optional"`
	Height  int64  `form:"height,optional"`
	Speed   int64  `form:"speed,optional"`
	Key     string `form:"key,optional"`
}

type SearchSelfoutResp struct {
	Current int64       `json:"current"`
	Size    int64       `json:"size"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"records"`
}
