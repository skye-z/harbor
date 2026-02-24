<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { systemApi } from '../plugins/api'
import { useContainerStore } from '../plugins/stores/containers'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'
import {
  ServerOutline,
  TrashOutline,
  RefreshOutline,
  SettingsOutline
} from '@vicons/ionicons5'

const message = useMessage()
const dialog = useDialog()
const containerStore = useContainerStore()
const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()

const loading = ref(false)
const pruneLoading = ref(false)
const pruneResult = ref<Record<string, number>>({})

const systemInfo = ref<any>(null)
const hostInfo = ref<any>(null)

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const loadSystemInfo = async () => {
  try {
    systemInfo.value = await systemApi.getSystemInfo()
    hostInfo.value = {
      hostname: systemInfo.value.name || 'Unknown',
      os: systemInfo.value.operating_system || 'Linux',
      architecture: systemInfo.value.architecture || 'Unknown',
      kernel: systemInfo.value.kernel_version || 'Unknown',
      cpus: systemInfo.value.n_cpu || 0,
      memory: systemInfo.value.mem_total || 0,
      osType: systemInfo.value.os_type || 'Unknown'
    }
  } catch (error: any) {
    message.error('加载系统信息失败: ' + error.message)
  }
}

const handlePrune = (type: 'containers' | 'images' | 'volumes' | 'networks' | 'all') => {
  const typeName = {
    containers: '容器',
    images: '镜像',
    volumes: '卷',
    networks: '网络',
    all: '全部资源'
  }[type]

  const pruneMethod = {
    containers: systemApi.pruneContainers,
    images: systemApi.pruneImages,
    volumes: systemApi.pruneVolumes,
    networks: systemApi.pruneNetworks,
    all: systemApi.pruneAll
  }[type]

  dialog.warning({
    title: '确认清理',
    content: `确定要清理未使用的${typeName}吗？此操作不可逆。`,
    positiveText: '清理',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        pruneLoading.value = true
        pruneResult.value = {}
        const result = await pruneMethod()
        pruneResult.value[type] = result.space_reclaimed
        if (result.details) {
          pruneResult.value = { ...pruneResult.value, ...result.details }
        }
        const spaceMB = result.space_reclaimed / 1024 / 1024
        const spaceGB = spaceMB / 1024
        const spaceText = spaceGB >= 1 ? `${spaceGB.toFixed(2)} GB` : `${spaceMB.toFixed(2)} MB`
        message.success(`${result.message}，释放空间: ${spaceText}`)

        await Promise.all([
          containerStore.fetchContainers(),
          imageStore.fetchImages(),
          volumeStore.fetchVolumes(),
          networkStore.fetchNetworks()
        ])
      } catch (error: any) {
        message.error('清理失败: ' + error.message)
      } finally {
        pruneLoading.value = false
      }
    }
  })
}

onMounted(() => {
  loadSystemInfo()
})
</script>

