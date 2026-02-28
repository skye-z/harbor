<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContainerStore } from '../plugins/stores/containers'
import { useMessage, NSelect, NButton, NIcon } from 'naive-ui'
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
const terminalRef = ref<InstanceType<typeof Terminal>>()

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
  terminalKey.value++
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
    <div class="toolbar">
      <div class="toolbar-left">
        <n-button quaternary size="small" @click="handleBack">
          <template #icon>
            <n-icon :component="ArrowBack" />
          </template>
        </n-button>
        <span class="container-name">{{ container?.names?.[0]?.replace(/^\//, '') || containerId }}</span>
      </div>
      <div class="toolbar-center">
        <n-select
          v-model:value="selectedShell"
          :options="shellOptions"
          size="small"
          style="width: 120px"
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
      </div>
      <div class="toolbar-right">
        <div class="status-indicator">
          <span
            class="status-dot"
            :class="isConnected ? 'connected' : 'disconnected'"
          ></span>
          <span class="status-text">{{ isConnected ? '已连接' : '未连接' }}</span>
        </div>
      </div>
    </div>

    <div class="terminal-area">
      <div v-if="!isConnected" class="placeholder">
        <n-icon size="48" :component="TerminalOutline" class="placeholder-icon" />
        <span>选择 Shell 后点击连接</span>
      </div>
      <Terminal
        v-else
        :key="terminalKey"
        :container-id="containerId"
        :shell="selectedShell"
        ref="terminalRef"
      />
    </div>
  </div>
</template>

<style scoped>
.terminal-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background: var(--terminal-bg, #1e1e1e);
  z-index: 100;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: var(--toolbar-bg, #252526);
  border-bottom: 1px solid var(--border-color, #3c3c3c);
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.container-name {
  color: var(--text-primary, #cccccc);
  font-weight: 500;
  font-size: 14px;
}

.toolbar-center {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  width: 100px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.connected {
  background: #0dbc79;
  box-shadow: 0 0 8px rgba(13, 188, 121, 0.6);
}

.status-dot.disconnected {
  background: #f14c4c;
}

.status-text {
  color: var(--text-secondary, #808080);
  font-size: 12px;
}

.terminal-area {
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
  color: var(--text-secondary, #808080);
}

.placeholder-icon {
  opacity: 0.5;
}

/* 明亮模式 */
[data-theme='light'] .terminal-page {
  --terminal-bg: #ffffff;
  --toolbar-bg: #f5f5f5;
  --border-color: #e0e0e0;
  --text-primary: #333333;
  --text-secondary: #666666;
}

/* 暗色模式 */
[data-theme='dark'] .terminal-page {
  --terminal-bg: #1e1e1e;
  --toolbar-bg: #252526;
  --border-color: #3c3c3c;
  --text-primary: #cccccc;
  --text-secondary: #808080;
}
</style>
