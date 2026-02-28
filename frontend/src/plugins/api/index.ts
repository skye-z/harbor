import { apiClient } from './client'
import type {
  Container,
  Image,
  Volume,
  Network
} from '../types'

export const authApi = {
  login: async (data: { username: string; password: string }) => {
    return apiClient.post('/user/login', data)
  }
}

export const containerApi = {
  list: async () => {
    return apiClient.get('/container/list')
  },
  get: async (id: string) => {
    return apiClient.get('/container/info', { params: { id } })
  },
  create: async (data: any) => {
    return apiClient.post('/container/create', data)
  },
  operation: async (id: string,action: string) => {
    return apiClient.get('/container/operation', { params: { id, action } })
  },
  logs: async (id: string, tail?: string) => {
    return apiClient.get('/container/logs', { params: { id, stdout: true, stderr: true, timestamps: true, tail: tail || '100' } })
  },
  stats: async (id: string) => {
    return apiClient.get('/container/stat', { params: { id } })
  },
  processes: async (id: string) => {
    return apiClient.get('/container/processes', { params: { id } })
  },
  rename: async (id: string, name: string) => {
    return apiClient.post('/container/rename', { id, name })
  },
  listFiles: async (id: string, path: string) => {
    return apiClient.get('/container/files', { params: { id, path } })
  },
  copyFrom: async (id: string, srcPath: string, dstPath?: string) => {
    return apiClient.get('/container/copy/from', { params: { id, src_path: srcPath, dst_path: dstPath } })
  },
  copyTo: async (id: string, srcPath: string, dstPath: string) => {
    return apiClient.post('/container/copy/to', { id, src_path: srcPath, dst_path: dstPath })
  },
  terminal: async (id: string) => {
    return apiClient.get('/container/terminal', { params: { id } })
  }
}

export const imageApi = {
  list: async () => {
    return apiClient.get('/image/list')
  },
  get: async (id: string) => {
    return apiClient.get('/image/inspect', { params: { id } })
  },
  search: async (query: string, limit?: number) => {
    return apiClient.get('/image/search', { params: { q: query, limit } })
  },
  pull: async (image: string, tag?: string) => {
    return apiClient.get('/image/pull', { params: { image, tag } })
  },
  getPullProgress: async () => {
    return apiClient.get('/image/pull/progress')
  },
  delete: async (id: string) => {
    return apiClient.get('/image/remove', { params: { id } })
  },
  build: async (data: { dockerfile: string; tag: string }) => {
    return apiClient.get('/image/build', { params: data })
  },
  tag: async (id: string, tag: string) => {
    return apiClient.get('/image/tag', { params: { id, tag } })
  },
  push: async (id: string, repo: string, tag: string) => {
    return apiClient.get('/image/push', { params: { id, repo, tag } })
  }
}

export const volumeApi = {
  list: async () => {
    return apiClient.get('/volume/list')
  },
  create: async (data: { name?: string; driver?: string }) => {
    return apiClient.get('/volume/create', { params: data })
  },
  delete: async (id: string) => {
    return apiClient.get('/volume/remove', { params: { id } })
  }
}

export const networkApi = {
  list: async () => {
    return apiClient.get('/network/list')
  },
  create: async (data: { name: string; driver?: string }) => {
    return apiClient.get('/network/create', { params: data })
  },
  delete: async (id: string) => {
    return apiClient.get('/network/remove', { params: { id } })
  },
  connect: async (id: string, container: string) => {
    return apiClient.get('/network/connect', { params: { id, container } })
  },
  disconnect: async (id: string, container: string) => {
    return apiClient.get('/network/disconnect', { params: { id, container } })
  }
}

export const systemApi = {
  getSystemInfo: async () => {
    return apiClient.get('/docker/info')
  },
  pruneContainers: async () => {
    return apiClient.get('/docker/prune/containers')
  },
  pruneImages: async () => {
    return apiClient.get('/docker/prune/images')
  },
  pruneVolumes: async () => {
    return apiClient.get('/docker/prune/volumes')
  },
  pruneNetworks: async () => {
    return apiClient.get('/docker/prune/networks')
  },
  pruneAll: async () => {
    return apiClient.get('/docker/prune/all')
  }
}

export const logApi = {
  getRecent: async (limit: number = 10, type: string = '') => {
    return apiClient.get('/logs/recent', { params: { limit, type } })
  }
}
