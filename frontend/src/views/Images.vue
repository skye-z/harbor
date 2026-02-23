<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useImageStore } from '../plugins/stores/images'
import { useMessage, useDialog } from 'naive-ui'
import {
  SearchOutline,
  CloudDownloadOutline,
  BuildOutline,
  RefreshOutline,
  Trash,
  Balloon
} from '@vicons/ionicons5'

const router = useRouter()
const imageStore = useImageStore()
const message = useMessage()
const dialog = useDialog()

const searchText = ref('')
const showBuildModal = ref(false)
const loading = ref(false)

const buildForm = ref({
  imageName: '',
  dockerfile: 'FROM alpine\nRUN apk add --no-cache curl'
})

const filteredImages = computed(() => {
  if (!searchText.value) return imageStore.images

  const search = searchText.value.toLowerCase()
  return imageStore.images.filter(i => {
    const tags = i.repo_tags?.join(' ') || ''
    const digest = i.repo_digests?.join(' ') || ''
    return (
      tags.includes(search) ||
      digest.includes(search) ||
      i.id.includes(search)
    )
  })
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor(diff / (1000 * 60))

  if (days > 0) return `${days} 天前`
  if (hours > 0) return `${hours} 小时前`
  if (minutes > 0) return `${minutes} 分钟前`
  return '刚刚'
}

const formatId = (id: string) => {
  return id.startsWith('sha256:') ? id.substring(7, 19) : id.substring(0, 12)
}

const getImageName = (image: any) => {
  if (image.repo_tags && image.repo_tags.length > 0) {
    return image.repo_tags[0].split(':')[0]
  }
  return formatId(image.id)
}

const copyId = (id: string) => {
  navigator.clipboard.writeText(id)
  message.success('已复制镜像ID')
}

const refresh = async () => {
  try {
    loading.value = true
    await imageStore.fetchImages()
  } catch (error: any) {
    message.error('刷新失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleDelete = (id: string) => {
  const image = imageStore.getImageById(id)
  const name = image?.repo_tags?.[0] || formatId(id)
  const d = dialog.warning({
    title: '确认删除',
    content: `确定要删除镜像 "${name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      d.destroy()
      try {
        await imageStore.deleteImage(id)
        message.success('镜像已删除')
        refresh()
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      }
    }
  })
}

onMounted(async () => {
  loading.value = true
  try {
    await imageStore.fetchImages()
  } catch (error: any) {
    message.error('加载镜像列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
})

const handleRefresh = async () => {
  loading.value = true
  try {
    await imageStore.fetchImages()
  } catch (error: any) {
    message.error('刷新失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const viewDetails = (id: string) => {
  router.push({ name: 'ImageDetail', params: { id } })
}
</script>

<template>
  <div class="images-page">
    <div class="view-header">
      <div class="filter-bar-header">
        <n-input v-model:value="searchText" placeholder="搜索镜像..." clearable style="width: 240px">
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
        <n-space>
          <n-button type="primary" @click="router.push('/images/pull')">
            <template #icon>
              <n-icon :component="CloudDownloadOutline" />
            </template>
            拉取
          </n-button>
          <n-button type="info" @click="showBuildModal = true">
            <template #icon>
              <n-icon :component="BuildOutline" />
            </template>
            构建
          </n-button>
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
      <n-grid x-gap="16" y-gap="16" cols="1 s:2 m:3 l:4 xl:5" responsive="screen" v-if="filteredImages.length > 0">
        <n-gi v-for="image in filteredImages" :key="image.id">
          <n-card hoverable class="image-card" style="cursor: pointer" @click="viewDetails(image.id)">
            <template #header>
              <div class="card-header">
                <div class="card-title-area">
                  <div class="card-title" :title="getImageName(image)">
                    {{ getImageName(image) }}
                  </div>
                  <div class="id-tags">
                    <n-tag size="small" :bordered="false" @click.stop="copyId(image.id)" style="cursor: pointer">
                      {{ formatId(image.id) }}
                    </n-tag>
                    <template v-if="image.repo_tags && image.repo_tags.length > 0">
                      <n-tag v-for="(tag, i) in (image.repo_tags || []).slice(0, 2)" type="info" :key="i" size="small"
                        :bordered="false" style="cursor: pointer">
                        {{ tag }}
                      </n-tag>
                      <n-tag v-if="image.repo_tags && image.repo_tags.length > 2" type="info" size="small" :bordered="false" style="cursor: pointer">
                        +{{ image.repo_tags.length - 2 }}
                      </n-tag>
                    </template>
                  </div>
                </div>
                <div class="card-meta">
                  <span class="meta-item">{{ formatTime(image.created) }}构建</span>
                </div>
              </div>
            </template>
            <template #action>
              <div style="display: flex;justify-content: space-between;align-items: center;">
                <div style="font-size: 12px;color: #999;">{{ formatSize(image.size) }}</div>
                <n-button size="small" quaternary circle type="error"
                  :disabled="image.containers > 0"
                  @click.stop="handleDelete(image.id)">
                  <template #icon>
                    <n-icon>
                      <Trash />
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-result v-else-if="!loading" status="404" title="这里什么都没有" description="在找啥呢? 要不要先拉取一个镜像试试" style="margin-top: 20vh;">
        <template #icon>
          <n-icon size="100">
            <Balloon />
          </n-icon>
        </template>
        <template #footer>
          <n-button type="primary" strong secondary @click="router.push('/images/pull')">
            <template #icon>
              <n-icon>
                <CloudDownloadOutline />
              </n-icon>
            </template>
            创建容器
          </n-button>
        </template>
      </n-result>
    </div>

    <n-modal v-model:show="showBuildModal" preset="card" title="构建镜像" style="width: 600px">
      <n-form :model="buildForm" label-placement="left" label-width="100px">
        <n-form-item label="镜像名称">
          <n-input v-model:value="buildForm.imageName" placeholder="myapp:latest" />
        </n-form-item>
        <n-form-item label="Dockerfile">
          <n-input v-model:value="buildForm.dockerfile" type="textarea" :autosize="{ minRows: 6, maxRows: 15 }"
            placeholder="FROM nginx:latest" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showBuildModal = false">取消</n-button>
          <n-button type="primary" @click="handleBuild" :loading="loading">
            构建
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.images-page {
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
}

.card-title-area {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.card-title {
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  font-weight: 600;
  font-size: 16px;
  max-width: 100%;
}

.id-tags {
  align-items: center;
  display: flex;
  gap: 4px;
}

.card-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  flex-shrink: 0;
  margin-left: 12px;
}

.meta-item {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}

.card-footer {
  padding: 8px 0 0 0;
  border-top: 1px solid var(--n-border-color);
  margin-top: 8px;
}

.size-text {
  font-size: 13px;
  color: var(--n-text-color-2);
}
</style>
