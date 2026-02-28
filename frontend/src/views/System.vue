<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { systemApi } from '../plugins/api'
import { useContainerStore } from '../plugins/stores/containers'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'
import {
  TrashOutline,
  PersonAddOutline,
  CreateOutline,
  ServerOutline
} from '@vicons/ionicons5'

const message = useMessage()
const dialog = useDialog()
const containerStore = useContainerStore()
const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()

const loading = ref(false)
const pruneLoading = ref(false)
const pruneResult = ref<Record<string, number>>({})

const systemInfo = ref<any>(null)
const hostInfo = ref<any>(null)

// 用户管理
const users = ref([
  { id: 1, username: 'admin', role: '管理员', status: 'active', createdAt: '2024-01-01' },
])
const showUserModal = ref(false)
const userForm = ref({ id: null as number | null, username: '', password: '', role: 'user' })
const isEditing = ref(false)

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const loadSystemInfo = async () => {
  try {
    systemInfo.value = await systemApi.getSystemInfo()
    hostInfo.value = {
      hostname: systemInfo.value.name || 'Unknown',
      os: systemInfo.value.operating_system || 'Linux',
      architecture: systemInfo.value.architecture || 'Unknown',
      kernel: systemInfo.value.kernel_version || 'Unknown',
      cpus: systemInfo.value.n_cpu || 0,
      memory: systemInfo.value.mem_total || 0,
      osType: systemInfo.value.os_type || 'Unknown'
    }
  } catch (error: any) {
    message.error('加载系统信息失败: ' + error.message)
  }
}

// 用户管理
const handleAddUser = () => {
  isEditing.value = false
  userForm.value = { id: null, username: '', password: '', role: 'user' }
  showUserModal.value = true
}

const handleEditUser = (user: any) => {
  isEditing.value = true
  userForm.value = { ...user, password: '' }
  showUserModal.value = true
}

const handleSaveUser = () => {
  if (!userForm.value.username) {
    message.error('请输入用户名')
    return
  }
  if (!isEditing.value && !userForm.value.password) {
    message.error('请输入密码')
    return
  }
  
  if (isEditing.value) {
    const index = users.value.findIndex(u => u.id === userForm.value.id)
    if (index !== -1) {
      users.value[index] = { ...users.value[index], ...userForm.value }
    }
    message.success('用户更新成功')
  } else {
    users.value.push({
      id: Date.now(),
      username: userForm.value.username,
      role: userForm.value.role === 'admin' ? '管理员' : '普通用户',
      status: 'active',
      createdAt: new Date().toISOString().split('T')[0]
    })
    message.success('用户创建成功')
  }
  showUserModal.value = false
}

