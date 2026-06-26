<template>
    <ProForm :columns="formColumns" :span="4" @ok="submit" :btn-position="'float-end'" ok-text="搜索">
        <template #suffix-btn>
            <el-button  @click="check" type="primary">检测搜索结果
            </el-button>
        </template>
    </ProForm>
</template>
<script lang="ts" setup>
import { createProForm, ProForm } from '@howuse/element-plus-form'
import { onMounted, ref } from 'vue';
import { useRequestRuleGet } from '@/api/tv';
import { unique } from 'howtools';
import { toLegacyMaps } from '@/utils/promptConfig'

const { execute: executeDoc, data: ruleData } = useRequestRuleGet()

const emit = defineEmits<{
    (e: 'search', args: any): void;
    (e: 'reset', args: any): void;
    (e: 'check'): void;
    (e: 'load'): void;
}>()

let group = {}
let topName = {}
const groupArray = ref<string[]>([])
const rightNameArray = ref<string[]>([])
const fakeNameArray = ref<string[]>([])

const { formColumns, setValue } = createProForm([
    {
        key: "group",
        label: "分组",
        value: "",
        is: 'ElSelect',
        data: groupArray,
        optionLabel: "index",
        optionValue: "index",
        props: {
            onChange(e) {
                rightNameArray.value = []
                rightNameArray.value.push(...(group[e]?.split("#") || []))
                const rName = rightNameArray.value[0]
                setValue("displayName", rName)

                fakeNameArray.value = unique(topName[rName]?.split("#") || [])
                setValue("fakeName", fakeNameArray.value[0])
            }
        }
    },
    {
        key: "displayName",
        label: "展示名称",
        value: "",
        is: 'ElSelect',
        data: rightNameArray,
        optionLabel: "index",
        optionValue: "index",
        props: {
            onChange(e) {
                fakeNameArray.value = unique(topName[e]?.split("#") || [])
                setValue("fakeName", fakeNameArray.value[0])
            }
        }
    },
    {
        key: "fakeName",
        label: "匹配名称",
        value: "",
        is: 'ElSelect',
        data: fakeNameArray,
        optionLabel: "index",
        optionValue: "index",
        props: {
            "filterable": true,
            "allow-create": true,
        }
    },
])

function submit(raw) {
    emit("search", raw)
}

onMounted(async () => {
    await executeDoc({})
    const maps = toLegacyMaps(ruleData.value ?? { groups: [] })
    group = maps.group
    topName = maps.name
    groupArray.value = Object.keys(group)
})

function check() {
    emit("check")
}
</script>
