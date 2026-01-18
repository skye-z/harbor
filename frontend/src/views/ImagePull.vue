<script setup lang="ts">
import { ref, computed, onUnmounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { imageApi } from '../plugins/api'
import { useMessage, NButton, NIcon, NInput, NSelect, NProgress, NTag, NAutoComplete } from 'naive-ui'
import {
  CloudDownloadOutline
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const registryOptions = [
  { label: '默认镜像源', value: '' },
  { label: 'Docker Hub', value: 'index.docker.io' },
  { label: 'Harbor 镜像', value: 'harbor' },
  { label: 'GitHub Container Registry', value: 'ghcr.io' }
]

const selectedRegistry = ref('')
const imageName = ref('')
const imageTag = ref('latest')

interface SearchResult {
  name: string
  description: string
  StarCount: number
  IsOfficial: boolean
}

const suggestions = ref<SearchResult[]>([])
const suggestionLoading = ref(false)

type PullStatus = 'idle' | 'pulling' | 'success' | 'error'
const pullStatus = ref<PullStatus>('idle')
const pullProgress = ref(0)
const pullLogs = ref<string[]>([])
const pullError = ref('')
let pollTimer: number | null = null
let searchTimer: number | null = null

const isPulling = computed(() => pullStatus.value === 'pulling')

const debouncedSearch = useDebounceFn(async (query: string) => {
  await loadSuggestions(query)
}, 300)

function useDebounceFn<T extends (...args: any[]) => any>(fn: T, delay: number): T {
  return ((...args: any[]) => {
    if (searchTimer) clearTimeout(searchTimer)
    searchTimer = window.setTimeout(() => fn(...args), delay)
  }) as T
}

const fullImageName = computed(() => {
  const registry = selectedRegistry.value ? selectedRegistry.value + '/' : ''
  return `${registry}${imageName.value}:${imageTag.value}`
})

interface LayerProgress {
  id: string
  status: string
  progress: string
  current: number
  total: number
}

interface PullProgressResponse {
  tag: string
  layers: LayerProgress[]
  percent: number
}

const loadSuggestions = async (query: string) => {
  if (!query || query.length < 1) {
    suggestions.value = []
    return
  }

  suggestionLoading.value = true
  try {
    const results = await imageApi.search(query, 10)
    suggestions.value = (results || []) as SearchResult[]
  } catch (e) {
    suggestions.value = []
  } finally {
    suggestionLoading.value = false
  }
}

const selectSuggestion = (item: SearchResult) => {
  imageName.value = item.name
}

const startPolling = () => {
  const seenLogs = new Set<string>()
  
  pollTimer = window.setInterval(async () => {
    try {
      const data: PullProgressResponse = await imageApi.getPullProgress()

      if (!data.tag) {
        clearInterval(pollTimer!)
        if (pullStatus.value !== 'error') {
          pullStatus.value = 'success'
          pullProgress.value = 100
          pullLogs.value.push('镜像拉取成功')
          message.success('镜像拉取成功')
        }
        return
      }

      pullProgress.value = data.percent

      data.layers.forEach(layer => {
        const logKey = `${layer.id}-${layer.status}-${layer.progress}`
        if (!seenLogs.has(logKey) && layer.status) {
          seenLogs.add(logKey)
          const progressStr = layer.progress ? ` ${layer.progress}` : ''
          pullLogs.value.push(`[${layer.id}] ${layer.status}${progressStr}`)
        }
      })
    } catch (e) {
      console.error('获取进度失败:', e)
    }
  }, 1000)
}

const renderSuggestionLabel = (option: any) => {
  const tags = []
  if (option.StarCount || option.star_count) {
    tags.push(h(NTag, { size: 'tiny', type: 'warning' }, { default: () => `★ ${option.StarCount || option.star_count}` }))
  }
  if (option.IsOfficial || option.is_official) {
    tags.push(h(NTag, { size: 'tiny', type: 'success' }, { default: () => '官方' }))
  }
  if (option.IsAutomated || option.is_automated) {
    tags.push(h(NTag, { size: 'tiny', type: 'info' }, { default: () => '自动' }))
  }

  return h('div', { style: 'padding: 8px;' }, [
    h('div', { style: 'font-size: 15px; font-weight: 600; margin-bottom: 4px;' }, option.name),
    tags.length > 0 ? h('div', { style: 'margin-bottom: 4px; display: flex; gap: 4px;' }, tags) : null,
    option.description ? h('div', { style: 'font-size: 12px; color: #888; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;' }, option.description) : null
  ])
}

const handlePull = async () => {
  if (!imageName.value) {
    message.error('请输入镜像名称')
    return
  }

  pullStatus.value = 'pulling'
  pullProgress.value = 0
  pullLogs.value = []
  pullError.value = ''

  pullLogs.value.push(`开始拉取镜像: ${fullImageName.value}`)

  startPolling()

  try {
    await imageApi.pull(fullImageName.value)
  } catch (error: any) {
    pullStatus.value = 'error'
    pullError.value = error.message || '拉取失败'
    pullLogs.value.push(`拉取失败: ${pullError.value}`)
    message.error('拉取失败: ' + pullError.value)
  }
}

const reset = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  imageName.value = ''
  imageTag.value = 'latest'
  pullStatus.value = 'idle'
  pullProgress.value = 0
  pullLogs.value = []
  pullError.value = ''
}

const preload = () => {
  const prefill = route.query.image as string
  if (prefill) {
    imageName.value = prefill
  }
}

preload()

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
</script>

<template>
  <div class="image-pull">
    <div class="page-header">
      <div class="view-header">
        <div class="title-group">
          <h1>拉取镜像</h1>
        </div>
      </div>
    </div>

    <n-grid x-gap="20" :cols="24">
      <n-gi :span="16">
        <n-card title="镜像配置" class="config-card">
          <n-form-item label="注册表">
            <n-select v-model:value="selectedRegistry" :options="registryOptions" placeholder="选择注册表" />
          </n-form-item>

          <n-form-item label="镜像名称">
            <n-auto-complete
              v-model:value="imageName"
              :options="suggestions.map(s => ({ label: s.name, value: s, key: s.name }))"
              :loading="suggestionLoading"
              placeholder="输入镜像名称搜索，如: nginx"
              :render-label="(option: any) => renderSuggestionLabel(option.value)"
              @select="(val: any) => selectSuggestion(val)"
              @update:value="debouncedSearch"
            />
          </n-form-item>

          <n-form-item label="版本标签">
            <n-input v-model:value="imageTag" placeholder="默认 latest" />
          </n-form-item>
        </n-card>

        <n-card class="config-card">
          <div class="preview-row">
            <span class="preview-label">完整镜像名</span>
            <n-tag type="info">{{ fullImageName }}</n-tag>
          </div>
        </n-card>

        <n-card v-if="pullStatus !== 'idle'" title="拉取进度" class="config-card">
          <div class="result-header">
            <span class="result-label">拉取状态</span>
            <n-tag :type="pullStatus === 'success' ? 'success' : pullStatus === 'error' ? 'error' : 'info'">
              {{ pullStatus === 'pulling' ? '拉取中...' : pullStatus === 'success' ? '成功' : '失败' }}
            </n-tag>
          </div>

          <n-progress
            v-if="pullStatus === 'pulling'"
            type="line"
            :percentage="pullProgress"
            :show-indicator="true"
            :height="8"
            border-radius="4"
          />

          <div class="logs-container" v-if="pullLogs.length > 0">
            <div class="log-item" v-for="(log, i) in pullLogs" :key="i">{{ log }}</div>
          </div>

          <n-alert v-if="pullStatus === 'error'" type="error" :title="pullError" style="margin-top: 12px" />
        </n-card>
      </n-gi>

      <n-gi :span="8">
        <n-card title="操作" class="action-card">
          <n-space vertical style="width: 100%">
            <n-button type="primary" block size="large" :loading="isPulling" :disabled="!imageName" @click="handlePull">
              <template #icon>
                <n-icon><CloudDownloadOutline /></n-icon>
              </template>
              拉取镜像
            </n-button>
            <n-button block @click="reset" :disabled="isPulling">重置</n-button>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<style scoped>
.image-pull {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 10px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
}

.config-card {
  margin-bottom: 20px;
}

.action-card {
  position: sticky;
  top: 60px;
}

.preview-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--n-color-popover);
  border-radius: 8px;
}

.preview-label {
  font-size: 13px;
  color: var(--n-text-color-3);
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.result-label {
  font-size: 14px;
  font-weight: 600;
}

.logs-container {
  margin-top: 12px;
  padding: 12px;
  background: var(--n-color-popover);
  border-radius: 8px;
  max-height: 200px;
  overflow-y: auto;
}

.log-item {
  font-size: 12px;
  font-family: monospace;
  color: var(--n-text-color-2);
  padding: 2px 0;
}
</style>
