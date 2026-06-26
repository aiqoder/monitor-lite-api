<template>
    <div class=" flex flex-col h-full">
        <div class=" flex justify-between bg-white px-4 items-center">
            <div class=" flex-1 w-0 text-xl font-semibold">
                流媒体 管理系统
                <span class=" text-sm text-red-900 ">v{{ version }}</span>
            </div>
            <div>
                <el-dropdown @command="handleCommand">
                    <span class="el-dropdown-link">
                        <el-avatar :size="45"
                            :src="'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item command="password">修改密码</el-dropdown-item>
                            <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
        </div>
        <div class=" flex-1 h-0 overflow-auto p-4">
            <div class=" flex gap-2 h-[calc(100vh-80px)]">
                <div class=" w-[200px] h-full bg-white" :style="{ 'background-image': `url(${bg})` }">
                    <el-menu ellipsis :default-active="activeIndex" :popper-offset="16" router @select="selectMenu">
                        <el-menu-item index="/tvlist">媒体维护</el-menu-item>
                        <el-menu-item index="/online">监控画面</el-menu-item>
                        <!-- <el-menu-item index="/epg">节目单</el-menu-item> -->
                        <el-menu-item index="/rule">频道分组</el-menu-item>
                        <!-- <el-menu-item index="/search">找媒体助手</el-menu-item> -->
                        <el-menu-item index="/subscriber">订阅池</el-menu-item>
                        <el-menu-item index="/setting">系统设置</el-menu-item>
                    </el-menu>
                </div>
                <el-card class=" flex-1 w-0">
                    <router-view></router-view>
                </el-card>
            </div>

        </div>
    </div>
</template>
<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import bg from "@/assets/sider-bg.png"
const version = import.meta.env.VITE_APP_VERSION || 'dev'

const router = useRouter()
const activeIndex = useSessionStorage("active-index", "/tvlist")

function selectMenu(menu) {
    activeIndex.value = menu
}

function handleCommand(cmd) {
    if (cmd === "password") {
        activeIndex.value = "/setting"
        router.push({ path: "/setting", query: { tab: "security" } })
        return
    }
    if (cmd === "logout") {
        sessionStorage.clear()
        ElMessage.success("已退出登录")
        router.replace("/login")
    }
}
</script>