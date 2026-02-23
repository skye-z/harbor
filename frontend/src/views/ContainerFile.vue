<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, NButton, NIcon, NInput, NDataTable, NUpload } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  ArrowBack,
  FolderOutline,
  DocumentOutline,
  RefreshOutline,
  DownloadOutline,
  CloudUploadOutline,
  ChevronUpOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()
const message = useMessage()

const containerId = computed(() => route.params.id as string)
const containerName = ref('')

const currentPath = ref('/')
const fileList = ref<any[]>([])
const fileLoading = ref(false)
const uploadLoading = ref(false)

const fileColumns: DataTableColumns<{
  name: string
  path: string
  type: string
  size: string
  mtime: string
}> = [
  {
    title: '名称',
    key: 'name',
    render(row: any) {
      return h('div', {
        style: { display: 'flex', alignItems: 'center', gap: '8px', cursor: row.is_dir ? 'pointer' : 'default' },
        onClick: () => handleNavigate(row)
      }, [
        h(NIcon, {
          component: row.is_dir ? FolderOutline : DocumentOutline,
          size: 16,
          color: row.is_dir ? '#dcb67a' : '#75beff'
        }),
        h('span', {
          style: { color: row.is_dir ? '#dcb67a' : '#d4d4d4' }
        }, row.name)
      ])
    }
  },
  { title: '大小', key: 'size', width: 100, ellipsis: { tooltip: true } },
  { title: '修改时间', key: 'mtime', width: 160, ellipsis: { tooltip: true } }
]

const loadContainerDetail = async () => {
  try {
    const data = await containerStore.getContainerInfo(containerId.value)
    if (data) {
      containerName.value = data.Name?.replace(/^\//, '') || containerId.value
    }
  } catch (error: any) {
    message.error('加载容器详情失败: ' + error.message)
  }
}

const loadFileList = async () => {
  if (!currentPath.value) {
    message.warning('请输入路径')
    return
  }
  try {
    fileLoading.value = true
    const result = await containerStore.listContainerFiles(containerId.value, currentPath.value)
    fileList.value = result || []
  } catch (error: any) {
    message.error('获取文件列表失败: ' + error.message)
    fileList.value = []
  } finally {
    fileLoading.value = false
  }
}

const handleNavigateUp = () => {
  if (currentPath.value === '/') return
  const parent = currentPath.value.split('/').slice(0, -1).join('/')
  currentPath.value = parent || '/'
  loadFileList()
}

const handleNavigate = (row: any) => {
  if (row.is_dir) {
    currentPath.value = row.path
    loadFileList()
  }
}

const handleDownload = async () => {
  if (!currentPath.value) {
    message.warning('请输入文件路径')
    return
  }
  try {
    fileLoading.value = true
    await containerStore.copyFromContainer(containerId.value, currentPath.value)
    message.success('文件下载成功')
  } catch (error: any) {
    message.error('下载失败: ' + error.message)
  } finally {
    fileLoading.value = false
  }
}

const handleUpload = async (options: any) => {
  if (!currentPath.value) {
    message.warning('请输入目标路径')
    return false
  }
  const file = options.file.file
  if (!file) return false

  try {
    uploadLoading.value = true
    await containerStore.copyToContainer(containerId.value, file.name, currentPath.value + '/' + file.name)
    message.success('文件上传成功')
    loadFileList()
  } catch (error: any) {
    message.error('上传失败: ' + error.message)
  } finally {
    uploadLoading.value = false
  }
  return false
}

const handleBack = () => {
  router.push({ name: 'ContainerDetail', params: { id: containerId.value } })
}

onMounted(async () => {
  await loadContainerDetail()
  if (!containerName.value) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
    return
  }
  await loadFileList()
})
</script>

<template>
  <div class="file-page">
    <div class="toolbar">
      <div class="toolbar-left">
        <n-button quaternary size="small" @click="handleBack">
          <template #icon>
            <n-icon :component="ArrowBack" />
          </template>
        </n-button>
        <span class="container-name">{{ containerName }}</span>
      </div>
      <div class="toolbar-center">
        <n-button quaternary size="small" @click="handleNavigateUp" :disabled="currentPath === '/'">
          <template #icon>
            <n-icon :component="ChevronUpOutline" />
          </template>
        </n-button>
        <n-input
          v-model:value="currentPath"
          placeholder="路径"
          size="small"
          style="width: 300px"
          @keyup.enter="loadFileList"
        />
        <n-button quaternary size="small" @click="loadFileList" :loading="fileLoading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
      <div class="toolbar-right">
        <n-button quaternary size="small" @click="handleDownload" :loading="fileLoading">
          <template #icon>
            <n-icon :component="DownloadOutline" />
          </template>
        </n-button>
        <n-upload :show-file-list="false" @before-upload="handleUpload" accept="*/*">
          <n-button quaternary size="small" :loading="uploadLoading">
            <template #icon>
              <n-icon :component="CloudUploadOutline" />
            </template>
          </n-button>
        </n-upload>
      </div>
    </div>

    <div class="file-area">
      <div v-if="fileList.length === 0 && !fileLoading" class="placeholder">
        <n-icon size="48" :component="FolderOutline" class="placeholder-icon" />
        <span>空目录</span>
      </div>
      <n-data-table
        v-else
        :columns="fileColumns"
        :data="fileList"
        :bordered="false"
        :loading="fileLoading"
        :pagination="false"
        size="small"
        class="file-table"
      />
    </div>
  </div>
</template>

<style scoped>
.file-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  z-index: 100;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
  flex-shrink: 0;
  gap: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.container-name {
  color: #cccccc;
  font-weight: 500;
  font-size: 14px;
}

.toolbar-center {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.file-area {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 8px 0;
}

.placeholder {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #808080;
}

.placeholder-icon {
  opacity: 0.5;
}

.file-table {
  background: transparent !important;
}

.file-table :deep(.n-data-table-th) {
  background: #252526 !important;
  color: #808080 !important;
  border-bottom: 1px solid #3c3c3c !important;
  font-weight: 500;
}

.file-table :deep(.n-data-table-td) {
  background: transparent !important;
  color: #d4d4d4 !important;
  border-bottom: 1px solid #2d2d2d !important;
}

.file-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #2a2d2e !important;
}

.file-table :deep(.n-data-table-empty) {
  padding: 48px 0;
}
</style>
