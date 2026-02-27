<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, NButton, NIcon, NInput, NScrollbar, NSpin } from 'naive-ui'
import {
  ArrowBack,
  FolderOutline,
  DocumentOutline,
  RefreshOutline,
  DownloadOutline,
  CloudUploadOutline,
  ChevronUpOutline,
  ChevronForwardOutline
} from '@vicons/ionicons5'

interface FileItem {
  name: string
  path: string
  type: 'directory' | 'file'
  size: string
  mtime: string
  is_dir: boolean
  loading?: boolean
}

interface Column {
  id: string
  path: string
  files: FileItem[]
  selectedIndex: number | null
}

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()
const message = useMessage()

const containerId = computed(() => route.params.id as string)
const containerName = ref('')

const columns = ref<Column[]>([])
const fileLoading = ref(false)
const uploadLoading = ref(false)
const columnsContainerRef = ref<HTMLElement>()
const windowWidth = ref(window.innerWidth)

// 监听窗口大小变化
const updateWindowWidth = () => {
  windowWidth.value = window.innerWidth
}

// 响应式栏数
const visibleColumnsCount = computed(() => {
  const width = windowWidth.value
  if (width >= 1300) return 3
  if (width >= 900) return 2
  return 1
})

// 获取当前可见的列
const visibleColumns = computed(() => {
  const count = visibleColumnsCount.value
  if (columns.value.length <= count) {
    return columns.value
  }
  return columns.value.slice(-count)
})

// 排序文件列表：文件夹在前按a-z排序，文件在后按a-z排序
const sortFiles = (files: FileItem[]): FileItem[] => {
  const filtered = files.filter(f => f.name !== '.' && f.name !== '..')
  const dirs = filtered.filter(f => f.is_dir).sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
  const fileItems = filtered.filter(f => !f.is_dir).sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
  return [...dirs, ...fileItems]
}

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

const loadColumnFiles = async (path: string, columnIndex: number) => {
  try {
    fileLoading.value = true
    const result = await containerStore.listContainerFiles(containerId.value, path)
    const sortedFiles = sortFiles(result || [])
    
    if (columnIndex === 0) {
      columns.value = [{
        id: Date.now().toString(),
        path,
        files: sortedFiles,
        selectedIndex: null
      }]
    } else if (columnIndex >= columns.value.length) {
      columns.value.push({
        id: Date.now().toString(),
        path,
        files: sortedFiles,
        selectedIndex: null
      })
    } else {
      columns.value = columns.value.slice(0, columnIndex).map((col, idx) => ({
        ...col,
        selectedIndex: idx === columnIndex - 1 ? 0 : col.selectedIndex
      }))
      columns.value.push({
        id: Date.now().toString(),
        path,
        files: sortedFiles,
        selectedIndex: null
      })
    }
  } catch (error: any) {
    message.error('获取文件列表失败: ' + error.message)
  } finally {
    fileLoading.value = false
  }
}

const handleItemClick = (columnIndex: number, itemIndex: number, item: FileItem) => {
  const actualColumnIndex = columns.value.length - visibleColumns.value.length + columnIndex
  
  columns.value[actualColumnIndex].selectedIndex = itemIndex
  
  if (item.is_dir) {
    item.loading = true
    columns.value = columns.value.slice(0, actualColumnIndex + 1)
    loadColumnFiles(item.path, actualColumnIndex + 1).finally(() => {
      item.loading = false
    })
  }
}

const handleBreadcrumbClick = (index: number) => {
  if (index < columns.value.length - 1) {
    columns.value = columns.value.slice(0, index + 1)
    columns.value[index].selectedIndex = null
  }
}

const handleNavigateUp = () => {
  if (columns.value.length === 0) return
  const currentColumn = columns.value[columns.value.length - 1]
  if (currentColumn.path === '/') return
  
  const parentPath = currentColumn.path.split('/').slice(0, -1).join('/') || '/'
  loadColumnFiles(parentPath, Math.max(0, columns.value.length - 1))
}

const handleRefresh = () => {
  if (columns.value.length > 0) {
    const lastColumn = columns.value[columns.value.length - 1]
    loadColumnFiles(lastColumn.path, columns.value.length - 1)
  }
}

const handleDownload = async (file: FileItem) => {
  if (file.is_dir) {
    message.warning('暂不支持下载文件夹')
    return
  }
  try {
    fileLoading.value = true
    await containerStore.copyFromContainer(containerId.value, file.path, file.name)
    message.success('文件下载成功')
  } catch (error: any) {
    message.error('下载失败: ' + error.message)
  } finally {
    fileLoading.value = false
  }
}

