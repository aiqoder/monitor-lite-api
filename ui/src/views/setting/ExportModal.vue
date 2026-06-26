<template>
  <div>
    <ProForm :columns="formColumns" :span="12" hidden-reset :display-expend="false" :btn-position="'center'"
      :ok-text="'保存'" @ok="save"></ProForm>
  </div>
  <div>
    <el-divider content-position="left">已创建链接</el-divider>

    <el-table :data="linkData">
      <el-table-column label="链接" prop="key">
          <template #default="{ row }">
            <a :href="`/cus/${row.key}`" target="_blank">/cus/{{ row.key }}</a>
          </template>
      </el-table-column>
      <el-table-column label="宽度" prop="width"></el-table-column>
      <el-table-column label="高度" prop="height"></el-table-column>
      <el-table-column label="响应速率" prop="speed"></el-table-column>
      <el-table-column label="操作" prop="action">
        <template #default="{ row }">
          <el-button @click="del(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <Pagination />
  </div>
</template>
<script setup lang="ts">
import { useAddSelfout, useDelSelfout, useSearchSelfout, useUpdateSelfout } from "@/api/selfout";
import { createProForm, ProForm } from "@howuse/element-plus-form"
import { ElMessage } from "element-plus";
import { onMounted } from "vue";

const { execute: executeAdd } = useAddSelfout()
const { execute: executeUpdate } = useUpdateSelfout()
const { execute: executeDel } = useDelSelfout()
const { search, data: linkData, load, Pagination } = useSearchSelfout()

onMounted(() => {
  search({})
})

const { formColumns } = createProForm([
  {
    key: "width",
    label: "最小宽度",
    value: 0,
    is: 'ElInputNumber',
    rules: [
      {
        required: true,
        message: '请输入最小宽度',
      },
    ],
  },
  {
    key: "height",
    label: "最小高度",
    value: 0,
    is: 'ElInputNumber',
    rules: [
      {
        required: true,
        message: '请输入最小高度',
      },
    ],
  },
  {
    key: "speed",
    label: "最小响应速率",
    value: 0,
    is: 'ElInputNumber',
    rules: [
      {
        required: true,
        message: '请输入最小响应速率',
      },
    ],
  },
  {
    key: "key",
    label: "链接标识",
    value: "",
    is: 'ElInput',
    rules: [
      {
        required: true,
        message: '请输入链接标识',
      },
    ],
  },
])

function save(raw) {
  if (raw.id) {
    executeUpdate({ data: raw }).then(res => {
      if (res.data.code == 500) {
        return ElMessage.error(res.data.msg)
      }
      ElMessage.success('保存成功')
      load()
    })
  } else {
    executeAdd({ data: raw }).then(res => {
      if (res.data.code == 500) {
        return ElMessage.error(res.data.msg)
      }
      ElMessage.success('保存成功')
      load()
    })
  }
}

function del(raw) {
  executeDel({
    params: {
      id: raw.id
    }
  }).then(res => {
    if (res.data.code == 500) {
      return ElMessage.error(res.data.msg)
    }
    ElMessage.success('删除成功')
    load()
  })
}
</script>