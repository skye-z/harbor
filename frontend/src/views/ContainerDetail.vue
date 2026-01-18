<template>
  <div class="container-detail">
    <div class="view-header">
      <div class="container-title-group">
        <div class="title-info">
          <h1>{{ container?.Name?.replace(/^\//, '') || containerId }}</h1>
          <n-tag :type="getStateType(container?.State?.Status)" round size="small" :bordered="false">
            <template #icon>
              <div :style="{ width: '6px', height: '6px', borderRadius: '50%', backgroundColor: 'currentColor', marginRight: '4px' }"></div>
            </template>
            {{ formatState(container?.State?.Status) }}
          </n-tag>
        </div>
      </div>
      <div class="header-actions">
        <n-button size="medium" type="warning" @click="handleStop" :disabled="container?.State?.Status !== 'running'" :loading="loading">
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
        <n-button size="medium" type="success" @click="handleStart" :disabled="container?.State?.Status === 'running'" :loading="loading">
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
        <n-button size="medium" type="primary" @click="openTerminal" :disabled="container?.State?.Status !== 'running'">
          <template #icon>
            <n-icon :component="TerminalOutline" />
          </template>
          终端
        </n-button>
      </div>
    </div>

    <n-space vertical :size="20">
      <n-grid :cols="3" :x-gap="16" :y-gap="16">
        <n-gi>
          <n-card title="CPU使用" size="small" :bordered="true">
            <div class="stat-card">
              <n-statistic tabular-nums>
                <template #prefix>
                  <n-icon :component="ServerOutline" />
                </template>
                {{ (stats?.cpu_percent || 0).toFixed(2) }}%
              </n-statistic>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="内存使用" size="small" :bordered="true">
            <div class="stat-card">
              <n-statistic tabular-nums>
                <template #prefix>
                  <n-icon :component="DiscOutline" />
                </template>
                {{ formatFileSize(stats?.memory_usage || 0) }}
              </n-statistic>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="内存限制" size="small" :bordered="true">
            <div class="stat-card">
              <n-statistic tabular-nums>
                <template #prefix>
                  <n-icon :component="HardwareChipOutline" />
                </template>
                {{ formatFileSize(stats?.memory_limit || 0) }}
              </n-statistic>
            </div>
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-card title="基本信息" size="small" :bordered="true">
            <n-descriptions :column="1" label-placement="left" bordered>
              <n-descriptions-item label="容器ID">
                <n-text code>{{ container?.Id?.substring(0, 12) || '-' }}</n-text>
              </n-descriptions-item>
              <n-descriptions-item label="容器名称">
                {{ container?.Name?.replace(/^\//, '') || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="镜像">
                {{ container?.Config?.Image || container?.Image || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">
                {{ formatDate(container?.Created) }}
              </n-descriptions-item>
              <n-descriptions-item label="启动命令">
                <n-text code>{{ container?.Path }} {{ (container?.Args || []).join(' ') }}</n-text>
              </n-descriptions-item>
              <n-descriptions-item label="工作目录">
                {{ container?.Config?.WorkingDir || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="重启次数">
                {{ container?.RestartCount || 0 }}
              </n-descriptions-item>
              <n-descriptions-item label="存储驱动">
                {{ container?.Driver || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="平台">
                {{ container?.Platform || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="进程ID">
                {{ container?.State?.Pid || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="运行时">
                {{ container?.HostConfig?.Runtime || '-' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="状态信息" size="small" :bordered="true">
            <n-descriptions :column="1" label-placement="left" bordered>
              <n-descriptions-item label="状态">
                <n-tag :type="getStateType(container?.State?.Status)" size="small">
                  {{ formatState(container?.State?.Status) }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="运行中">
                <n-tag :type="container?.State?.Running ? 'success' : 'error'" size="small">
                  {{ container?.State?.Running ? '是' : '否' }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="退出码">
                {{ container?.State?.ExitCode ?? '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="启动时间">
                {{ formatDate(container?.State?.StartedAt) }}
              </n-descriptions-item>
              <n-descriptions-item label="结束时间">
                {{ formatDate(container?.State?.FinishedAt) }}
              </n-descriptions-item>
              <n-descriptions-item label="重启策略">
                {{ container?.HostConfig?.RestartPolicy?.Name || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="网络模式">
                {{ container?.HostConfig?.NetworkMode || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="自动删除">
                {{ container?.HostConfig?.AutoRemove ? '是' : '否' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>
      </n-grid>

      <n-card title="端口映射" size="small" :bordered="true">
        <n-data-table
          v-if="portBindings.length > 0"
          :columns="portColumns"
          :data="portBindings"
          :bordered="false"
          :pagination="{ pageSize: 10 }"
        />
        <n-empty v-else description="暂无端口映射" />
      </n-card>

      <n-card title="挂载点" size="small" :bordered="true">
        <n-data-table
          v-if="mounts.length > 0"
          :columns="mountColumns"
          :data="mounts"
          :bordered="false"
          :pagination="{ pageSize: 10 }"
        />
        <n-empty v-else description="暂无挂载点" />
      </n-card>

      <n-card title="环境变量" size="small" :bordered="true">
        <n-data-table
          v-if="envVars.length > 0"
          :columns="envColumns"
          :data="envVars"
          :bordered="false"
          :pagination="{ pageSize: 10 }"
        />
        <n-empty v-else description="暂无环境变量" />
      </n-card>

      <n-card title="网络信息" size="small" :bordered="true">
        <n-data-table
          v-if="networkInfos.length > 0"
          :columns="networkColumns"
          :data="networkInfos"
          :bordered="false"
          :pagination="{ pageSize: 10 }"
        />
        <n-empty v-else description="暂无网络信息" />
      </n-card>

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
        <div class="logs-container" v-if="logs.length > 0">
          <pre class="logs-content"><code
            v-for="(line, index) in logs"
            :key="index"
            :class="getLogLineClass(line)"
          >{{ line }}</code></pre>
        </div>
        <n-empty v-else description="暂无日志" />
      </n-card>

      <n-card title="文件管理" :bordered="true">
        <template #header-extra>
          <n-space>
            <n-button size="small" @click="handleNavigateUp" :disabled="filePath === '/'">
              <template #icon>
                <n-icon :component="FolderOpenOutline" />
              </template>
              返回上级
            </n-button>
            <n-input v-model:value="currentPath" placeholder="输入容器内路径，如 /app" style="width: 300px" />
            <n-button size="small" type="primary" @click="loadFileList" :loading="fileLoading">
              <template #icon>
                <n-icon :component="SearchOutline" />
              </template>
              浏览
            </n-button>
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
          输入容器内文件/目录路径，点击浏览查看目录内容，点击下载可将文件下载到本地。点击上传可选择本地文件上传到容器当前路径。
        </n-alert>
        <n-data-table
          v-if="fileList.length > 0"
          :columns="fileColumns"
          :data="fileList"
          :bordered="false"
          :pagination="{ pageSize: 20 }"
        />
        <n-empty v-else description="暂无文件信息，请在上方输入路径并点击浏览" />
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
import { useMessage, useDialog, NButton, NIcon, NCard, NDataTable, NTag, NText, NGrid, NGi, NEmpty, NInput, NSpace, NAlert, NUpload, NStatistic, NDescriptions, NDescriptionsItem } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { formatContainerState, getContainerStateType } from '../plugins/utils/container'
import {
  PlayOutline,
  StopOutline,
  RefreshOutline,
  TrashOutline,
  TerminalOutline,
  DownloadOutline,
  CloudUploadOutline,
  FolderOutline,
  FolderOpenOutline,
  DocumentOutline,
  HardwareChipOutline,
  DiscOutline,
  ServerOutline,
  SearchOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()
const message = useMessage()
const dialog = useDialog()

const containerId = computed(() => route.params.id as string)

const container = ref<any>(null)
const containerLoading = ref(true)

const stats = ref<any>(null)
const logs = ref<string[]>([])

const getLogLineClass = (line: string) => {
  const lower = line.toLowerCase()
  if (lower.includes('error') || lower.includes('err') || lower.includes('fatal') || lower.includes('exception')) return 'log-error'
  if (lower.includes('warn')) return 'log-warn'
  if (lower.includes('info')) return 'log-info'
  if (lower.includes('debug')) return 'log-debug'
  return ''
}
const processes = ref<any[]>([])
const loading = ref(false)
const logsLoading = ref(false)
const processesLoading = ref(false)
const currentPath = ref('/')
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

const getStateType = (state: string) => {
  if (state === 'running') return 'success'
  if (state === 'paused') return 'warning'
  return 'default'
}

const formatState = (state: string) => {
  if (state === 'running') return '运行中'
  if (state === 'exited') return '已停止'
  if (state === 'paused') return '已暂停'
  if (state === 'created') return '已创建'
  if (state === 'restarting') return '重启中'
  if (state === 'removing') return '删除中'
  if (state === 'dead') return '已终止'
  return state || '未知'
}

const portBindings = computed(() => {
  const bindings = container.value?.HostConfig?.PortBindings || {}
  const result = []
  for (const [port, configs] of Object.entries(bindings)) {
    const [containerPort, protocol] = port.split('/')
    const configArray = configs as Array<{ HostIp: string; HostPort: string }>
    result.push({
      port,
      containerPort,
      protocol: protocol?.toUpperCase(),
      hostIp: configArray[0]?.HostIp || '',
      hostPort: configArray[0]?.HostPort || ''
    })
  }
  return result
})

const mounts = computed(() => {
  return container.value?.Mounts || []
})

const envVars = computed(() => {
  const env = container.value?.Config?.Env || []
  return env.map((e: string) => {
    const [key, ...valueParts] = e.split('=')
    return { key, value: valueParts.join('=') }
  })
})

const networkInfos = computed(() => {
  const networks = container.value?.NetworkSettings?.Networks || {}
  return Object.entries(networks).map(([name, info]: [string, any]) => ({
    network: name,
    ipAddress: info?.IPAddress || '-',
    gateway: info?.Gateway || '-',
    macAddress: info?.MacAddress || '-',
    prefixLen: info?.IPPrefixLen || '-'
  }))
})

const portColumns: DataTableColumns<any> = [
  { title: '容器端口', key: 'containerPort', width: 120 },
  { title: '协议', key: 'protocol', width: 80 },
  { title: '宿主机端口', key: 'hostPort', width: 120 },
  { title: '宿主机IP', key: 'hostIp', width: 150 }
]

const mountColumns: DataTableColumns<any> = [
  { title: '类型', key: 'Type', width: 80 },
  { title: '源路径', key: 'Source', ellipsis: { tooltip: true } },
  { title: '目标路径', key: 'Destination', ellipsis: { tooltip: true } },
  { title: '模式', key: 'Mode', width: 80 },
  {
    title: '读写',
    key: 'RW',
    width: 80,
    render: (row: any) => row.RW ? '是' : '否'
  }
]

const envColumns: DataTableColumns<any> = [
  { title: '变量名', key: 'key', width: 200 },
  { title: '变量值', key: 'value', ellipsis: { tooltip: true } }
]

const networkColumns: DataTableColumns<any> = [
  { title: '网络名称', key: 'network', width: 120 },
  { title: 'IP地址', key: 'ipAddress', width: 150 },
  { title: '网关', key: 'gateway', width: 150 },
  { title: 'MAC地址', key: 'macAddress', width: 180 },
  { title: '前缀长度', key: 'prefixLen', width: 100 }
]

const processColumns: DataTableColumns<any> = [
  { title: 'USER', key: '0', width: 100 },
  { title: 'PID', key: '1', width: 80 },
  { title: 'PPID', key: '2', width: 80 },
  { title: 'C', key: '3', width: 60 },
  { title: 'STIME', key: '4', width: 80 },
  { title: 'TTY', key: '5', width: 60 },
  { title: 'TIME', key: '6', width: 100 },
  { title: 'CMD', key: '7', ellipsis: { tooltip: true } }
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
        style: { display: 'flex', alignItems: 'center', gap: '8px', cursor: row.is_dir ? 'pointer' : 'default', color: row.is_dir ? '#2080f0' : 'inherit' },
        onClick: () => handleNavigate(row)
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

const loadContainerDetail = async () => {
  containerLoading.value = true
  try {
    const data = await containerStore.getContainerInfo(containerId.value)
    container.value = data?.Container || data
  } catch (error: any) {
    message.error('获取容器详情失败: ' + error.message)
  } finally {
    containerLoading.value = false
  }
}

const loadStats = async () => {
  try {
    const data = await containerStore.getContainerStats(containerId.value)
    stats.value = {
      cpu_percent: data?.cpu ? data.cpu * 100 : 0,
      memory_usage: data?.memory || 0,
      memory_limit: data?.memory_limit || 0
    }
  } catch (error: any) {
    console.error('加载统计信息失败:', error)
  }
}

const loadLogs = async () => {
  try {
    logsLoading.value = true
    const result = await containerStore.getContainerLogs(containerId.value)
    const logText = result?.logs || result || ''
    logs.value = logText.split('\n').filter((line: string) => line.length > 0)
  } catch (error: any) {
    message.error('加载日志失败: ' + error.message)
  } finally {
    logsLoading.value = false
  }
}

const viewFullLogs = () => {
  router.push({ name: 'ContainerLogs', params: { id: containerId.value } })
}

const loadProcesses = async () => {
  try {
    processesLoading.value = true
    const result = await containerStore.getContainerProcesses(containerId.value)
    const processesData = result?.Processes || []
    const titles = result?.Titles || []
    
    processes.value = processesData.map((row: string[]) => {
      const obj: any = {}
      row.forEach((cell, index) => {
        obj[index] = cell
      })
      return obj
    })
  } catch (error: any) {
    message.error('加载进程列表失败: ' + error.message)
  } finally {
    processesLoading.value = false
  }
}

const handleStop = async () => {
  dialog.warning({
    title: '确认停止容器',
    content: `确定要停止容器 "${container.value?.Name?.replace(/^\//, '')}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        await containerStore.operationContainer(containerId.value, 'stop')
        message.success('容器已停止')
        await loadContainerDetail()
      } catch (error: any) {
        message.error('停止容器失败: ' + error.message)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleStart = async () => {
  try {
    loading.value = true
    await containerStore.operationContainer(containerId.value, 'start')
    message.success('容器已启动')
    await loadContainerDetail()
  } catch (error: any) {
    message.error('启动容器失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleRestart = async () => {
  dialog.warning({
    title: '确认重启容器',
    content: `确定要重启容器 "${container.value?.Name?.replace(/^\//, '')}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        await containerStore.operationContainer(containerId.value, 'restart')
        message.success('容器已重启')
        await loadContainerDetail()
      } catch (error: any) {
        message.error('重启容器失败: ' + error.message)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleDelete = async () => {
  dialog.warning({
    title: '确认删除容器',
    content: `确定要删除容器 "${container.value?.Name?.replace(/^\//, '')}" 吗？此操作不可恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        await containerStore.operationContainer(containerId.value, 'remove')
        message.success('容器已删除')
        router.push({ name: 'Containers' })
      } catch (error: any) {
        message.error('删除容器失败: ' + error.message)
      } finally {
        loading.value = false
      }
    }
  })
}

const openTerminal = () => {
  router.push({ name: 'ContainerTerminal', params: { id: containerId.value } })
}

const loadFileList = async () => {
  if (!currentPath.value) {
    message.warning('请输入路径')
    return
  }
  try {
    fileLoading.value = true
    const result = await containerStore.listContainerFiles(containerId.value, currentPath.value)
    fileList.value = result || []
  } catch (error: any) {
    message.error('获取文件列表失败: ' + error.message)
    fileList.value = []
  } finally {
    fileLoading.value = false
  }
}

const handleNavigateUp = () => {
  if (currentPath.value === '/') return
  const parent = currentPath.value.split('/').slice(0, -1).join('/')
  currentPath.value = parent || '/'
  loadFileList()
}

const handleNavigate = (row: any) => {
  if (row.is_dir) {
    currentPath.value = row.path
    loadFileList()
  }
}

const handleDownload = async () => {
  if (!currentPath.value) {
    message.warning('请输入文件路径')
    return
  }
  try {
    fileLoading.value = true
    await containerStore.copyFromContainer(containerId.value, currentPath.value)
    message.success('文件下载成功')
  } catch (error: any) {
    message.error('下载失败: ' + error.message)
  } finally {
    fileLoading.value = false
  }
}

const handleUpload = async (options: any) => {
  if (!currentPath.value) {
    message.warning('请输入目标路径')
    return false
  }
  const file = options.file.file
  if (!file) return false

  try {
    uploadLoading.value = true
    await containerStore.copyToContainer(containerId.value, file.name, currentPath.value + '/' + file.name)
    message.success('文件上传成功')
    loadFileList()
  } catch (error: any) {
    message.error('上传失败: ' + error.message)
  } finally {
    uploadLoading.value = false
  }
  return false
}

onMounted(async () => {
  await loadContainerDetail()
  if (!container.value?.Id) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
    return
  }
  await loadStats()
  await loadLogs()
  setTimeout(() => {
    loadFileList()
  }, 1000)
  if (container.value?.State?.Status === 'running') {
    setTimeout(() => {
      loadProcesses()
    }, 3000)
  }
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

.stat-card {
  text-align: center;
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

.log-debug {
  color: #8b949e;
}
</style>
