<template>
  <div class="log-viewer">
    <n-space vertical :size="16">
      <n-space>
        <n-input v-model:value="searchQuery" placeholder="搜索日志..." clearable style="width: 300px" @input="handleSearch">
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
        <n-select v-model:value="logLevel" :options="levelOptions" placeholder="日志级别" style="width: 150px"
          @update:value="handleFilter" />
        <n-button @click="handleClear" quaternary>
          <template #icon>
            <n-icon :component="TrashOutline" />
          </template>
          清空
        </n-button>
        <n-switch v-model:value="autoScroll">
          <template #checked>自动滚动</template>
          <template #unchecked>手动滚动</template>
        </n-switch>
        <n-switch v-model:value="followLogs" @update:value="handleFollowChange">
          <template #checked>实时</template>
          <template #unchecked>暂停</template>
        </n-switch>
      </n-space>

      <div ref="logContainer" class="log-container" @scroll="handleScroll">
        <div v-for="(log, index) in filteredLogs" :key="index" :class="['log-line', getLogLevelClass(log)]">
          <span class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</span>
          <span class="log-content" v-html="highlightSearch(log.content)"></span>
        </div>
        <n-empty v-if="filteredLogs.length === 0" description="暂无日志" />
      </div>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { SearchOutline, TrashOutline } from '@vicons/ionicons5'

interface Props {
  containerId: string
}

interface LogLine {
  timestamp: number
  content: string
  level: 'error' | 'warning' | 'info' | 'debug'
}

const props = defineProps<Props>()

const logs = ref<LogLine[]>([])
const searchQuery = ref('')
const logLevel = ref<'all' | 'error' | 'warning' | 'info' | 'debug'>('all')
const autoScroll = ref(true)
const followLogs = ref(true)
const logContainer = ref<HTMLElement>()

let ws: WebSocket | null = null

const levelOptions = [
  { label: '全部', value: 'all' },
  { label: '错误', value: 'error' },
  { label: '警告', value: 'warning' },
  { label: '信息', value: 'info' },
  { label: '调试', value: 'debug' }
]

const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    if (logLevel.value !== 'all' && log.level !== logLevel.value) {
      return false
    }
    if (searchQuery.value && !log.content.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      return false
    }
    return true
  })
})

const getLogLevelClass = (log: LogLine) => {
  return `log-level-${log.level}`
}

const detectLogLevel = (content: string): LogLine['level'] => {
  const lowerContent = content.toLowerCase()
  if (lowerContent.includes('error') || lowerContent.includes('fail') || lowerContent.includes('exception')) {
    return 'error'
  } else if (lowerContent.includes('warning') || lowerContent.includes('warn')) {
    return 'warning'
  } else if (lowerContent.includes('debug')) {
    return 'debug'
  }
  return 'info'
}

const formatTimestamp = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', { hour12: false }) + '.' + date.getMilliseconds().toString().padStart(3, '0')
}

const highlightSearch = (content: string) => {
  if (!searchQuery.value) {
    return content
  }
  const regex = new RegExp(`(${searchQuery.value})`, 'gi')
  return content.replace(regex, '<mark>$1</mark>')
}

const handleSearch = () => {
  if (autoScroll.value) {
    scrollToBottom()
  }
}

const handleFilter = () => {
  if (autoScroll.value) {
    scrollToBottom()
  }
}

const handleClear = () => {
  logs.value = []
}

const handleScroll = () => {
  if (logContainer.value) {
    const { scrollTop, scrollHeight, clientHeight } = logContainer.value
    autoScroll.value = scrollTop + clientHeight >= scrollHeight - 100
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const handleFollowChange = (value: boolean) => {
  if (value) {
    startFollowLogs()
  } else {
    stopFollowLogs()
  }
}

const startFollowLogs = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'logs_stop'
    }))
    ws.close()
    ws = null
  } else if (ws) {
    ws.close()
    ws = null
  }

  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws${token ? `?token=${encodeURIComponent(token)}` : ''}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    ws?.send(JSON.stringify({
      type: 'logs',
      data: {
        container_id: props.containerId,
        tail: '100',
        follow: true
      }
    }))
  }

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)

      switch (message.type) {
        case 'logs_started':
          console.log('Logs streaming started')
          break
        case 'logs_update':
          const content = atob(message.data.data)
          const newLog: LogLine = {
            timestamp: Date.now(),
            content,
            level: detectLogLevel(content)
          }
          logs.value.push(newLog)
          if (autoScroll.value) {
            scrollToBottom()
          }
          break
        case 'logs_error':
          console.error('Logs error:', message.data.error)
          break
        case 'logs_end':
          console.log('Logs streaming completed')
          break
      }
    } catch (error) {
      console.error('Error parsing logs message:', error)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  ws.onclose = () => {
    console.log('WebSocket disconnected')
  }
}

const stopFollowLogs = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'logs_stop'
    }))
    ws.close()
    ws = null
  } else if (ws) {
    ws.close()
    ws = null
  }
}

const loadInitialLogs = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws${token ? `?token=${encodeURIComponent(token)}` : ''}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    ws?.send(JSON.stringify({
      type: 'logs',
      data: {
        container_id: props.containerId,
        tail: '100',
        follow: false
      }
    }))
  }

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)

      switch (message.type) {
        case 'logs_started':
          console.log('Logs streaming started')
          break
        case 'logs_update':
          const content = atob(message.data.data)
          const newLog: LogLine = {
            timestamp: Date.now(),
            content,
            level: detectLogLevel(content)
          }
          logs.value.push(newLog)
          break
        case 'logs_end':
          console.log('Logs streaming completed')
          scrollToBottom()
          if (followLogs.value) {
            startFollowLogs()
          }
          break
        case 'logs_error':
          console.error('Logs error:', message.data.error)
          break
      }
    } catch (error) {
      console.error('Error parsing logs message:', error)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  ws.onclose = () => {
    console.log('WebSocket disconnected')
  }
}

onMounted(() => {
  loadInitialLogs()
})

onBeforeUnmount(() => {
  stopFollowLogs()
})
</script>

<style scoped>
.log-viewer {
  width: 100%;
  height: 500px;
}

.log-container {
  flex: 1;
  overflow-y: auto;
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  line-height: 1.6;
}

.log-line {
  padding: 4px 0;
  border-bottom: 1px solid #333;
}

.log-timestamp {
  color: #888;
  margin-right: 12px;
}

.log-content {
  color: #e5e5e5;
}

.log-level-error .log-content {
  color: #f14c4c;
}

.log-level-warning .log-content {
  color: #d29922;
}

.log-level-info .log-content {
  color: #4ec9b0;
}

.log-level-debug .log-content {
  color: #6a9955;
}

:deep(mark) {
  background: #ffd700;
  color: #000;
  padding: 0 2px;
  border-radius: 2px;
}
</style>
