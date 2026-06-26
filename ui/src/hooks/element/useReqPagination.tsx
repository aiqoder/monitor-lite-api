import type {HowAxiosRequestConfig, HowExRequestOptions} from "@howuse/axios";
import {useDefAxiosRequest} from "@/utils/http";
import {usePagination} from "./usePagination";
import {computed, ComputedRef, MaybeRef, unref, watch, type Component} from "vue";
import {ElPagination} from "element-plus";
import {isDefined} from "@vueuse/shared";
import {merge} from "lodash-es";
import 'element-plus/theme-chalk/src/pagination.scss';
export type Page<T> = {
    current: number; //-
    pages: number; //-
    records: T[];
    searchCount: boolean; //-
    size: number; //-
    total: number; //-
}

/**
 *用法
 * const {
 *     Pagination,
 *     search,
 *     data,
 * } = useReqPagination({});
 * <template>
 *     <Pagination />
 * </template>
 * @param conf
 * @param options
 */
export function useReqPagination<T extends <T>(ref: MaybeRef<T> | ComputedRef<T>) => T>(conf: HowAxiosRequestConfig, options?: HowExRequestOptions & {size?: number}) {
    const {data, execute} = useDefAxiosRequest<Page<T & {idx: number}>>(conf, options);
    const {changeSize, changeCurrent, size, current, total} = usePagination(options?.size);

    let lastConf: HowAxiosRequestConfig;

    function load() {
       return execute(merge({
            [conf.method?.toLocaleLowerCase() == "post" ? 'data': 'params']: {
                current: unref(current),
                size: unref(size),
            },
        }, lastConf));
    }

    watch(data, (data) => {
        current.value = isDefined(data?.current) ? data?.current :1;
        total.value = isDefined(data?.total) ? data?.total : 0;
        if(isDefined(data?.size)) {
            size.value = data?.size;
        }
        if(unref(data)?.records && unref(data)?.records.length) {
            unref(data)?.records.forEach((o, idx) => {
                o.idx = (current.value -1) * size.value + idx + 1;
            })
        }
    })

    function search(conf: HowAxiosRequestConfig) {
        lastConf = conf;
        current.value = 1;
        return load();
    }

    /**
     * 修改页面大小修改数据
     */
    function handleChangeSize(_size: number) {
        changeSize(_size);
        return load();
    }

    /**
     * 修改页面页码
     */
    function handleChangeCurrent(_current: number) {
        changeCurrent(_current);
        return load();
    }
    const Pagination: Component = {
        render() {
            return <div class=" w-full py-5 flex justify-center default">
                 <ElPagination total={unref(total)}  
                    layout="total,sizes, prev, pager, next, jumper, ->, slot" 
                    currentPage={unref(current)} 
                    pageSize={unref(size)} 
                    onUpdate:current-page={handleChangeCurrent} 
                    onUpdate:page-size={handleChangeSize}
                    pageSizes={ [10, 20, 30, 40, 50, 100, 500, 1000] }
                 />
            </div>
        }
    }

    const getQuery = () => lastConf


    return {
        load,
        search,
        data: computed<T[]>({
            get() {
                return unref(data)?.records as []
            },
            set(val) {
                if(isDefined(unref(data)?.records)) {
                    // @ts-ignore
                    unref(data).records = val
                }
            }
        }),
        handleChangeSize,
        handleChangeCurrent,
        size,
        current,
        total,
        Pagination,
        getQuery
    };
}
