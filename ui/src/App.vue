<script setup lang="ts">
import { RouterView } from 'vue-router'
import zhCn from 'element-plus/dist/locale/zh-cn'
import { useIntervalFn, useWebSocket } from '@vueuse/core';
import { h, watchEffect } from 'vue';
import { ElMessageBox, ElNotification } from 'element-plus';
// 旧版本浏览器websocket必须写入全路径才可正确连接websocket
const { data } = useWebSocket(`ws://${location.host}/v1/tv/tips`, {
  autoReconnect: true,
})

const msgObj = {}
const msgObjTimeOut = {}

// 3s 没更新，表示没有继续推送消息立即删除
useIntervalFn(() => {
  for (const id in msgObj) {
    if (msgObjTimeOut[id] - Date.now() > 3000) {
      delete msgObj[id]
      delete msgObjTimeOut[id]
    }
  }
}, 3000)

watchEffect(() => {
  if (!data.value) return
  // 没有登录，不显示该消息弹出
  if(!sessionStorage.getItem("auth")) return
  const msg = data.value
  const [id, content] = msg.split("#")

  let displayTxt = "";
  msgObj[id] = content
  msgObjTimeOut[id] = Date.now()

  for (const id in msgObj) {
    displayTxt += (id + "：" + msgObj[id] + "\n")
  }

  const newDate = sessionStorage.getItem("msg-tip-time")
  if (newDate == null || new Date() > new Date(Number(newDate))) {
    const msg = document.getElementById("tip-message")
    if (msg) {
      msg.innerText = displayTxt
      return
    }
    const notice = ElNotification({
      title: '消息提示',
      message: h("div", { style: { color: 'red' }, id: "tip-message" }, displayTxt),
      showClose: false,
      position: 'bottom-right',
      duration: 0,

      onClick() {
        ElMessageBox({
          title: "是否10分钟后再提示",
          showClose: true,
          showCancelButton: true,
        }).then(() => {
          sessionStorage.setItem("msg-tip-time", Date.now() + 10 * 60 * 1000 + "")
          notice.close()
        })
      },
    })
  }
})
</script>

<template>
  <el-config-provider :locale="zhCn">
    <RouterView />
  </el-config-provider>
</template>