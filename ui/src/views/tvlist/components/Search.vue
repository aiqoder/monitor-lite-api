<template>
    <ProForm :columns="formColumns" :span="5" @ok="submit" :btn-position="'float-end'">

    </ProForm>

    <div class=" border border-dashed rounded-md border-gray-400 p-2">
        <el-button plain @click="importM3u">
            <el-icon>
                <UploadFilled />
            </el-icon>文件导入
        </el-button>
        <el-button plain @click="openTxtVisible">
            <el-icon>
                <Document />
            </el-icon>文本导入
        </el-button>
        <el-dropdown class=" mx-2" @command="batchDelete">
            <el-button type="danger">
                <el-icon>
                    <DeleteFilled />
                </el-icon>删除
            </el-button>
            <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item command="select">删除选中</el-dropdown-item>
                    <el-dropdown-item command="delloseEfficacy">删除失效</el-dropdown-item>
                    <el-dropdown-item command="clear">清空所有</el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>
        <el-button @click="emptyGroup" type="danger">
            <el-icon>
                <DocumentDelete />
            </el-icon>清分组
        </el-button>
        <el-button @click="updateGroup" type="primary">
            <el-icon>
                <Refresh />
            </el-icon>AI 分组
        </el-button>
        <el-button @click="checkAll" type="primary">
            <el-icon>
                <Checked />
            </el-icon>检测
        </el-button>
    </div>

    <el-dialog title="文本导入" v-model="txtVisible">
        <el-form ref="ruleFormRef" :model="formModel" label-width="auto" status-icon>
            <el-form-item label="文本" prop="text" :rules="[
                { required: true, message: '请输入文本' }
            ]">
                <el-input placeholder="仅支持单行文本导入" spellcheck="false" type="textarea" v-model="formModel.text"
                    show-word-limit :autosize="{ minRows: 20, maxRows: 20 }" :resize="'none'"
                    :maxlength="50000"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="saveFile()">保存</el-button>
        </template>
    </el-dialog>
</template>
<script lang="ts" setup>
import { createProForm, ProForm } from '@howuse/element-plus-form'
import { UploadFilled, Document, DeleteFilled, Refresh, DocumentDelete, Checked } from "@element-plus/icons-vue"
import { loadLocalFile } from 'howtools';
import { readerHandleM3u, readerHandleTxt, readerSuffixM3u, readerSuffixTxt } from '@/utils/file';
import { useToggle } from '@vueuse/core';
import { onMounted, reactive, ref } from 'vue';
import { useRequestPix, useRequestTvBatchUpdate } from '@/api/tv';
import { ElMessage } from 'element-plus';

const { execute: addTv } = useRequestTvBatchUpdate() // 更新tv列表
const { execute: executePix, data: pixData } = useRequestPix()
executePix({})

const formModel = reactive({
    text: "",
    mode: "txt",
})

const [txtVisible, toggle] = useToggle()

const emit = defineEmits<{
    (e: 'search', args: any): void;
    (e: 'reset', args: any): void;
    (e: 'insert'): void;
    (e: 'load'): void;
    (e: 'batch-delete', args: "select" | "clear"): void;
    (e: 'empty-group'): void;
    (e: 'update-group'): void;
    (e: 'check'): void;
}>()

const { formColumns, formModel: formModel2 } = createProForm([
    {
        label: "原始名称",
        key: "name",
        value: "",
        is: "ElInput",
    },
    {
        label: "展示名称",
        key: "displayName",
        value: "",
        is: "ElInput",
    },
    {
        label: "分组",
        key: "group",
        value: "",
        is: "ElInput",
    },
    {
        key: "url",
        label: "链接",
        value: "",
        is: 'ElInput',
    },
    {
        key: "pix",
        label: "分辨率",
        value: [],
        is: 'ElSelect',
        data: pixData,
        optionLabel: 'index',
        optionValue: 'index',
        props: {
            clearable: true,
            multiple: true,
            "collapse-tags": true,
            "max-collapse-tags": "2",
        }
    },
])

function submit(raw) {
    const pix = raw.pix
    const widths = []
    const heights = []
    for (const p of pix) {
        const [w, h] = p.split("X")
        if (w) {
            widths.push(w)
            heights.push(h)
        }
    }
    emit("search", {
        ...raw,
        pix: undefined,
        width: widths.join(","),
        height: heights.join(",")
    })
}

onMounted(() => {
    submit(formModel2.value)
})

function importM3u() {
    loadLocalFile({ accept: ".m3u,.txt" }).then((files) => {
        const file = files[0]
        if (file.name.endsWith(".txt")) {
            readerSuffixTxt(file).then(res => {
                const txt = res.map(item => {
                    return [item.name, item.url].join(",")
                }).join("\n")

                formModel.text = txt
                formModel.mode = "txt"
                toggle(true)
            })
        } else {
            readerSuffixM3u(file).then(res => {
                const txt = res.map(item => {
                    return [item.name, item.url].join(",")
                }).join("\n")

                formModel.text = txt
                formModel.mode = "txt"
                toggle(true)
            })
        }
    })
}

function openTxtVisible() {
    toggle(true)
}

function batchDelete(t: "select" | "clear") {
    emit("batch-delete", t)
}

function emptyGroup() {
    emit("empty-group")
}

function updateGroup() {
    emit("update-group")
}

function checkAll() {
    emit("check")
}

async function saveFile() {
    const mode = formModel.mode
    ElMessage.warning({
        message: "导入数据中，请勿关闭页面...，否则数据将丢失",
        duration: 99999,
        type: "warning"
    })

    if (formModel.text.includes("#EXTINF:")) {
        formModel.mode = "m3u"
    } else {
        formModel.mode = "txt"
    }

    if (mode == 'txt') {
        const tvs = readerHandleTxt(formModel.text).map(item => { return { ...item, fail: false } })
        await addTv({ data: { tvs } })
        ElMessage.success(`导入${tvs.length}条数据`)
    }
    if (mode == 'm3u') {
        const tvs = readerHandleM3u(formModel.text).map(item => { return { ...item, fail: false } })
        await addTv({ data: { tvs } })
        ElMessage.success(`导入${tvs.length}条数据`)
    }

    setTimeout(() => {
        toggle(false)
        formModel.text = ""
        ElMessage.closeAll()
    }, 1000)
}
</script>