const handleDeleteUser = (user: any) => {
  // 检查是否是最后一个管理员
  if (user.role === '管理员') {
    const adminCount = users.value.filter(u => u.role === '管理员').length
    if (adminCount <= 1) {
      message.error('系统必须保留至少一个管理员')
      return
    }
  }
  
  dialog.warning({
    title: '确认删除',
    content: `确定要删除用户 "${user.username}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      users.value = users.value.filter(u => u.id !== user.id)
      message.success('用户已删除')
    }
  })
}

// 资源清理
const handlePrune = (type: 'containers' | 'images' | 'volumes' | 'networks' | 'logs' | 'cache' | 'all') => {
  const typeName = {
    containers: '容器',
    images: '镜像',
    volumes: '存储卷',
    networks: '网络',
    logs: '日志',
    cache: '缓存',
    all: '全部资源'
  }[type]

  const pruneMethod = {
    containers: systemApi.pruneContainers,
    images: systemApi.pruneImages,
    volumes: systemApi.pruneVolumes,
    networks: systemApi.pruneNetworks,
    logs: async () => ({ space_reclaimed: 0, message: '日志清理功能开发中' }),
    cache: async () => ({ space_reclaimed: 0, message: '缓存清理功能开发中' }),
    all: systemApi.pruneAll
  }[type]

  dialog.warning({
    title: '确认清理',
    content: `确定要清理${typeName}吗？此操作不可逆。`,
    positiveText: '清理',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        pruneLoading.value = true
        pruneResult.value = {}
        const result = await pruneMethod()
        pruneResult.value[type] = result.space_reclaimed
        if (result.details) {
          pruneResult.value = { ...pruneResult.value, ...result.details }
        }
        const spaceMB = result.space_reclaimed / 1024 / 1024
        const spaceGB = spaceMB / 1024
        const spaceText = spaceGB >= 1 ? `${spaceGB.toFixed(2)} GB` : `${spaceMB.toFixed(2)} MB`
        message.success(`${result.message}，释放空间: ${spaceText}`)

        await Promise.all([
          containerStore.fetchContainers(),
          imageStore.fetchImages(),
          volumeStore.fetchVolumes(),
          networkStore.fetchNetworks()
        ])
      } catch (error: any) {
        message.error('清理失败: ' + error.message)
      } finally {
        pruneLoading.value = false
      }
    }
  })
}

onMounted(() => {
  loadSystemInfo()
})
</script>

<template>
  <div class="system-page">
    <!-- 顶部4张统计卡片 -->
    <n-grid x-gap="16" y-gap="16" cols="2 s:2 m:4" responsive="screen" style="margin-top: 10px;">
      <n-gi>
        <n-card class="stat-card" :bordered="false">
          <div class="stat-info">
            <div class="stat-label">CPU</div>
            <div class="stat-value">{{ hostInfo?.cpus || 0 }} 核</div>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card" :bordered="false">
          <div class="stat-info">
            <div class="stat-label">内存</div>
            <div class="stat-value">{{ formatSize(hostInfo?.memory || 0) }}</div>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card" :bordered="false">
          <div class="stat-info">
            <div class="stat-label">架构</div>
            <div class="stat-value">{{ hostInfo?.architecture || 'Unknown' }}</div>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card" :bordered="false">
          <div class="stat-info">
            <div class="stat-label">Docker版本</div>
            <div class="stat-value">{{ systemInfo?.server_version || 'Unknown' }}</div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 用户管理卡片 -->
    <n-card title="用户管理" :bordered="false" style="margin-top: 20px;">
      <template #header-extra>
        <n-button type="primary" size="small" @click="handleAddUser">
          <template #icon>
            <n-icon :component="PersonAddOutline" />
          </template>
          新增用户
        </n-button>
      </template>
      <n-list>
        <n-list-item v-for="user in users" :key="user.id">
          <n-thing :title="user.username">
            <template #avatar>
              <n-avatar :size="38">{{ user.username[0].toUpperCase() }}</n-avatar>
            </template>
            <template #header-extra>
              <n-space align="center">
                <n-button size="small" quaternary @click="handleEditUser(user)">
                  <template #icon>
                    <n-icon :component="CreateOutline" />
                  </template>
                </n-button>
                <n-button size="small" quaternary type="error" @click="handleDeleteUser(user)">
                  <template #icon>
                    <n-icon :component="TrashOutline" />
                  </template>
                </n-button>
              </n-space>
            </template>
            <template #description>
              <n-space align="center" :wrap="false">
                <n-tag :type="user.role === '管理员' ? 'error' : 'default'" size="small">{{ user.role }}</n-tag>
                <n-tag :type="user.status === 'active' ? 'success' : 'warning'" size="small">
                  {{ user.status === 'active' ? '正常' : '禁用' }}
                </n-tag>
              </n-space>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-card>

    <!-- 资源清理卡片 -->
    <n-card title="资源清理" :bordered="false" style="margin-top: 20px;">
      <n-grid x-gap="16" y-gap="16" cols="2 s:3 m:6" responsive="screen">
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('containers')" :loading="pruneLoading">
            清理容器
          </n-button>
        </n-gi>
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('images')" :loading="pruneLoading">
            清理镜像
          </n-button>
        </n-gi>
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('volumes')" :loading="pruneLoading">
            清理存储
          </n-button>
        </n-gi>
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('networks')" :loading="pruneLoading">
            清理网络
          </n-button>
        </n-gi>
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('logs')" :loading="pruneLoading">
            清理日志
          </n-button>
        </n-gi>
        <n-gi>
          <n-button class="prune-btn" size="large" @click="handlePrune('cache')" :loading="pruneLoading">
            清理缓存
          </n-button>
        </n-gi>
      </n-grid>
    </n-card>

    <!-- 用户编辑/新增弹窗 -->
    <n-modal v-model:show="showUserModal" :title="isEditing ? '编辑用户' : '新增用户'" preset="card" style="width: 400px;">
      <n-form :model="userForm" label-placement="left" label-width="80px">
        <n-form-item label="用户名">
          <n-input v-model:value="userForm.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码">
          <n-input v-model:value="userForm.password" type="password" placeholder="请输入密码" />
        </n-form-item>
        <n-form-item label="角色">
          <n-select v-model:value="userForm.role" :options="[
            { label: '管理员', value: 'admin' },
            { label: '普通用户', value: 'user' }
          ]" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showUserModal = false">取消</n-button>
          <n-button type="primary" @click="handleSaveUser">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.system-page {
  padding: 0 10px 10px 10px;
  max-width: 1400px;
  margin: 0 auto;
}

.view-header {
  margin-bottom: 0;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
}

.stat-card {
  transition: all 0.3s ease;
  cursor: default;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-info {
  padding: 8px 0;
}

.stat-label {
  font-size: 14px;
  color: var(--n-text-color-3);
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--n-text-color-1);
}

.prune-btn {
  width: 100%;
  height: 80px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.create-time {
  font-size: 12px;
  color: #999;
  margin-left: 8px;
}
</style>
