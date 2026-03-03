<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { volumeApi, containerApi } from '../plugins/api'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useContainerStore } from '../plugins/stores/containers'
import {
  ArrowBackOutline,
  TrashOutline,
  DownloadOutline,
  FolderOutline,
  LinkOutline,
  UnlinkOutline,
  RefreshOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const volumeStore = useVolumeStore()
const containerStore = useContainerStore()

const volumeId = route.params.id as string
const loading = ref(false)
const exportLoading = ref(false)
const fileTreeLoading = ref(false)
const activeTab = ref('info')

const volume = computed(() => {
  return volumeStore.volumes.find(v => v.id === volumeId)
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
const selectedFile = ref<string>('')

const loadVolumeDetail = async () => {
  loading.value = true
  try {
    await volumeStore.fetchVolumes()
    if (!volume.value) {
      message.error('数据卷不存在')
      router.push({ name: 'Storage' })
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
        await volumeApi.delete(volumeId)
        message.success('数据卷已删除')
        router.push({ name: 'Storage' })
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
    const response = await volumeApi.export(volumeId)
    // 创建下载链接
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

// 加载文件树
const loadFileTree = async (path: string = '/') => {
  fileTreeLoading.value = true
  try {
    const files = await volumeApi.listFiles(volumeId, path)
    fileTree.value = files
  } catch (error: any) {
    message.error('加载文件列表失败: ' + error.message)
  } finally {
    fileTreeLoading.value = false
  }
}

// 解除挂载
const handleUnmount = async (containerId: string) => {
  try {
    await containerApi.unmountVolume(containerId, volume.value?.name)
    message.success('已解除挂载')
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('解除挂载失败: ' + error.message)
  }
}

onMounted(() => {
  loadVolumeDetail()
  containerStore.fetchContainers()
})
</script>

<template>
  <div class="volume-detail-page">
    <div class="page-header">
      <n-button text @click="router.push({ name: 'Storage' })">
        <template #icon>
          <n-icon :component="ArrowBackOutline" />
        </template>
        返回
      </n-button>
      <h1 class="page-title">{{ volume?.name || '数据卷详情' }}</h1>
      <n-space>
        <n-button @click="handleExport" :loading="exportLoading">
          <template #icon>
            <n-icon :component="DownloadOutline" />
          </template>
          导出
        </n-button>
        <n-button type="error" @click="handleDelete">
          <template #icon>
            <n-icon :component="TrashOutline" />
          </template>
          删除
        </n-button>
      </n-space>
    </div>

    <n-tabs v-model:value="activeTab" type="line" style="margin-top: 20px;">
      <n-tab-pane name="info" tab="基本信息">
        <n-card :bordered="false">
          <n-descriptions bordered :column="2">
            <n-descriptions-item label="名称">{{ volume?.name }}</n-descriptions-item>
            <n-descriptions-item label="驱动">{{ volume?.driver }}</n-descriptions-item>
            <n-descriptions-item label="挂载点">{{ volume?.mountpoint }}</n-descriptions-item>
            <n-descriptions-item label="创建时间">
              <n-time v-if="volume?.created_at" :time="new Date(volume.created_at).getTime()" type="datetime" />
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-tab-pane>

      <n-tab-pane name="containers" tab="容器挂载">
        <n-card :bordered="false">
          <n-list v-if="mountedContainers.length > 0">
            <n-list-item v-for="container in mountedContainers" :key="container.id">
              <n-thing :title="container.names[0]?.replace(/^\//, '') || container.id">
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
                <n-button size="small" quaternary type="error" @click="handleUnmount(container.id)">
                  <template #icon>
                    <n-icon :component="UnlinkOutline" />
                  </template>
                  解除挂载
                </n-button>
              </template>
            </n-list-item>
          </n-list>
          <n-empty v-else description="暂无容器挂载此数据卷" />
        </n-card>
      </n-tab-pane>

      <n-tab-pane name="files" tab="文件浏览">
        <n-card :bordered="false">
          <n-spin :show="fileTreeLoading">
            <n-tree
              :data="fileTree"
              :render-prefix="() => h(NIcon, { component: FolderOutline })"
              selectable
              @update:selected-keys="(keys) => selectedFile = keys[0] as string"
            />
          </n-spin>
          <n-empty v-if="fileTree.length === 0 && !fileTreeLoading" description="暂无文件数据" />
        </n-card>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<style scoped>
.volume-detail-page {
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
