<template>
  <div class="storage">
    <div class="page-header">
      <div class="view-header">
        <div class="title-group">
          <h1>存储管理</h1>
        </div>
      </div>
    </div>

    <n-grid x-gap="20" :cols="24">
      <n-gi :span="24">
        <n-card class="config-card">
          <n-row :gutter="16">
            <n-col :span="6">
              <n-statistic label="总容量" :value="formatSize(storage.total)" />
            </n-col>
            <n-col :span="6">
              <n-statistic label="已使用" :value="formatSize(storage.used)" />
            </n-col>
            <n-col :span="6">
              <n-statistic label="可用空间" :value="formatSize(storage.available)" />
            </n-col>
            <n-col :span="6">
              <n-statistic label="使用率" :value="storage.usagePercent" suffix="%" />
            </n-col>
          </n-row>
          <div style="margin-top: 16px">
            <n-progress
              type="line"
              :percentage="storage.usagePercent"
              :status="storage.usagePercent > 80 ? 'error' : storage.usagePercent > 60 ? 'warning' : 'success'"
              :height="12"
              :show-indicator="true"
            />
          </div>
        </n-card>
      </n-gi>

      <n-gi :span="12">
        <n-card title="数据卷" class="config-card">
          <template #header-extra>
            <n-space>
              <n-button type="primary" size="small" @click="showVolumeModal = true">
                <template #icon>
                  <n-icon :component="AddOutline" />
                </template>
                新建
              </n-button>
              <n-button quaternary size="small" @click="refreshVolumes" :loading="loadingVolumes">
                <template #icon>
                  <n-icon :component="RefreshOutline" />
                </template>
              </n-button>
            </n-space>
          </template>
          <n-scrollbar style="max-height: 500px">
            <n-space vertical :size="12">
              <n-card v-for="volume in volumes" :key="volume.id" size="small" hoverable>
                <div class="volume-header">
                  <div class="volume-title">
                    <n-icon :component="ServerOutline" :size="18" color="#2080f0" />
                    <span>{{ volume.name }}</span>
                  </div>
                  <n-space>
                    <n-button text size="small" type="primary" @click="handleInspectVolume(volume.id)">
                      <template #icon>
                        <n-icon :component="SettingsOutline" :size="16" />
                      </template>
                    </n-button>
                    <n-button text size="small" type="error" @click="handleDeleteVolume(volume.id)">
                      <template #icon>
                        <n-icon :component="TrashOutline" :size="16" />
                      </template>
                    </n-button>
                  </n-space>
                </div>
                <n-divider style="margin: 8px 0" />
                <div class="volume-info">
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">驱动</n-text>
                    <n-tag size="small" :bordered="false" type="info">{{ volume.driver }}</n-tag>
                  </div>
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">挂载点</n-text>
                    <n-text code style="font-size: 12px">{{ volume.mountpoint }}</n-text>
                  </div>
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">创建时间</n-text>
                    <n-text>{{ volume.created_at }}</n-text>
                  </div>
                </div>
              </n-card>
              <n-empty v-if="volumes.length === 0" description="暂无数据卷" size="medium" />
            </n-space>
          </n-scrollbar>
        </n-card>
      </n-gi>

      <n-gi :span="12">
        <n-card title="网络" class="config-card">
          <template #header-extra>
            <n-space>
              <n-button type="primary" size="small" @click="showNetworkModal = true">
                <template #icon>
                  <n-icon :component="AddOutline" />
                </template>
                新建
              </n-button>
              <n-button quaternary size="small" @click="refreshNetworks" :loading="loadingNetworks">
                <template #icon>
                  <n-icon :component="RefreshOutline" />
                </template>
              </n-button>
            </n-space>
          </template>
          <n-scrollbar style="max-height: 500px">
            <n-space vertical :size="12">
              <n-card v-for="network in networks" :key="network.id" size="small" hoverable>
                <div class="network-header">
                  <div class="network-title">
                    <n-icon :component="GlobeOutline" :size="18" color="#18a058" />
                    <span>{{ network.name }}</span>
                  </div>
                  <n-space>
                    <n-button text size="small" type="primary" @click="handleInspectNetwork(network.id)">
                      <template #icon>
                        <n-icon :component="SettingsOutline" :size="16" />
                      </template>
                    </n-button>
                    <n-button text size="small" type="error" @click="handleDeleteNetwork(network.id)">
                      <template #icon>
                        <n-icon :component="TrashOutline" :size="16" />
                      </template>
                    </n-button>
                  </n-space>
                </div>
                <n-divider style="margin: 8px 0" />
                <div class="network-info">
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">驱动</n-text>
                    <n-tag size="small" :bordered="false" :type="network.driver === 'bridge' ? 'info' : 'default'">
                      {{ network.driver }}
                    </n-tag>
                  </div>
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">子网</n-text>
                    <n-text code style="font-size: 12px">{{ network.subnet || '-' }}</n-text>
                  </div>
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">网关</n-text>
                    <n-text code style="font-size: 12px">{{ network.gateway || '-' }}</n-text>
                  </div>
                  <div class="info-item">
                    <n-text depth="3" style="font-size: 12px">创建时间</n-text>
                    <n-text>{{ network.created_at }}</n-text>
                  </div>
                </div>
              </n-card>
              <n-empty v-if="networks.length === 0" description="暂无网络" size="medium" />
            </n-space>
          </n-scrollbar>
        </n-card>
      </n-gi>
    </n-grid>

    <n-modal v-model:show="showVolumeModal" preset="card" title="新建数据卷" style="width: 500px">
      <n-form :model="volumeForm" label-placement="left" label-width="100px">
        <n-form-item label="卷名称">
          <n-input v-model:value="volumeForm.name" placeholder="my-volume" />
        </n-form-item>
        <n-form-item label="驱动">
          <n-select v-model:value="volumeForm.driver" :options="driverOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showVolumeModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateVolume" :loading="loadingVolumes">创建</n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal v-model:show="showNetworkModal" preset="card" title="新建网络" style="width: 500px">
      <n-form :model="networkForm" label-placement="left" label-width="100px">
        <n-form-item label="网络名称">
          <n-input v-model:value="networkForm.name" placeholder="my-network" />
        </n-form-item>
        <n-form-item label="驱动">
          <n-select v-model:value="networkForm.driver" :options="networkDriverOptions" />
        </n-form-item>
        <n-form-item label="子网">
          <n-input v-model:value="networkForm.subnet" placeholder="172.20.0.0/16" />
        </n-form-item>
        <n-form-item label="网关">
          <n-input v-model:value="networkForm.gateway" placeholder="172.20.0.1" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showNetworkModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateNetwork" :loading="loadingNetworks">创建</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { volumeApi, networkApi, systemApi } from '../plugins/api'
