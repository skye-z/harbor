<template>
  <div class="container-create">
    <div class="page-header">
      <div class="view-header">
        <div class="title-group">
          <h1>创建容器</h1>
        </div>
      </div>
    </div>

    <n-form ref="formRef" :model="formData" :rules="formRules" :disabled="imageLoading || creating"
      label-placement="top">
      <n-grid x-gap="20" :cols="24">
        <n-gi :span="16">
          <n-card title="基本配置" class="config-card">
            <n-form-item label="镜像" path="image">
              <n-select v-model:value="formData.image" :options="imageOptions" :disabled="preselectedImage"
                placeholder="选择镜像" filterable clearable />
            </n-form-item>
            <n-form-item label="容器名称" path="name">
              <n-input v-model:value="formData.name" placeholder="输入容器名称" />
            </n-form-item>
          </n-card>

          <n-card title="环境变量" class="config-card">
            <n-form-item v-for="(env, index) in formData.envs" :key="index">
              <n-space style="width: 100%">
                <n-input v-model:value="env.key" placeholder="变量名" style="flex: 1" />
                <n-input v-model:value="env.value" placeholder="变量值" style="flex: 1" />
                <n-button type="error" ghost circle @click="removeEnv(index)">
                  <template #icon>
                    <n-icon :component="TrashOutline" />
                  </template>
                </n-button>
              </n-space>
            </n-form-item>
            <n-button type="primary" ghost @click="addEnv">
              <template #icon>
                <n-icon :component="AddOutline" />
              </template>
              添加环境变量
            </n-button>
          </n-card>

          <n-card title="存储配置" class="config-card">
            <n-form-item :show-feedback="false">
              <n-space style="width: 100%" vertical>
                <div v-for="(mount, index) in formData.volumes" :key="index" class="mount-item">
                  <n-space>
                    <n-select v-model:value="mount.type" :options="mountTypeOptions" placeholder="类型"
                      style="width: 100px" @update:value="() => mount.volumeId = ''" />
                    <template v-if="mount.type === 'volume'">
                      <n-select v-model:value="mount.volumeId" :options="volumeOptions" placeholder="选择卷"
                        style="width: 200px" />
                    </template>
                    <template v-else>
                      <n-input v-model:value="mount.hostPath" placeholder="宿主机路径" style="width: 200px" />
                    </template>
                    <n-input v-model:value="mount.containerPath" placeholder="容器内路径" />
                    <n-select v-model:value="mount.mode" :options="mountModeOptions" placeholder="模式"
                      style="width: 80px" />
                    <n-button type="error" ghost circle @click="removeVolume(index)">
                      <template #icon>
                        <n-icon :component="TrashOutline" />
                      </template>
                    </n-button>
                  </n-space>
                </div>
            <n-button type="primary" ghost @click="addVolume">
              <template #icon>
                <n-icon :component="AddOutline" />
              </template>
              添加挂载
            </n-button>
              </n-space>
            </n-form-item>
          </n-card>

          <n-card title="网络配置" class="config-card">
            <n-form-item label="网络模式">
              <n-select v-model:value="formData.networkMode" :options="networkModeOptions" placeholder="选择网络模式" />
            </n-form-item>
            <n-form-item label="端口映射" :show-feedback="false">
              <n-space style="width: 100%" vertical>
                <div v-for="(port, index) in formData.portBindings" :key="index" class="port-item">
                  <n-space align="center">
                    <n-input-number v-model:value="port.hostPort" :min="1" :max="65535" placeholder="1-65535"
                      style="width: 150px" :show-button="false" />
                    <span style="font-weight: bold; color: var(--n-text-color-3);">:</span>
                    <n-input-number v-model:value="port.containerPort" :min="1" :max="65535" placeholder="1-65535"
                      style="width: 150px" :show-button="false" />
                    <n-select v-model:value="port.protocol" :options="protocolOptions" placeholder="协议"
                      style="width: 100px" />
                    <n-button type="error" ghost circle @click="removePort(index)">
                      <template #icon>
                        <n-icon :component="TrashOutline" />
                      </template>
                    </n-button>
                  </n-space>
                </div>
                <n-button type="primary" ghost @click="addPort">
                  <template #icon>
                    <n-icon :component="AddOutline" />
                  </template>
                  添加端口映射
                </n-button>
              </n-space>
            </n-form-item>
          </n-card>

          <n-card class="config-card">
            <n-collapse>
              <n-collapse-item title="更多配置" name="more">
                <n-grid :cols="2" x-gap="20">
                  <n-gi>
                    <n-form-item label="内存限制 (MB)">
                      <n-input-number v-model:value="formData.memory" :min="0" placeholder="例如: 512"
                        style="width: 100%" :show-button="false" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="CPU 权重">
                      <n-input-number v-model:value="formData.cpuShares" :min="0" :max="1024" placeholder="例如: 512"
                        style="width: 100%" :show-button="false" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="重启策略">
                      <n-select v-model:value="formData.restartPolicy" :options="restartPolicyOptions"
                        placeholder="选择重启策略" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="工作目录">
                      <n-input v-model:value="formData.workingDir" placeholder="容器内工作目录" />
                    </n-form-item>
                  </n-gi>
                </n-grid>
                <n-form-item label="命令">
                  <n-input v-model:value="formData.command" placeholder="默认命令 (可选)" />
                </n-form-item>
                <n-form-item label="自动启动">
                  <n-switch v-model:value="formData.autoStart" />
                </n-form-item>
              </n-collapse-item>
            </n-collapse>
          </n-card>
        </n-gi>

        <n-gi :span="8">
          <n-card title="操作" class="action-card">
            <n-space vertical style="width: 100%">
              <n-button type="primary" block size="large" :loading="creating" @click="handleCreate">
                创建容器
              </n-button>
              <n-button block :loading="imageLoading" @click="handleReset">
                重置表单
              </n-button>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import {
  TrashOutline,
  AddOutline
} from '@vicons/ionicons5'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'
import { useContainerStore } from '../plugins/stores/containers'
import { containerApi, imageApi } from '../plugins/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()
const containerStore = useContainerStore()

