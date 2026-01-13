<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, NTag, NButton, NSpace, NText, NIcon, NGrid, NGi, NCard, NDataTable, NRadioGroup, NRadioButton, NInput, NSelect } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useContainerStore } from '../plugins/stores/containers'
import {
  TrashOutline,
  SearchOutline,
  SettingsOutline,
  StopCircleOutline,
  PlayCircleOutline,
  RefreshOutline,
  ListOutline,
  GridOutline
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const containerStore = useContainerStore()
const containers = computed(() => containerStore.containers)

const searchText = ref('')
const statusFilter = ref('all')
const loading = ref(false)
const viewMode = ref<'card'>('card')

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
  return container.names?.[0]?.replace(/^\//, '') || container.id
}

const formatId = (id: string) => {
  return id?.substring(0, 12) || ''
}

const handleStart = async (id: string) => {
  try {
    await containerStore.startContainer(id)
    message.success('容器启动成功')
  } catch (error: any) {
    message.error('启动失败：' + error.message)
  }
}

const handleStop = async (id: string) => {
  try {
    await containerStore.stopContainer(id)
    message.success('容器已停止')
  } catch (error: any) {
    message.error('停止失败：' + error.message)
  }
}

const handleDelete = async (id: string) => {
  const name = getContainerName({ names: containers.value.find(c => c.id === id)?.names })
  dialog.warning({
    title: '确认删除',
    content: `确定要删除容器 "${name || id}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await containerStore.deleteContainer(id)
        message.success('容器已删除')
      } catch (error: any) {
        message.error('删除失败：' + error.message)
      }
    }
  })
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

const rowProps = (row: any) => {
  return {
    style: 'cursor: pointer;',
    onClick: () => viewDetails(row.id)
  }
}

const createColumns = (): DataTableColumns<any> => {
  return [
    {
      title: '名称/ID',
      key: 'names',
      render(row) {
        return h('div', { class: 'container-identity' }, [
          h('div', { class: 'container-name' }, row.names?.[0]?.replace(/^\//, '') || row.id),
          h(NText, { depth: 3, style: 'font-size: 12px; font-family: monospace' }, { default: () => formatId(row.id) })
        ])
      }
    },
    {
      title: '镜像',
      key: 'image',
      render(row) {
        return h(NText, { depth: 3, style: 'font-size: 12px' }, { default: () => row.image || '-' })
      }
    },
    {
      title: '状态',
      key: 'state',
      render(row) {
        return h(
          'div',
          { style: 'display: flex; flex-direction: column; gap: 4px;' },
          [
            h(
              NTag,
              {
                type: row.state === 'running' ? 'success' : 'default',
                size: 'small',
                bordered: false,
                round: true
              },
              {
                default: () => row.state === 'running' ? 'Running' : row.state,
                icon: () => h('div', {
                  style: {
                    width: '6px',
                    height: '6px',
                    borderRadius: '50%',
                    backgroundColor: 'currentColor',
                    marginRight: '4px'
                  }
                })
              }
            ),
            row.status ? h(NText, { depth: 3, style: 'font-size: 12px;' }, { default: () => row.status }) : null
          ]
        )
      }
    },
    {
      title: '端口映射',
      key: 'ports',
      render(row) {
        if (!row.ports || row.ports.length === 0) return '-'
        return h(NSpace, { size: 4, wrap: true }, {
          default: () => row.ports.map((p: any) =>
            h(NTag, { size: 'small', bordered: true, style: 'font-size: 12px' }, {
              default: () => `${p.public_port}:${p.private_port}/${p.type}`
            })
          )
        })
      }
    },
    {
      title: '操作',
      key: 'actions',
      fixed: 'right',
      width: 120,
      render(row) {
        return h(NSpace, { size: 'small', onClick: (e: MouseEvent) => e.stopPropagation() }, {
          default: () => [
            h(NButton, {
              size: 'small',
              quaternary: true,
              circle: true,
              type: 'info',
              onClick: (e: MouseEvent) => {
                e.stopPropagation()
                viewDetails(row.id)
              }
            }, { icon: () => h(NIcon, { component: SettingsOutline }) }),
            h(NButton, {
              size: 'small',
              quaternary: true,
              circle: true,
              type: row.state === 'running' ? 'warning' : 'success',
              onClick: (e: MouseEvent) => {
                e.stopPropagation()
                row.state === 'running' ? handleStop(row.id) : handleStart(row.id)
              }
            }, { icon: () => h(NIcon, { component: row.state === 'running' ? StopCircleOutline : PlayCircleOutline }) }),
            h(NButton, {
              size: 'small',
              quaternary: true,
              circle: true,
              type: 'error',
              onClick: (e: MouseEvent) => {
                e.stopPropagation()
                handleDelete(row.id)
              }
            }, { icon: () => h(NIcon, { component: TrashOutline }) })
          ]
        })
      }
    }
  ]
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
          <n-radio-group v-model:value="viewMode">
            <n-radio-button value="table">
              <n-icon :component="ListOutline" />
            </n-radio-button>
            <n-radio-button value="card">
              <n-icon :component="GridOutline" />
            </n-radio-button>
          </n-radio-group>
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
      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4" responsive="screen">
        <n-gi v-for="container in filteredContainers" :key="container.id">
          <n-card hoverable class="container-card" style="cursor: pointer" @click="viewDetails(container.id)">
            <template #header>
              <div class="card-header">
                <div class="card-title" :title="getContainerName(container)">
                  {{ getContainerName(container) }}
                </div>
                <n-tag
                  :type="container.state === 'running' ? 'success' : 'default'"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ container.state }}
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

            <n-space vertical size="small">
              <div class="info-item">
                <span class="label">ID:</span>
                <span class="value monospace">{{ formatId(container.id) }}</span>
              </div>
              <div class="info-item">
                <span class="label">镜像:</span>
                <span class="value" :title="container.image">{{ container.image }}</span>
              </div>
              <div class="info-item" v-if="container.status">
                <span class="label">时长:</span>
                <span class="value">{{ container.status }}</span>
              </div>
              <div class="info-item">
                <span class="label">端口:</span>
                <div class="value">
                  <n-space size="small" v-if="container.ports && container.ports.length > 0">
                    <n-tag v-for="(p, i) in container.ports.slice(0, 2)" :key="i" size="tiny" :bordered="true">
                      {{ p.public_port }}:{{ p.private_port }}
                    </n-tag>
                    <n-tag v-if="container.ports.length > 2" size="tiny" :bordered="true">+{{ container.ports.length - 2 }}</n-tag>
                  </n-space>
                  <span v-else>-</span>
                </div>
              </div>
            </n-space>

            <template #action>
              <div class="card-actions">
                <n-button size="small" quaternary circle type="info" @click.stop="viewDetails(container.id)">
                  <template #icon><n-icon :component="SettingsOutline" /></template>
                </n-button>
                <n-button
                  size="small"
                  quaternary
                  circle
                  :type="container.state === 'running' ? 'warning' : 'success'"
                  @click.stop="container.state === 'running' ? handleStop(container.id) : handleStart(container.id)"
                >
                  <template #icon>
                    <n-icon :component="container.state === 'running' ? StopCircleOutline : PlayCircleOutline" />
                  </template>
                </n-button>
                <n-button size="small" quaternary circle type="error" @click.stop="handleDelete(container.id)">
                  <template #icon><n-icon :component="TrashOutline" /></template>
                </n-button>
              </div>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<style scoped>
.containers-page {
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.card-title {
  font-weight: 600;
  font-size: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
