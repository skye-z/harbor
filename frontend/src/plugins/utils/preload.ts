// 资源预加载工具
// 提前加载关键资源，提升用户体验

/**
 * 预加载图片资源
 * @param src 图片 URL
 * @returns Promise
 */
export function preloadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve(img)
    img.onerror = () => reject(new Error(`Failed to load image: ${src}`))
    img.src = src
  })
}

/**
 * 预加载多个图片资源
 * @param urls 图片 URL 数组
 */
export async function preloadImages(urls: string[]): Promise<void> {
  const promises = urls.map(url => preloadImage(url))
  try {
    await Promise.all(promises)
  } catch (error) {
    console.warn('Some images failed to preload:', error)
  }
}

/**
 * 预加载脚本文件
 * @param src 脚本 URL
 */
export function preloadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const script = document.createElement('link')
    script.rel = 'preload'
    script.as = 'script'
    script.href = src
    script.onload = () => resolve()
    script.onerror = () => reject(new Error(`Failed to load script: ${src}`))
    document.head.appendChild(script)
  })
}

/**
 * 预加载样式文件
 * @param href 样式 URL
 */
export function preloadStyle(href: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const link = document.createElement('link')
    link.rel = 'preload'
    link.as = 'style'
    link.href = href
    link.onload = () => resolve()
    link.onerror = () => reject(new Error(`Failed to load style: ${href}`))
    document.head.appendChild(link)
  })
}

/**
 * DNS 预解析
 * @param hostnames 主机名数组
 */
export function prefetchDNS(hostnames: string[]): void {
  hostnames.forEach(hostname => {
    const link = document.createElement('link')
    link.rel = 'dns-prefetch'
    link.href = `//${hostname}`
    document.head.appendChild(link)
  })
}

/**
 * 预连接资源
 * @param urls URL 数组
 */
export function prefetchConnections(urls: string[]): void {
  urls.forEach(url => {
    const link = document.createElement('link')
    link.rel = 'preconnect'
    link.href = url
    document.head.appendChild(link)
  })
}

/**
 * 路由级别的预加载策略
 * @param currentRoute 当前路由
 */
export function preloadRouteResources(currentRoute: string): void {
  const routePreloadMap: Record<string, string[]> = {
    Dashboard: [],
    Containers: [],
    Images: [],
    Storage: [],
    System: []
  }

  const resources = routePreloadMap[currentRoute] || []

  // 预连接 API
  if (resources.length > 0) {
    prefetchConnections(['/api'])
  }
}

/**
 * 优化首屏加载
 * 预加载关键资源，提升首屏加载速度
 */
export function optimizeFirstScreenLoad(): void {
  // DNS 预解析
  prefetchDNS(['api.harbor.com', 'cdn.jsdelivr.net'])

  // 预连接关键域名
  prefetchConnections(['/api'])

  // 预加载关键字体（如果使用自定义字体）
  // preloadStyle('/fonts/custom-font.woff2')

  // 预加载关键图标
  // preloadImages(['/icons/icon.png'])
}

/**
 * 监听网络状态，优化加载策略
 */
export function setupNetworkAwareLoading(): void {
  const connection = (navigator as any).connection

  if (!connection) return

  // 根据网络速度调整预加载策略
  if (connection.effectiveType === '4g' || connection.effectiveType === 'wifi') {
    // 快速网络，积极预加载
    optimizeFirstScreenLoad()
  } else if (connection.effectiveType === '2g' || connection.effectiveType === '3g') {
    // 慢速网络，谨慎预加载
    console.log('Slow network detected, using conservative preload strategy')
  }

  // 监听网络变化
  connection.addEventListener('change', () => {
    console.log('Network changed:', connection.effectiveType)
    if (connection.effectiveType === '4g' || connection.effectiveType === 'wifi') {
      optimizeFirstScreenLoad()
    }
  })
}

/**
 * 懒加载图片（Intersection Observer）
 * @param selector 图片选择器
 * @param rootMargin 根边距
 */
export function lazyLoadImages(
  selector: string = 'img[data-src]',
  rootMargin: string = '50px'
): IntersectionObserver | null {
  if (!('IntersectionObserver' in window)) {
    // 不支持 IntersectionObserver，直接加载
    document.querySelectorAll(selector).forEach(img => {
      const image = img as HTMLImageElement
      if (image.dataset.src) {
        image.src = image.dataset.src
        delete image.dataset.src
      }
    })
    return null
  }

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const image = entry.target as HTMLImageElement
        if (image.dataset.src) {
          image.src = image.dataset.src
          delete image.dataset.src
          observer.unobserve(image)
        }
      }
    })
  }, {
    rootMargin,
    threshold: 0.01
  })

  // 开始观察所有懒加载图片
  setTimeout(() => {
    document.querySelectorAll(selector).forEach(img => {
      observer.observe(img)
    })
  }, 100)

  return observer
}

/**
 * 清理预加载的资源
 */
export function cleanupPrefetch(): void {
  // 移除所有预加载的链接
  document.querySelectorAll('link[rel="dns-prefetch"], link[rel="preconnect"]').forEach(link => {
    link.remove()
  })
}
