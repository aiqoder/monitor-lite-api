import { useReqPagination } from "@/hooks/element/useReqPagination";
import { useAxiosRequest, useDefAxiosRequest } from "@/utils/http";

export function useRequestTvPage() {
    return useReqPagination({ url: "/v1/tv/page" })
}

export function useRequestPix() {
    return useDefAxiosRequest({ url: "/v1/tv/pix" }, { defaultVal: [] })
}

// axios.get(`${searchUrl.getUrl}/v1/tv/json`, { params: { tvName: value, mode: unref(searchMode) } })
export function useRequestTvJson() {
    return useDefAxiosRequest({ url: "/v1/tv/json" }, { defaultVal: [] })
}

export function useRequestTvCheck() {
    return useDefAxiosRequest({ url: "/v1/tv/check", method: "post" }, { defaultVal: [] })
}

export function useRequestTvUpdate() {
    return useDefAxiosRequest({ url: "/v1/tv/update", method: "post" })
}

export function useRequestTvBatchUpdate() {
    return useDefAxiosRequest({ url: "/v1/tv/batchupdate", method: "post" })
}

export function useRequestTvBatchDelete() {
    return useDefAxiosRequest({ url: "/v1/tv/batchdelete", method: "post" })
}


export function useRequestTvLoseEfficacy() {
    return useDefAxiosRequest({ url: "/v1/tv/deleteLoseEfficacy", method: "post" })
}

export interface RuleGroup {
    name: string
    channels: string[]
}

export interface RuleConfig {
    groups: RuleGroup[]
}

export function useRequestRuleGet() {
    return useAxiosRequest<RuleConfig>({ url: "/v1/tv/rule/get" }, { defaultVal: { groups: [] } })
}

export function useRequestRuleUpdate() {
    return useDefAxiosRequest<any>({ url: "/v1/tv/rule/update", method: "post" })
}

export function useRequestIdentify() {
    return useDefAxiosRequest<any>({ url: "/v1/tv/identify", method: "get" })
}

export function useRequestFind() {
    return useAxiosRequest<any>({ url: "/v1/setting/find", method: "get" })
}

export function useRequestUpdate() {
    return useAxiosRequest<any>({ url: "/v1/setting/update", method: "post" })
}

export function useRequestChangePassword() {
    return useDefAxiosRequest<any>({ url: "/v1/setting/changePassword", method: "post" })
}

export function useRequestAiModels() {
    return useAxiosRequest<{ models: string[] }>({ url: "/v1/setting/aiModels", method: "get" }, { defaultVal: { models: [] } })
}

export function useRequestEmptyGroup() {
    return useAxiosRequest<any>({ url: "/v1/tv/emptyGroup", method: "post" })
}

export function useRequestUpdateGroup() {
    return useAxiosRequest<any>({ url: "/v1/tv/updateGroup", method: "post" })
}

export function useRequestCheckAll() {
    return useAxiosRequest<any>({ url: "/v1/tv/checkAll", method: "post" })
}

// 订阅接口
export function useRequestSubscribers() {
    return useAxiosRequest<any>({ url: "/v1/subscriber/subscribers", method: "get" }, { defaultVal: [] })
}

export function useRequestSubscriberUpdate() {
    return useAxiosRequest<any>({ url: "/v1/subscriber/update", method: "post" })
}

export function useRequestSubscriberDel() {
    return useAxiosRequest<any>({ url: "/v1/subscriber/delete/{id}", method: "get" })
}

export function useRequestSubscriberGrab() {
    return useAxiosRequest<any>({ url: "/v1/subscriber/grab/{id}", method: "get" })
}

//EPG
export function useRequestEpgList() {
    return useAxiosRequest<any>({ url: "/v1/epg/epgList", method: "get" }, { defaultVal: [] })
}

// 查询播放列表
export function useVideoSuper() {
    return useAxiosRequest<string>({ url: "/v1/tv/super", method: "get" })
}