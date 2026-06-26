<template>
    <div class=" flex justify-center items-center h-full">

        <div class=" bg-white lg:w-[400px] w-full p-4 rounded-md shadow-lg">
            <div class=" flex gap-2 cursor-pointer justify-center">
                <div class=" border-b-2 text-xl" :class="{ active: mode == 'email' }" @click="mode = 'email'">流媒体 管理系统 
                    <span class=" text-sm text-red-900 ">v{{ version }}</span>
                </div>
            </div>
            <Login class="mt-8" :mode="mode" />
        </div>

        <div class=" fixed bottom-2 text-xl bg-yellow-50">
            用于海康、大华等监控设备，RTSP协议监控画面自动检测，接入守护，批量管理。
        </div>
    </div>
</template>
<script setup lang="ts">
import Login from "@/components/Login.vue"
import { onMounted, ref } from "vue";
import { version } from "../../package.json"
const mode = ref<"email" | "wechat">("email")
</script>
<style scoped lang="scss">
.active {
    @apply border-b-orange-600 font-semibold;
}
</style>