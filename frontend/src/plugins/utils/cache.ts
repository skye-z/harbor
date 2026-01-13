// 缓存管理工具
// 提供统一的缓存接口，支持内存缓存和持久化存储

interface CacheData<T> {
  data: T
  timestamp: number
  ttl: number // 缓存生存时间（毫秒）
}

class CacheManager {
  private memoryCache = new Map<string, CacheData<any>>()
  private static instance: CacheManager

  private constructor() {}

  // 单例模式
  static getInstance(): CacheManager {
    if (!CacheManager.instance) {
      CacheManager.instance = new CacheManager()
    }
    return CacheManager.instance
  }

  /**
   * 设置缓存
   * @param key 缓存键
   * @param value 缓存值
   * @param ttl 生存时间（毫秒），默认 5 分钟
   * @param persist 是否持久化到 localStorage，默认 false
   */
  set<T>(key: string, value: T, ttl: number = 5 * 60 * 1000, persist: boolean = false): void {
    const cacheData: CacheData<T> = {
      data: value,
      timestamp: Date.now(),
      ttl
    }

    // 内存缓存
    this.memoryCache.set(key, cacheData)

    // 持久化缓存
    if (persist) {
      try {
        localStorage.setItem(key, JSON.stringify(cacheData))
      } catch (e) {
        console.warn(`Failed to persist cache for key: ${key}`, e)
      }
    }
  }

  /**
   * 获取缓存
   * @param key 缓存键
   * @returns 缓存值，如果不存在或已过期则返回 null
   */
  get<T>(key: string): T | null {
    // 先尝试从内存缓存获取
    const memCache = this.memoryCache.get(key) as CacheData<T>
    if (memCache) {
      if (!this.isExpired(memCache)) {
        return memCache.data
      }
      this.memoryCache.delete(key)
    }

    // 尝试从 localStorage 获取
    const localCache = this.getFromLocalStorage<T>(key)
    if (localCache) {
      // 同步到内存缓存
      this.memoryCache.set(key, localCache)
      if (!this.isExpired(localCache)) {
        return localCache.data
      }
      // 如果已过期，删除 localStorage 中的缓存
      this.remove(key)
    }

    return null
  }

  /**
   * 删除指定键的缓存
   * @param key 缓存键
   */
  remove(key: string): void {
    this.memoryCache.delete(key)
    try {
      localStorage.removeItem(key)
    } catch (e) {
      console.warn(`Failed to remove cache for key: ${key}`, e)
    }
  }

  /**
   * 清空所有缓存
   */
  clear(): void {
    this.memoryCache.clear()
    try {
      localStorage.clear()
    } catch (e) {
      console.warn('Failed to clear localStorage', e)
    }
  }

  /**
   * 清空指定前缀的缓存
   * @param prefix 缓存键前缀
   */
  clearByPrefix(prefix: string): void {
    // 清空内存缓存
    for (const key of this.memoryCache.keys()) {
      if (key.startsWith(prefix)) {
        this.memoryCache.delete(key)
      }
    }

    // 清空 localStorage 缓存
    try {
      const keys = Object.keys(localStorage)
      for (const key of keys) {
        if (key.startsWith(prefix)) {
          localStorage.removeItem(key)
        }
      }
    } catch (e) {
      console.warn('Failed to clear caches by prefix', e)
    }
  }

  /**
   * 从 localStorage 获取缓存
   * @param key 缓存键
   * @returns 缓存数据
   */
  private getFromLocalStorage<T>(key: string): CacheData<T> | null {
    try {
      const value = localStorage.getItem(key)
      if (value) {
        return JSON.parse(value) as CacheData<T>
      }
    } catch (e) {
      console.warn(`Failed to parse cache for key: ${key}`, e)
    }
    return null
  }

  /**
   * 检查缓存是否过期
   * @param cache 缓存数据
   * @returns 是否过期
   */
  private isExpired<T>(cache: CacheData<T>): boolean {
    return Date.now() - cache.timestamp > cache.ttl
  }

  /**
   * 获取缓存统计信息
   */
  getStats(): {
    memoryCacheSize: number
    localStorageSize: number
  } {
    let localStorageSize = 0
    try {
      localStorageSize = Object.keys(localStorage).length
    } catch (e) {
      console.warn('Failed to get localStorage size', e)
    }

    return {
      memoryCacheSize: this.memoryCache.size,
      localStorageSize
    }
  }
}

// 导出单例实例
export const cacheManager = CacheManager.getInstance()

// 导出缓存键常量
export const CACHE_KEYS = {
  // 容器缓存
  CONTAINERS_LIST: 'cache:containers:list',
  CONTAINER_DETAIL: 'cache:container:detail:',

  // 镜像缓存
  IMAGES_LIST: 'cache:images:list',
  IMAGE_DETAIL: 'cache:image:detail:',

  // 网络缓存
  NETWORKS_LIST: 'cache:networks:list',
  NETWORK_DETAIL: 'cache:network:detail:',

  // 卷缓存
  VOLUMES_LIST: 'cache:volumes:list',
  VOLUME_DETAIL: 'cache:volume:detail:',

  // 系统信息缓存
  SYSTEM_INFO: 'cache:system:info'
}

// 默认缓存生存时间（毫秒）
export const CACHE_TTL = {
  SHORT: 60 * 1000,        // 1 分钟
  MEDIUM: 5 * 60 * 1000,   // 5 分钟
  LONG: 15 * 60 * 1000,     // 15 分钟
  VERY_LONG: 60 * 60 * 1000 // 1 小时
}
