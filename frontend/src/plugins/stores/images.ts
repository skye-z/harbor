import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { imageApi } from '../api'
import type { Image } from '../../types'

export const useImageStore = defineStore('images', () => {
  const images = ref<Image[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalImages = computed(() => images.value.length)

  const totalSize = computed(() => {
    return images.value.reduce((sum, img) => sum + (img.size || 0), 0)
  })

  const fetchImages = async () => {
    loading.value = true
    error.value = null
    try {
      const data = await imageApi.list()
      images.value = data
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch images'
      throw err
    } finally {
      loading.value = false
    }
  }

  const getImageById = (id: string) => {
    return images.value.find(i => i.id === id)
  }

  const pullImage = async (image: string, tag?: string) => {
    loading.value = true
    error.value = null
    try {
      await imageApi.pull(image, tag)
      await fetchImages()
    } catch (err: any) {
      error.value = err.message || 'Failed to pull image'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteImage = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      await imageApi.delete(id)
      images.value = images.value.filter(i => i.id !== id)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete image'
      throw err
    } finally {
      loading.value = false
    }
  }

  const buildImage = async (name: string, dockerfile: string) => {
    loading.value = true
    error.value = null
    try {
      await imageApi.build({ dockerfile, tag: name })
      await fetchImages()
    } catch (err: any) {
      error.value = err.message || 'Failed to build image'
      throw err
    } finally {
      loading.value = false
    }
  }

  const tagImage = async (id: string, tag: string) => {
    loading.value = true
    error.value = null
    try {
      await imageApi.tag(id, tag)
      await fetchImages()
    } catch (err: any) {
      error.value = err.message || 'Failed to tag image'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    images,
    loading,
    error,
    totalImages,
    totalSize,
    fetchImages,
    getImageById,
    pullImage,
    deleteImage,
    buildImage,
    tagImage
  }
})
