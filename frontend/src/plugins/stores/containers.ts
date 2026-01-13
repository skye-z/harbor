import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { containerApi } from '../api'
import type { Container } from '../../types'

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

  const startContainer = async (id: string) => {
    await containerApi.start(id)
    await fetchContainers()
  }

  const stopContainer = async (id: string) => {
    await containerApi.stop(id)
    await fetchContainers()
  }

  const restartContainer = async (id: string) => {
    await containerApi.restart(id)
    await fetchContainers()
  }

  const deleteContainer = async (id: string) => {
    await containerApi.delete(id)
    containers.value = containers.value.filter(c => c.id !== id)
  }

  const getContainerLogs = async (id: string) => {
    return await containerApi.logs(id)
  }

  const getContainerStats = async (id: string) => {
    return await containerApi.stats(id)
  }

  const getContainerProcesses = async (id: string) => {
    return await containerApi.processes(id)
  }

  const getContainerInfo = async (id: string) => {
    return await containerApi.get(id)
  }

  const renameContainer = async (id: string, name: string) => {
    return await containerApi.rename(id, name)
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
    startContainer,
    stopContainer,
    restartContainer,
    deleteContainer,
    getContainerLogs,
    getContainerStats,
    getContainerProcesses,
    getContainerInfo,
    renameContainer,
    copyFromContainer,
    copyToContainer
  }
})
