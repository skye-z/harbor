import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { containerApi } from '../api'
import type { Container } from '../../types'

const CACHE_PREFIX = 'container_detail_'
const CACHE_EXPIRE = 5 * 60 * 1000

interface CacheItem {
  data: any
  timestamp: number
}

const getCache = (id: string): any | null => {
  try {
    const key = CACHE_PREFIX + id
    const cached = localStorage.getItem(key)
    if (cached) {
      const item: CacheItem = JSON.parse(cached)
      const now = Date.now()
      if (now - item.timestamp < CACHE_EXPIRE) {
        return item.data
      }
      localStorage.removeItem(key)
    }
  } catch (e) {
    console.error('Failed to read cache:', e)
  }
  return null
}

const setCache = (id: string, data: any) => {
  try {
    const key = CACHE_PREFIX + id
    const item: CacheItem = {
      data,
      timestamp: Date.now()
    }
    localStorage.setItem(key, JSON.stringify(item))
  } catch (e) {
    console.error('Failed to write cache:', e)
  }
}

export const useContainerStore = defineStore('containers', () => {
  const containers = ref<Container[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const runningContainers = computed(() =>
    containers.value.filter(c => c.state === 'running').length
  )

  const totalContainers = computed(() => containers.value.length)

  const fetchContainers = async () => {
    loading.value = true
    error.value = null
    try {
      const data = await containerApi.list()
      containers.value = data || []
    } catch (err: any) {
      error.value = err.message || '获取容器列表失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const getContainerById = (id: string) => {
    return containers.value.find(c => c.id === id)
  }

  const operationContainer = async (id: string,action: string) => {
    await containerApi.operation(id,action)
    await fetchContainers()
  }

  const getContainerLogs = async (id: string, tail?: string) => {
    return await containerApi.logs(id, tail)
  }

  const getContainerStats = async (id: string) => {
    return await containerApi.stats(id)
  }

  const getContainerProcesses = async (id: string) => {
    return await containerApi.processes(id)
  }

  const getContainerInfo = async (id: string, useCache = true) => {
    if (useCache) {
      const cached = getCache(id)
      if (cached) {
        return cached
      }
    }
    const data = await containerApi.get(id)
    setCache(id, data)
    return data
  }

  const getContainerInfoCached = (id: string): any | null => {
    return getCache(id)
  }

  const refreshContainerInfo = async (id: string) => {
    const data = await containerApi.get(id)
    setCache(id, data)
    return data
  }

  const renameContainer = async (id: string, name: string) => {
    return await containerApi.rename(id, name)
  }

  const listContainerFiles = async (id: string, path: string) => {
    return await containerApi.listFiles(id, path)
  }

  const copyFromContainer = async (id: string, srcPath: string, dstPath?: string) => {
    return await containerApi.copyFrom(id, srcPath, dstPath)
  }

  const copyToContainer = async (id: string, srcPath: string, dstPath: string) => {
    return await containerApi.copyTo(id, srcPath, dstPath)
  }

  return {
    containers,
    loading,
    error,
    runningContainers,
    totalContainers,
    fetchContainers,
    getContainerById,
    operationContainer,
    getContainerLogs,
    getContainerStats,
    getContainerProcesses,
    getContainerInfo,
    getContainerInfoCached,
    refreshContainerInfo,
    renameContainer,
    listContainerFiles,
    copyFromContainer,
    copyToContainer
  }
})
