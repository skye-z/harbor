<template>
  <Transition name="loading" mode="out-in">
    <div v-if="visible" class="page-loading-container">
      <div class="loading-backdrop">
        <div class="loading-content">
          <n-spin :size="48">
            <template #description>
              <div class="loading-text">{{ text }}</div>
            </template>
          </n-spin>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  visible?: boolean
  text?: string
}>()

const localVisible = ref(props.visible)
const localText = ref(props.text || '加载中...')

watch(() => props.visible, (newVal) => {
  localVisible.value = newVal
})

watch(() => props.text, (newVal) => {
  if (newVal) {
    localText.value = newVal
  }
})
</script>

<style scoped>
.loading-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  transition: background-color 0.3s ease;
}

[data-theme='dark'] .loading-backdrop {
  background-color: rgba(0, 0, 0, 0.7);
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.loading-text {
  font-size: 14px;
  color: #666;
  margin-top: 12px;
  transition: color 0.3s ease;
}

[data-theme='dark'] .loading-text {
  color: #999;
}

.loading-enter-active,
.loading-leave-active {
  transition: opacity 0.3s ease;
}

.loading-enter-from,
.loading-leave-to {
  opacity: 0;
}

.loading-enter-to,
.loading-leave-from {
  opacity: 1;
}

.loading-content :deep(.n-spin) {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
