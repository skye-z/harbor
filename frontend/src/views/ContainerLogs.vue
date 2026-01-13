<template>
  <div class="logs-page">
    <div class="page-header">
      <div class="header-left">
        <n-button quaternary circle @click="router.back()">
          <template #icon>
            <n-icon :component="ArrowBackOutline" />
          </template>
        </n-button>
        <div class="header-info">
          <h1>容器日志</h1>
          <n-tag :type="container?.state === 'running' ? 'success' : 'default'" size="small" round :bordered="false">
            {{ container?.names?.[0]?.replace(/^\//, '') || containerId }}
          </n-tag>
        </div>
      </div>
      <div class="header-actions">
        <n-space>
          <n-select
            v-model:value="tailLines"
            :options="tailOptions"
            style="width: 140px"
            size="small"
          />
          <n-input
            v-model:value="searchText"
            placeholder="过滤日志..."
            clearable
            style="width: 200px"
            size="small"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
          </n-input>
          <n-button @click="toggleAutoRefresh" :type="autoRefresh ? 'warning' : 'default'" size="small">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            {{ autoRefresh ? '停止刷新' : '自动刷新' }}
          </n-button>
          <n-button @click="loadLogs" :loading="loading" size="small">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            刷新
          </n-button>
        </n-space>
      </div>
    </div>

    <n-card class="logs-card" content-style="padding: 0; display: flex; flex-direction: column; height: 100%;">
      <div class="logs-toolbar">
        <n-space align="center">
          <n-tag :type="statusType" size="small">{{ statusText }}</n-tag>
          <n-text depth="3">共 {{ filteredLogs.length }} 行</n-text>
          <n-divider vertical />
          <n-checkbox v-model:checked="showTimestamp" size="small">显示时间戳</n-checkbox>
          <n-checkbox v-model:checked="follow" size="small" :disabled="autoRefresh">跟随</n-checkbox>
        </n-space>
      </div>
      
      <div class="logs-container" ref="logsContainerRef">
        <div v-if="loading && !logs.length" class="logs-loading">
          <n-spin size="small" />
          <span>加载中...</span>
        </div>
        <div v-else-if="!logs.length" class="logs-empty">
          <n-empty description="暂无日志" />
        </div>
        <pre v-else class="logs-content"><code
          v-for="(line, index) in filteredLogs"
          :key="index"
          :class="getLineClass(line)"
        >{{ line }}</code></pre>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { ArrowBackOutline, SearchOutline, RefreshOutline } from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()

const containerId = computed(() => route.params.id as string)
const container = computed(() => containerStore.getContainerById(containerId.value))

const logs = ref<string[]>([])
const loading = ref(false)
const searchText = ref('')
const tailLines = ref('100')
const showTimestamp = ref(true)
const follow = ref(false)
const autoRefresh = ref(false)
const logsContainerRef = ref<HTMLElement>()

let autoRefreshTimer: number | null = null

const tailOptions = [
  { label: '最后50行', value: '50' },
  { label: '最后100行', value: '100' },
  { label: '最后200行', value: '200' },
  { label: '最后500行', value: '500' },
  { label: '全部', value: 'all' }
]

const statusType = computed(() => {
  if (!container.value) return 'default'
  return container.value.state === 'running' ? 'success' : 'warning'
})

const statusText = computed(() => {
  if (!container.value) return '未知'
  return container.value.state === 'running' ? '运行中' : '已停止'
})

const filteredLogs = computed(() => {
  if (!searchText.value) return logs.value
  const keyword = searchText.value.toLowerCase()
  return logs.value.filter(line => line.toLowerCase().includes(keyword))
})

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    autoRefreshTimer = window.setInterval(() => {
      loadLogs()
    }, 3000)
  } else if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

const getLineClass = (line: string) => {
  if (line.toLowerCase().includes('error') || line.toLowerCase().includes('err')) return 'log-error'
  if (line.toLowerCase().includes('warn')) return 'log-warn'
  if (line.toLowerCase().includes('info')) return 'log-info'
  return ''
}

const loadLogs = async () => {
  if (!containerId.value) return
  
  loading.value = true
  try {
    const result = await containerStore.getContainerLogs(containerId.value)
    const logText = result?.logs || result || ''
    logs.value = logText.split('\n').filter(line => line.length > 0)
    
    if (follow.value || autoRefresh.value) {
      await nextTick()
      scrollToBottom()
    }
  } catch (error: any) {
    console.error('加载日志失败:', error)
  } finally {
    loading.value = false
  }
}

const scrollToBottom = () => {
  if (logsContainerRef.value) {
    logsContainerRef.value.scrollTop = logsContainerRef.value.scrollHeight
  }
}

onMounted(async () => {
  if (!containerStore.containers.length) {
    await containerStore.fetchContainers()
  }
  if (!container.value) {
    router.push({ name: 'Containers' })
    return
  }
  await loadLogs()
})

onUnmounted(() => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
  }
})

watch(tailLines, () => {
  loadLogs()
})
</script>

<style scoped>
.logs-page {
  height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-info h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
}

.logs-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.logs-toolbar {
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logs-container {
  flex: 1;
  overflow: auto;
  background: #0d1117;
  padding: 12px 16px;
}

.logs-loading,
.logs-empty {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #8b949e;
}

.logs-content {
  margin: 0;
  font-family: 'Fira Code', Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.logs-content code {
  display: block;
}

.log-error {
  color: #ff7b72;
}

.log-warn {
  color: #d29922;
}

.log-info {
  color: #58a6ff;
}
</style>
