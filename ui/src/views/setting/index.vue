<template>
    <div class="overflow-y-auto overflow-x-hidden h-[80vh]">
        <el-tabs v-model="activeTab" class="setting-tabs">
            <el-tab-pane label="通用设置" name="general">
                <ProForm
                    :columns="generalColumns"
                    :span="24"
                    :btn-position="'center'"
                    :display-expend="false"
                    hidden-reset
                    ok-text="保存"
                    @ok="submit"
                >
                    <template #plusKey="{ data }">
                        <div class="flex w-full flex-col">
                            <ElInput v-model:model-value="data.value" />
                            <div class="font-bold">
                                {{ local }}/v1/tv/w/{{ data.value }}
                                <el-button type="primary" link @click="handleSelfOut">【自建输出策略】</el-button>
                            </div>
                        </div>
                    </template>
                    <template #autoDelCount="{ data }">
                        <div class="flex w-full flex-col">
                            <ElInput v-model:model-value="data.value" />
                            <div class="font-bold text-red-900">
                                填入一个数字，失败次数达到这个数字就会删除。填入1，表示检测失败立即删除。
                            </div>
                        </div>
                    </template>
                    <template #blackList="{ data }">
                        <div class="flex w-full flex-col">
                            <ElInput
                                v-model:model-value="data.value"
                                v-bind="{
                                    type: 'textarea',
                                    autosize: { minRows: 10, maxRows: 10 },
                                    placeholder: 'xxx#yyy#zzz',
                                    spellcheck: 'false',
                                    resize: 'none',
                                }"
                            />
                            <div class="font-bold text-red-900">
                                填入此处的链接将不会在生成的链接当中显示出来，一般填入域名或者IP即可。
                            </div>
                        </div>
                    </template>
                </ProForm>
            </el-tab-pane>

            <el-tab-pane label="账号安全" name="security">
                <el-form :model="pwdForm" label-width="120px" class="pwd-form">
                    <el-form-item label="原密码">
                        <ElInput
                            v-model="pwdForm.oldPassword"
                            type="text"
                            autocomplete="off"
                            placeholder="请输入当前密码"
                            class="pwd-field"
                            :class="{ 'pwd-masked': !pwdVisible.old }"
                        >
                            <template #suffix>
                                <el-button link type="primary" @click="pwdVisible.old = !pwdVisible.old">
                                    <el-icon><View v-if="!pwdVisible.old" /><Hide v-else /></el-icon>
                                </el-button>
                            </template>
                        </ElInput>
                    </el-form-item>
                    <el-form-item label="新密码">
                        <ElInput
                            v-model="pwdForm.newPassword"
                            type="text"
                            autocomplete="off"
                            placeholder="至少 4 位"
                            class="pwd-field"
                            :class="{ 'pwd-masked': !pwdVisible.new }"
                        >
                            <template #suffix>
                                <el-button link type="primary" @click="pwdVisible.new = !pwdVisible.new">
                                    <el-icon><View v-if="!pwdVisible.new" /><Hide v-else /></el-icon>
                                </el-button>
                            </template>
                        </ElInput>
                    </el-form-item>
                    <el-form-item label="确认新密码">
                        <ElInput
                            v-model="pwdForm.confirmPassword"
                            type="text"
                            autocomplete="off"
                            placeholder="再次输入新密码"
                            class="pwd-field"
                            :class="{ 'pwd-masked': !pwdVisible.confirm }"
                        >
                            <template #suffix>
                                <el-button link type="primary" @click="pwdVisible.confirm = !pwdVisible.confirm">
                                    <el-icon><View v-if="!pwdVisible.confirm" /><Hide v-else /></el-icon>
                                </el-button>
                            </template>
                        </ElInput>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" :loading="pwdLoading" @click="changePassword">修改密码</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>

            <el-tab-pane label="AI设置" name="ai">
                <ProForm
                    :columns="aiColumns"
                    :span="24"
                    :btn-position="'center'"
                    :display-expend="false"
                    hidden-reset
                    ok-text="保存"
                    @ok="submit"
                >
                    <template #aiBaseUrl="{ data }">
                        <ElInput
                            v-model:model-value="data.value"
                            placeholder="OpenAI 兼容 API，如 https://api.openai.com/v1"
                            @blur="fetchModels"
                        />
                    </template>
                    <template #aiApiKey="{ data }">
                        <ElInput
                            v-model:model-value="data.value"
                            type="text"
                            autocomplete="off"
                            autocapitalize="off"
                            spellcheck="false"
                            name="ai-api-key-config"
                            placeholder="sk-..."
                            class="api-key-field"
                            :class="{ 'api-key-masked': !apiKeyVisible }"
                            data-1p-ignore
                            data-lpignore="true"
                            @blur="fetchModels"
                        >
                            <template #suffix>
                                <el-button link type="primary" class="api-key-toggle" @click="apiKeyVisible = !apiKeyVisible">
                                    <el-icon>
                                        <View v-if="!apiKeyVisible" />
                                        <Hide v-else />
                                    </el-icon>
                                </el-button>
                            </template>
                        </ElInput>
                    </template>
                    <template #aiModel="{ data }">
                        <div class="flex w-full items-center gap-2">
                            <el-select
                                v-model="data.value"
                                filterable
                                allow-create
                                default-first-option
                                clearable
                                placeholder="选择或输入模型名称"
                                :loading="modelsLoading"
                                class="flex-1"
                            >
                                <el-option v-for="m in modelOptions" :key="m" :label="m" :value="m" />
                            </el-select>
                            <el-button :loading="modelsLoading" @click="fetchModels">刷新</el-button>
                        </div>
                    </template>
                </ProForm>
            </el-tab-pane>
        </el-tabs>
    </div>
    <el-dialog title="自建输出" v-model="outVisible">
        <ExportModal />
    </el-dialog>
