<template>
  <div class="resource-topology">
    <v-chart ref="chartRef" :option="topologyOption" style="height: calc(100vh - 200px); max-height: 500px;" @click="handleNodeClick" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { GraphChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import { useContainerStore } from '../plugins/stores/containers'
import { useImageStore } from '../plugins/stores/images'
import { useVolumeStore } from '../plugins/stores/volumes'
import { useNetworkStore } from '../plugins/stores/networks'

const router = useRouter()

use([
  CanvasRenderer,
  GraphChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const containerStore = useContainerStore()
const imageStore = useImageStore()
const volumeStore = useVolumeStore()
const networkStore = useNetworkStore()

const nodes = ref<any[]>([])
const links = ref<any[]>([])
const chartRef = ref<any>(null)

// 处理节点点击事件
const handleNodeClick = (params: any) => {
  if (params.dataType === 'node') {
    const node = params.data
    const nodeId = node.id
    
    if (nodeId.startsWith('container-')) {
      const containerId = nodeId.replace('container-', '')
      router.push({ name: 'ContainerDetail', params: { id: containerId } })
    } else if (nodeId.startsWith('image-')) {
      const imageId = nodeId.replace('image-', '')
      router.push({ name: 'ImageDetail', params: { id: imageId } })
    } else if (nodeId.startsWith('volume-')) {
      const volumeId = nodeId.replace('volume-', '')
      router.push({ name: 'VolumeDetail', params: { id: volumeId } })
    } else if (nodeId.startsWith('network-')) {
      const networkId = nodeId.replace('network-', '')
      router.push({ name: 'NetworkDetail', params: { id: networkId } })
    }
  }
}

const topologyOption = computed<EChartsOption>(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (params: any) => {
      if (params.dataType === 'edge') {
        const link = params.data
        const containerNode = nodes.value.find(n => n.id === link.source)
        const networkNode = nodes.value.find(n => n.id === link.target)
        
        if (link.name && containerNode && networkNode) {
          return `<div style="padding: 8px;">
            <div style="font-weight: bold; margin-bottom: 4px;">容器 → 网络</div>
            <div>容器: ${containerNode.name}</div>
            <div>网络: ${networkNode.name}</div>
            <div>IP地址: ${link.name || 'N/A'}</div>
          </div>`
        }
        return ''
      } else if (params.dataType === 'node') {
        const node = params.data
        let info = `<div style="padding: 8px;">
          <div style="font-weight: bold; margin-bottom: 4px;">${node.name}</div>`
        
        if (node.category === 0) {
          info += `<div>类型: 镜像</div>`
          info += `<div>大小: ${formatSize(node.value)}</div>`
        } else if (node.category === 1) {
          info += `<div>类型: 容器</div>`
          const container = containerStore.containers.find(c => c.id === node.id.replace('container-', ''))
          if (container) {
            info += `<div>镜像: ${container.image}</div>`
            info += `<div>状态: ${container.state}</div>`
            if (container.network_settings?.networks) {
              const networks = Object.entries(container.network_settings.networks).map(([name, config]: [string, any]) => {
                return `${name}: ${config.ip_address || 'N/A'}`
              })
              if (networks.length > 0) {
                info += `<div style="margin-top: 4px;">网络: ${networks.join(', ')}</div>`
              }
            }
          }
        } else if (node.category === 2) {
          info += `<div>类型: 卷</div>`
          const volume = volumeStore.volumes.find(v => v.id === node.id.replace('volume-', ''))
          if (volume) {
            info += `<div>驱动: ${volume.driver}</div>`
          }
        } else if (node.category === 3) {
          info += `<div>类型: 网络</div>`
          const network = networkStore.networks.find(n => n.id === node.id.replace('network-', ''))
          if (network) {
            info += `<div>驱动: ${network.driver}</div>`
            info += `<div>子网: ${network.subnet || 'N/A'}</div>`
          }
        }
        
        info += '</div>'
        return info
      }
      return ''
    }
  },
  legend: {
    data: ['容器', '镜像', '卷', '网络']
  },
  series: [
    {
      type: 'graph',
      layout: 'force',
      data: nodes.value,
      links: links.value,
      categories: [
        { name: '镜像' },
        { name: '容器' },
        { name: '卷' },
        { name: '网络' }
      ],
      roam: true,
      label: {
        show: true,
        position: 'right',
        formatter: '{b}',
        textBorderColor: '#fff',
        textBorderWidth: 1,
        textShadowColor: '#fff',
        textShadowBlur: 2
      },
      labelLayout: {
        hideOverlap: true
      },
      edgeSymbol: ['circle', 'arrow'],
      edgeSymbolSize: [4, 10],
      scaleLimit: {
        min: 0.4,
        max: 2
      },
      lineStyle: {
        color: 'source',
        curveness: 0.3
      },
      emphasis: {
        focus: 'adjacency',
        lineStyle: {
          width: 10
        }
      },
      force: {
        repulsion: 1000,
        edgeLength: 80,
        gravity: 0.1
      }
    }
  ]
}))

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

