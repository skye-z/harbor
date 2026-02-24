<template>
  <div class="dashboard">
    <n-spin :show="loading">
      <n-grid x-gap="10" y-gap="10" item-responsive :cols="24">
        <n-gi span="12 800:6">
          <n-el tag="div" class="card stat-card">
            <div class="stat-content">
              <div class="stat-label">容器</div>
              <div class="stat-value">
                <n-text type="success">{{ containerStats.running }}</n-text>
                <n-text depth="3"> · </n-text>
                <n-text type="error">{{ containerStats.error }}</n-text>
                <n-text depth="3"> · </n-text>
                <n-text depth="3">{{ containerStats.stopped }}</n-text>
              </div>
            </div>
          </n-el>
        </n-gi>
        <n-gi span="12 800:6">
          <n-el tag="div" class="card stat-card" hoverable>
            <div class="stat-content">
              <div class="stat-label">镜像</div>
              <div class="stat-value">
                {{ imageStore.totalImages }}
                <n-text depth="3" style="font-size: 14px"> {{ formatSize(imageTotalSize) }}</n-text>
              </div>
            </div>
          </n-el>
        </n-gi>
        <n-gi span="12 800:6">
          <n-el tag="div" class="card stat-card" hoverable>
            <div class="stat-content">
              <div class="stat-label">存储</div>
              <div class="stat-value">
                <n-text :type="storageType">{{ formatSize(storageAvailable) }}</n-text>
              </div>
            </div>
          </n-el>
        </n-gi>
        <n-gi span="12 800:6">
          <n-el tag="div" class="card stat-card" hoverable>
            <div class="stat-content">
              <div class="stat-label">版本</div>
              <div class="stat-value">
                {{ dockerVersion }}
              </div>
            </div>
          </n-el>
        </n-gi>
        <n-gi span="24 570:12 800:8 1100:6">
          <n-el tag="div" class="card logs-card">
            <n-scrollbar style="height: calc(100vh - 116px)">
              <n-timeline>
                <n-timeline-item v-for="log in recentLogs" :key="log.id" :type="(log.type as any)" :title="log.action"
                  :time="log.time" />
              </n-timeline>
            </n-scrollbar>
            <div class="logs-more">
              <n-button strong secondary block type="primary" @click="goToLogs">
                查看更多
              </n-button>
            </div>
          </n-el>
        </n-gi>
        <n-gi span="24 570:12 800:16 1100:18">
          <n-el tag="div" class="card topology-card">
            <ResourceTopology />
          </n-el>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'
import { useUserStore } from '../plugins/stores/user'
import { systemApi, logApi } from '../plugins/api'
import ResourceTopology from '../components/ResourceTopology.vue'

const router = useRouter()
const containerStore = useContainerStore()
const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()
const userStore = useUserStore()

const loading = ref(false)
const dockerVersion = ref('Unknown')
const storageTotal = ref(0)
const storageAvailable = ref(0)

const recentLogs = ref<any[]>([])

const loadRecentLogs = async () => {
  try {
    const logs = await logApi.getRecent(10, '')
    recentLogs.value = logs.map((log: any) => ({
      id: log.id,
      action: log.message || log.action,
      time: new Date(log.created_at).toLocaleString('zh-CN', { hour12: false }),
      type: log.level === 'error' ? 'error' : log.level === 'warning' ? 'warning' : 'info'
    }))
  } catch (error: any) {
    console.error('Failed to load logs:', error)
    recentLogs.value = []
  }
}

const goToLogs = () => {
  router.push({ name: 'PlatformLogs' })
}

const containerStats = computed(() => {
  const running = containerStore.containers.filter(c => c.state === 'running').length
  const error = containerStore.containers.filter(c => c.state === 'dead' || c.state === 'restarting').length
  const stopped = containerStore.containers.filter(c => c.state === 'exited' || c.state === 'created').length
  return { running, error, stopped }
})

const imageTotalSize = computed(() => {
  return imageStore.images.reduce((total, img) => total + (img.size || 0), 0)
})

const storagePercentage = computed(() => {
  if (storageTotal.value === 0) return 0
  return (storageAvailable.value / storageTotal.value) * 100
})

const storageType = computed(() => {
  const percentage = storagePercentage.value
  if (percentage < 10) return 'error'
  if (percentage < 30) return 'warning'
  return 'success'
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      containerStore.fetchContainers(),
      imageStore.fetchImages(),
      volumeStore.fetchVolumes(),
      networkStore.fetchNetworks(),
      loadSystemInfo(),
      loadRecentLogs()
    ])
  } catch (err: any) {
    console.error('Failed to fetch data:', err)
  } finally {
    loading.value = false
  }
})

const loadSystemInfo = async () => {
  try {
    const info = await systemApi.getSystemInfo()
    dockerVersion.value = info.server_version || info.version || 'Unknown'

    storageTotal.value = info.disk_total || 0
    storageAvailable.value = info.disk_available || 0
  } catch (err) {
    console.error('Failed to load system info:', err)
    dockerVersion.value = 'Unknown'
    storageTotal.value = 0
    storageAvailable.value = 0
  }
}

</script>

<style scoped>
.dashboard {
  padding: 5px 10px 10px 10px;
}

.card {
  background-color: var(--base-color);
  border-radius: 12px;
  padding: 15px;
}

.stat-card {
  transition: all 0.3s ease;
  cursor: default;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-label {
  font-size: 14px;
  color: var(--n-text-color-3);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  line-height: 1.2;
}

.logs-card {
  max-height: 530px !important;
  height: calc(100vh - 140px);
  flex-direction: column;
  position: relative;
  display: flex;
  padding: 0;
}

.logs-card :deep(.n-timeline) {
  padding: 12px;
  padding-bottom: 60px;
  overflow-y: auto;
  flex: 1;
}

.logs-more {
  padding: 10px;
}

.logs-card,
.topology-card {
  max-height: 500px;
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--n-border-color);
  font-weight: 600;
  font-size: 14px;
}

.topology-card :deep(.n-card__content) {
  padding: 0;
}
</style>