</template>

<script lang="ts" setup>
import { createProForm, ProForm } from '@howuse/element-plus-form'
import { ElInput, ElMessage, ElSwitch, ElTimePicker } from 'element-plus'
import { View, Hide } from '@element-plus/icons-vue'
import { useRequestFind, useRequestUpdate, useRequestAiModels, useRequestChangePassword } from '@/api/tv'
import { onMounted, reactive, ref, unref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToggle } from '@vueuse/core'
import ExportModal from './ExportModal.vue'

const { execute: executeFind, data: settingData } = useRequestFind()
const { execute: executeUpdate } = useRequestUpdate()
const { execute: executeAiModels } = useRequestAiModels()
const { execute: executeChangePassword } = useRequestChangePassword()
const router = useRouter()
const route = useRoute()
const [outVisible, toggleOut] = useToggle()
const apiKeyVisible = ref(false)
const pwdLoading = ref(false)
const pwdVisible = reactive({ old: false, new: false, confirm: false })
const pwdForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
})
const modelOptions = ref<string[]>([])
const modelsLoading = ref(false)
const activeTab = ref('general')

function syncTabFromRoute() {
    const tab = route.query.tab
    const tabName = Array.isArray(tab) ? tab[0] : tab
    if (tabName === 'security' || tabName === 'ai' || tabName === 'general') {
        activeTab.value = tabName
    }
}

watch(() => route.query.tab, syncTabFromRoute)

