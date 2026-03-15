<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog, NTag, NButton, NSpace, NIcon } from 'naive-ui'
import { networkApi } from '../plugins/api'
import { useNetworkStore } from '../plugins/stores/networks'
import { useContainerStore } from '../plugins/stores/containers'
import {
  ArrowBackOutline,
  TrashOutline,
  LinkOutline,
  UnlinkOutline,
  RefreshOutline,
  GlobeOutline,
  CubeOutline,
  GitNetworkOutline,
  AddOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const networkStore = useNetworkStore()
const containerStore = useContainerStore()

const networkName = route.params.name as string
const loading = ref(false)
const showConnectModal = ref(false)
const selectedContainer = ref<string>('')
const connectLoading = ref(false)

const network = computed(() => {
  return networkStore.networks.find(n => n.name === networkName)
})

// 获取网络ID（用于API调用）
const networkId = computed(() => network.value?.id || network.value?.name)

// 获取接入该网络的容器
const connectedContainers = computed(() => {
  return containerStore.containers.filter(container => {
    if (container.network_settings?.networks) {
      return Object.keys(container.network_settings.networks).includes(network.value?.name || '')
    }
    return false
  })
})

// 获取未接入该网络的容器
const availableContainers = computed(() => {
  return containerStore.containers.filter(container => {
    if (container.network_settings?.networks) {
      return !Object.keys(container.network_settings.networks).includes(network.value?.name || '')
    }
    return true
  })
})

const loadNetworkDetail = async () => {
  loading.value = true
  try {
    await networkStore.fetchNetworks()
    await containerStore.fetchContainers()
    if (!network.value) {
      message.error('网络不存在')
      router.push({ name: 'Connect' })
      return
    }
  } catch (error: any) {
    message.error('加载网络信息失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 删除网络
const handleDelete = () => {
  if (connectedContainers.value.length > 0) {
    message.error('该网络有容器接入，请先断开所有容器')
    return
  }

  dialog.warning({
    title: '确认删除',
    content: `确定要删除网络 "${network.value?.name}" 吗？此操作不可逆。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await networkApi.delete(networkId.value)
        message.success('网络已删除')
        router.push({ name: 'Connect' })
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      }
    }
  })
}

// 断开容器连接
const handleDisconnect = async (containerId: string) => {
  try {
    await networkApi.disconnectContainer(networkId.value, containerId)
    message.success('已断开连接')
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('断开连接失败: ' + error.message)
  }
}

// 连接容器
const handleConnect = async () => {
  if (!selectedContainer.value) {
    message.error('请选择要连接的容器')
    return
  }

  try {
    connectLoading.value = true
    await networkApi.connectContainer(networkId.value, selectedContainer.value)
    message.success('容器已接入网络')
    showConnectModal.value = false
    selectedContainer.value = ''
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('连接失败: ' + error.message)
  } finally {
    connectLoading.value = false
  }
}

// 获取容器在网络中的IP
const getContainerIP = (container: any) => {
  const networkConfig = container.network_settings?.networks?.[network.value?.name || '']
  return networkConfig?.ip_address || 'N/A'
}

// 获取容器MAC地址
const getContainerMAC = (container: any) => {
  const networkConfig = container.network_settings?.networks?.[network.value?.name || '']
  return networkConfig?.mac_address || 'N/A'
}

onMounted(() => {
  loadNetworkDetail()
})
</script>

<template>
  <div class="network-detail">
    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ network?.name || networkName }}</h1>
          <div class="header-actions">
            <n-button type="primary" ghost @click="showConnectModal = true">
              <template #icon>
                <n-icon :component="LinkOutline" />
              </template>
              接入容器
            </n-button>
            <n-button type="error" ghost @click="handleDelete">
              <template #icon>
                <n-icon :component="TrashOutline" />
              </template>
              删除
            </n-button>
          </div>
        </div>
        <div class="subtitle-text">
          <n-tag :type="network?.driver === 'bridge' ? 'info' : 'default'" size="small">
            {{ network?.driver || 'bridge' }}
          </n-tag>
          <span style="margin-left: 8px;">{{ network?.id || '-' }}</span>
        </div>
      </div>
    </div>

    <n-spin :show="loading">
      <n-grid item-responsive x-gap="10" cols="24">
        <!-- 左侧：基本信息 -->
        <n-gi span="24 760:9 900:8">
          <n-card size="small" title="基本信息" style="margin-bottom: 10px;">
            <template #header-extra>
              <n-tag :type="network?.scope === 'local' ? 'success' : 'warning'" size="small" :bordered="false">
                {{ network?.scope || 'local' }}
              </n-tag>
            </template>
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item label="ID">
                <n-ellipsis style="max-width: 200px;">
                  {{ network?.id || '-' }}
                </n-ellipsis>
              </n-descriptions-item>
              <n-descriptions-item label="名称">
                {{ network?.name || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="驱动">
                {{ network?.driver || 'bridge' }}
              </n-descriptions-item>
              <n-descriptions-item label="范围">
                {{ network?.scope || 'local' }}
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">
                <n-time v-if="network?.created" :time="new Date(network.created).getTime()" type="datetime" />
                <span v-else>-</span>
              </n-descriptions-item>
            </n-descriptions>
          </n-card>

          <n-card size="small" title="网络配置" style="margin-bottom: 10px;">
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item label="子网">
                {{ network?.subnet || network?.ipam?.config?.[0]?.subnet || 'N/A' }}
              </n-descriptions-item>
              <n-descriptions-item label="网关">
                {{ network?.gateway || network?.ipam?.config?.[0]?.gateway || 'N/A' }}
              </n-descriptions-item>
              <n-descriptions-item label="IP范围">
                {{ network?.ipam?.config?.[0]?.ip_range || 'N/A' }}
              </n-descriptions-item>
              <n-descriptions-item label="内部网络">
                {{ network?.internal ? '是' : '否' }}
              </n-descriptions-item>
              <n-descriptions-item label="可接入">
                {{ network?.attachable ? '是' : '否' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>

        <!-- 右侧：接入容器 -->
        <n-gi span="24 760:15 900:16">
          <n-card size="small" title="接入容器" style="margin-bottom: 10px;">
            <template #header-extra>
              <n-tag type="info" size="small">{{ connectedContainers.length }} 个容器</n-tag>
            </template>
            <n-list v-if="connectedContainers.length > 0">
              <n-list-item v-for="container in connectedContainers" :key="container.id">
                <n-thing :title="container.names[0]?.replace(/^\//, '') || container.id">
                  <template #avatar>
                    <n-avatar>
                      <n-icon :component="CubeOutline" />
                    </n-avatar>
                  </template>
                  <template #description>
                    <n-space vertical size="small">
                      <n-space>
                        <n-tag size="small" :type="container.state === 'running' ? 'success' : 'default'">
                          {{ container.state }}
                        </n-tag>
                        <span style="font-size: 12px; color: #999;">{{ container.image }}</span>
                      </n-space>
                      <n-space>
                        <n-tag size="small" type="info">IP: {{ getContainerIP(container) }}</n-tag>
                        <n-tag size="small" type="warning">MAC: {{ getContainerMAC(container) }}</n-tag>
                      </n-space>
                    </n-space>
                  </template>
                </n-thing>
                <template #suffix>
                  <n-space>
                    <n-button size="small" quaternary type="primary" @click="router.push({ name: 'ContainerDetail', params: { id: container.id } })">
                      查看
                    </n-button>
                    <n-button size="small" quaternary type="error" @click="handleDisconnect(container.id)">
                      断开
                    </n-button>
                  </n-space>
                </template>
              </n-list-item>
            </n-list>
            <n-empty v-else description="暂无容器接入此网络" />
          </n-card>

          <n-card size="small" title="标签">
            <n-space v-if="network?.labels && Object.keys(network.labels).length > 0" wrap>
              <n-tag v-for="(value, key) in network.labels" :key="key" size="small">
                {{ key }}: {{ value }}
              </n-tag>
            </n-space>
            <n-empty v-else description="暂无标签" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>

    <!-- 接入容器弹窗 -->
    <n-modal v-model:show="showConnectModal" preset="card" title="接入容器到网络" style="width: 500px;">
      <n-form label-placement="left" label-width="100px">
        <n-form-item label="选择容器">
          <n-select v-model:value="selectedContainer" placeholder="请选择要接入的容器"
            :options="availableContainers.map(c => ({ label: c.names[0]?.replace(/^\//, '') || c.id, value: c.id }))" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showConnectModal = false">取消</n-button>
          <n-button type="primary" @click="handleConnect" :loading="connectLoading">确认</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.network-detail {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 10px;
}

.title-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.view-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.view-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.subtitle-text {
  font-size: 14px;
  color: var(--n-text-color-3);
  display: flex;
  align-items: center;
}
</style>
