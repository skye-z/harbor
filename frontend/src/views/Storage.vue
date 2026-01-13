<template>
  <div class="storage">
    <n-grid x-gap="10" y-gap="10" :cols="24">
      <n-gi :span="24">
        <n-el tag="div" class="card overview-card">
          <n-row :gutter="16" style="padding: 16px 20px 0">
            <n-col :span="6">
              <n-statistic label="总容量">
                {{ formatSize(storage.total) }}
              </n-statistic>
            </n-col>
            <n-col :span="6">
              <n-statistic label="已使用">
                {{ formatSize(storage.used) }}
              </n-statistic>
            </n-col>
            <n-col :span="6">
              <n-statistic label="可用空间">
                {{ formatSize(storage.available) }}
              </n-statistic>
            </n-col>
            <n-col :span="6">
              <n-statistic label="使用率" :value="storage.usagePercent" suffix="%" />
            </n-col>
          </n-row>
          <div style="padding: 0 20px 16px">
            <n-progress
              type="line"
              :percentage="storage.usagePercent"
              :status="storage.usagePercent > 80 ? 'error' : 'success'"
              :height="12"
              :show-indicator="true"
            />
          </div>
        </n-el>
      </n-gi>

      <n-gi :span="12">
        <n-el tag="div" class="card resource-card">
          <div class="card-header">
            <n-space align="center">
              <n-text strong>数据卷</n-text>
            </n-space>
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
          </div>
          <div class="card-content">
            <n-scrollbar style="max-height: 500px">
              <n-space vertical :size="12">
                <n-el tag="div" v-for="volume in volumes" :key="volume.id" class="volume-card" hoverable>
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
                      <n-text depth="3" style="font-size: 12px">容量</n-text>
                      <n-text>{{ volume.size ? formatSize(volume.size) : '-' }}</n-text>
                    </div>
                    <div class="info-item">
                      <n-text depth="3" style="font-size: 12px">已使用</n-text>
                      <n-text>{{ volume.used ? formatSize(volume.used) : '-' }}</n-text>
                    </div>
                    <div class="info-item">
                      <n-text depth="3" style="font-size: 12px">创建时间</n-text>
                      <n-text>{{ volume.created_at }}</n-text>
                    </div>
                  </div>
                </n-el>
                <n-empty v-if="volumes.length === 0" description="暂无数据卷" size="medium" />
              </n-space>
            </n-scrollbar>
          </div>
        </n-el>
      </n-gi>

      <n-gi :span="12">
        <n-el tag="div" class="card resource-card">
          <div class="card-header">
            <n-space align="center">
              <n-text strong>网络</n-text>
            </n-space>
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
          </div>
          <div class="card-content">
            <n-scrollbar style="max-height: 500px">
              <n-space vertical :size="12">
                <n-el tag="div" v-for="network in networks" :key="network.id" class="network-card" hoverable>
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
                      <n-tag
                        size="small"
                        :bordered="false"
                        :type="network.driver === 'bridge' ? 'info' : 'default'"
                      >
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
                </n-el>
                <n-empty v-if="networks.length === 0" description="暂无网络" size="medium" />
              </n-space>
            </n-scrollbar>
          </div>
        </n-el>
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
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { volumeApi, networkApi } from '../plugins/api'
import {
  ServerOutline,
  CloudOutline,
  CheckmarkCircleOutline,
  PieChartOutline,
  AddOutline,
  TrashOutline,
  SettingsOutline,
  RefreshOutline,
  GlobeOutline
} from '@vicons/ionicons5'
import type { Volume, Network } from '../types'

const message = useMessage()

const showVolumeModal = ref(false)
const showNetworkModal = ref(false)
const loadingVolumes = ref(false)
const loadingNetworks = ref(false)

const storage = ref({
  total: 500 * 1024 * 1024 * 1024,
  used: 375 * 1024 * 1024 * 1024,
  available: 125 * 1024 * 1024 * 1024,
  usagePercent: 75
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

const handleInspectVolume = (_id: string) => {
  message.info('查看详情功能开发中')
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

const handleInspectNetwork = (_id: string) => {
  message.info('查看详情功能开发中')
}

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

onMounted(() => {
  refreshVolumes()
  refreshNetworks()
})
</script>

<style scoped>
.storage {
  padding: 5px 10px 10px 10px;
}

.card {
  border: 1px solid var(--n-border-color);
  background-color: var(--n-card-color);
  border-radius: 12px;
  overflow: hidden;
}

.overview-card,
.resource-card {
  border-radius: 12px;
}

.volume-card,
.network-card {
  transition: all 0.3s ease;
  border-radius: 8px;
  border: 1px solid var(--n-border-color);
  background-color: var(--n-card-color);
  padding: 12px 16px;
}

.volume-card:hover,
.network-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--n-border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-content {
  padding: 16px 20px;
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
  color: var(--n-text-color-1);
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
