import { type HowAxiosRequestConfig, type HowExRequestOptions, createAxios } from "@howuse/axios"
import type { AxiosResponse, AxiosError } from "axios";
import { ElLoading, ElMessage } from "element-plus";

export interface DefResponse<T = any> {
  code: 0 | 1 | 500 | 200 | 1001;
  msg: string;
  data?: T;
  success: boolean,
}

// 创建实例
const { server: instance, useAxiosRequest, useBlobDownload } = createAxios({
  timeout: 1000 * 60 * 2,
  headers: {
    "Content-Type": "application/json",
  },
});
// 请求拦截
instance.interceptors.request.use(
  (config: any) => {
    config.headers["session-id"] = sessionStorage.getItem("auth")
    // 检测链接不进行loading
    if (config.url.includes("/check")) {
      return config;
    }
    ElLoading.service({ fullscreen: true, text: "程序正在拼命加载中...", background: "rgba(0, 0, 0, 0.8)" });
    return config;
  },
  (err) => {
    console.log(err);
  }
);
// 响应拦截
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    ElLoading.service().close();
    return response
  },
  (error: AxiosError) => {
    ElLoading.service().close();
    if (error.code == 'ERR_BAD_RESPONSE') {
      ElMessage.error("您的服务器网络存在异常，请检查")
    } else {
      ElMessage.error(error.response?.data || "系统正在拼命中...")
      if (error.response?.status == 401) {
        location.href = "/admin/login"
      }
    }

    return error
  }
);

function useDefAxiosRequest<T = any>(config: HowAxiosRequestConfig, options?: HowExRequestOptions) {
  const {
    response,
    data,
    error,
    edata,
    execute,
    aborted,
    abort,
    finished: isFinished,
    loading: isLoading,
  } = useAxiosRequest<DefResponse<T>>(config, options)

  return {
    response,
    data,
    error,
    edata,
    execute: execute as unknown as (config: HowAxiosRequestConfig, options?: HowExRequestOptions) => Promise<DefResponse<T>>,
    aborted,
    abort,
    isFinished,
    isLoading
  }
}

export { useDefAxiosRequest, useBlobDownload, useAxiosRequest }
export default instance;
