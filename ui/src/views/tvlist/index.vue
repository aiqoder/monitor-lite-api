<template>
    <div class=" mb-2">
        <Search @search="searchHeader" @batch-delete="batchDelete" @empty-group="emptyGroup" @update-group="updateGroup"
            @check="checkAll" />
    </div>
    <div class=" h-[calc(100vh-320px)]">
        <el-auto-resizer>
            <template #default="{ height, width }">
                <el-table :data="data" border @selection-change="handleSelectionChange" :height="height"
                    ref="multipleTableRef">
                    <el-table-column type="selection" width="55" />
                    <el-table-column prop="name" label="原始名称" width="180" />
                    <el-table-column prop="displayName" label="展示名称" width="180" />
                    <el-table-column prop="url" label="链接" show-overflow-tooltip min-width="180" />
                    <el-table-column prop="pix" label="分辨率" width="180">
                        <template #default="{ row }">
                            {{ row.width }}X{{ row.height }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="speed" label="瞬时速度" width="100" />
                    <el-table-column prop="group" label="分组" width="100" show-overflow-tooltip />
                    <el-table-column prop="failCount" label="失败次数" width="100" />
                    <el-table-column prop="weight" label="权重" width="100" />
                    <el-table-column prop="updateTime" label="更新时间" width="180" />
                    <el-table-column prop="action" label="操作" width="100" fixed="right">
                        <template #default="{ row }">
                            <el-button type="primary" @click="handleEdit(row)">编辑</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </template>
        </el-auto-resizer>
        <Pagination />
    </div>
    <ProFormDialog v-model:visible="editVisible" title="编辑">
        <Edit :data="currentEdit" @ok="saveEdit" />
    </ProFormDialog>

    <el-dialog title="提示" v-model="allCheckVisible">
        <div>
            此过程无法终止，直到全部检查完，通常您不需要进行此操作，后台会自动进行。如果您是云服务器，可能无法正常运行该功能。<br />注意：需要安装ffmpeg<br />每个链接检测时间少则100ms，多则9秒，请自行估计时长
        </div>
        <template #footer>
            <el-button type="primary" @click="handleCheckAll('fail')">检查失败的媒体</el-button>
            <el-button type="primary" @click="handleCheckAll('pix0')">检查分辨率0X0的媒体</el-button>
            <el-button type="primary" @click="handleCheckAll('all')">检查所有媒体</el-button>
            <el-button type="primary" @click="handleCheckAll('select')">检查选中的媒体</el-button>
        </template>
    </el-dialog>
</template>
<script lang="ts" setup>
    import { useRequestCheckAll, useRequestEmptyGroup, useRequestTvBatchDelete, useRequestTvLoseEfficacy, useRequestTvPage, useRequestTvUpdate, useRequestUpdateGroup } from '@/api/tv';
    import { ref } from 'vue';
    import Edit from "./components/Edit.vue"
    import Search from "./components/Search.vue"
    import { ProFormDialog } from "@howuse/element-plus-form"
    import { useToggle } from '@vueuse/core';
    import { unref } from 'vue';
    import { ElMessage, ElMessageBox, type TableInstance } from 'element-plus';

    const { search, load, data, Pagination } = useRequestTvPage()
    const { execute: executeUpdate } = useRequestTvUpdate()
    const { execute: executeDel } = useRequestTvBatchDelete()
    const { execute: executeLoseEfficacy } = useRequestTvLoseEfficacy()
    const { execute: executeEmpty } = useRequestEmptyGroup()
    const { execute: executeUpdateGroup } = useRequestUpdateGroup()
    const { execute: executeCheckAll } = useRequestCheckAll()
    const [allCheckVisible, toggleAllCheck] = useToggle()

    function searchHeader(query) {
        search({ params: query })
    }

    const [editVisible, toggleEdit] = useToggle()

    const currentEdit = ref({})
    function handleEdit(row: any) {
        currentEdit.value = row
        toggleEdit(true)
    }

    function saveEdit(raw) {
        executeUpdate({ data: { ...unref(currentEdit), ...raw, group: raw.group.join(","), fail: false } }).then(() => {
            load()
            ElMessage.success('保存成功')
            toggleEdit(false)
        })

    }
    const multipleTableRef = ref<TableInstance>()
    const multipleSelection = ref([])
    const handleSelectionChange = (val) => {
        multipleSelection.value = val
    }

    function batchDelete(type: "select" | "clear" | "delloseEfficacy") {
        if (type == "select") {
            if (unref(multipleSelection).length == 0) {
                ElMessage.error("未选中任何数据")
                return
            }

            ElMessageBox.confirm("确认要删除吗？", "警告", {
                type: "error"
            }).then(async () => {
                const ids = unref(multipleSelection).map(item => item.id)
                await executeDel({
                    data: {
                        ids,
                    }
                })
                multipleSelection.value = []
                load()
            })
        } else if (type == "clear") {
            ElMessageBox.confirm("注意，此操作将会清空所有数据，确认要执行吗？", "警告", {
                type: "error"
            }).then(async () => {
                await executeDel({
                    data: {
                        ids: ["-1"],
                    }
                })
                multipleSelection.value = []
                load()
            })
        }else if(type == "delloseEfficacy") {
            ElMessageBox.confirm("注意，此操作清空失败次数大于0的数据，确认要执行吗？", "警告", {
                type: "error"
            }).then(async () => {
                await executeLoseEfficacy({
                    data: {
                  
                    }
                })
                multipleSelection.value = []
                load()
            })
        }
    }

    function emptyGroup() {
        ElMessageBox.confirm("确认要清空所有分组和展示名称吗？", "警告", {
            type: "error"
        }).then(async () => {
            executeEmpty().then(() => {
                load()
            })
        })
    }

    function updateGroup() {
        ElMessageBox.confirm("确认要进行 AI 分组吗？只会处理没有分组的数据", "提示", {
            type: "warning"
        }).then(async () => {
            await executeUpdateGroup()
            ElMessage.success("AI 分组完成")
            load()
        })
    }

    function checkAll() {
        toggleAllCheck()
    }

    async function handleCheckAll(type) {
        if (type == "select") {
            if (unref(multipleSelection).length == 0) {
                ElMessage.error("未选中任何数据")
                return
            }
            await executeCheckAll({ data: { type, extra: unref(multipleSelection).map(item => item.id).join(",") } })
        } else {
            await executeCheckAll({ data: { type } })
        }
        toggleAllCheck(false)
        ElMessage.success("检测命令已下发")
        load()
    }
</script>