import {
  ServerOutline,
  AddOutline,
  TrashOutline,
  SettingsOutline,
  RefreshOutline,
  GlobeOutline
} from '@vicons/ionicons5'
import type { Volume, Network } from '../types'

const message = useMessage()
const dialog = useDialog()

const showVolumeModal = ref(false)
const showNetworkModal = ref(false)
const loadingVolumes = ref(false)
const loadingNetworks = ref(false)
const loadingStorage = ref(false)

const storageData = ref({
  disk_total: 0,
  disk_used: 0,
  disk_available: 0
})

const storage = computed(() => {
  const total = storageData.value.disk_total || 0
  const used = storageData.value.disk_used || 0
  const available = storageData.value.disk_available || 0
  const usagePercent = total > 0 ? Math.round((used / total) * 100) : 0
  return {
    total,
    used,
    available,
    usagePercent
  }
})

const volumeForm = ref({
  name: '',
  driver: 'local'
})

const networkForm = ref({
  name: '',
  driver: 'bridge',
  subnet: '',
  gateway: ''
})

const driverOptions = [
  { label: 'Local', value: 'local' },
  { label: 'NFS', value: 'nfs' },
  { label: 'SMB', value: 'smb' }
]

const networkDriverOptions = [
  { label: 'Bridge', value: 'bridge' },
  { label: 'Overlay', value: 'overlay' },
  { label: 'Macvlan', value: 'macvlan' }
]

