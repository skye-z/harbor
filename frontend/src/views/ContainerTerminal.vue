<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, NCard, NSpace, NSelect, NButton, NIcon, NPageHeader, NSpin } from 'naive-ui'
import { ArrowBack, TerminalOutline, PowerOutline } from '@vicons/ionicons5'
import Terminal from '../components/Terminal.vue'

const route = useRoute()
const router = useRouter()
const containerStore = useContainerStore()
const message = useMessage()

const containerId = computed(() => route.params.id as string)
const container = computed(() => containerStore.getContainerById(containerId.value))

const selectedShell = ref('/bin/sh')
const isConnected = ref(false)
const terminalKey = ref(0)

const shellOptions = [
  { label: '/bin/sh', value: '/bin/sh' },
  { label: '/bin/bash', value: '/bin/bash' },
  { label: '/bin/zsh', value: '/bin/zsh' },
  { label: 'sh', value: 'sh' },
  { label: 'bash', value: 'bash' }
]

const handleBack = () => {
  router.push({ name: 'ContainerDetail', params: { id: containerId.value } })
}

const handleConnect = () => {
  if (!containerId.value) return
  isConnected.value = true
  terminalKey.value++
}

const handleDisconnect = () => {
  isConnected.value = false
}

onMounted(async () => {
  if (!containerStore.containers.length) {
    await containerStore.fetchContainers()
  }
  if (!container.value) {
    message.error('容器不存在')
    router.push({ name: 'Containers' })
  }
})
</script>

<template>
  <div class="terminal-page">
    <div class="page-header">
      <div class="title-group">
        <div class="view-header">
          <h1>{{ container?.names?.[0]?.replace(/^\//, '') || containerId }} - 终端</h1>
          <div class="header-actions">
            <n-button size="medium" @click="handleBack">
              <template #icon>
                <n-icon :component="ArrowBack" />
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

    <n-card class="mt-4 terminal-card" content-style="padding: 0; height: 100%; display: flex; flex-direction: column;">
      <div class="terminal-controls p-4 border-b border-gray-700 bg-gray-900 flex items-center justify-between">
        <n-space align="center">
          <span class="text-gray-300">Shell:</span>
          <n-select
            v-model:value="selectedShell"
            :options="shellOptions"
            size="small"
            style="width: 150px"
            :disabled="isConnected"
          />
          <n-button
            v-if="!isConnected"
            type="primary"
            size="small"
            @click="handleConnect"
          >
            <template #icon>
              <n-icon :component="TerminalOutline" />
            </template>
            连接
          </n-button>
          <n-button
            v-else
            type="error"
            size="small"
            @click="handleDisconnect"
          >
            <template #icon>
              <n-icon :component="PowerOutline" />
            </template>
            断开
          </n-button>
        </n-space>

        <div class="status-indicator flex items-center gap-2">
          <div
            class="w-2 h-2 rounded-full transition-colors duration-300"
            :class="isConnected ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]' : 'bg-red-500'"
          ></div>
          <span class="text-xs text-gray-400">
            {{ isConnected ? '已连接' : '未连接' }}
          </span>
        </div>
      </div>

      <div class="terminal-wrapper flex-1 bg-[#1e1e1e] relative">
        <div v-if="!isConnected" class="absolute inset-0 flex items-center justify-center text-gray-500">
          <n-space vertical align="center">
            <n-icon size="48" :component="TerminalOutline" />
            <span>请选择 Shell 并点击连接</span>
          </n-space>
        </div>

        <Terminal
          v-if="isConnected"
          :key="terminalKey"
          :container-id="containerId"
          :shell="selectedShell"
        />
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.terminal-page {
  height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;
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

.terminal-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.09);
}

.terminal-wrapper {
  min-height: 0;
  overflow: hidden;
}

.mt-4 {
  margin-top: 16px;
}
</style>
