<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { logApi } from '../plugins/api'
import { NButton, NIcon, NSelect } from 'naive-ui'
import { ArrowBack, SearchOutline, RefreshOutline, StopOutline, PlayOutline, SettingsOutline, CubeOutline, ImagesOutline, ServerOutline, GlobeOutline } from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

const logs = ref<any[]>([])
const loading = ref(false)
const searchText = ref('')
const logType = ref('')
const limit = ref('50')

let autoRefreshTimer: number | null = null

const typeOptions = [
  { label: '全部', value: '' },
  { label: '容器', value: 'container' },
  { label: '镜像', value: 'image' },
  { label: '卷', value: 'volume' },
  { label: '网络', value: 'network' },
  { label: '系统', value: 'system' }
]

const limitOptions = [
  { label: '20条', value: '20' },
  { label: '50条', value: '50' },
  { label: '100条', value: '100' },
  { label: '200条', value: '200' }
]

const filteredLogs = computed(() => {
  if (!searchText.value) return logs.value
  const keyword = searchText.value.toLowerCase()
  return logs.value.filter((log: any) => {
    const message = (log.message || log.action || '').toLowerCase()
    const target = (log.target || '').toLowerCase()
    return message.includes(keyword) || target.includes(keyword)
  })
})

const getLogTypeIcon = (type: string) => {
  switch (type) {
    case 'container': return CubeOutline
    case 'image': return ImagesOutline
    case 'volume': return ServerOutline
    case 'network': return GlobeOutline
    default: return SettingsOutline
  }
}

const getLogTypeColor = (type: string) => {
  switch (type) {
    case 'container': return '#2080f0'
    case 'image': return '#f0a020'
    case 'volume': return '#18a058'
    case 'network': return '#d03050'
    case 'system': return '#a0a0a0'
    default: return '#a0a0a0'
  }
}

const getLogTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    container: '容器',
    image: '镜像',
    volume: '卷',
    network: '网络',
    system: '系统'
  }
  return map[type] || type
}

const loadLogs = async () => {
  loading.value = true
  try {
    const data = await logApi.getRecent(parseInt(limit.value), logType.value)
    logs.value = Array.isArray(data) ? data : []
  } catch (error: any) {
    console.error('加载日志失败:', error)
  } finally {
    loading.value = false
  }
}

const handleBack = () => {
  router.push({ name: 'Dashboard' })
}

onMounted(async () => {
  await loadLogs()
})

onUnmounted(() => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
  }
})

watch([logType, limit], () => {
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
        <span class="container-name">平台日志</span>
      </div>
      <div class="toolbar-center">
        <n-select
          v-model:value="logType"
          :options="typeOptions"
          size="small"
          style="width: 100px"
        />
        <n-select
          v-model:value="limit"
          :options="limitOptions"
          size="small"
          style="width: 90px"
        />
        <n-input
          v-model:value="searchText"
          placeholder="搜索..."
          clearable
          size="small"
          style="width: 140px"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" size="14" />
          </template>
        </n-input>
      </div>
      <div class="toolbar-right">
        <n-button quaternary size="small" @click="loadLogs" :loading="loading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
    </div>

    <div class="logs-area">
      <div v-if="loading && !logs.length" class="placeholder">
        <span>加载中...</span>
      </div>
      <div v-else-if="!logs.length" class="placeholder">
        <span>暂无日志</span>
      </div>
      <div v-else class="logs-list">
        <div
          v-for="(log, index) in filteredLogs"
          :key="log.id || index"
          class="log-item"
        >
          <div class="log-header">
            <div class="log-type">
              <n-icon :component="getLogTypeIcon(log.type)" :color="getLogTypeColor(log.type)" size="16" />
              <span>{{ getLogTypeLabel(log.type) }}</span>
            </div>
            <span class="log-time">{{ new Date(log.created_at).toLocaleString('zh-CN', { hour12: false }) }}</span>
          </div>
          <div class="log-content">
            {{ log.message || log.action }}
          </div>
          <div v-if="log.target" class="log-target">
            目标: {{ log.target }}
          </div>
          <div v-if="log.username" class="log-user">
            操作人: {{ log.username }}
          </div>
        </div>
      </div>
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

.logs-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-item {
  background: #2a2d2e;
  border-radius: 6px;
  padding: 12px;
  border: 1px solid #3c3c3c;
  transition: all 0.2s;
}

.log-item:hover {
  border-color: #4a4a4a;
  background: #2d2f30;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.log-type {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  color: #cccccc;
}

.log-time {
  font-size: 12px;
  color: #808080;
  font-family: monospace;
}

.log-content {
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 6px;
}

.log-target,
.log-user {
  font-size: 12px;
  color: #808080;
  margin-top: 2px;
}
</style>