const handleUpload = async (options: any) => {
  if (columns.value.length === 0) {
    message.warning('请先选择目录')
    return false
  }
  const currentPath = columns.value[columns.value.length - 1].path
  const file = options.file.file
  if (!file) return false

  try {
    uploadLoading.value = true
    await containerStore.copyToContainer(containerId.value, file.name, currentPath + '/' + file.name)
    message.success('文件上传成功')
    handleRefresh()
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

const getFileIconColor = (item: FileItem) => {
  if (item.is_dir) return '#dcb67a'
  const ext = item.name.split('.').pop()?.toLowerCase()
  const colorMap: Record<string, string> = {
    'js': '#f7df1e',
    'ts': '#3178c6',
    'vue': '#42b883',
    'css': '#264de4',
    'html': '#e34c26',
    'json': '#888888',
    'md': '#ffffff',
    'go': '#00add8',
    'py': '#3776ab',
    'java': '#b07219',
    'sh': '#89e051',
    'yml': '#cb171e',
    'yaml': '#cb171e',
    'dockerfile': '#2496ed',
    'sql': '#336791'
  }
  return colorMap[ext || ''] || '#75beff'
}

onMounted(async () => {
  window.addEventListener('resize', updateWindowWidth)
  await loadContainerDetail()
  if (!containerName.value) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
    return
  }
  await loadColumnFiles('/', 0)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateWindowWidth)
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
        <n-button quaternary size="small" @click="handleNavigateUp" :disabled="columns.length <= 1 && columns[0]?.path === '/'">
          <template #icon>
            <n-icon :component="ChevronUpOutline" />
          </template>
        </n-button>
        <div class="breadcrumb">
          <span 
            v-for="(col, index) in columns" 
            :key="col.id"
            class="breadcrumb-item"
            :class="{ active: index === columns.length - 1 }"
            @click="handleBreadcrumbClick(index)"
          >
            {{ col.path === '/' ? '根目录' : col.path.split('/').pop() }}
            <n-icon v-if="index < columns.length - 1" :component="ChevronForwardOutline" size="12" />
          </span>
        </div>
        <n-button quaternary size="small" @click="handleRefresh" :loading="fileLoading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
      <div class="toolbar-right">
        <n-upload :show-file-list="false" @before-upload="handleUpload" accept="*/*">
          <n-button quaternary size="small" :loading="uploadLoading">
            <template #icon>
              <n-icon :component="CloudUploadOutline" />
            </template>
          </n-button>
        </n-upload>
      </div>
    </div>

    <div class="file-browser">
      <div v-if="columns.length === 0" class="placeholder">
        <n-icon size="48" :component="FolderOutline" class="placeholder-icon" />
        <span>加载中...</span>
      </div>
      <div v-else class="columns-container" ref="columnsContainerRef">
        <div 
          v-for="(column, colIndex) in visibleColumns" 
          :key="column.id"
          class="column"
          :class="{ 
            'column-1': visibleColumnsCount === 1,
            'column-2': visibleColumnsCount === 2,
            'column-3': visibleColumnsCount === 3
          }"
        >
          <n-scrollbar class="column-scroll">
            <div v-if="column.files.length === 0" class="column-empty">
              <span>空目录</span>
            </div>
            <div
              v-for="(item, itemIndex) in column.files"
              :key="item.path"
              class="file-item"
              :class="{ 
                selected: column.selectedIndex === itemIndex,
                'is-directory': item.is_dir
              }"
              @click="handleItemClick(colIndex, itemIndex, item)"
            >
              <n-icon 
                :component="item.is_dir ? FolderOutline : DocumentOutline" 
                size="18"
                :color="getFileIconColor(item)"
                class="file-icon"
              />
              <span class="file-name" :title="item.name">{{ item.name }}</span>
              <n-spin
                v-if="item.is_dir && item.loading"
                size="small"
                class="folder-loading"
              />
              <n-button
                v-else-if="!item.is_dir"
                quaternary
                size="tiny"
                class="download-btn"
                @click.stop="handleDownload(item)"
              >
                <template #icon>
                  <n-icon :component="DownloadOutline" size="14" />
                </template>
              </n-button>
              <span class="file-size">{{ item.size }}</span>
            </div>
          </n-scrollbar>
        </div>
      </div>
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
  flex-shrink: 0;
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
  flex: 1;
  justify-content: center;
  min-width: 0;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #808080;
  overflow: hidden;
  flex: 1;
  justify-content: center;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.2s;
  white-space: nowrap;
}

.breadcrumb-item:hover {
  background: #3c3c3c;
  color: #cccccc;
}

.breadcrumb-item.active {
  color: #cccccc;
  font-weight: 500;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.file-browser {
  flex: 1;
  min-height: 0;
  overflow: hidden;
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

.columns-container {
  display: flex;
  height: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.column {
  flex-shrink: 0;
  border-right: 1px solid #3c3c3c;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.column:last-child {
  border-right: none;
}

@media (max-width: 899px) {
  .column {
    width: 100% !important;
  }
}

@media (min-width: 900px) and (max-width: 1299px) {
  .column {
    width: 50% !important;
  }
}

@media (min-width: 1300px) {
  .column {
    width: 33.333% !important;
  }
}

.column-scroll {
  flex: 1;
}

.column-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #808080;
  font-size: 13px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid #2d2d2d;
}

.file-item:hover {
  background: #2a2d2e;
}

.file-item.selected {
  background: #094771;
}

.file-item.selected:hover {
  background: #0d5a8f;
}

.file-icon {
  flex-shrink: 0;
}

.file-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  color: #d4d4d4;
}

.file-item.is-directory .file-name {
  color: #dcb67a;
}

.download-btn {
  opacity: 0;
  transition: opacity 0.2s;
  flex-shrink: 0;
  padding: 2px !important;
  height: 24px !important;
  width: 24px !important;
}

.file-item:hover .download-btn {
  opacity: 1;
}

.folder-loading {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-size {
  font-size: 11px;
  color: #808080;
  flex-shrink: 0;
  min-width: 60px;
  text-align: right;
}
</style>
