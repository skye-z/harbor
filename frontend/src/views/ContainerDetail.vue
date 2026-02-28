<template>
  <div class="container-detail">

    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ container?.Name?.replace(/^\//, '') || containerId }}</h1>
          <div class="header-actions">
            <n-button size="medium" type="warning" @click="handleStop"
              :disabled="container?.State?.Status !== 'running'" :loading="loading">
              <template #icon>
                <n-icon :component="StopOutline" />
              </template>
              停止
            </n-button>
            <n-button size="medium" type="info" @click="handleRestart"
              :disabled="container?.State?.Status !== 'running'" :loading="loading">
              <template #icon>
                <n-icon :component="RefreshOutline" />
              </template>
              重启
            </n-button>
            <n-button size="medium" type="success" @click="handleStart"
              :disabled="container?.State?.Status === 'running'" :loading="loading">
              <template #icon>
                <n-icon :component="PlayOutline" />
              </template>
              启动
            </n-button>
            <n-button size="medium" type="error" @click="handleDelete"
              :disabled="container?.State?.Status === 'running'" :loading="loading">
              <template #icon>
                <n-icon :component="TrashOutline" />
              </template>
              删除
            </n-button>
          </div>
        </div>
        <div class="subtitle-text" v-if="container">
          {{ container?.Path }} {{ (container?.Args || []).join(' ') }}
        </div>
      </div>
    </div>

    <n-grid item-responsive x-gap="10" cols="24">
      <n-grid-item span="24 760:9 900:8">
        <n-card size="small" title="基本信息" style="margin-bottom: 10px;">
          <template #header-extra>
            <n-tag :type="getStateType(container?.State?.Status)" round size="small" :bordered="false">
              <template #icon>
                <div
                  :style="{ width: '6px', height: '6px', borderRadius: '50%', backgroundColor: 'currentColor', marginRight: '4px' }">
                </div>
              </template>
              {{ formatState(container?.State?.Status) }}
            </n-tag>
          </template>
          <n-descriptions :column="1" label-placement="left">
            <n-descriptions-item label="编号">
              {{ container?.Id?.substring(0, 12) || '-' }}
            </n-descriptions-item>
            <n-descriptions-item label="镜像">
              <div style="cursor: pointer;" @click="goToImage(container?.Image)">
                {{ container?.Config?.Image || container?.Image || '-' }}
              </div>
            </n-descriptions-item>
            <n-descriptions-item label="平台">
              {{
                (container?.ImageManifestDescriptor?.platform?.os + '/' +
                  container?.ImageManifestDescriptor?.platform?.architecture)
                || '-' }}
            </n-descriptions-item>
            <n-descriptions-item label="重启次数">
              {{ container?.RestartCount || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="重启策略">
              {{ container?.HostConfig?.RestartPolicy?.Name || '-' }}
            </n-descriptions-item>
            <n-descriptions-item label="工作目录">
              {{ container?.Config?.WorkingDir || '-' }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
        <div style="margin-bottom: 10px;">
          <n-button-group style="width: 100%;">
            <n-button size="medium" type="primary" style="width: 33.3%;" @click="openTerminal"
              :disabled="container?.State?.Status !== 'running'">
              <template #icon>
                <n-icon :component="TerminalOutline" />
              </template>
              终端
            </n-button>
            <n-button size="medium" type="primary" style="width: 33.3%;" @click="openLogs"
              :disabled="container?.State?.Status !== 'running'">
              <template #icon>
                <n-icon :component="ListOutline" />
              </template>
              日志
            </n-button>
            <n-button size="medium" type="primary" style="width: 33.3%;" @click="openFile"
              :disabled="container?.State?.Status !== 'running'">
              <template #icon>
                <n-icon :component="FolderOutline" />
              </template>
              文件
            </n-button>
          </n-button-group>
        </div>
        <n-card size="small" style="margin-bottom: 10px;">
          <n-descriptions :column="2" label-placement="left">
            <n-descriptions-item label="CPU">
              {{ (stats?.cpu_percent || 0).toFixed(2) }}%
            </n-descriptions-item>
            <n-descriptions-item label="内存">
              {{ formatFileSize(stats?.memory_usage || 0) }}
            </n-descriptions-item>
            <n-descriptions-item label="上传">
              {{ formatFileSize(stats?.networkTx || 0) }}
            </n-descriptions-item>
            <n-descriptions-item label="下载">
              {{ formatFileSize(stats?.networkRx || 0) }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
        <n-card size="small" style="margin-bottom: 10px;">
          <n-space v-if="displayPorts.length" wrap>
            <template v-for="item in displayPorts" :key="item.key">
              <n-tag v-if="item.hasBinding" size="small" :bordered="false" type="success" style="cursor: pointer"
                @click.stop="openPort(item.port)">
                {{ item.port }}
                <template #icon>
                  <n-icon :component="OpenOutline" />
                </template>
              </n-tag>
              <n-tag v-else size="small" :bordered="false">
                {{ item.port }}
              </n-tag>
            </template>
          </n-space>
          <n-tag v-else size="small" :bordered="false">
            -
          </n-tag>
        </n-card>
      </n-grid-item>
      <n-grid-item span="24 760:15 900:16">
        <n-card size="small" title="启动信息" style="margin-bottom: 10px;">
          <n-descriptions :column="2" label-placement="left">
            <n-descriptions-item label="创建时间" content-style="width: 140px">
              <n-time v-if="container?.Created" :time="new Date(container?.Created)" type="datetime" />
            </n-descriptions-item>
            <n-descriptions-item label="进程ID" content-style="width: 140px">
              {{ container?.State?.Pid || '-' }}
            </n-descriptions-item>
            <n-descriptions-item label="启动时间" content-style="width: 140px">
              <n-time v-if="container?.State?.StartedAt" :time="new Date(container?.State?.StartedAt)"
                type="datetime" />
            </n-descriptions-item>
            <n-descriptions-item label="退出码" content-style="width: 140px">
              {{ container?.State?.ExitCode ?? '-' }}
            </n-descriptions-item>
            <n-descriptions-item label="停止时间" content-style="width: 140px">
              <n-time v-if="container?.State?.FinishedAt" :time="new Date(container?.State?.FinishedAt)"
                type="datetime" />
            </n-descriptions-item>
            <n-descriptions-item label="自动删除" content-style="width: 140px">
              {{ container?.HostConfig?.AutoRemove ? '是' : '否' }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
        <n-card v-if="mounts.length" size="small" title="存储映射" style="margin-bottom: 10px;">
          <template #header-extra>
            <div style="color: #999;">{{ '由 ' + container?.Driver + ' 驱动' || '-' }}</div>
          </template>
          <n-space vertical :size="8">
            <div v-for="(mount, index) in mounts" :key="index">
              <n-button-group style="width: 100%;">
                <n-button size="small" style="width: calc(50% - 50px);" tertiary @click="copyText(mount.Source)">{{
                  mount.Source }}</n-button>
                <n-button size="small" style="width: 100px;" strong secondary :type="mount.RW ? 'success' : 'default'"
                  @click="copyText(`docker cp ${mount.Source} ${container?.Name?.replace(/^\//, '')}:${mount.Destination}`)">
                  {{ mount.RW ? '读写' : '只读' }}
                </n-button>
                <n-button size="small" style="width: calc(50% - 50px);" tertiary @click="copyText(mount.Destination)">{{
                  mount.Destination }}</n-button>
              </n-button-group>
            </div>
          </n-space>
        </n-card>
        <n-card size="small" title="环境变量" style="margin-bottom: 10px;">
          <n-space v-if="container?.Config?.Env?.length" wrap>
            <n-tag v-for="(env, index) in container?.Config?.Env" :key="index" size="small" style="cursor: pointer"
              @click="copyText(env.split('=').slice(1).join('='))">
              {{ env.split('=')[0] }}={{ env.split('=').slice(1).join('=').substring(0, 20) }}{{
                env.split('=').slice(1).join('=').length > 20 ? '...' : '' }}
            </n-tag>
          </n-space>
          <n-empty v-else description="暂无环境变量" />
        </n-card>
        <n-card v-if="networkInfos.length" size="small" title="网络信息">
          <div v-for="item in networkInfos" style="border-radius: 4px; margin-bottom: 8px;">
            <div style="margin-bottom: 4px;">
              <n-tag size="small" type="success">{{ item.network }}</n-tag>
              {{ item.gateway }}/{{ item.prefixLen }}
            </div>
            <div>
              <n-tag size="small">分配地址: {{ item.ipAddress }}</n-tag>
              <n-tag v-if="item.macAddress !== '-'" size="small" style="margin-left: 8px;">MAC: {{ item.macAddress
                }}</n-tag>
            </div>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>
    <div style="margin-top: 20px; text-align: center;">
      <n-button @click="showRawData = true">显示原始数据</n-button>
    </div>
    <n-modal v-model:show="showRawData" preset="card" title="容器原始数据" style="width: 800px; max-width: 90%;">
      <n-scrollbar style="max-height: 600px">
        <n-code :code="rawData" language="json" word-wrap show-line-numbers />
      </n-scrollbar>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, useDialog, NIcon, NModal, NCode } from 'naive-ui'
import {
  PlayOutline,
  StopOutline,
  RefreshOutline,
  TrashOutline,
  TerminalOutline,
  FolderOutline,
  OpenOutline,
  ListOutline
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

const loading = ref(false)

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const copyText = (text: string) => {
  navigator.clipboard.writeText(text)
  message.success('已复制')
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

const openPort = (port: number) => {
  window.open(`http://localhost:${port}`, '_blank')
}

const displayPorts = computed(() => {
  const ports = container.value?.NetworkSettings?.Ports || {}
  const portSet = new Set<string>()
  const result: { port: number; key: string; hasBinding: boolean }[] = []

  for (const [key, value] of Object.entries(ports)) {
    const portNum = key.split('/')[0]
    if (portSet.has(portNum)) continue
    portSet.add(portNum)

    const hasBinding = value !== null && Array.isArray(value) && value.length > 0 && value[0]?.HostPort
    result.push({
      port: Number(portNum),
      key,
      hasBinding
    })
  }
  // 按端口号排序，确保显示顺序稳定
  return result.sort((a, b) => a.port - b.port)
})

const mounts = computed(() => {
  return container.value?.Mounts || []
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

const loadContainerDetail = async () => {
  try {
    const data = await containerStore.refreshContainerInfo(containerId.value)
    container.value = data
  } catch (error: any) {
    message.error('获取容器详情失败: ' + error.message)
  }
}

const loadStats = async () => {
  try {
    const data = await containerStore.getContainerStats(containerId.value)
    stats.value = {
      cpu_percent: data?.cpu ? data.cpu * 100 : 0,
      memory_usage: data?.memory || 0,
      memory_limit: data?.memory_limit || 0,
      networkRx: data?.networkRx || 0,
      networkTx: data?.networkTx || 0
    }
  } catch (error: any) {
    console.error('加载统计信息失败:', error)
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
        containerStore.operationContainer(containerId.value, 'stop')
        message.success('容器已停止')
        loadContainerDetail()
      } catch (error: any) {
        message.error('停止容器失败: ' + error.message)
      } finally {
        setTimeout(() => {
          loading.value = false
        }, 500);
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
        containerStore.operationContainer(containerId.value, 'restart')
        message.success('容器已重启')
        loadContainerDetail()
      } catch (error: any) {
        message.error('重启容器失败: ' + error.message)
      } finally {
        setTimeout(() => {
          loading.value = false
        }, 500);
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
        containerStore.operationContainer(containerId.value, 'remove')
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

const openLogs = () => {
  router.push({ name: 'ContainerLogs', params: { id: containerId.value } })
}

const openFile = () => {
  router.push({ name: 'ContainerFile', params: { id: containerId.value } })
}

const goToImage = (imageName: string) => {
  if (!imageName) return
  router.push({ name: 'ImageDetail', params: { id: imageName.substring(0, 12) } })
}

const showRawData = ref(false)
const rawData = computed(() => JSON.stringify(container.value, null, 2))

onMounted(async () => {
  containerLoading.value = true

  const cached = containerStore.getContainerInfoCached(containerId.value)
  if (cached) {
    container.value = cached
    containerLoading.value = false
    loadStats()
  }

  try {
    const data = await containerStore.refreshContainerInfo(containerId.value)
    container.value = data
    if (!cached) {
      loadStats()
    }
  } catch (error: any) {
    if (!cached) {
      message.error('获取容器详情失败: ' + error.message)
      router.push({ name: 'Containers' })
      return
    }
  } finally {
    containerLoading.value = false
  }

  if (!container.value?.Id) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
  }
})
</script>

<style scoped>
.container-detail {
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
  align-items: center;
}


.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-group .subtitle {
  color: var(--n-text-color-3);
  font-weight: 400;
}

.subtitle-text {
  margin-top: 8px;
  color: var(--n-text-color-3);
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
</style>