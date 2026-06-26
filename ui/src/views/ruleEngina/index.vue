<template>
    <div class="prompt-engine" v-loading="loading">
        <el-alert type="info" :closable="false" show-icon class="tip">
            按分组维护规范名称对照表，切换左侧标签编辑各分组；频道名称每行一个，也支持 # 分隔。
        </el-alert>

        <div class="toolbar">
            <el-input
                v-model="keyword"
                size="small"
                clearable
                placeholder="搜索分组或频道"
                class="search-input"
            />
            <el-button size="small" @click="addGroup">添加分组</el-button>
            <span class="count-tip">{{ displayGroups.length }} 个分组 · {{ totalChannels }} 条频道</span>
            <el-button type="primary" size="small" @click="save">保存</el-button>
        </div>

        <el-empty v-if="displayGroups.length === 0" description="暂无分组，点击「添加分组」开始维护" />
        <div v-else class="rule-layout">
            <aside class="group-nav">
                <button
                    v-for="group in displayGroups"
                    :key="group.id"
                    type="button"
                    class="group-nav-item"
                    :class="{ active: activeGroupId === group.id }"
                    @click="activeGroupId = group.id"
                >
                    <span class="group-nav-name">{{ group.name.trim() || '未命名分组' }}</span>
                    <span v-if="channelCount(group) > 0" class="group-nav-count">{{ channelCount(group) }}</span>
                </button>
            </aside>

            <main v-if="activeGroup" class="group-content">
                <div class="group-meta">
                    <span class="field-label">分组名称</span>
                    <el-autocomplete
                        v-model="activeGroup.name"
                        size="small"
                        clearable
                        :fetch-suggestions="queryGroupSuggestions"
                        placeholder="如 央视频道"
                        class="group-name-input"
                    />
                    <span class="channel-count">{{ channelCount(activeGroup) }} 个频道</span>
                    <el-button link type="danger" @click="removeGroup(activeGroup)">删除分组</el-button>
                </div>
                <div class="channel-editor">
                    <span class="field-label">规范名称</span>
                    <el-input
                        v-model="activeGroup.channelsText"
                        type="textarea"
                        :autosize="{ minRows: 10, maxRows: 28 }"
                        placeholder="每行一个，例如：&#10;CCTV-1综合&#10;CCTV-2财经"
                        spellcheck="false"
                    />
                </div>
            </main>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, computed, ref, watch } from 'vue'
import { useRequestRuleGet, useRequestRuleUpdate, type RuleGroup } from '@/api/tv'
import { ElMessage } from 'element-plus'
import {
    groupOptionsFromConfig,
    groupsToEdits,
    editsToGroups,
    parseChannelNames,
    createGroupEdit,
    type GroupEdit,
} from '@/utils/promptConfig'

const { execute: executeDoc, data: ruleData } = useRequestRuleGet()
const { execute: executeUpdate } = useRequestRuleUpdate()

const groups = ref<GroupEdit[]>([])
const loading = ref(false)
const keyword = ref('')
const activeGroupId = ref<number | null>(null)

const groupOptions = computed(() =>
    groupOptionsFromConfig({ groups: editsToGroups(groups.value) })
)

const displayGroups = computed(() => {
    const kw = keyword.value.trim().toLowerCase()
    if (!kw) {
        return groups.value
    }
    return groups.value.filter(g =>
        g.name.toLowerCase().includes(kw) ||
        g.channelsText.toLowerCase().includes(kw)
    )
})

const activeGroup = computed(() =>
    displayGroups.value.find(g => g.id === activeGroupId.value) ?? null
)

const totalChannels = computed(() =>
    groups.value.reduce((sum, g) => sum + channelCount(g), 0)
)

watch(displayGroups, (list) => {
    if (list.length === 0) {
        activeGroupId.value = null
        return
    }
    if (!list.some(g => g.id === activeGroupId.value)) {
        activeGroupId.value = list[0].id
    }
})

onMounted(async () => {
    loading.value = true
    try {
        await executeDoc({})
        groups.value = groupsToEdits(ruleData.value?.groups ?? [])
        if (groups.value.length > 0) {
            activeGroupId.value = groups.value[0].id
        }
    } finally {
        loading.value = false
    }
})

function channelCount(group: GroupEdit) {
    return parseChannelNames(group.channelsText).length
}

function queryGroupSuggestions(queryString: string, cb: (results: { value: string }[]) => void) {
    const q = queryString.trim().toLowerCase()
    const results = groupOptions.value
        .filter(g => !q || g.toLowerCase().includes(q))
        .map(g => ({ value: g }))
    cb(results)
}

function addGroup() {
    const group = createGroupEdit()
    groups.value.unshift(group)
    keyword.value = ''
    activeGroupId.value = group.id
}

function removeGroup(group: GroupEdit) {
    const idx = groups.value.indexOf(group)
    if (idx >= 0) {
        groups.value.splice(idx, 1)
    }
}

async function save() {
    const payload: RuleGroup[] = editsToGroups(groups.value)
    if (payload.length === 0) {
        ElMessage.warning('请至少维护一个分组和频道名称')
        return
    }
    await executeUpdate({ data: { groups: payload } })
    ElMessage.success('保存成功')
}
</script>

<style lang="scss" scoped>
.prompt-engine {
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: calc(100vh - 180px);
    overflow: hidden;
}

.tip {
    flex-shrink: 0;
}

.toolbar {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
}

.search-input {
    width: 220px;
}

.count-tip {
    flex: 1;
    color: var(--el-text-color-secondary);
    font-size: 13px;
}

.rule-layout {
    flex: 1;
    min-height: 0;
    display: flex;
    border: 1px solid var(--el-border-color-light);
    border-radius: 4px;
    overflow: hidden;
    background: var(--el-bg-color);
}

.group-nav {
    flex-shrink: 0;
    width: 200px;
    overflow-y: auto;
    overscroll-behavior: contain;
    border-right: 1px solid var(--el-border-color-light);
    background: var(--el-fill-color-light);
    padding: 6px 0;
}

.group-nav-item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 10px 14px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    font-size: 13px;
    color: var(--el-text-color-regular);
    transition: background-color 0.15s, color 0.15s;

    &:hover {
        background: var(--el-fill-color);
    }

    &.active {
        background: var(--el-bg-color);
        color: var(--el-color-primary);
        font-weight: 500;
        box-shadow: inset 3px 0 0 var(--el-color-primary);
    }
}

.group-nav-name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.group-nav-count {
    flex-shrink: 0;
    min-width: 20px;
    padding: 0 6px;
    border-radius: 10px;
    background: var(--el-fill-color-darker);
    color: var(--el-text-color-secondary);
    font-size: 12px;
    line-height: 20px;
    text-align: center;

    .active & {
        background: var(--el-color-primary-light-8);
        color: var(--el-color-primary);
    }
}

.group-content {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    overflow-y: auto;
}

.group-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
}

.field-label {
    flex-shrink: 0;
    width: 72px;
    color: var(--el-text-color-secondary);
    font-size: 13px;
}

.group-name-input {
    width: 220px;
}

.channel-count {
    flex: 1;
    color: var(--el-text-color-secondary);
    font-size: 13px;
}

.channel-editor {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    flex: 1;
    min-height: 0;

    :deep(.el-textarea) {
        flex: 1;
    }
}
</style>