const formRef = ref()
const creating = ref(false)
const imageLoading = ref(false)
const originalImageDetail = ref<any>(null)

const preselectedImage = computed(() => !!route.query.image)

interface EnvVar {
  key: string
  value: string
}

interface VolumeMount {
  type: 'volume' | 'bind'
  volumeId: string
  hostPath: string
  containerPath: string
  mode: string
}

interface PortBinding {
  hostPort: number | null
  containerPort: number | null
  protocol: string
}

const formData = reactive({
  image: null as string | null,
  name: '',
  envs: [] as EnvVar[],
  volumes: [] as VolumeMount[],
  networkMode: 'bridge',
  portBindings: [] as PortBinding[],
  memory: null as number | null,
  cpuShares: null as number | null,
  restartPolicy: 'unless-stopped',
  workingDir: '',
  command: '',
  autoStart: true
})

const formRules = {
  image: {
    required: true,
    message: '请选择镜像',
    trigger: 'change'
  },
  name: {
    required: true,
    message: '请输入容器名称',
    trigger: 'blur'
  }
}

const imageOptions = computed(() =>
  imageStore.images.map(img => ({
    label: img.repo_tags?.[0] || img.id,
    value: img.repo_tags?.[0] || img.id
  }))
)

const volumeOptions = computed(() =>
  volumeStore.volumes.map(v => ({
    label: v.name,
    value: v.id
  }))
)

const mountModeOptions = [
  { label: '读写', value: 'rw' },
  { label: '只读', value: 'ro' }
]

const mountTypeOptions = [
  { label: '存储卷', value: 'volume' },
  { label: '路径', value: 'bind' }
]

const networkModeOptions = [
  { label: 'Bridge', value: 'bridge' },
  { label: 'Host', value: 'host' },
  { label: 'None', value: 'none' },
  { label: 'Container', value: 'container' }
]

const protocolOptions = [
  { label: 'TCP', value: 'tcp' },
  { label: 'UDP', value: 'udp' }
]

const restartPolicyOptions = [
  { label: '始终重启', value: 'always' },
  { label: '除非手动停止', value: 'unless-stopped' },
  { label: '失败时重启', value: 'on-failure' },
  { label: '不重启', value: 'no' }
]

const addEnv = () => {
  formData.envs.push({ key: '', value: '' })
}

const removeEnv = (index: number) => {
  formData.envs.splice(index, 1)
}

const addVolume = () => {
  formData.volumes.push({ type: 'volume', volumeId: '', hostPath: '', containerPath: '', mode: 'rw' })
}

const removeVolume = (index: number) => {
  formData.volumes.splice(index, 1)
}

const addPort = () => {
  formData.portBindings.push({ hostPort: null, containerPort: null, protocol: 'tcp' })
}

const removePort = (index: number) => {
  formData.portBindings.splice(index, 1)
}

const fillFromImageDetail = (imageDetail: any) => {
  if (!imageDetail) return

  originalImageDetail.value = JSON.parse(JSON.stringify(imageDetail))

  if (imageDetail.Config?.Env) {
    formData.envs = imageDetail.Config.Env.map((env: string) => {
      const [key, ...valueParts] = env.split('=')
      return {
        key,
        value: valueParts.join('=')
      }
    })
  }

  if (imageDetail.Config?.ExposedPorts) {
    formData.portBindings = Object.keys(imageDetail.Config.ExposedPorts).map(port => {
      const [containerPort, protocol] = port.split('/')
      return {
        hostPort: null,
        containerPort: parseInt(containerPort),
        protocol: protocol.toLowerCase()
      }
    })
  }

  if (imageDetail.Config?.WorkingDir) {
    formData.workingDir = imageDetail.Config.WorkingDir
  }

  if (imageDetail.Config?.Cmd) {
    formData.command = imageDetail.Config.Cmd.join(' ')
  }
}

