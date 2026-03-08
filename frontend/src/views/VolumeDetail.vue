<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog, NTag, NButton, NSpace, NIcon } from 'naive-ui'
import { volumeApi } from '../plugins/api'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useContainerStore } from '../plugins/stores/containers'
import {
  ArrowBackOutline,
  TrashOutline,
  DownloadOutline,
  FolderOutline,
  LinkOutline,
  UnlinkOutline,
  RefreshOutline,
  ServerOutline,
  CubeOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const volumeStore = useVolumeStore()
const containerStore = useContainerStore()

const volumeName = route.params.name as string
const loading = ref(false)
const exportLoading = ref(false)
const fileTreeLoading = ref(false)

const volume = computed(() => {
  // 通过 name 查找数据卷
  return volumeStore.volumes.find(v => v.name === volumeName)
})

// 获取挂载该卷的容器
const mountedContainers = computed(() => {
  return containerStore.containers.filter(container => {
    if (container.mounts && Array.isArray(container.mounts)) {
      return container.mounts.some((mount: any) => mount.name === volume.value?.name)
    }
    return false
  })
})

// 文件树数据
const fileTree = ref<any[]>([])

const loadVolumeDetail = async () => {
  loading.value = true
  try {
    await volumeStore.fetchVolumes()
    await containerStore.fetchContainers()
    if (!volume.value) {
      message.error('数据卷不存在')
      router.push({ name: 'Connect' })
      return
    }
  } catch (error: any) {
    message.error('加载数据卷信息失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 删除数据卷
const handleDelete = () => {
  if (mountedContainers.value.length > 0) {
    message.error('该数据卷正在被容器使用，请先解除挂载')
    return
  }

  dialog.warning({
    title: '确认删除',
    content: `确定要删除数据卷 "${volume.value?.name}" 吗？此操作不可逆。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await volumeApi.delete(volumeName)
        message.success('数据卷已删除')
        router.push({ name: 'Connect' })
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      }
    }
  })
}

// 打包导出数据卷
const handleExport = async () => {
  try {
    exportLoading.value = true
    const response = await volumeApi.export(volumeName)
    const blob = new Blob([response], { type: 'application/x-tar' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${volume.value?.name}.tar`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    message.success('导出成功')
  } catch (error: any) {
    message.error('导出失败: ' + error.message)
  } finally {
    exportLoading.value = false
  }
}

// 格式化大小
const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

onMounted(() => {
  loadVolumeDetail()
})
</script>

<template>
  <div class="volume-detail">
    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ volume?.name || volumeName }}</h1>
          <div class="header-actions">
            <n-button type="primary" ghost @click="handleExport" :loading="exportLoading">
              <template #icon>
                <n-icon :component="DownloadOutline" />
              </template>
              导出
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
          {{ volume?.mountpoint || '-' }}
        </div>
      </div>
    </div>

    <n-spin :show="loading">
      <n-grid item-responsive x-gap="10" cols="24">
        <!-- 左侧：基本信息 -->
        <n-gi span="24 760:9 900:8">
          <n-card size="small" title="基本信息" style="margin-bottom: 10px;">
            <template #header-extra>
              <n-tag type="info" size="small" :bordered="false">
                {{ volume?.driver || 'local' }}
              </n-tag>
            </template>
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item label="ID">
                {{ volume?.id || volume?.name || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="名称">
                {{ volume?.name || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="驱动">
                {{ volume?.driver || 'local' }}
              </n-descriptions-item>
              <n-descriptions-item label="挂载点">
                <n-ellipsis style="max-width: 200px;">
                  {{ volume?.mountpoint || '-' }}
                </n-ellipsis>
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">
                <n-time v-if="volume?.created_at" :time="new Date(volume.created_at).getTime()" type="datetime" />
                <span v-else>-</span>
              </n-descriptions-item>
              <n-descriptions-item label="范围">
                {{ volume?.scope || 'local' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>

          <n-card size="small" title="使用统计" style="margin-bottom: 10px;">
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item label="挂载容器数">
                {{ mountedContainers.length }}
              </n-descriptions-item>
              <n-descriptions-item label="引用次数">
                {{ volume?.usage_data?.ref_count || 0 }}
              </n-descriptions-item>
              <n-descriptions-item label="大小">
                {{ formatSize(volume?.usage_data?.size || 0) }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>

        <!-- 右侧：挂载容器 -->
        <n-gi span="24 760:15 900:16">
          <n-card size="small" title="挂载容器" style="margin-bottom: 10px;">
            <n-list v-if="mountedContainers.length > 0">
              <n-list-item v-for="container in mountedContainers" :key="container.id">
                <n-thing :title="container.names[0]?.replace(/^\//, '') || container.id">
                  <template #avatar>
                    <n-avatar>
                      <n-icon :component="CubeOutline" />
                    </n-avatar>
                  </template>
                  <template #description>
                    <n-space>
                      <n-tag size="small" :type="container.state === 'running' ? 'success' : 'default'">
                        {{ container.state }}
                      </n-tag>
                      <span style="font-size: 12px; color: #999;">{{ container.image }}</span>
                    </n-space>
                  </template>
                </n-thing>
                <template #suffix>
                  <n-button size="small" quaternary type="primary" @click="router.push({ name: 'ContainerDetail', params: { id: container.id } })">
                    查看
                  </n-button>
                </template>
              </n-list-item>
            </n-list>
            <n-empty v-else description="暂无容器挂载此数据卷" />
          </n-card>

          <n-card size="small" title="标签">
            <n-space v-if="volume?.labels && Object.keys(volume.labels).length > 0" wrap>
              <n-tag v-for="(value, key) in volume.labels" :key="key" size="small">
                {{ key }}: {{ value }}
              </n-tag>
            </n-space>
            <n-empty v-else description="暂无标签" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<style scoped>
.volume-detail {
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
  word-break: break-all;
}
</style>
