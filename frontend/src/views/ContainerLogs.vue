<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { NButton, NIcon, NSelect, NInput, NCheckbox } from 'naive-ui'
import { ArrowBack, SearchOutline, RefreshOutline, StopOutline, PlayOutline } from '@vicons/ionicons5'

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
  { label: '50行', value: '50' },
  { label: '100行', value: '100' },
  { label: '200行', value: '200' },
  { label: '500行', value: '500' },
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

const handleBack = () => {
  router.push({ name: 'ContainerDetail', params: { id: containerId.value } })
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

<template>
  <div class="logs-page">
    <div class="toolbar">
      <div class="toolbar-left">
        <n-button quaternary size="small" @click="handleBack">
          <template #icon>
            <n-icon :component="ArrowBack" />
          </template>
        </n-button>
        <span class="container-name">{{ container?.names?.[0]?.replace(/^\//, '') || containerId }}</span>
        <span class="status-indicator">
          <span class="status-dot" :class="statusType"></span>
          <span class="status-text">{{ statusText }}</span>
        </span>
      </div>
      <div class="toolbar-center">
        <n-select
          v-model:value="tailLines"
          :options="tailOptions"
          size="small"
          style="width: 90px"
        />
        <n-input
          v-model:value="searchText"
          placeholder="过滤..."
          clearable
          size="small"
          style="width: 140px"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" size="14" />
          </template>
        </n-input>
        <n-checkbox v-model:checked="showTimestamp" size="small">时间</n-checkbox>
        <n-checkbox v-model:checked="follow" size="small" :disabled="autoRefresh">跟随</n-checkbox>
      </div>
      <div class="toolbar-right">
        <n-button
          quaternary
          size="small"
          @click="toggleAutoRefresh"
          :type="autoRefresh ? 'warning' : 'default'"
        >
          <template #icon>
            <n-icon :component="autoRefresh ? StopOutline : PlayOutline" />
          </template>
        </n-button>
        <n-button quaternary size="small" @click="loadLogs" :loading="loading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
    </div>

    <div class="logs-area" ref="logsContainerRef">
      <div v-if="loading && !logs.length" class="placeholder">
        <span>加载中...</span>
      </div>
      <div v-else-if="!logs.length" class="placeholder">
        <span>暂无日志</span>
      </div>
      <pre v-else class="logs-content"><code
        v-for="(line, index) in filteredLogs"
        :key="index"
        :class="getLineClass(line)"
      >{{ line }}</code></pre>
    </div>
  </div>
</template>

<style scoped>
.logs-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  z-index: 100;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
  flex-shrink: 0;
  gap: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.container-name {
  color: #cccccc;
  font-weight: 500;
  font-size: 14px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.success {
  background: #0dbc79;
  box-shadow: 0 0 8px rgba(13, 188, 121, 0.6);
}

.status-dot.warning {
  background: #e5e510;
}

.status-text {
  color: #808080;
  font-size: 12px;
}

.toolbar-center {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.logs-area {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 8px 12px;
}

.placeholder {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #808080;
}

.logs-content {
  margin: 0;
  font-family: 'Fira Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  color: #d4d4d4;
}

.logs-content code {
  display: block;
}

.log-error {
  color: #f14c4c;
}

.log-warn {
  color: #e5e510;
}

.log-info {
  color: #3794ff;
}
</style>
