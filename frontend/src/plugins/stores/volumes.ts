import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { volumeApi } from '../api'
import type { Volume } from '../../types'
import { cacheManager, CACHE_KEYS, CACHE_TTL } from '../utils/cache'

export const useVolumeStore = defineStore('volumes', () => {
  const volumes = ref<Volume[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalVolumes = computed(() => volumes.value.length)

  const fetchVolumes = async (forceRefresh: boolean = false) => {
    // 检查缓存
    if (!forceRefresh) {
      const cached = cacheManager.get<Volume[]>(CACHE_KEYS.VOLUMES_LIST)
      if (cached) {
        volumes.value = cached
        return
      }
    }

    loading.value = true
    error.value = null
    try {
      const data = await volumeApi.list()
      volumes.value = data
      // 缓存数据，有效期10分钟
      cacheManager.set(CACHE_KEYS.VOLUMES_LIST, data, CACHE_TTL.LONG)
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch volumes'
      throw err
    } finally {
      loading.value = false
    }
  }

  const getVolumeById = (id: string) => {
    return volumes.value.find(v => v.id === id)
  }

  const createVolume = async (data: any) => {
    loading.value = true
    error.value = null
    try {
      const volume = await volumeApi.create(data)
      volumes.value.push(volume)
      // 清空卷列表缓存
      cacheManager.remove(CACHE_KEYS.VOLUMES_LIST)
      return volume
    } catch (err: any) {
      error.value = err.message || 'Failed to create volume'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteVolume = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      await volumeApi.delete(id)
      volumes.value = volumes.value.filter(v => v.id !== id)
      // 清空卷列表缓存
      cacheManager.remove(CACHE_KEYS.VOLUMES_LIST)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete volume'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    volumes,
    loading,
    error,
    totalVolumes,
    fetchVolumes,
    getVolumeById,
    createVolume,
    deleteVolume
  }
})
