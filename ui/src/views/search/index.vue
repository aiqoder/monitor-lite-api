<template>
    <div class=" mb-2">
        <Search @search="searchHeader" @check="checkTv" />
    </div>
    <div class=" h-[calc(100vh-200px)]">
        <el-auto-resizer>
            <template #default="{ height, width }">
                <div>
                    <el-table :data="dataSource" border stripe :height="height">
                        <el-table-column prop="name" label="节目名称" width="150px" />
                        <el-table-column prop="url" label="链接地址" />
                        <el-table-column prop="result" label="检测结果" width="150px" />
                    </el-table>
                </div>
            </template>
        </el-auto-resizer>
    </div>

</template>
<script setup lang="ts">
    import { ElAlert, ElMessage, ElNotification } from "element-plus";
    import Search from "./components/Search.vue"
    import { useRequestTvCheck, useRequestTvJson } from "@/api/tv";
    import { computed } from "vue";
    import { unref } from "vue";
    const { execute: executeSearch, data } = useRequestTvJson()
    const { execute: executeCheck } = useRequestTvCheck()
    const dataSource = computed({
        get() {
            return data.value?.data || []
        },

        set(v: any) {
            data.value = { data: v, code: 200, success: true, msg: "success" };
        }
    })


    function searchHeader(raw: any) {
        if(!raw.fakeName) {
            ElMessage.error("请输入或者选择匹配名称")
            return
        }

        executeSearch({ params: { tvName: raw.fakeName, mode: "re" } })
    }

    async function checkTv() {
        ElNotification.warning("检测中...，请勿刷新或者关闭页面")
        for (const [index, tv] of (unref(data).data || []).entries()) {
            ElMessage.warning(`正在检测第${index + 1}个节目，剩余${unref(data).data?.length - index}个`)
            dataSource.value[index].result = "检测中..."
            await executeCheck({
                data: {
                    url: tv.url,
                    name: tv.name,
                }
            }).then((res: any) => {
                if (res.status > 200) {
                    // ElMessage.error(res.response.data)
                    dataSource.value[index].result = "检测失败"
                    return;
                }
                dataSource.value[index].result = "检测成功，已入库"
            }).catch(() => {
                dataSource.value[index].result = "检测失败"
            })
        }
    }
</script>