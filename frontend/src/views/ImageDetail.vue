<template>
  <div class="image-detail">
    <div class="page-header">
      <n-breadcrumb>
        <n-breadcrumb-item @click="$router.push({ name: 'Images' })">Images</n-breadcrumb-item>
        <n-breadcrumb-item>{{ imageDisplayName }}</n-breadcrumb-item>
      </n-breadcrumb>

      <div class="view-header">
        <div class="title-group">
          <h1>
            {{ imageName }}
            <span class="subtitle" v-if="imageTag">: {{ imageTag }}</span>
            <span class="id-tag">{{ shortImageId }}</span>
          </h1>
          <div class="subtitle-text">
            {{ image?.os }}/{{ image?.architecture }}
          </div>
        </div>
        <div class="header-actions">
          <n-button type="primary" @click="handleCreateContainer">
            <template #icon>
              <n-icon :component="PlayCircleOutline" />
            </template>
            创建容器
          </n-button>
          <n-button type="error" ghost @click="handleDelete">
            <template #icon>
              <n-icon :component="TrashOutline" />
            </template>
            删除镜像
          </n-button>
        </div>
      </div>
    </div>

    <n-grid x-gap="24" cols="24">
      <n-gi span="16">
        <n-card title="镜像详情" size="small">
          <n-descriptions :column="2" label-placement="left">
            <n-descriptions-item label="ID">
              <n-text code>{{ image?.id?.substring(0, 20) }}...</n-text>
            </n-descriptions-item>
            <n-descriptions-item label="创建时间">
              <n-time v-if="image" :time="new Date(image.created * 1000)" type="date" />
            </n-descriptions-item>
            <n-descriptions-item label="大小">
              {{ formatSize(image?.size || 0) }}
            </n-descriptions-item>
            <n-descriptions-item label="架构">
              {{ image?.os }} / {{ image?.architecture }}
            </n-descriptions-item>
            <n-descriptions-item label="标签">
              <n-space>
                <n-tag v-for="tag in image?.repo_tags" :key="tag" type="info" size="small">
                  {{ tag }}
                </n-tag>
                <n-tag v-if="!image?.repo_tags?.length" type="warning" size="small">
                  none
                </n-tag>
              </n-space>
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-gi>

      <n-gi span="8">
        <n-card title="标签管理" size="small">
          <template #header-extra>
            <n-button size="small" type="primary" @click="showTagModal = true">
              添加标签
            </n-button>
          </template>
          <n-space vertical v-if="image?.repo_tags?.length">
            <n-tag
              v-for="tag in image?.repo_tags"
              :key="tag"
              type="info"
              closable
              @close="handleRemoveTag(tag)"
            >
              {{ tag }}
            </n-tag>
          </n-space>
          <n-empty v-else description="暂无标签" />
        </n-card>
      </n-gi>
    </n-grid>

    <n-modal v-model:show="showTagModal" preset="card" title="添加标签" style="width: 400px">
      <n-form :model="tagForm" label-placement="left" label-width="80px">
        <n-form-item label="标签">
          <n-input v-model:value="tagForm.tag" placeholder="myrepo:v1.0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showTagModal = false">取消</n-button>
          <n-button type="primary" @click="handleAddTag" :loading="tagLoading">添加</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useImageStore } from '../plugins/stores/images'
import { useMessage, useDialog } from 'naive-ui'
import {
  PlayCircleOutline,
  TrashOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const imageStore = useImageStore()
const message = useMessage()
const dialog = useDialog()

const showTagModal = ref(false)
const tagForm = ref({ tag: '' })
const tagLoading = ref(false)

const imageId = computed(() => route.params.id as string)

const image = computed(() => imageStore.getImageById(imageId.value))

const imageDisplayName = computed(() => {
  const img = image.value
  return img?.repo_tags?.[0] || img?.id || 'Unknown'
})

const imageName = computed(() => {
  const name = imageDisplayName.value
  if (name.includes(':')) {
    return name.split(':')[0]
  }
  return name
})

const imageTag = computed(() => {
  const name = imageDisplayName.value
  if (name.includes(':')) {
    return name.split(':')[1]
  }
  return ''
})

const shortImageId = computed(() => {
  const id = imageId.value
  return id.startsWith('sha256:') ? id.substring(7, 19) : id.substring(0, 12)
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

const handleCreateContainer = () => {
  router.push({ name: 'ContainerCreate', query: { image: imageDisplayName.value } })
}

const handleDelete = () => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除镜像 "${imageDisplayName.value}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await imageStore.deleteImage(imageId.value)
        message.success('镜像已删除')
        router.push({ name: 'Images' })
      } catch (error: any) {
        message.error('删除失败: ' + error.message)
      }
    }
  })
}

const handleAddTag = async () => {
  if (!tagForm.value.tag) {
    message.error('请输入标签')
    return
  }

  try {
    tagLoading.value = true
    await imageStore.tagImage(imageId.value, tagForm.value.tag)
    message.success('标签添加成功')
    showTagModal.value = false
    tagForm.value.tag = ''
    await imageStore.fetchImages()
  } catch (error: any) {
    message.error('添加标签失败: ' + error.message)
  } finally {
    tagLoading.value = false
  }
}

const handleRemoveTag = async (tag: string) => {
  dialog.warning({
    title: '确认移除',
    content: `确定要移除标签 "${tag}" 吗？`,
    positiveText: '移除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await imageStore.deleteImage(tag)
        message.success('标签已移除')
        await imageStore.fetchImages()
      } catch (error: any) {
        message.error('移除标签失败: ' + error.message)
      }
    }
  })
}

onMounted(async () => {
  await imageStore.fetchImages()
  if (!image.value) {
    message.error('镜像不存在')
    router.push({ name: 'Images' })
  }
})
</script>

<style scoped>
.image-detail {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-top: 12px;
}

.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-group .subtitle {
  color: var(--n-text-color-3);
  font-weight: 400;
}

.title-group .id-tag {
  font-family: monospace;
  background: var(--n-color-modal);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--n-text-color-3);
  border: 1px solid var(--n-border-color);
}

.subtitle-text {
  margin-top: 8px;
  color: var(--n-text-color-3);
}

.header-actions {
  display: flex;
  gap: 12px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
