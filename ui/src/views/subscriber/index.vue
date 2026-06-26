<template>
    <div class=" mb-2">
        <Search @add="addSubscriber" />
    </div>
    <div>
        <el-table :data="data" border height="75vh">
            <el-table-column prop="name" label="名称" width="180" />
            <el-table-column prop="url" label="链接" />
            <el-table-column prop="count" label="读取数量" width="180" />
            <el-table-column prop="checkTime" label="读取时间" width="180" />
            <el-table-column prop="action" label="操作" width="270" fixed="right">
                <template #default="{ row }">
                    <el-popconfirm title="是否立即读取，此操作无法撤销" @confirm="confirmGrab(row)">
                        <template #reference><el-button type="primary">立即读取</el-button></template>
                    </el-popconfirm>

                    <el-button type="primary" @click="handleEdit(row)">编辑</el-button>
                    <el-button type="danger" @click="batchDelete(row)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
    </div>
    <ProFormDialog v-model:visible="editVisible" title="编辑">
        <Edit :data="currentEdit" @ok="saveEdit" />
    </ProFormDialog>
</template>
<script lang="ts" setup>
import { useRequestSubscriberDel, useRequestSubscriberGrab, useRequestSubscribers, useRequestSubscriberUpdate } from '@/api/tv';
import { onMounted, ref } from 'vue';
import Edit from "./components/Edit.vue"
import Search from "./components/Search.vue"
import { ProFormDialog } from "@howuse/element-plus-form"
import { useToggle } from '@vueuse/core';
import { unref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const { execute: search, data } = useRequestSubscribers()
const { execute: executeUpdate } = useRequestSubscriberUpdate()
const { execute: executeDel } = useRequestSubscriberDel()
const { execute: executeGrab } = useRequestSubscriberGrab()

onMounted(() => {
    search({})
})

const [editVisible, toggleEdit] = useToggle()

const currentEdit = ref({})
function handleEdit(row: any) {
    currentEdit.value = row
    toggleEdit(true)
}

function saveEdit(raw) {
    executeUpdate({ data: { url: raw.url, name: raw.name, id: unref(currentEdit).id } }).then(() => {
        search()
        ElMessage.success('保存成功')
        toggleEdit(false)
    })

}

function batchDelete(raw) {
    ElMessageBox.confirm("确认要删除吗？", "警告", {
        type: "error"
    }).then(async () => {
        await executeDel({
            path: {
                id: raw.id,
            }
        })
        search()
    })
}

async function confirmGrab(raw) {
    await executeGrab({
        path: {
            id: raw.id,
        }
    })
    search()
}

function addSubscriber() {
    toggleEdit(true)
}
</script>