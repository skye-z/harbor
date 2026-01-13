<template>
  <div class="images">
    <div class="view-header">
      <div class="header-content">
        <h1>镜像列表</h1>
        <div class="subtitle">管理本地容器镜像和仓库</div>
      </div>
      <div class="header-actions">
        <n-input
          v-model:value="searchText"
          placeholder="搜索镜像名称或ID..."
          clearable
          style="width: 300px"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
        <n-button type="primary" @click="showPullModal = true">
          <template #icon>
            <n-icon :component="CloudDownloadOutline" />
          </template>
          拉取镜像
        </n-button>
        <n-button type="info" @click="showBuildModal = true">
          <template #icon>
            <n-icon :component="BuildOutline" />
          </template>
          构建镜像
        </n-button>
        <n-button quaternary circle @click="refresh" :loading="loading">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
    </div>

    <n-card :bordered="true" class="table-card">
      <n-data-table
        :columns="columns"
        :data="filteredImages"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
      />
    </n-card>

    <n-modal v-model:show="showPullModal" preset="card" title="拉取镜像" style="width: 500px">
      <n-form :model="pullForm" label-placement="left" label-width="100px">
        <n-form-item label="镜像名称">
          <n-input v-model:value="pullForm.image" placeholder="nginx:latest" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showPullModal = false">取消</n-button>
          <n-button type="primary" @click="handlePull" :loading="pulling">
            拉取
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal v-model:show="showBuildModal" preset="card" title="构建镜像" style="width: 600px">
      <n-form :model="buildForm" label-placement="left" label-width="100px">
        <n-form-item label="镜像名称">
          <n-input v-model:value="buildForm.imageName" placeholder="myapp:latest" />
        </n-form-item>
        <n-form-item label="Dockerfile">
          <n-input
            v-model:value="buildForm.dockerfile"
            type="textarea"
            :autosize="{ minRows: 6, maxRows: 15 }"
            placeholder="FROM nginx:latest"
          />
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

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useImageStore } from '../plugins/stores/images'
import { useMessage, useDialog, NTag, NButton, NSpace, NText, NTime, NIcon } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  SearchOutline,
  CloudDownloadOutline,
  BuildOutline,
  RefreshOutline,
  TrashOutline,
  SettingsOutline,
  CubeOutline
} from '@vicons/ionicons5'

const router = useRouter()
const imageStore = useImageStore()
const message = useMessage()
const dialog = useDialog()

const searchText = ref('')
const showPullModal = ref(false)
const showBuildModal = ref(false)
const loading = ref(false)
const pulling = ref(false)

const pullForm = ref({
  image: ''
})

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

const columns: DataTableColumns<any> = [
  {
    title: '名称与标签',
    key: 'repo_tags',
    render(row) {
      return h('div', { class: 'image-name-cell' }, [
        h('div', { class: 'image-icon' }, h(NIcon, { component: CubeOutline, size: 18 })),
        h('div', { class: 'image-details' }, [
          h('div', { style: 'font-weight: 600;' },
            row.repo_tags && row.repo_tags.length > 0
              ? row.repo_tags[0].split(':')[0]
              : (row.id.startsWith('sha256:') ? row.id.substring(7, 19) : row.id.substring(0, 12))
          ),
          h('div', { style: 'display: flex; gap: 4px; flex-wrap: wrap;' },
            (row.repo_tags && row.repo_tags.length > 0)
              ? row.repo_tags.slice(0, 3).map((tag: string) => {
                  const version = tag.split(':')[1] || 'latest'
                  return h(NTag, { size: 'small', bordered: true, style: 'font-family: monospace; font-size: 12px' }, { default: () => version })
                })
              : [h(NTag, { size: 'small', type: 'warning', bordered: true }, { default: () => 'none' })]
          )
        ])
      ])
    }
  },
  {
    title: '镜像ID',
    key: 'id',
    render(row) {
      return h(NText, { depth: 3, style: 'font-family: monospace; font-size: 12px' }, { default: () => row.id.startsWith('sha256:') ? row.id.substring(7, 19) : row.id.substring(0, 12) })
    }
  },
  {
    title: '创建时间',
    key: 'created',
    render(row) {
      return h(NTime, { time: new Date(row.created * 1000), type: 'relative' })
    }
  },
  {
    title: '大小',
    key: 'size',
    render(row) {
      return h(NText, { depth: 3 }, { default: () => formatSize(row.size) })
    }
  },
  {
    title: '操作',
    key: 'actions',
    fixed: 'right',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            quaternary: true,
            circle: true,
            type: 'info',
            onClick: () => viewDetails(row.id)
          }, { icon: () => h(NIcon, { component: SettingsOutline }) }),
          h(NButton, {
            size: 'small',
            quaternary: true,
            circle: true,
            type: 'error',
            onClick: () => handleDelete(row.id)
          }, { icon: () => h(NIcon, { component: TrashOutline }) })
        ]
      })
    }
  }
]

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

const handlePull = async () => {
  if (!pullForm.value.image) {
    message.error('请输入镜像名称')
    return
  }

  try {
    pulling.value = true
    await imageStore.pullImage(pullForm.value.image)
    showPullModal.value = false
    pullForm.value.image = ''
    message.success('开始拉取镜像')
    refresh()
  } catch (error: any) {
    message.error('拉取失败: ' + error.message)
  } finally {
    pulling.value = false
  }
}

const handleBuild = async () => {
  if (!buildForm.value.imageName) {
    message.error('请输入镜像名称')
    return
  }

  try {
    loading.value = true
    await imageStore.buildImage(buildForm.value.imageName, buildForm.value.dockerfile)
    showBuildModal.value = false
    message.success('开始构建镜像')
    refresh()
  } catch (error: any) {
    message.error('构建失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const viewDetails = (id: string) => {
  router.push({ name: 'ImageDetail', params: { id } })
}

const handleDelete = (id: string) => {
  const image = imageStore.getImageById(id)
  dialog.warning({
    title: '确认删除',
    content: `确定要删除镜像 "${image?.repo_tags?.[0] || image?.id}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
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
  await refresh()
})
</script>

<style scoped>
.images {
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
}

.subtitle {
  color: var(--n-text-color-3);
  font-size: 14px;
  margin-top: 4px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.table-card {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.image-name-cell) {
  display: flex;
  align-items: center;
  gap: 12px;
}

:deep(.image-icon) {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  background-color: var(--n-color-modal);
  border: 1px solid var(--n-border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--n-text-color-3);
}

:deep(.image-details) {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
</style>
