<template>
  <div class="container-detail">
    <div class="view-header">
      <div class="container-title-group">
        <div class="title-info">
          <h1>{{ container?.names?.[0]?.replace(/^\//, '') || containerId }}</h1>
          <n-tag :type="container?.state === 'running' ? 'success' : 'default'" round size="small" :bordered="false">
            <template #icon>
              <div :style="{ width: '6px', height: '6px', borderRadius: '50%', backgroundColor: 'currentColor', marginRight: '4px' }"></div>
            </template>
            {{ container?.state === 'running' ? 'RUNNING' : container?.state?.toUpperCase() }}
          </n-tag>
        </div>
      </div>
      <div class="header-actions">
        <n-button size="medium" type="warning" @click="handleStop" :disabled="container?.state !== 'running'" :loading="loading">
          <template #icon>
            <n-icon :component="StopOutline" />
          </template>
          停止
        </n-button>
        <n-button size="medium" type="info" @click="handleRestart" :loading="loading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
          重启
        </n-button>
        <n-button size="medium" type="success" @click="handleStart" :disabled="container?.state === 'running'" :loading="loading">
          <template #icon>
            <n-icon :component="PlayOutline" />
          </template>
          启动
        </n-button>
        <n-button size="medium" type="error" @click="handleDelete" :loading="loading">
          <template #icon>
            <n-icon :component="TrashOutline" />
          </template>
          删除
        </n-button>
        <n-button size="medium" type="primary" @click="openTerminal" :disabled="container?.state !== 'running'">
          <template #icon>
            <n-icon :component="TerminalOutline" />
          </template>
          终端
        </n-button>
      </div>
    </div>

    <n-space vertical :size="24">
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-card title="实时统计" size="small" :bordered="true">
            <n-grid :cols="2" :y-gap="8">
              <n-gi>
                <div class="mini-stat">
                  <span class="mini-stat-label">CPU</span>
                  <span class="mini-stat-value">{{ (stats?.cpu_percent || 0).toFixed(2) }}%</span>
                </div>
              </n-gi>
              <n-gi>
                <div class="mini-stat">
                  <span class="mini-stat-label">MEM</span>
                  <span class="mini-stat-value">{{ formatFileSize(stats?.memory_usage || 0) }}</span>
                </div>
              </n-gi>
              <n-gi>
                <div class="mini-stat">
                  <span class="mini-stat-label">内存使用</span>
                  <span class="mini-stat-value">{{ formatFileSize(stats?.memory_usage || 0) }}</span>
                </div>
              </n-gi>
              <n-gi>
                <div class="mini-stat">
                  <span class="mini-stat-label">内存限制</span>
                  <span class="mini-stat-value">{{ formatFileSize(stats?.memory_limit || 0) }}</span>
                </div>
              </n-gi>
            </n-grid>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="容器信息" size="small" :bordered="true">
            <div class="info-list compact">
              <div class="info-item">
                <span class="info-label">容器 ID</span>
                <span class="info-value">
                  <n-text code>{{ container?.id?.substring(0, 12) || '-' }}</n-text>
                </span>
              </div>
              <div class="info-item">
                <span class="info-label">镜像</span>
                <span class="info-value">{{ container?.image || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">创建时间</span>
                <span class="info-value">{{ formatDate(container?.created) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">命令</span>
                <span class="info-value">{{ container?.command || '-' }}</span>
              </div>
            </div>
          </n-card>
        </n-gi>
      </n-grid>

      <n-card title="容器日志" :bordered="true">
        <template #header-extra>
          <n-space>
            <n-button size="small" @click="loadLogs" :loading="logsLoading">
              <template #icon>
                <n-icon :component="RefreshOutline" />
              </template>
              刷新
            </n-button>
            <n-button size="small" type="primary" @click="viewFullLogs">
              查看完整日志
            </n-button>
          </n-space>
        </template>
        <div class="logs-container" v-if="logs">
          <pre class="logs-content">{{ logs }}</pre>
        </div>
        <n-empty v-else description="暂无日志" />
      </n-card>

      <n-card title="文件管理" :bordered="true">
        <template #header-extra>
          <n-space>
            <n-input v-model:value="filePath" placeholder="输入容器内路径，如 /app" style="width: 300px" />
            <n-button size="small" type="primary" @click="handleDownload" :loading="fileLoading">
              <template #icon>
                <n-icon :component="DownloadOutline" />
              </template>
              下载
            </n-button>
            <n-upload :show-file-list="false" @before-upload="handleUpload" accept="*/*">
              <n-button size="small" type="info" :loading="uploadLoading">
                <template #icon>
                  <n-icon :component="CloudUploadOutline" />
                </template>
                上传
              </n-button>
            </n-upload>
          </n-space>
        </template>
        <n-alert type="info" :show-icon="false" style="margin-bottom: 16px">
          输入容器内文件/目录路径，点击下载可将文件下载到本地。点击上传可选择本地文件上传到容器当前路径。
        </n-alert>
        <n-data-table
          v-if="fileList.length > 0"
          :columns="fileColumns"
          :data="fileList"
          :bordered="false"
        />
        <n-empty v-else description="暂无文件信息，请在上方输入路径并点击下载或上传" />
      </n-card>

      <n-card title="进程列表" :bordered="true">
        <template #header-extra>
          <n-button size="small" @click="loadProcesses" :loading="processesLoading">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            刷新
          </n-button>
        </template>
        <n-data-table
          v-if="processes.length > 0"
          :columns="processColumns"
          :data="processes"
          :bordered="false"
          :pagination="{ pageSize: 20 }"
        />
        <n-empty v-else description="暂无进程信息" />
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, useDialog, NButton, NIcon, NCard, NDataTable, NTag, NText, NGrid, NGi, NEmpty, NInput, NSpace, NAlert, NUpload } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  PlayOutline,
  StopOutline,
  RefreshOutline,
  TrashOutline,
  TerminalOutline,
  DownloadOutline,
  CloudUploadOutline,
  FolderOutline,
  DocumentOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()
const message = useMessage()
const dialog = useDialog()

const containerId = computed(() => route.params.id as string)

const container = computed(() => containerStore.getContainerById(containerId.value))

const stats = ref<any>(null)
const logs = ref<string>('')
const processes = ref<any[]>([])
const loading = ref(false)
const logsLoading = ref(false)
const processesLoading = ref(false)
const filePath = ref('/')
const fileList = ref<any[]>([])
const fileLoading = ref(false)
const uploadLoading = ref(false)

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (timestamp: number | string | undefined) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

const processColumns: DataTableColumns<{
  user: string
  pid: string
  cpu: string
  mem: string
  vsz: string
  rss: string
  tty: string
  stat: string
  start: string
  time: string
  command: string
}> = [
  { title: 'USER', key: 'user', width: 100 },
  { title: 'PID', key: 'pid', width: 80 },
  { title: 'CPU%', key: 'cpu', width: 80 },
  { title: 'MEM%', key: 'mem', width: 80 },
  { title: 'VSZ', key: 'vsz', width: 100 },
  { title: 'RSS', key: 'rss', width: 100 },
  { title: 'TTY', key: 'tty', width: 80 },
  { title: 'STAT', key: 'stat', width: 80 },
  { title: 'START', key: 'start', width: 100 },
  { title: 'TIME', key: 'time', width: 100 },
  {
    title: 'COMMAND',
    key: 'command',
    ellipsis: { tooltip: true }
  }
]

const fileColumns: DataTableColumns<{
  name: string
  path: string
  type: string
  size: string
  mtime: string
}> = [
  {
    title: '名称',
    key: 'name',
    render(row: any) {
      return h('div', {
        style: { display: 'flex', alignItems: 'center', gap: '8px' }
      }, [
        h(NIcon, {
          component: row.is_dir ? FolderOutline : DocumentOutline,
          size: 18,
          color: row.is_dir ? '#f0a020' : '#2080f0'
        }),
        h('span', row.name)
      ])
    }
  },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '类型', key: 'type', width: 100 },
  { title: '大小', key: 'size', width: 100 },
  { title: '修改时间', key: 'mtime', width: 150 }
]

const loadStats = async () => {
  try {
    const data = await containerStore.getContainerStats(containerId.value)
    stats.value = {
      cpu_percent: data?.cpu ? data.cpu * 100 : 0,
      memory_usage: data?.memory || 0
    }
  } catch (error: any) {
    console.error('加载统计信息失败:', error)
  }
}

const loadLogs = async () => {
  try {
    logsLoading.value = true
    const result = await containerStore.getContainerLogs(containerId.value)
    logs.value = result?.logs || result || ''
  } catch (error: any) {
    message.error('加载日志失败: ' + error.message)
  } finally {
    logsLoading.value = false
  }
}

const viewFullLogs = () => {
  router.push({ name: 'ContainerLogs', params: { id: containerId.value } })
}

const openTerminal = () => {
  router.push({ name: 'ContainerTerminal', params: { id: containerId.value } })
}

const handleDownload = async () => {
  if (!filePath.value) {
    message.warning('请输入文件路径')
    return
  }
  try {
    fileLoading.value = true
    await containerApi.copyFrom(containerId.value, filePath.value)
    message.success('文件下载成功')
    // 添加下载记录到列表
    fileList.value.unshift({
      name: filePath.value.split('/').pop() || filePath.value,
      path: filePath.value,
      type: 'file',
      size: '-',
      mtime: new Date().toLocaleString(),
      is_dir: false
    })
  } catch (error: any) {
    message.error('下载失败: ' + error.message)
  } finally {
    fileLoading.value = false
  }
}

const handleUpload = async (options: any) => {
  if (!filePath.value) {
    message.warning('请输入目标路径')
    return false
  }
  const file = options.file.file
  if (!file) return false

  try {
    uploadLoading.value = true
    // 先上传到服务器临时目录，再复制到容器
    const formData = new FormData()
    formData.append('file', file)
    // 使用上传API
    await containerApi.copyTo(containerId.value, file.name, filePath.value + '/' + file.name)
    message.success('文件上传成功')
    // 添加上传记录到列表
    fileList.value.unshift({
      name: file.name,
      path: filePath.value + '/' + file.name,
      type: 'file',
      size: formatFileSize(file.size),
      mtime: new Date().toLocaleString(),
      is_dir: false
    })
  } catch (error: any) {
    message.error('上传失败: ' + error.message)
  } finally {
    uploadLoading.value = false
  }
  return false
}

onMounted(async () => {
  await containerStore.fetchContainers()
  if (!container.value?.id) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
    return
  }
  await Promise.all([loadStats(), loadLogs(), loadProcesses()])
})
</script>

<style scoped>
.container-detail {
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.container-title-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-info h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.mini-stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mini-stat-label {
  font-size: 11px;
  color: var(--n-text-color-3);
  text-transform: uppercase;
}

.mini-stat-value {
  font-size: 16px;
  font-weight: 700;
}

.info-list {
  display: flex;
  flex-direction: column;
}

.info-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--n-border-color);
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: var(--n-text-color-3);
}

.info-value {
  font-weight: 500;
  text-align: right;
  max-width: 60%;
  word-break: break-all;
}

.logs-container {
  background: var(--n-color-modal);
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow: auto;
}

.logs-content {
  margin: 0;
  font-family: monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
