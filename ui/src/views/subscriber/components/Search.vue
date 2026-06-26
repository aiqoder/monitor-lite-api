<template>
    <div class=" border border-dashed rounded-md border-gray-400 p-2">
        <el-button plain @click="addSubscriber">
            <el-icon>
                <Plus />
            </el-icon>新增订阅
        </el-button>
    </div>
</template>
<script lang="ts" setup>
import { Plus } from "@element-plus/icons-vue"


const emit = defineEmits<{
    (e: 'add', args: any): void;
}>()

function addSubscriber() {
    emit("add", "")
}
</script>