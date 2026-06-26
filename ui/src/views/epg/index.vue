<template>
    <div>
        <el-alert :closable="false">
            <div>EPG接口（每日早8点定时采集）：{{ local }}/v1/epg/diyp</div>
        </el-alert>
    </div>
    <el-tabs type="border-card" @tab-change="tabChange" v-model="defaultTab">
        <el-tab-pane :label="group" :name="group" v-for="group in groupArray">
            <div class="h-[76vh] flex">
                <div class=" w-[256px] overflow-auto">
                    <div class="bg-[#f5f7fa] px-4 py-2 cursor-pointer"
                        :class="{ '!bg-white border-red-400 border-l-2': activeRightName == name }"
                        v-for="name in rightNameArray" @click="changeChannel(name)">
                        <div class=" inline-flex ">
                            <img :src="`/logo/${name}.png`" width="48px" height="48px" class=" bg-[#a855f7] rounded"/>
                            <span> {{ name }}</span>
                        </div>
                    </div>
                </div>
                <div class=" w-0 flex-1 px-3">
                    <div class=" flex justify-around  gap-4">
                        <div v-for="d in datetimeArray" class=" flex-1 w-0">
                            <div class=" bg-purple-500 text-white text-center px-2 py-1 rounded-md">
                                <div>{{ d.week }}</div>
                                <div>{{ d.date }}</div>
                                <!-- <div>{{ d.time }}</div> -->
                            </div>
                            <div class=" pt-4 flex flex-col gap-2  h-[65vh] overflow-auto">
                                <div v-for="epg in getEpgData(d)"
                                    class=" text-sm flex flex-col border-[1px] w-full p-1 rounded">
                                    <span class=" inline-flex gap-2 mb-1">
                                        <span>{{ epg.start }}</span>
                                        <span>{{ epg.title }}</span>
                                    </span>
                                    <div class="overflow-hidden text-ellipsis line-clamp-2 text-gray-500">描述：{{ epg.desc
                                        || "未知" }}</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </el-tab-pane>
    </el-tabs>
</template>

<script lang="ts" setup>
import { useRequestEpgList, useRequestRuleGet } from '@/api/tv';
import { dayjs } from 'element-plus';
import { onMounted, ref, unref } from 'vue';
import { toLegacyMaps } from '@/utils/promptConfig'
const local = location.origin
const { execute: executeDoc, data: ruleData } = useRequestRuleGet()
const { execute: executeEpgList, data: epgData } = useRequestEpgList()

let group = {}
let topName = {}
const groupArray = ref<string[]>([])
const rightNameArray = ref<string[]>([])
const activeRightName = ref("")
const defaultTab = ref("")
const datetimeArray = ref<{ date: string, time: string, week: string }[]>([])

onMounted(async () => {
    await executeDoc({})
    const maps = toLegacyMaps(ruleData.value ?? { groups: [] })
    group = maps.group
    topName = maps.name
    groupArray.value = Object.keys(group)
    tabChange(unref(groupArray)[0])

    // 生成一星期节目表日历
    const daysOfWeekCN = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'];
    for (const index of [-2, -1, 0, 1]) {
        datetimeArray.value.unshift({
            date: dayjs().add(index, 'd').format("YYYY-MM-DD"),
            time: dayjs().add(index, 'd').format("HH:mm:ss"),
            week: index == 0 ? "今天" : daysOfWeekCN[dayjs().subtract(index, 'd').day()],
        })
    }
})


function tabChange(tabName) {
    defaultTab.value = tabName
    setTimeout(() => {
        rightNameArray.value = group[tabName].split("#") || []
        activeRightName.value = unref(rightNameArray)[0]

        executeEpgList({
            params: {
                channel: activeRightName.value
            }
        })
    }, 500)
}

function changeChannel(name) {
    activeRightName.value = name
    executeEpgList({
        params: {
            channel: name
        }
    })
}

function getEpgData({ date }) {
    return epgData.value?.filter(item => item.date == date)
}
</script>

<style>
.demo-tabs>.el-tabs__content {
    padding: 32px;
    color: #6b778c;
    font-size: 32px;
    font-weight: 600;
}

.demo-tabs .custom-tabs-label .el-icon {
    vertical-align: middle;
}

.demo-tabs .custom-tabs-label span {
    vertical-align: middle;
    margin-left: 4px;
}
</style>