<template>
  <div class="system-page">
    <div class="view-header">
      <h1 class="page-title">系统信息</h1>
    </div>

    <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen" style="margin-top: 10px;">
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">主机名</div>
          <div class="info-value">{{ hostInfo?.hostname || 'Loading...' }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">操作系统</div>
          <div class="info-value">{{ hostInfo?.os || 'Loading...' }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">内核版本</div>
          <div class="info-value">{{ hostInfo?.kernel || 'Loading...' }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">架构</div>
          <div class="info-value">{{ hostInfo?.architecture || 'Loading...' }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">CPU核心数</div>
          <div class="info-value large">{{ hostInfo?.cpus || 0 }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">总内存</div>
          <div class="info-value">{{ formatSize(hostInfo?.memory || 0) }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">Docker版本</div>
          <div class="info-value">{{ systemInfo?.server_version || 'Loading...' }}</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="info-card" :bordered="false">
          <div class="info-label">API版本</div>
          <div class="info-value">1.43</div>
        </n-card>
      </n-gi>
    </n-grid>

    <div style="margin-top: 30px;">
      <h2 style="margin: 0 0 10px 0; font-size: 18px; font-weight: 600;">资源清理</h2>
      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:4" responsive="screen">
        <n-gi>
          <n-card class="prune-card" :bordered="false" hoverable>
            <div class="prune-content">
              <n-icon :component="SettingsOutline" size="28" color="#2080f0" style="margin-bottom: 8px;" />
              <div class="prune-label">清理容器</div>
              <div class="prune-desc">清理所有已停止的容器</div>
              <n-space justify="space-between" align="center" style="margin-top: 12px;">
                <n-text v-if="pruneResult.containers" type="success" style="font-size: 12px;">
                  -{{ formatSize(pruneResult.containers) }}
                </n-text>
                <n-button size="small" type="warning" @click="handlePrune('containers')" :loading="pruneLoading">
                  <template #icon>
                    <n-icon :component="TrashOutline" size="14" />
                  </template>
                  执行
                </n-button>
              </n-space>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="prune-card" :bordered="false" hoverable>
            <div class="prune-content">
              <n-icon :component="SettingsOutline" size="28" color="#f0a020" style="margin-bottom: 8px;" />
              <div class="prune-label">清理镜像</div>
              <div class="prune-desc">清理未使用的镜像</div>
              <n-space justify="space-between" align="center" style="margin-top: 12px;">
                <n-text v-if="pruneResult.images" type="success" style="font-size: 12px;">
                  -{{ formatSize(pruneResult.images) }}
                </n-text>
                <n-button size="small" type="warning" @click="handlePrune('images')" :loading="pruneLoading">
                  <template #icon>
                    <n-icon :component="TrashOutline" size="14" />
                  </template>
                  执行
                </n-button>
              </n-space>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="prune-card" :bordered="false" hoverable>
            <div class="prune-content">
              <n-icon :component="SettingsOutline" size="28" color="#18a058" style="margin-bottom: 8px;" />
              <div class="prune-label">清理卷</div>
              <div class="prune-desc">清理未使用的卷</div>
              <n-space justify="space-between" align="center" style="margin-top: 12px;">
                <n-text v-if="pruneResult.volumes" type="success" style="font-size: 12px;">
                  -{{ formatSize(pruneResult.volumes) }}
                </n-text>
                <n-button size="small" type="warning" @click="handlePrune('volumes')" :loading="pruneLoading">
                  <template #icon>
                    <n-icon :component="TrashOutline" size="14" />
                  </template>
                  执行
                </n-button>
              </n-space>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="prune-card" :bordered="false" hoverable>
            <div class="prune-content">
              <n-icon :component="SettingsOutline" size="28" color="#d03050" style="margin-bottom: 8px;" />
              <div class="prune-label">清理网络</div>
              <div class="prune-desc">清理未使用的网络</div>
              <n-space justify="space-between" align="center" style="margin-top: 12px;">
                <n-text v-if="pruneResult.networks" type="success" style="font-size: 12px;">
                  -{{ formatSize(pruneResult.networks) }}
                </n-text>
                <n-button size="small" type="warning" @click="handlePrune('networks')" :loading="pruneLoading">
                  <template #icon>
                    <n-icon :component="TrashOutline" size="14" />
                  </template>
                  执行
                </n-button>
              </n-space>
            </div>
          </n-card>
        </n-gi>
        <n-gi :span="24">
          <n-card class="prune-card danger" :bordered="false" hoverable>
            <div class="prune-content">
              <n-icon :component="TrashOutline" size="28" color="#d03050" style="margin-bottom: 8px;" />
              <div class="prune-label danger">一键清理全部</div>
              <div class="prune-desc">清理所有未使用的资源</div>
              <n-space justify="space-between" align="center" style="margin-top: 12px;">
                <n-text v-if="pruneResult.all" type="success" style="font-size: 12px; font-weight: 600;">
                  总释放: {{ formatSize(pruneResult.all) }}
                </n-text>
                <n-button size="small" type="error" @click="handlePrune('all')" :loading="pruneLoading">
                  <template #icon>
                    <n-icon :component="TrashOutline" size="14" />
                  </template>
                  清理全部
                </n-button>
              </n-space>
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<style scoped>
.system-page {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.view-header {
  margin-bottom: 0;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
}

.info-card {
  transition: all 0.3s ease;
  cursor: default;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.info-label {
  font-size: 14px;
  color: var(--n-text-color-3);
  margin-bottom: 8px;
}

.info-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--n-text-color-1);
}

.info-value.large {
  font-size: 32px;
}

.prune-card {
  transition: all 0.3s ease;
  cursor: default;
}

.prune-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.prune-card.danger {
  background: rgba(208, 48, 80, 0.05);
  border: 1px solid rgba(208, 48, 80, 0.2);
}

.prune-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.prune-label {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.prune-label.danger {
  color: #d03050;
}

.prune-desc {
  font-size: 12px;
  color: var(--n-text-color-3);
}
</style>
