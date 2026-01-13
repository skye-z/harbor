import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { systemApi } from '../api'

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const systemInfo = ref<any>(null)

  const isLoading = computed(() => loading.value)

  const fetchSystemInfo = async () => {
    loading.value = true
    error.value = null
    try {
      const data = await systemApi.getSystemInfo()
      systemInfo.value = data
      return data
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch system info'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    systemInfo,
    isLoading,
    fetchSystemInfo
  }
})

