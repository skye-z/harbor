<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog, NTag, NButton, NSpace, NIcon } from 'naive-ui'
import { networkApi, containerApi } from '../plugins/api'
import { useNetworkStore } from '../plugins/stores/networks'
import { useContainerStore } from '../plugins/stores/containers'
import {
  ArrowBackOutline,
  TrashOutline,
  LinkOutline,
  UnlinkOutline,
  RefreshOutline,
  GlobeOutline,
  ServerOutline,
  SettingsOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const networkStore = useNetworkStore()
const containerStore = useContainerStore()

const networkId = route.params.id as string
const loading = ref(false)
const activeTab = ref('info')
const showConnectModal = ref(false)
const selectedContainer = ref<string>('')
const connectLoading = ref(false)

const network = computed(() => {
  return networkStore.networks.find(n => n.id === networkId)
})

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
    if (!network.value) {
      message.error('网络不存在')
      router.push({ name: 'Storage' })
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
        await networkApi.delete(networkId)
        message.success('网络已删除')
        router.push({ name: 'Storage' })
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      }
    }
  })
}

// 断开容器连接
const handleDisconnect = async (containerId: string) => {
  try {
    await networkApi.disconnectContainer(networkId, containerId)
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
    await networkApi.connectContainer(networkId, selectedContainer.value)
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

onMounted(() => {
  loadNetworkDetail()
  containerStore.fetchContainers()
})
</script>

<template>
  <div class="network-detail-page">
    <div class="page-header">
      <n-button text @click="router.push({ name: 'Storage' })">
        <template #icon>
          <n-icon :component="ArrowBackOutline" />
        </template>
        返回
      </n-button>
      <h1 class="page-title">{{ network?.name || '网络详情' }}</h1>
      <n-button type="error" @click="handleDelete">
        <template #icon>
          <n-icon :component="TrashOutline" />
        </template>
        删除
      </n-button>
    </div>

    <n-tabs v-model:value="activeTab" type="line" style="margin-top: 20px;">
      <n-tab-pane name="info" tab="基本信息">
        <n-card :bordered="false">
          <n-descriptions bordered :column="2">
            <n-descriptions-item label="名称">{{ network?.name }}</n-descriptions-item>
            <n-descriptions-item label="驱动">{{ network?.driver }}</n-descriptions-item>
            <n-descriptions-item label="ID">{{ network?.id }}</n-descriptions-item>
            <n-descriptions-item label="范围">{{ network?.scope }}</n-descriptions-item>
            <n-descriptions-item label="子网">{{ network?.subnet || 'N/A' }}</n-descriptions-item>
            <n-descriptions-item label="网关">{{ network?.gateway || 'N/A' }}</n-descriptions-item>
            <n-descriptions-item label="创建时间">
              <n-time v-if="network?.created" :time="new Date(network.created).getTime()" type="datetime" />
            </n-descriptions-item>
            <n-descriptions-item label="接入容器数">{{ connectedContainers.length }}</n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-tab-pane>

      <n-tab-pane name="containers" tab="容器接入">
        <n-card :bordered="false">
          <template #header-extra>
            <n-button type="primary" size="small" @click="showConnectModal = true">
              <template #icon>
                <n-icon :component="LinkOutline" />
              </template>
              接入容器
            </n-button>
          </template>
          <n-list v-if="connectedContainers.length > 0">
            <n-list-item v-for="container in connectedContainers" :key="container.id">
              <n-thing :title="container.names[0]?.replace(/^\//, '') || container.id">
                <template #description>
                  <n-space>
                    <n-tag size="small" :type="container.state === 'running' ? 'success' : 'default'">
                      {{ container.state }}
                    </n-tag>
                    <span style="font-size: 12px; color: #999;">{{ container.image }}</span>
                    <n-tag size="small" type="info">IP: {{ getContainerIP(container) }}</n-tag>
                  </n-space>
                </template>
              </n-thing>
              <template #suffix>
                <n-button size="small" quaternary type="error" @click="handleDisconnect(container.id)">
                  <template #icon>
                    <n-icon :component="UnlinkOutline" />
                  </template>
                  断开
                </n-button>
              </template>
            </n-list-item>
          </n-list>
          <n-result v-else status="404" title="暂无容器接入" description="点击右上角按钮接入容器" style="margin-top: 10vh;">
            <template #icon>
              <n-icon size="80">
                <GlobeOutline />
              </n-icon>
            </template>
          </n-result>
        </n-card>
      </n-tab-pane>

      <n-tab-pane name="traffic" tab="流量监控">
        <n-card :bordered="false">
          <n-result status="info" title="功能开发中" description="网络流量监控功能即将上线" style="margin-top: 10vh;">
            <template #icon>
              <n-icon size="80">
                <SettingsOutline />
              </n-icon>
            </template>
          </n-result>
        </n-card>
      </n-tab-pane>
    </n-tabs>

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
.network-detail-page {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 10px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  flex: 1;
}
</style>