const buildTopology = () => {
  const newNodes: any[] = []
  const newLinks: any[] = []

  imageStore.images.forEach((image) => {
    newNodes.push({
      id: `image-${image.id}`,
      name: image.repo_tags?.[0] || image.id,
      category: 0,
      symbolSize: 20,
      value: image.size,
      itemStyle: {
        color: '#5470c6'
      }
    })
  })

  networkStore.networks.forEach((network) => {
    newNodes.push({
      id: `network-${network.id}`,
      name: network.name,
      category: 3,
      symbolSize: 15,
      itemStyle: {
        color: '#91cc75'
      }
    })
  })

  volumeStore.volumes.forEach((volume) => {
    newNodes.push({
      id: `volume-${volume.id}`,
      name: volume.name,
      category: 2,
      symbolSize: 15,
      itemStyle: {
        color: '#fac858'
      }
    })
  })

  containerStore.containers.forEach((container) => {
    newNodes.push({
      id: `container-${container.id}`,
      name: container.names[0]?.replace(/^\//, '') || container.id,
      category: 1,
      symbolSize: 25,
      itemStyle: {
        color: '#ee6666'
      }
    })

    const imageNode = newNodes.find(n => {
      const tagName = n.name
      if (tagName === container.image) return true
      if (tagName.startsWith(container.image + ':')) return true
      if (tagName === container.image + ':latest') return true
      return false
    })
    if (imageNode) {
      newLinks.push({
        source: imageNode.id,
        target: `container-${container.id}`,
        lineStyle: {
          color: '#5470c6'
        }
      })
    }

    if (container.network_settings && container.network_settings.networks) {
      Object.keys(container.network_settings.networks).forEach((networkName) => {
        const networkNode = newNodes.find(n => n.name === networkName)
        if (networkNode) {
          const networkConfig = container.network_settings.networks[networkName]
          newLinks.push({
            source: `container-${container.id}`,
            target: networkNode.id,
            name: networkConfig.ip_address || 'N/A',
            lineStyle: {
              color: '#91cc75'
            }
          })
        }
      })
    }

    if (container.mount && Array.isArray(container.mount) && container.mount.length > 0) {
      container.mount.forEach((mount: any) => {
        if (mount.name) {
          const volumeNode = newNodes.find(n => n.name === mount.name)
          if (volumeNode) {
            newLinks.push({
              source: `container-${container.id}`,
              target: volumeNode.id,
              lineStyle: {
                color: '#fac858'
              }
            })
          }
        }
      })
    }
  })

  nodes.value = newNodes
  links.value = newLinks
}

onMounted(async () => {
  try {
    await Promise.all([
      containerStore.fetchContainers(),
      imageStore.fetchImages(),
      volumeStore.fetchVolumes(),
      networkStore.fetchNetworks()
    ])
    buildTopology()
  } catch (error) {
    console.error('Failed to fetch data for topology:', error)
  }
})
</script>

<style scoped>
.resource-topology {
  width: 100%;
  height: 100%;
}
</style>
