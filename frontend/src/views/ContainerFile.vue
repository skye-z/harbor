<template>
  <div class="container-file">
    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ containerName }} - 文件管理</h1>
          <div class="header-actions">
            <n-button size="medium" @click="router.back()">
              <template #icon>
                <n-icon :component="ArrowBackOutline" />
              </template>
              返回
            </n-button>
          </div>
        </div>
        <div class="subtitle-text">
          容器ID: {{ containerId }}
        </div>
      </div>
    </div>

    <n-card>
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="handleNavigateUp" :disabled="currentPath === '/'">
            <template #icon>
              <n-icon :component="FolderOpenOutline" />
            </template>
            返回上级
          </n-button>
          <n-input v-model:value="currentPath" placeholder="输入容器内路径，如 /app" style="width: 300px" @keyup.enter="loadFileList" />
          <n-button size="small" type="primary" @click="loadFileList" :loading="fileLoading">
            <template #icon>
              <n-icon :component="SearchOutline" />
            </template>
            浏览
          </n-button>
          <n-button size="small" type="primary" @click="handleDownload" :loading="fileLoading">
            <template #icon>
              <n-icon :component="DownloadOutline" />
            </template>
            下载
          </n-button>
          <n-upload :show-file-list="false" @before-upload="handleUpload" accept="*/*">
            <n-button size="small" type="info" :loading="uploadLoading">
              <template #icon>
                <n-icon :component="CloudUploadOutline" />
              </template>
              上传
            </n-button>
          </n-upload>
        </n-space>
      </template>
      <n-alert type="info" :show-icon="false" style="margin-bottom: 16px">
        输入容器内文件/目录路径，点击浏览查看目录内容，点击下载可将文件下载到本地。点击上传可选择本地文件上传到容器当前路径。
      </n-alert>
      <n-data-table v-if="fileList.length > 0" :columns="fileColumns" :data="fileList" :bordered="false"
        :pagination="{ pageSize: 20 }" />
      <n-empty v-else description="暂无文件信息，请在上方输入路径并点击浏览" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, NIcon } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  ArrowBackOutline,
  FolderOutline,
  FolderOpenOutline,
  DocumentOutline,
  SearchOutline,
  DownloadOutline,
  CloudUploadOutline
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
          style: { display: 'flex', alignItems: 'center', gap: '8px', cursor: row.is_dir ? 'pointer' : 'default', color: row.is_dir ? '#2080f0' : 'inherit' },
          onClick: () => handleNavigate(row)
        }, [
          h(NIcon, {
            component: row.is_dir ? FolderOutline : DocumentOutline,
            size: 18,
            color: row.is_dir ? '#f0a020' : '#2080f0'
          }),
          h('span', row.name)
        ])
      }
    },
    { title: '路径', key: 'path', ellipsis: { tooltip: true } },
    { title: '类型', key: 'type', width: 100 },
    { title: '大小', key: 'size', width: 100 },
    { title: '修改时间', key: 'mtime', width: 150 }
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

<style scoped>
.container-file {
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
  align-items: center;
}

.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-group .subtitle-text {
  color: var(--n-text-color-3);
  font-size: 14px;
  margin-top: 4px;
}
</style>
