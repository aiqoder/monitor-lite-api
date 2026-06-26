<template>
    <ProForm :columns="formColumns" :span="24" :btn-position="'center'" :display-expend="false" hidden-reset
        ok-text="保存" />
</template>
<script setup lang="ts">
import { createProForm, ProForm } from "@howuse/element-plus-form"
import { syncRefs } from "@vueuse/core";
import { ref } from "vue";
const props = defineProps<{ data: {} }>()

const { formColumns, formModel } = createProForm([
    {
        key: "name",
        label: "订阅名称",
        value: "",
        is: 'ElInput',
        rules: [
            { required: true, message: "请输入" },
            { type: "string", min: 2, message: "至少输入两个字符" }
        ]
    },
    {
        key: "url",
        label: "链接",
        value: "",
        is: 'ElInput',
        rules: [
            { required: true, message: "请输入" },
        ]
    },
])

syncRefs(() => props.data, formModel)
</script>