const handleImageChange = async (imageName: string | null) => {
  if (!imageName) {
    return
  }

  imageLoading.value = true
  try {
    const imageDetail = await imageApi.get(imageName)
    fillFromImageDetail(imageDetail)
  } catch (error: any) {
    console.error('获取镜像详情失败:', error)
  } finally {
    imageLoading.value = false
  }
}

watch(() => formData.image, handleImageChange)

const handleReset = () => {
  if (preselectedImage.value && route.query.image) {
    formData.image = route.query.image as string
  } else {
    formData.image = null
  }
  formData.name = ''
  formData.volumes = []
  formData.memory = null
  formData.cpuShares = null
  formData.restartPolicy = 'unless-stopped'
  formData.autoStart = true

  if (originalImageDetail.value) {
    fillFromImageDetail(originalImageDetail.value)
  } else {
    formData.envs = []
    formData.portBindings = []
    formData.workingDir = ''
    formData.command = ''
  }
}

const handleCreate = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  dialog.warning({
    title: '确认创建容器',
    content: `确定要创建容器 "${formData.name}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await doCreateContainer()
    }
  })
}

const doCreateContainer = async () => {
  creating.value = true
  try {
    const envs: Record<string, string> = {}
    formData.envs.forEach(e => {
      if (e.key && e.value) {
        envs[e.key] = e.value
      }
    })

    const exposedPorts: Record<string, Record<string, string>> = {}
    const portBindings: Record<string, Array<{ HostPort: string }>> = {}
    formData.portBindings.forEach(p => {
      if (p.hostPort && p.containerPort) {
        const portKey = `${p.containerPort}/${p.protocol.toUpperCase()}`
        exposedPorts[portKey] = {}
        portBindings[portKey] = [{ HostPort: String(p.hostPort) }]
      }
    })

    const hostConfig: any = {
      NetworkMode: formData.networkMode,
      RestartPolicy: {
        Name: formData.restartPolicy
      }
    }

    if (formData.memory && formData.memory > 0) {
      hostConfig.Memory = formData.memory * 1024 * 1024
    }

    if (formData.cpuShares && formData.cpuShares > 0) {
      hostConfig.CPUShares = formData.cpuShares
    }

    if (Object.keys(portBindings).length > 0) {
      hostConfig.PortBindings = portBindings
    }

    const binds: string[] = []
    formData.volumes.forEach(m => {
      if (m.containerPath) {
        if (m.type === 'volume' && m.volumeId && m.containerPath) {
          const volume = volumeStore.volumes.find(v => v.id === m.volumeId)
          if (volume) {
            binds.push(`${volume.name}:${m.containerPath}:${m.mode}`)
          }
        } else if (m.type === 'bind' && m.hostPath && m.containerPath) {
          binds.push(`${m.hostPath}:${m.containerPath}:${m.mode}`)
        }
      }
    })

    if (binds.length > 0) {
      hostConfig.Binds = binds
    }

    // 构建符合 Docker API 格式的配置
    const dockerConfig: any = {
      Image: formData.image,
      name: formData.name,
      Env: envs,
      Cmd: formData.command ? formData.command.split(' ') : undefined,
      WorkingDir: formData.workingDir || undefined,
      Tty: true,
      OpenStdin: true,
      AutoRemove: !formData.autoStart
    }

    // 添加 ExposedPorts
    if (Object.keys(exposedPorts).length > 0) {
      dockerConfig.ExposedPorts = exposedPorts
    }

    // 构建 HostConfig，使用 Docker API 标准字段名
    const dockerHostConfig: any = {
      NetworkMode: formData.networkMode,
      RestartPolicy: {
        Name: formData.restartPolicy
      },
      AutoRemove: !formData.autoStart
    }

    if (formData.memory && formData.memory > 0) {
      dockerHostConfig.Memory = formData.memory * 1024 * 1024
    }

    if (formData.cpuShares && formData.cpuShares > 0) {
      dockerHostConfig.CpuShares = formData.cpuShares
    }

    if (Object.keys(portBindings).length > 0) {
      dockerHostConfig.PortBindings = portBindings
    }

    if (binds.length > 0) {
      dockerHostConfig.Binds = binds
    }

    dockerConfig.HostConfig = dockerHostConfig

    await containerApi.create(dockerConfig)
    message.success('容器创建成功')
    router.push({ name: 'Containers' })
  } catch (error: any) {
    message.error('创建失败: ' + error.message)
  } finally {
    creating.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    imageStore.fetchImages(),
    volumeStore.fetchVolumes(),
    networkStore.fetchNetworks()
  ])

  if (preselectedImage.value && route.query.image) {
    const imageName = route.query.image as string
    formData.image = imageName
  }
})
</script>

<style scoped>
.container-create {
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
  align-items: flex-start;
}

.title-group h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
}

.config-card {
  margin-bottom: 20px;
}

.action-card {
  position: sticky;
  top: 60px;
}

.mount-item,
.port-item {
  padding: 10px 10px 10px 0;
  background: var(--n-color-hover);
  border-radius: 6px;
}

.config-card:deep(.n-form-item.n-form-item--top-labelled) {
  grid-template-rows: none;
}
</style>
