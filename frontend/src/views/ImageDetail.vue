<template>
  <div class="image-detail">
    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ shortImageId }}</h1>
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
        <div class="subtitle-text">
          <div style="margin-bottom: 10px;" v-if="image?.Config?.Labels?.['org.opencontainers.image.description']">
            {{ image.Config.Labels['org.opencontainers.image.description'] }}
          </div>
          <n-tag v-for="tag in image?.RepoTags" :key="tag" style="margin-right: 10px;" type="info" closable
            @close="handleRemoveTag(tag)">
            {{ tag }}
          </n-tag>
          <n-button size="small" type="primary" @click="showTagModal = true">
            添加标签
          </n-button>
        </div>
      </div>
    </div>

    <n-spin :show="loading">
      <n-grid item-responsive x-gap="10" cols="24">
        <n-gi span="24 760:9 900:8">
          <n-card size="small" title="基本信息" style="margin-bottom: 10px;">
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item label="名称">
                {{ image?.Config?.Labels?.['org.opencontainers.image.title'] || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="厂商">
                {{ image?.Config?.Labels?.['org.opencontainers.image.vendor'] || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="版本">
                {{ image?.Config?.Labels?.['org.opencontainers.image.version'] || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="架构">
                {{ image?.Architecture }}/{{ image?.Os }}
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">
                <n-time v-if="image?.Created" :time="new Date(image.Created)" type="datetime" />
              </n-descriptions-item>
              <n-descriptions-item label="大小">
                {{ formatSize(image?.Size || 0) }}
              </n-descriptions-item>
              <n-descriptions-item label="工作目录">
                {{ image?.Config?.WorkingDir || '-' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>

          <n-card size="small" title="预设端口">
            <n-space v-if="image?.Config?.ExposedPorts && Object.keys(image?.Config?.ExposedPorts).length" wrap>
              <n-tag v-for="(port, key) in image?.Config?.ExposedPorts" :key="key" size="small" type="warning">
                {{ key }}
              </n-tag>
            </n-space>
            <n-empty v-else description="暂无暴露端口" />
          </n-card>
        </n-gi>
        <n-gi span="24 760:15 900:16">

          <n-card size="small" title="更多信息" style="margin-bottom: 10px;">
            <n-descriptions :column="2" label-placement="left">
              <n-descriptions-item label="许可证">
                {{ image?.Config?.Labels?.['org.opencontainers.image.licenses'] || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="文档">
                <n-button v-if="image?.Config?.Labels?.['org.opencontainers.image.documentation']" text type="primary"
                  @click="openLink(image.Config.Labels['org.opencontainers.image.documentation'])">
                  查看
                </n-button>
                <span v-else>-</span>
              </n-descriptions-item>
              <n-descriptions-item label="源码">
                <n-button v-if="image?.Config?.Labels?.['org.opencontainers.image.source']" text type="primary"
                  @click="openLink(image.Config.Labels['org.opencontainers.image.source'])">
                  查看
                </n-button>
                <span v-else>-</span>
              </n-descriptions-item>
              <n-descriptions-item label="官网">
                <n-button v-if="image?.Config?.Labels?.['org.opencontainers.image.url']" text type="primary"
                  @click="openLink(image.Config.Labels['org.opencontainers.image.url'])">
                  查看
                </n-button>
                <span v-else>-</span>
              </n-descriptions-item>
            </n-descriptions>
          </n-card>


          <n-card size="small" title="环境变量" style="margin-bottom: 10px;">
            <n-space v-if="image?.Config?.Env?.length" wrap>
              <n-tag v-for="(env, index) in image?.Config?.Env" :key="index" size="small">
                {{ env }}
              </n-tag>
            </n-space>
            <n-empty v-else description="暂无环境变量" />
          </n-card>
          <n-card size="small" title="默认命令">
            <n-code :code="image?.Config?.Cmd?.join(' ') || '-'" language="bash" />
          </n-card>
          <n-card size="small" style="margin-top: 10px;">
            <template #header>
              <span>使用此镜像的容器 ({{ usingContainers.length }})</span>
            </template>
            <n-empty v-if="usingContainers.length === 0" description="暂无容器使用此镜像" />
            <n-list v-else size="small">
              <n-list-item v-for="c in usingContainers" :key="c.id" style="cursor: pointer"
                @click="viewContainer(c.id)">
                <n-thing :title="getContainerName(c)">
                  <template #avatar>
                    <n-tag :type="getStateType(c.state)" size="small" round>
                      {{ formatState(c.state) }}
                    </n-tag>
                  </template>
                  <template #header-extra>
                    <n-tag size="small" :bordered="false">{{ formatId(c.id) }}</n-tag>
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>

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
import { useContainerStore } from '../plugins/stores/containers'
import { imageApi } from '../plugins/api'
import { useMessage, useDialog } from 'naive-ui'
import {
  PlayCircleOutline,
  TrashOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const imageStore = useImageStore()
const containerStore = useContainerStore()
const message = useMessage()
const dialog = useDialog()

const showTagModal = ref(false)
const tagForm = ref({ tag: '' })
const tagLoading = ref(false)
const imageDetail = ref<any>(null)
const loading = ref(true)

const imageId = computed(() => route.params.id as string)

const image = computed(() => imageDetail.value)

const stateMap: Record<string, string> = {
  running: '运行',
  exited: '停止',
  created: '创建',
  paused: '暂停',
  restarting: '重启',
  removing: '移除',
  dead: '异常'
}

const formatState = (state: string) => {
  return stateMap[state] || state || '未知'
}

const getStateType = (state: string) => {
  if (state === 'running') return 'success'
  if (state === 'paused') return 'warning'
  return 'default'
}

const usingContainers = computed(() => {
  if (!image.value) return []
  const imgId = imageId.value
  return containerStore.containers.filter((c: any) => {
    const fullId = image.value?.id || ''
    const shortId = fullId.startsWith('sha256:') ? fullId.substring(7) : fullId
    return c.image_id === imgId || c.image_id === fullId ||
      c.image_id?.startsWith(shortId) ||
      c.image?.includes(imgId) || c.image?.includes(fullId)
  })
})

const formatId = (id: string) => {
  return id?.substring(0, 12) || ''
}

const getContainerName = (container: any) => {
  return container.names?.[0]?.replace(/^\//, '') || formatId(container.id)
}

const viewContainer = (id: string) => {
  router.push({ name: 'ContainerDetail', params: { id } })
}

const imageDisplayName = computed(() => {
  const img = image.value
  return img?.RepoTags?.[0] || img?.Id || 'Unknown'
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
  const id = image.value?.Id || imageId.value || ''
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
  router.push({ name: 'ContainerCreate', query: { image: imageDisplayName.value, imageId: imageId.value } })
}

const openLink = (url: string) => {
  if (url) {
    window.open(url, '_blank')
  }
}

const handleDelete = () => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除镜像 "${imageDisplayName.value}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await imageApi.delete(imageId.value)
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
    await imageApi.tag(imageId.value, tagForm.value.tag)
    message.success('标签添加成功')
    showTagModal.value = false
    tagForm.value.tag = ''
    const res = await imageApi.get(imageId.value)
    imageDetail.value = res
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
        await imageApi.delete(tag)
        message.success('标签已移除')
        const res = await imageApi.get(imageId.value)
        imageDetail.value = res
      } catch (error: any) {
        message.error('移除标签失败: ' + error.message)
      }
    }
  })
}

onMounted(async () => {
  loading.value = true
  try {
    imageDetail.value = await imageApi.get(imageId.value)
    await containerStore.fetchContainers()
  } catch (error: any) {
    message.error('获取镜像详情失败: ' + error.message)
    router.push({ name: 'Images' })
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.image-detail {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 10px;
}

.view-header {
  justify-content: space-between;
  align-items: flex-start;
  display: flex;
}

.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  align-items: center;
  display: flex;
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