const volumes = ref<Volume[]>([])
const networks = ref<Network[]>([])

const refreshStorage = async () => {
  loadingStorage.value = true
  try {
    const info = await systemApi.getSystemInfo()
    storageData.value = {
      disk_total: info.disk_total || 0,
      disk_used: info.disk_used || 0,
      disk_available: info.disk_available || 0
    }
  } catch (error: any) {
    console.error('获取磁盘信息失败:', error)
  } finally {
    loadingStorage.value = false
  }
}

const refreshVolumes = async () => {
  loadingVolumes.value = true
  try {
    volumes.value = await volumeApi.list()
  } catch (error: any) {
    message.error('刷新失败: ' + error.message)
  } finally {
    loadingVolumes.value = false
  }
}

const refreshNetworks = async () => {
  loadingNetworks.value = true
  try {
    networks.value = await networkApi.list()
  } catch (error: any) {
    message.error('刷新失败: ' + error.message)
  } finally {
    loadingNetworks.value = false
  }
}

const handleCreateVolume = async () => {
  if (!volumeForm.value.name) {
    message.error('请输入卷名称')
    return
  }

  loadingVolumes.value = true
  try {
    await volumeApi.create(volumeForm.value)
    showVolumeModal.value = false
    message.success('数据卷创建成功')
    volumeForm.value = { name: '', driver: 'local' }
    await refreshVolumes()
  } catch (error: any) {
    message.error('创建失败: ' + error.message)
  } finally {
    loadingVolumes.value = false
  }
}

const handleDeleteVolume = async (id: string) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除该数据卷吗？此操作不可恢复。',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      loadingVolumes.value = true
      try {
        await volumeApi.delete(id)
        message.success('数据卷删除成功')
        await refreshVolumes()
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      } finally {
        loadingVolumes.value = false
      }
    }
  })
}

const handleInspectVolume = (id: string) => {
  const volume = volumes.value.find(v => v.id === id)
  if (volume) {
    dialog.info({
      title: '数据卷详情',
      content: JSON.stringify(volume, null, 2),
      positiveText: '确定'
    })
  }
}

const handleCreateNetwork = async () => {
  if (!networkForm.value.name) {
    message.error('请输入网络名称')
    return
  }

  loadingNetworks.value = true
  try {
    await networkApi.create(networkForm.value)
    showNetworkModal.value = false
    message.success('网络创建成功')
    networkForm.value = { name: '', driver: 'bridge', subnet: '', gateway: '' }
    await refreshNetworks()
  } catch (error: any) {
    message.error('创建失败: ' + error.message)
  } finally {
    loadingNetworks.value = false
  }
}

const handleDeleteNetwork = async (id: string) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除该网络吗？此操作不可恢复。',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      loadingNetworks.value = true
      try {
        await networkApi.delete(id)
        message.success('网络删除成功')
        await refreshNetworks()
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      } finally {
        loadingNetworks.value = false
      }
    }
  })
}

const handleInspectNetwork = (id: string) => {
  const network = networks.value.find(n => n.id === id)
  if (network) {
    dialog.info({
      title: '网络详情',
      content: JSON.stringify(network, null, 2),
      positiveText: '确定'
    })
  }
}

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

onMounted(() => {
  refreshStorage()
  refreshVolumes()
  refreshNetworks()
})
</script>

<style scoped>
.storage {
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

.volume-header,
.network-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.volume-title,
.network-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 14px;
}

.volume-info,
.network-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
