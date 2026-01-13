import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { networkApi } from '../api'
import type { Network } from '../../types'
import { cacheManager, CACHE_KEYS, CACHE_TTL } from '../utils/cache'

export const useNetworkStore = defineStore('networks', () => {
  const networks = ref<Network[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalNetworks = computed(() => networks.value.length)

  const fetchNetworks = async (forceRefresh: boolean = false) => {
    // 检查缓存
    if (!forceRefresh) {
      const cached = cacheManager.get<Network[]>(CACHE_KEYS.NETWORKS_LIST)
      if (cached) {
        networks.value = cached
        return
      }
    }

    loading.value = true
    error.value = null
    try {
      const data = await networkApi.list()
      networks.value = data
      // 缓存数据，有效期10分钟
      cacheManager.set(CACHE_KEYS.NETWORKS_LIST, data, CACHE_TTL.LONG)
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch networks'
      throw err
    } finally {
      loading.value = false
    }
  }

  const getNetworkById = (id: string) => {
    return networks.value.find(n => n.id === id)
  }

  const createNetwork = async (data: any) => {
    loading.value = true
    error.value = null
    try {
      const network = await networkApi.create(data)
      networks.value.push(network)
      // 清空网络列表缓存
      cacheManager.remove(CACHE_KEYS.NETWORKS_LIST)
      return network
    } catch (err: any) {
      error.value = err.message || 'Failed to create network'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteNetwork = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      await networkApi.delete(id)
      networks.value = networks.value.filter(n => n.id !== id)
      // 清空网络列表缓存
      cacheManager.remove(CACHE_KEYS.NETWORKS_LIST)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete network'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    networks,
    loading,
    error,
    totalNetworks,
    fetchNetworks,
    getNetworkById,
    createNetwork,
    deleteNetwork
  }
})