const generalColumns = [
    {
        key: 'unknownGroup',
        label: '播放列表显示未分组',
        value: false,
        is: ElSwitch,
        span: 6,
        props: {
            'active-value': '1',
            'inactive-value': '0',
        },
    },
    {
        key: 'autoEpg',
        label: '自动更新EPG',
        value: false,
        is: ElSwitch,
        span: 6,
        props: {
            'active-value': '1',
            'inactive-value': '0',
        },
    },
    {
        key: 'autoCheck',
        label: '自动守护',
        value: false,
        is: ElSwitch,
        span: 6,
        props: {
            'active-value': '1',
            'inactive-value': '0',
        },
    },
    {
        key: 'plusKey',
        label: 'DIYP链接',
        value: '',
        is: 'ElSelect',
    },
    {
        key: 'autoDelCount',
        label: '自动删除数量',
        value: '',
        is: 'ElSelect',
        rules: [{ pattern: /^[12345]$/, message: '只能输入1到5', trigger: 'change' }],
    },
    {
        key: 'subscriberTime',
        label: '订阅自动读取时间(每天)',
        value: '00:00',
        is: ElTimePicker,
        props: {
            style: { width: '300px' },
            format: 'HH:mm',
            valueFormat: 'HH:mm',
        },
    },
    {
        key: 'blackList',
        label: '黑名单',
        value: '',
        is: 'ElInput',
        props: {
            type: 'textarea',
            autosize: { minRows: 10, maxRows: 10 },
            placeholder: 'xxx#yyy#zzz',
            spellcheck: 'false',
            resize: 'none',
        },
        rules: [{ pattern: /[0-9a-zA-Z#\.:\/]+$/, message: '只能包含数字、字母、#、.、:、/', trigger: 'change' }],
    },
]

const aiColumns = [
    {
        key: 'aiBaseUrl',
        label: 'AI 接口地址',
        value: 'https://api.openai.com/v1',
        is: 'ElInput',
    },
    {
        key: 'aiApiKey',
        label: 'AI API Key',
        value: '',
        is: 'ElInput',
    },
    {
        key: 'aiModel',
        label: 'AI 模型',
        value: 'gpt-4o-mini',
        is: 'ElSelect',
    },
]

const { getValue, setValue } = createProForm([...generalColumns, ...aiColumns])

onMounted(() => {
    syncTabFromRoute()
    executeFind({}).then((res) => {
        const settings = res.data.settings
        for (const setting of settings) {
            if (setting.key === 'password' || setting.key === 'username' || setting.key === 'autoGroup') {
                continue
            }
            setValue(setting.key, setting.value)
        }
        fetchModels()
    })
})

async function fetchModels() {
    const baseUrl = getValue('aiBaseUrl')
    const apiKey = getValue('aiApiKey')
    if (!baseUrl || !apiKey) {
        return
    }
    modelsLoading.value = true
    try {
        const res = await executeAiModels({
            params: { baseUrl, apiKey },
        })
        modelOptions.value = res.data?.models || []
    } catch {
        ElMessage.warning('获取模型列表失败，可手动输入模型名称')
    } finally {
        modelsLoading.value = false
    }
}

const local = location.origin

async function submit(raw) {
    const data = unref(settingData).settings
    const allKeys = [...generalColumns, ...aiColumns].map((item) => item.key)
    const merged = { ...raw }
    for (const key of allKeys) {
        if (merged[key] === undefined) {
            merged[key] = getValue(key)
        }
    }

    const settings = [...data]
        .map((item) => {
            item.value = merged[item.key]
            return item
        })
        .filter((item) => !!item.value && item.key !== 'password' && item.key !== 'username')

    await executeUpdate({ data: { settings } })
    executeFind({}).then((res) => {
        res.data.settings.forEach((item: any) => {
            sessionStorage.setItem(item.key, item.value)
        })
    })
    ElMessage.success('保存成功')
}

async function changePassword() {
    if (!pwdForm.oldPassword || !pwdForm.newPassword || !pwdForm.confirmPassword) {
        ElMessage.warning('请完整填写密码信息')
        return
    }
    if (pwdForm.newPassword !== pwdForm.confirmPassword) {
        ElMessage.warning('两次输入的新密码不一致')
        return
    }
    pwdLoading.value = true
    try {
        await executeChangePassword({
            data: {
                oldPassword: pwdForm.oldPassword,
                newPassword: pwdForm.newPassword,
            },
        })
        ElMessage.success('密码修改成功，请重新登录')
        pwdForm.oldPassword = ''
        pwdForm.newPassword = ''
        pwdForm.confirmPassword = ''
        sessionStorage.clear()
        setTimeout(() => router.replace('/login'), 800)
    } catch {
        ElMessage.error('密码修改失败，请检查原密码是否正确')
    } finally {
        pwdLoading.value = false
    }
}

function handleSelfOut() {
    toggleOut(true)
}
</script>

<style scoped>
.setting-tabs :deep(.el-tabs__content) {
    padding-top: 12px;
}

.api-key-field.api-key-masked :deep(.el-input__inner) {
    -webkit-text-security: disc;
    text-security: disc;
}

.api-key-toggle {
    padding: 0 4px;
}

.pwd-form {
    max-width: 560px;
    padding-bottom: 24px;
}

.pwd-field.pwd-masked :deep(.el-input__inner) {
    -webkit-text-security: disc;
    text-security: disc;
}
</style>
