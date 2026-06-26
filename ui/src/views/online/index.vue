<template>
    <div class="flex gap-2 h-[80vh]">
        <div class=" w-[280px] overflow-auto">
            <el-collapse>
                <el-collapse-item :title="`【${g}】`" :name="g" v-for="(arr, g) in groups" :key="g">
                    <div class=" round py-2 px-1 border-b  cursor-pointer flex justify-between" v-for="a in arr"
                        @click="handleDeboucePlay(a, g)"
                        :class="{ 'bg-blue-500 text-white': currentName == `${g}#${a.name}` }">
                        <span>{{ a.name }}</span>
                        <select v-model="a.select" @click.stop class=" w-12 bg-purple-400"
                            :onchange="(e) => handleDeboucePlay2(e, a, g)" size="small">
                            <option :value="u" v-for="(u, index) in a.url.split('#')">通道{{ index
                                }}</option>
                        </select>
                    </div>
                </el-collapse-item>
            </el-collapse>
        </div>
        <div class=" flex-1 w-0">
            <video ref="videoElement" class=" w-full object-fill" style="height: -webkit-fill-available;" controls
                autoplay></video>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useVideoSuper } from "@/api/tv";
import { isUrl } from "@/utils/file";
import mpegts from "th-mpegts.js"
import { ref, unref } from "vue";

// -----------------列表

type IGroup = {
    group: {
        name: string;
        url: string;
        select: string;
    }[]
}


const groups = ref<IGroup>()
const currentName = ref<string>("")
const { execute: executeVideoData, data: videoData } = useVideoSuper()
executeVideoData().then(res => {
    const manifest = res.data || ""

    const xgroup: any = {}

    const results = manifest?.split("\n") || [];
    let lastGroup: undefined | string = undefined // 用于记录导入的节目分组内容

    for (const item of results) {
        // 检查分组内容，并设置最后一次的分组名称
        if (item.trim().endsWith("#genre#")) {
            lastGroup = item.trim().split(",")[0]
            if (lastGroup && !xgroup[lastGroup]) xgroup[lastGroup] = []
        }
        const [name, url] = item.replace(/\r/g, "").split(",");
        if (!name || !url) continue;

        if (isUrl(url) && lastGroup) {
            xgroup[lastGroup].push({ name: name, url: url, select: url.split("#")[0] })
        }
    }
    groups.value = xgroup
})

// -------------------播放器
const videoElement = ref()
const mpegPlayer = ref()
function justifyPlayTime() {
    if (!document.hidden) {
        // 显示
        let video = unref(videoElement)
        let buffered = video.buffered
        if (buffered.length > 0) {
            let end = buffered.end(0)
            if (end - video.currentTime > 0.15) {
                video.currentTime = end - 0.1
            }
        }
    } else {
        // 隐藏
    }
}

function distroyMpegts() {
    if (mpegPlayer.value) {
        mpegPlayer.value.pause()
        mpegPlayer.value.unload()
        mpegPlayer.value.detachMediaElement()
        mpegPlayer.value.destroy()
        mpegPlayer.value = null

        document.removeEventListener('visibilitychange', justifyPlayTime)
    }
}

function play(url) {
    distroyMpegts()
    if (mpegts.getFeatureList().mseLivePlayback) {
        mpegPlayer.value = mpegts.createPlayer({
            type: 'flv',  // could also be mpegts, m2ts, flv
            isLive: true,
            url: `ws://${location.host}/v1/video/play?url=${encodeURIComponent(url)}`,
        });
        mpegPlayer.value.attachMediaElement(unref(videoElement));
        mpegPlayer.value.load();
        mpegPlayer.value.play();
    }
}

function handleDeboucePlay(a, g) {
    if (currentName.value == `${g}#${a.name}`) {
        return
    }

    currentName.value = `${g}#${a.name}`
    const urls = a.url.split("#")
    play(urls[0])
}

function handleDeboucePlay2(e, a, g) {
     currentName.value = `${g}#${a.name}`
    play(e.target.value)
}
</script>