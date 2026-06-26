<template>
  <div class=" flex flex-col gap-2 w-full">
    <ElForm :rules="rules" ref="formRef" :model="modelRef">
      <ElFormItem prop="username">
        <ElInput placeholder="请输入用户名" v-model:model-value="modelRef.username" size="large"></ElInput>
      </ElFormItem>
      <ElFormItem prop="password">
        <ElInput placeholder="密码" type="password" v-model:model-value="modelRef.password" size="large">
        </ElInput>
      </ElFormItem>
      <p class=" text-right inline-flex w-full gap-2">
        <ElButton type="primary" class=" flex-1" @click="handleValidateButtonClick">登录</ElButton>
      </p>
    </ElForm>
  </div>
</template>

<script lang="ts" setup>
  import { type FormInstance, type FormRules, ElMessage } from "element-plus";

  import {ref, unref } from "vue";
  import { useRequestFind, useRequestIdentify } from "@/api/tv"

  const { execute: executeFind } = useRequestFind()
  const { execute } = useRequestIdentify()

  interface ModelType {
    username: string | null;
    password: string | null;
  }

  const formRef = ref<FormInstance>();

  const modelRef = ref<ModelType>({
    username: "admin",
    password: "admin123",
  });

  const rules: FormRules = {
    username: [
      {
        required: true,
        message: "请输入用户名",
      },
    ],
    password: [
      {
        required: true,
        message: "请输入密码",
      },
    ],
  };

  async function handleValidateButtonClick(e: MouseEvent) {
    formRef.value?.validate((valid: boolean) => {
      if (!valid) {
        return
      }
      execute({
        params: {
          username: unref(modelRef).username,
          password: unref(modelRef).password,
        }
      }).then((res) => {
        if (res.data?.token) {
          ElMessage.success("登录成功")
          const redirectUrl = sessionStorage.getItem("redirect_url") || "/admin"
          sessionStorage.clear()
          sessionStorage.setItem("login", "1")
          sessionStorage.setItem("auth", res.data.token)
          sessionStorage.setItem("rule_password", res.data.password)

          executeFind({}).then((res) => {
            res.data.settings.forEach((item: any) => {
              sessionStorage.setItem(item.key, item.value)
              sessionStorage.setItem(`${item.key}_id`, item.id)
              location.href = redirectUrl
            });
          })
        } else {
          ElMessage.error("登录失败")
        }
      }).catch(() => {
        ElMessage.error("用户名或密码错误")
      })
    })
  }
</script>
