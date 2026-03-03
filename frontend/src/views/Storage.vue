<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useMessage, useDialog, NButton, NSpace, NIcon } from 'naive-ui'
import { volumeApi, networkApi, systemApi } from '../plugins/api'
import {
  ServerOutline,
  Add,
  Trash,
  SettingsOutline,
  RefreshOutline,
  GlobeOutline,
  Balloon
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

<template>
  <div class="storage-page">
    <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen" style="margin-top: 10px;">
      <n-gi v-for="(stat, key) in [
        { label: '总容量', value: formatSize(storage.total), color: '#2080f0' },
        { label: '已使用', value: formatSize(storage.used), color: '#f0a020' },
        { label: '可用空间', value: formatSize(storage.available), color: '#18a058' },
        { label: '使用率', value: storage.usagePercent + '%', color: storage.usagePercent > 80 ? '#d03050' : '#18a058' }
      ]" :key="key">
        <n-card class="stat-card" :bordered="false">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value" :style="{ color: stat.color }">{{ stat.value }}</div>
        </n-card>
      </n-gi>
    </n-grid>

    <div style="margin-top: 20px;">
      <div style="margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center;">
        <h2 style="margin: 0; font-size: 18px; font-weight: 600;">数据卷</h2>
        <n-space>
          <n-button size="small" @click="showVolumeModal = true">
            <template #icon>
              <n-icon :component="Add" />
            </template>
            新建
          </n-button>
          <n-button size="small" quaternary @click="refreshVolumes" :loading="loadingVolumes">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
          </n-button>
        </n-space>
      </div>

      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen" v-if="volumes.length > 0">
        <n-gi v-for="volume in volumes" :key="volume.id">
          <n-card hoverable class="volume-card" :bordered="false" style="cursor: pointer" @click="router.push({ name: 'VolumeDetail', params: { id: volume.id } })">
            <template #header>
              <div>
                <div class="card-header">
                  <div class="card-title">
                    {{ volume.name }}
                  </div>
                  <n-tag size="small" :bordered="false" type="info">{{ volume.driver }}</n-tag>
                </div>
                <div style="font-size: 12px;color: #999;">{{ volume.mountpoint }}</div>
              </div>
            </template>
            <template #action>
              <div style="display: flex;align-items: center;justify-content: space-between;">
                <n-time v-if="volume.created_at" :time="new Date(volume.created_at).getTime()" type="relative"
                  style="font-size: 12px;color: #999;" />
                <n-button size="small" quaternary circle type="error" @click.stop="handleDeleteVolume(volume.id)">
                  <template #icon>
                    <n-icon>
                      <Trash />
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-result v-else status="404" title="这里什么都没有" style="margin-top: 15px;">
      </n-result>
    </div>

    <div style="margin-top: 30px;">
      <div style="margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center;">
        <h2 style="margin: 0; font-size: 18px; font-weight: 600;">网络</h2>
        <n-space>
          <n-button size="small" @click="showNetworkModal = true">
            <template #icon>
              <n-icon :component="Add" />
            </template>
            新建
          </n-button>
          <n-button size="small" quaternary @click="refreshNetworks" :loading="loadingNetworks">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
          </n-button>
        </n-space>
      </div>

      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen" v-if="networks.length > 0">
        <n-gi v-for="network in networks" :key="network.id">
          <n-card hoverable class="network-card" :bordered="false">
            <template #header>
              <div class="card-header">
                <div>
                  <div class="card-title">
                    {{ network.name }}
                  </div>
                </div>
                <n-tag size="small" :bordered="false" :type="network.driver === 'bridge' ? 'info' : 'default'">
                  {{ network.driver }}
                </n-tag>
              </div>
            </template>
            <div class="network-info">
              <div class="info-item">
                <span class="label">子网</span>
                <span class="value">{{ network.subnet || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">网关</span>
                <span class="value">{{ network.gateway || '-' }}</span>
              </div>
            </div>
            <template #action>
              <n-button class="del-btn" size="small" quaternary circle type="error"
                @click="handleDeleteNetwork(network.id)">
                <template #icon>
                  <n-icon>
                    <Trash />
                  </n-icon>
                </template>
              </n-button>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-result v-else status="404" title="这里什么都没有" style="margin-top: 15px;">
      </n-result>
    </div>

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

<style scoped>
.storage-page {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.view-header {
  margin-bottom: 0;
}

.filter-bar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
  display: flex;
  align-items: center;
}

.volume-card,
.network-card {
  transition: all 0.3s ease;
}

.volume-card:hover,
.network-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.volume-info,
.network-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.info-item .label {
  color: var(--n-text-color-3);
}

.info-item .value {
  font-family: monospace;
}

.del-btn {
  float: right;
}
</style>
