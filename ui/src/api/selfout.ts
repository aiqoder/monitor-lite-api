import { useReqPagination } from "@/hooks/element/useReqPagination";
import { useDefAxiosRequest } from "@/utils/http";

// 新增
export function useAddSelfout() {
    return useDefAxiosRequest({ url: "/v1/selfout/addSelfout", method: 'POST' });
}
// 更新
export function useUpdateSelfout() {
    return useDefAxiosRequest({ url: "/v1/selfout/updateSelfout", method: 'POST' });
}
// 删除
export function useDelSelfout() {
    return useDefAxiosRequest({ url: "/v1/selfout/delSelfout", method: 'DELETE' });
}
// 查询
export function useGetSelfoutById() {
    return useDefAxiosRequest({ url: "/v1/selfout/getSelfoutById", method: 'GET' });
}
// 分页查询
export function useSearchSelfout() {
    return useReqPagination({ url: "/v1/selfout/searchSelfout", method: 'GET' });
}
