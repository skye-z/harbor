<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { systemApi } from '../plugins/api'
import { useUserStore } from '../plugins/stores/user'
import { useContainerStore } from '../plugins/stores/containers'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'
import {
  ServerOutline,
  ShieldCheckmarkOutline,
  TrashOutline,
  TrashBinOutline,
  CheckmarkCircleOutline,
  LayersOutline,
  ConstructOutline
} from '@vicons/ionicons5'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()
const containerStore = useContainerStore()
const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()

const loading = ref(false)
const pruneLoading = ref(false)
const pruneResult = ref<Record<string, number>>({})

const systemInfo = ref<any>(null)
const hostInfo = ref<any>(null)
const versionInfo = ref({
  version: 'v1.0.0',
  edition: 'community',
  editionName: '社区版'
})

const licenseInfo = ref({
  level: 'community',
  levelName: '社区版',
  licensee: 'Community User',
  issueDate: '2024-01-01',
  expireDate: '9999-12-31',
  features: ['基础 Docker 管理', '容器管理', '镜像管理', '存储管理', '网络管理'],
  allFeatures: ['基础 Docker 管理', '容器管理', '镜像管理', '存储管理', '网络管理', '镜像编辑', '高级监控', '企业级支持']
})

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

const handlePrune = (type: 'containers' | 'images' | 'volumes' | 'networks' | 'all') => {
  const typeName = {
    containers: '容器',
    images: '镜像',
    volumes: '卷',
    networks: '网络',
    all: '全部资源'
  }[type]

  dialog.warning({
    title: '确认清理',
    content: `确定要清理未使用的${typeName}吗？此操作不可逆。`,
    positiveText: '清理',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        pruneLoading.value = true
        pruneResult.value = {}
        const result = await systemApi.prune(type)
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
    <n-space vertical :size="32">
      
      <!-- Host Info Section -->
      <section>
        <div class="section-header">
          <div class="icon-wrapper cyan">
            <n-icon :component="ServerOutline" />
          </div>
          <h2 class="section-title">宿主信息</h2>
        </div>
        
        <n-grid cols="1 s:2 m:2 l:4" responsive="screen" :x-gap="16" :y-gap="16">
          <!-- Hostname -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label">主机名</div>
              <div class="card-value">{{ hostInfo?.hostname || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- OS -->
          <n-gi span="1 s:2 m:2 l:2">
            <n-card class="data-card" :bordered="false">
              <div class="flex justify-between items-start">
                <div>
                  <div class="card-label">操作系统</div>
                  <div class="card-value">{{ hostInfo?.os || 'Loading...' }}</div>
                </div>
                <div class="text-right hidden-mobile">
                  <div class="card-label">系统类型</div>
                  <n-tag size="small" :bordered="false" type="info">{{ hostInfo?.osType || 'Unknown' }}</n-tag>
                </div>
              </div>
            </n-card>
          </n-gi>

          <!-- Kernel -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label">内核版本</div>
              <div class="card-value truncate">{{ hostInfo?.kernel || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- Architecture -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label">架构</div>
              <div class="card-value">{{ hostInfo?.architecture || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- CPU -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="flex justify-between items-center">
                <div>
                  <div class="card-label">CPU核心数</div>
                  <div class="card-value large">{{ hostInfo?.cpus || 0 }}</div>
                </div>
              </div>
            </n-card>
          </n-gi>

          <!-- Memory -->
          <n-gi span="1 s:2 m:2 l:2">
            <n-card class="data-card" :bordered="false">
              <div class="flex justify-between items-center mb-2">
                <div class="card-label">总内存</div>
                <div class="card-value">{{ formatSize(hostInfo?.memory || 0) }}</div>
              </div>
              <n-progress
                type="line"
                :percentage="100"
                :show-indicator="false"
                processing
              />
              <div class="flex justify-between mt-2 text-xs text-gray-400">
                <span>0 GB</span>
                <span>MAX</span>
              </div>
            </n-card>
          </n-gi>
        </n-grid>
      </section>

      <!-- Engine Info Section -->
      <section>
        <div class="section-header">
          <div class="icon-wrapper purple">
            <n-icon :component="LayersOutline" />
          </div>
          <h2 class="section-title">引擎信息</h2>
        </div>

        <n-grid cols="1 s:2 m:3 l:4" responsive="screen" :x-gap="16" :y-gap="16">
          <!-- Docker Version -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Docker 版本</div>
              <div class="card-value-mono">{{ systemInfo?.server_version || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- API Version -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">API 版本</div>
              <div class="card-value-mono">1.43</div>
              <div class="text-xs text-gray-400 mt-1">Minimum: 1.12</div>
            </n-card>
          </n-gi>

          <!-- Go Version -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Go 版本</div>
              <div class="card-value-mono">go1.20.5</div>
            </n-card>
          </n-gi>

          <!-- Git Commit -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Git Commit</div>
              <div class="card-value-mono">a61e2b4</div>
            </n-card>
          </n-gi>

          <!-- Build Time -->
          <n-gi span="1 s:2 m:2 l:2">
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">构建时间</div>
              <div class="card-value-mono text-base">2023-07-14T22:27:35.000000000+00:00</div>
            </n-card>
          </n-gi>

          <!-- Experimental -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="flex justify-between items-center h-full">
                <div class="flex flex-col">
                  <span class="card-label-small mb-1">实验特性</span>
                  <span class="font-semibold">{{ systemInfo?.experimental ? 'Enabled' : 'Disabled' }}</span>
                </div>
                <n-switch :value="systemInfo?.experimental" disabled size="small" />
              </div>
            </n-card>
          </n-gi>

          <!-- Storage Driver -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">存储驱动</div>
              <div class="card-value-mono">{{ systemInfo?.driver || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- Log Driver -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">日志驱动</div>
              <div class="card-value-mono">{{ systemInfo?.logging_driver || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- Cgroup Driver -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Cgroup 驱动</div>
              <div class="card-value-mono">{{ systemInfo?.cgroup_driver || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- Cgroup Version -->
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Cgroup 版本</div>
              <div class="card-value-mono">{{ systemInfo?.cgroup_version || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- Docker Root -->
          <n-gi span="1 s:2 m:2 l:2">
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">Docker 根目录</div>
              <div class="card-value-mono text-base">{{ systemInfo?.docker_root_dir || 'Loading...' }}</div>
            </n-card>
          </n-gi>

          <!-- ID -->
          <n-gi span="1 s:2 m:2 l:2">
            <n-card class="data-card" :bordered="false">
              <div class="card-label-small mb-2">ID</div>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <div class="card-value-mono text-sm truncate">{{ systemInfo?.id || 'Loading...' }}</div>
                </template>
                {{ systemInfo?.id }}
              </n-tooltip>
            </n-card>
          </n-gi>

          <!-- Status Cards Row -->
          <n-gi span="1 s:2 m:3 l:4">
             <n-grid cols="2 m:4" responsive="screen" :x-gap="16" :y-gap="16">
                <n-gi>
                   <n-card class="data-card" :bordered="false" size="small">
                      <div class="flex justify-between items-center">
                         <span class="card-label-small">调试模式</span>
                         <n-tag size="small" :type="systemInfo?.debug ? 'warning' : 'default'">{{ systemInfo?.debug ? 'Enabled' : 'Disabled' }}</n-tag>
                      </div>
                   </n-card>
                </n-gi>
                <n-gi>
                   <n-card class="data-card" :bordered="false" size="small">
                      <div class="flex justify-between items-center">
                         <span class="card-label-small">IPv4 转发</span>
                         <n-tag size="small" :type="systemInfo?.ipv4_forwarding ? 'success' : 'default'">{{ systemInfo?.ipv4_forwarding ? 'Enabled' : 'Disabled' }}</n-tag>
                      </div>
                   </n-card>
                </n-gi>
                <n-gi>
                   <n-card class="data-card" :bordered="false" size="small">
                      <div class="flex justify-between items-center">
                         <span class="card-label-small">Live Restore</span>
                         <n-tag size="small" :type="systemInfo?.live_restore_enabled ? 'success' : 'default'">{{ systemInfo?.live_restore_enabled ? 'Enabled' : 'Disabled' }}</n-tag>
                      </div>
                   </n-card>
                </n-gi>
                <n-gi>
                   <n-card class="data-card" :bordered="false" size="small">
                      <div class="flex justify-between items-center">
                         <span class="card-label-small">OOM Kill</span>
                         <n-tag size="small" :type="systemInfo?.oom_kill_disable ? 'success' : 'warning'">{{ systemInfo?.oom_kill_disable ? 'False' : 'True' }}</n-tag>
                      </div>
                   </n-card>
                </n-gi>
             </n-grid>
          </n-gi>
        </n-grid>
      </section>

      <!-- Resource Cleanup Section -->
      <section>
        <div class="section-header">
          <div class="icon-wrapper orange">
            <n-icon :component="ConstructOutline" />
          </div>
          <h2 class="section-title">资源清理</h2>
        </div>

        <n-alert type="warning" class="mb-4">
          清理操作将永久删除未使用的资源，此操作不可逆。请谨慎操作！
        </n-alert>

        <n-grid cols="1 s:2 m:4" responsive="screen" :x-gap="16" :y-gap="16">
          <n-gi>
            <n-card class="data-card" :bordered="false">
              <n-space vertical :size="12">
                <span class="card-label">清理容器</span>
                <n-text depth="3" class="text-xs">清理所有已停止的容器</n-text>
                <n-space justify="space-between" align="center">
                  <n-button type="warning" secondary size="small" @click="handlePrune('containers')" :loading="pruneLoading">
                    <template #icon>
                      <n-icon :component="TrashOutline" />
                    </template>
                    执行清理
                  </n-button>
                  <n-text v-if="pruneResult.containers" type="success" class="text-xs">
                    -{{ formatSize(pruneResult.containers) }}
                  </n-text>
                </n-space>
              </n-space>
            </n-card>
          </n-gi>

          <n-gi>
            <n-card class="data-card" :bordered="false">
              <n-space vertical :size="12">
                <span class="card-label">清理镜像</span>
                <n-text depth="3" class="text-xs">清理未使用的镜像</n-text>
                <n-space justify="space-between" align="center">
                  <n-button type="warning" secondary size="small" @click="handlePrune('images')" :loading="pruneLoading">
                    <template #icon>
                      <n-icon :component="TrashOutline" />
                    </template>
                    执行清理
                  </n-button>
                  <n-text v-if="pruneResult.images" type="success" class="text-xs">
                    -{{ formatSize(pruneResult.images) }}
                  </n-text>
                </n-space>
              </n-space>
            </n-card>
          </n-gi>

          <n-gi>
            <n-card class="data-card" :bordered="false">
              <n-space vertical :size="12">
                <span class="card-label">清理卷</span>
                <n-text depth="3" class="text-xs">清理未使用的卷</n-text>
                <n-space justify="space-between" align="center">
                  <n-button type="warning" secondary size="small" @click="handlePrune('volumes')" :loading="pruneLoading">
                    <template #icon>
                      <n-icon :component="TrashOutline" />
                    </template>
                    执行清理
                  </n-button>
                  <n-text v-if="pruneResult.volumes" type="success" class="text-xs">
                    -{{ formatSize(pruneResult.volumes) }}
                  </n-text>
                </n-space>
              </n-space>
            </n-card>
          </n-gi>

          <n-gi>
            <n-card class="data-card" :bordered="false">
              <n-space vertical :size="12">
                <span class="card-label">清理网络</span>
                <n-text depth="3" class="text-xs">清理未使用的网络</n-text>
                <n-space justify="space-between" align="center">
                  <n-button type="warning" secondary size="small" @click="handlePrune('networks')" :loading="pruneLoading">
                    <template #icon>
                      <n-icon :component="TrashOutline" />
                    </template>
                    执行清理
                  </n-button>
                  <n-text v-if="pruneResult.networks" type="success" class="text-xs">
                    -{{ formatSize(pruneResult.networks) }}
                  </n-text>
                </n-space>
              </n-space>
            </n-card>
          </n-gi>

          <n-gi span="1 s:2 m:4">
            <n-card class="data-card error-border" :bordered="false">
              <n-space justify="space-between" align="center">
                 <div class="flex flex-col">
                    <span class="card-label text-error">一键清理全部</span>
                    <n-text depth="3" class="text-xs">清理所有未使用的容器、镜像、卷和网络</n-text>
                 </div>
                 <n-space align="center">
                    <n-text v-if="pruneResult.all" type="success" class="font-bold">
                      总释放: {{ formatSize(pruneResult.all) }}
                    </n-text>
                    <n-button type="error" @click="handlePrune('all')" :loading="pruneLoading">
                      <template #icon>
                        <n-icon :component="TrashBinOutline" />
                      </template>
                      清理全部
                    </n-button>
                 </n-space>
              </n-space>
            </n-card>
          </n-gi>
        </n-grid>
      </section>

    </n-space>
  </div>
</template>

<style scoped>
.system-page {
  padding: 5px 10px 10px 10px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.icon-wrapper {
  padding: 8px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-wrapper.cyan {
  background-color: rgba(6, 182, 212, 0.1);
  color: #06b6d4;
  border: 1px solid rgba(6, 182, 212, 0.2);
}

.icon-wrapper.purple {
  background-color: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  border: 1px solid rgba(168, 85, 247, 0.2);
}

.icon-wrapper.green {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.icon-wrapper.orange {
  background-color: rgba(249, 115, 22, 0.1);
  color: #f97316;
  border: 1px solid rgba(249, 115, 22, 0.2);
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.data-card {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 12px;
}

.data-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.error-border {
  border: 1px solid rgba(239, 68, 68, 0.2);
  background-color: rgba(239, 68, 68, 0.02);
}

.card-label {
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  color: var(--n-text-color-3);
  margin-bottom: 4px;
  letter-spacing: 0.05em;
}

.text-error {
  color: #d03050;
}

.card-label-small {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.card-value {
  font-size: 18px;
  font-family: 'JetBrains Mono', monospace;
  font-weight: 500;
  color: var(--n-text-color-1);
}

.card-value.large {
  font-size: 30px;
  font-weight: 700;
}

.card-value-mono {
  font-size: 18px;
  font-family: 'JetBrains Mono', monospace;
  color: var(--n-text-color-1);
}

.truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hidden-mobile {
  display: none;
}

@media (min-width: 640px) {
  .hidden-mobile {
    display: block;
  }
}
</style>