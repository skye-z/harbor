<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NTag, NButton, NSpace, NIcon } from 'naive-ui'
import { useContainerStore } from '../plugins/stores/containers'
import { formatContainerState, getContainerStateType } from '../plugins/utils/container'
import {
  Trash,
  SearchOutline,
  Stop,
  Play,
  Pause,
  RefreshOutline,
  Balloon,
  Add
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const containerStore = useContainerStore()
const containers = computed(() => containerStore.containers)

const searchText = ref('')
const statusFilter = ref('all')
const loading = ref(false)
const loadingStates = ref<Record<string, boolean>>({})

const statusOptions = [
  { label: '全部', value: 'all' },
  { label: '运行中', value: 'running' },
  { label: '已停止', value: 'exited' },
  { label: '已创建', value: 'created' }
]

const filteredContainers = computed(() => {
  let result = containers.value
  if (statusFilter.value !== 'all') {
    result = result.filter(c => c.state === statusFilter.value)
  }
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(c =>
      c.names?.some(name => name?.toLowerCase().includes(search)) ||
      c.id?.toLowerCase().includes(search)
    )
  }
  return result
})

const viewDetails = (id: string) => {
  router.push({ name: 'ContainerDetail', params: { id } })
}

const getContainerName = (container: any) => {
  return formatContainerName(container.names?.[0]?.replace(/^\//, '') || container.id)
}

function formatContainerName(name: string) {
  if (typeof name !== 'string' || name.trim().length === 0) {
    return name || '';
  }

  const replacedStr = name.replace(/[_-]/g, ' ');
  const formattedStr = replacedStr.split(/\s+/)
    .filter(word => word.length > 0)
    .map(word => {
      return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
    })
    .join(' ');

  return formattedStr;
}

const formatId = (id: string) => {
  return id?.substring(0, 12) || ''
}

const copyId = (id: string) => {
  navigator.clipboard.writeText(id)
  message.success('已复制容器ID')
}

const openPort = (port: number) => {
  window.open(`http://localhost:${port}`, '_blank')
}

const formatPorts = (ports: any[]) => {
  if (!ports || ports.length === 0) return []
  const portMap = new Map<string, { public_port: number; private_port: number; types: string[]; ips: string[] }>()
  for (const port of ports) {
    const key = `${port.public_port}:${port.private_port}`
    if (!portMap.has(key)) {
      portMap.set(key, {
        public_port: port.public_port,
        private_port: port.private_port,
        types: [],
        ips: []
      })
    }
    if (port.type && !portMap.get(key)!.types.includes(port.type)) {
      portMap.get(key)!.types.push(port.type)
    }
    if (port.ip && !port.ip.startsWith('invalid') && !portMap.get(key)!.ips.includes(port.ip)) {
      portMap.get(key)!.ips.push(port.ip)
    }
  }
  return Array.from(portMap.values())
}

const handleOperation = async (id: string, action: string, btnType: string) => {
  const key = `${id}-${btnType}`
  loadingStates.value[key] = true
  try {
    await containerStore.operationContainer(id, action)
  } catch (error: any) {
    message.error('操作失败：' + error.message)
  } finally {
    loadingStates.value[key] = false
  }
}

const handleDelete = async (id: string) => {
  const name = getContainerName({ names: containers.value.find(c => c.id === id)?.names })
  const d = dialog.warning({
    title: '确认删除',
    content: `确定要删除容器 "${name || id}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      d.destroy()
      loadingStates.value[`${id}-delete`] = true
      try {
        await containerStore.operationContainer(id, 'remove')
        message.success('容器已删除')
      } catch (error: any) {
        message.error('删除失败：' + error.message)
      } finally {
        loadingStates.value[`${id}-delete`] = false
      }
    }
  })
}

const handleCreateContainer = () => {
  router.push({ name: 'ContainerCreate' })
}

onMounted(async () => {
  loading.value = true
  try {
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('加载容器列表失败：' + error.message)
  } finally {
    loading.value = false
  }
})

const handleRefresh = async () => {
  loading.value = true
  try {
    await containerStore.fetchContainers()
    message.success('刷新成功')
  } catch (error: any) {
    message.error('刷新失败：' + error.message)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('加载容器列表失败：' + error.message)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="containers-page">
    <div class="view-header">
      <div class="filter-bar-header">
        <n-space :size="16">
          <n-input v-model:value="searchText" placeholder="搜索容器..." clearable style="width: 240px">
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
          </n-input>
          <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="状态筛选" style="width: 150px" />
        </n-space>
        <n-space>
          <n-button @click="handleRefresh" :loading="loading">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            刷新
          </n-button>
        </n-space>
      </div>
    </div>

    <div class="content-area">
      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen" v-if="filteredContainers.length > 0">
        <n-gi v-for="container in filteredContainers" :key="container.id">
          <n-card hoverable class="container-card" style="cursor: pointer" @click="viewDetails(container.id)">
            <template #header>
              <div class="card-header">
                <div>
                  <div class="card-title" :title="getContainerName(container)">
                    {{ getContainerName(container) }}
                  </div>
                  <n-space size="small">
                    <n-tag size="small" :bordered="false" @click.stop="copyId(container.id)" style="cursor: pointer">
                      {{ formatId(container.id) }}
                    </n-tag>
                    <template v-for="(port, index) in formatPorts(container.ports)" :key="index">
                      <n-tag v-if="port.public_port !== 0" size="small" :bordered="false" type="success" style="cursor: pointer" @click.stop="openPort(port.public_port)">
                        {{ port.public_port }}
                      </n-tag>
                    </template>
                    <n-tag v-if="!container.ports || container.ports.length === 0" size="small" :bordered="false">
                      -
                    </n-tag>
                  </n-space>
                </div>
                <n-tag
                  :type="getContainerStateType(container.state)"
                  size="small" round :bordered="false">
                  {{ formatContainerState(container.state) }}
                  <template #icon>
                    <div :style="{
                      width: '6px',
                      height: '6px',
                      borderRadius: '50%',
                      backgroundColor: 'currentColor',
                      marginRight: '0'
                    }"></div>
                  </template>
                </n-tag>
              </div>
            </template>
            <template #action>
              <n-button class="del-btn" size="small" quaternary circle type="error"
                :disabled="container.state === 'running' || loadingStates[`${container.id}-delete`]"
                :loading="loadingStates[`${container.id}-delete`]" @click.stop="handleDelete(container.id)">
                <template #icon>
                  <n-icon>
                    <Trash />
                  </n-icon>
                </template>
              </n-button>
              <n-button-group>
                <n-button size="small" tertiary :type="container.state === 'running' ? 'error' : 'success'"
                  :disabled="container.state === 'paused' || Object.values(loadingStates).some(v => v)"
                  :loading="loadingStates[`${container.id}-startstop`]"
                  @click.stop="container.state === 'running' ? handleOperation(container.id, 'stop', 'startstop') : handleOperation(container.id, 'start', 'startstop')">
                  <template #icon>
                    <n-icon :component="container.state === 'running' ? Stop : Play" />
                  </template>
                </n-button>
                <n-button size="small" tertiary type="warning"
                  :disabled="container.state !== 'running' && container.state !== 'paused' || Object.values(loadingStates).some(v => v)"
                  :loading="loadingStates[`${container.id}-pause`]"
                  @click.stop="container.state === 'running' ? handleOperation(container.id, 'pause', 'pause') : handleOperation(container.id, 'unpause', 'pause')">
                  <template #icon>
                    <n-icon>
                      <component :is="container.state === 'running' ? Pause : Play" />
                    </n-icon>
                  </template>
                </n-button>
              </n-button-group>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-result v-else status="404" title="这里什么都没有" description="在找啥呢? 要不要先创建一个容器试试" style="margin-top: 20vh;">
        <template #icon>
          <n-icon size="100">
            <Balloon />
          </n-icon>
        </template>
        <template #footer>
          <n-button type="primary" strong secondary @click="handleCreateContainer">
            <template #icon>
              <n-icon>
                <Add />
              </n-icon>
            </template>
            创建容器
          </n-button>
        </template>
      </n-result>
    </div>
  </div>
</template>

<style scoped>
.containers-page {
  padding: 0 10px 10px 10px;
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
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

.content-area {
  width: 100%;
}

.table-view {
  background-color: var(--n-color-modal);
  border-radius: 8px;
}

:deep(.container-name) {
  font-weight: 600;
  font-size: 14px;
}

:deep(.container-identity) {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.card-header {
  justify-content: space-between;
  display: flex;
  width: 100%;
}

.card-title {
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  font-weight: 600;
  font-size: 16px;
  max-width: 70%;
}

.info-item {
  display: flex;
  font-size: 13px;
  line-height: 1.6;
}

.info-item .label {
  color: var(--n-text-color-3);
  width: 40px;
  flex-shrink: 0;
}

.info-item .value {
  color: var(--n-text-color-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.info-item .value.monospace {
  font-family: monospace;
}

.del-btn {
  float: right;
}
</style>
