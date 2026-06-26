<template>
    <ProForm :columns="formColumns" :span="24" :btn-position="'center'" :display-expend="false" hidden-reset
        ok-text="保存" />
</template>
<script setup lang="ts">
import { useRequestRuleGet } from "@/api/tv";
import { createProForm, ProForm } from "@howuse/element-plus-form"
import { onMounted, ref, watch } from "vue";
import { toLegacyMaps } from '@/utils/promptConfig'

let group = {}
const groupArray = ref<string[]>([])
const rightNameArray = ref<string[]>([])

const props = defineProps<{ data: [] }>()
const { execute: executeDoc, data: ruleData } = useRequestRuleGet()

const { formColumns, formModel, setValue } = createProForm([
    {
        key: "group",
        label: "分组",
        value: [],
        is: 'ElSelect',
        data: groupArray,
        optionLabel: "index",
        optionValue: "index",
        props: {
            multiple: true,
            onChange(e) {
                rightNameArray.value = []
                for (const groupName of e) {
                    rightNameArray.value.push(...(group[groupName]?.split("#") || []))
                }
                setValue("displayName", rightNameArray.value[0])
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
    },
    {
        key: "weight",
        label: "权重",
        value: "",
        is: 'ElInputNumber',
    },
])

watch(() => props.data, (val) => {
    formModel.value = {...val, group: !val.group ? [] :val.group.split(",")}
},
    {
        immediate: true,
    })

onMounted(async () => {
    await executeDoc({})
    const { group: g } = toLegacyMaps(ruleData.value ?? { groups: [] })
    group = g
    groupArray.value = Object.keys(group)
    setTimeout(() => {
        rightNameArray.value = group[formModel.value?.group[0]]?.split("#") || []
    }, 500)
})
</script>
