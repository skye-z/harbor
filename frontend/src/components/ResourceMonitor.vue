<template>
  <div class="resource-monitor">
    <n-space vertical :size="16">
      <n-card title="CPU使用率">
        <v-chart :option="cpuChartOption" style="height: 200px" />
      </n-card>
      <n-card title="内存使用率">
        <v-chart :option="memoryChartOption" style="height: 200px" />
      </n-card>
      <n-card title="网络流量">
        <v-chart :option="networkChartOption" style="height: 200px" />
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

interface Props {
  containerId: string
}

const props = defineProps<Props>()

const maxDataPoints = 60

const cpuData = ref<number[]>(new Array(maxDataPoints).fill(0))
const memoryData = ref<number[]>(new Array(maxDataPoints).fill(0))
const networkRxData = ref<number[]>(new Array(maxDataPoints).fill(0))
const networkTxData = ref<number[]>(new Array(maxDataPoints).fill(0))
const timeLabels = ref<string[]>([])

let ws: WebSocket | null = null

const cpuChartOption = computed<EChartsOption>(() => ({
  tooltip: {
    trigger: 'axis'
  },
  xAxis: {
    type: 'category',
    data: timeLabels.value
  },
  yAxis: {
    type: 'value',
    max: 100,
    axisLabel: {
      formatter: '{value}%'
    }
  },
  series: [
    {
      name: 'CPU使用率',
      type: 'line',
      data: cpuData.value,
      smooth: true,
      areaStyle: {}
    }
  ]
}))

const memoryChartOption = computed<EChartsOption>(() => ({
  tooltip: {
    trigger: 'axis'
  },
  xAxis: {
    type: 'category',
    data: timeLabels.value
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      formatter: (value: number) => `${(value / 1024 / 1024).toFixed(2)} MB`
    }
  },
  series: [
    {
      name: '内存使用',
      type: 'line',
      data: memoryData.value,
      smooth: true,
      areaStyle: {}
    }
  ]
}))

const networkChartOption = computed<EChartsOption>(() => ({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['网络接收', '网络发送']
  },
  xAxis: {
    type: 'category',
    data: timeLabels.value
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      formatter: (value: number) => `${(value / 1024).toFixed(2)} KB`
    }
  },
  series: [
    {
      name: '网络接收',
      type: 'line',
      data: networkRxData.value,
      smooth: true
    },
    {
      name: '网络发送',
      type: 'line',
      data: networkTxData.value,
      smooth: true
    }
  ]
}))

const updateTimeLabels = () => {
  const now = new Date()
  const timeStr = now.toLocaleTimeString('zh-CN')
  timeLabels.value.push(timeStr)
  if (timeLabels.value.length > maxDataPoints) {
    timeLabels.value.shift()
  }
}

const updateChartData = (stats: any) => {
  updateTimeLabels()

  cpuData.value.push(stats.cpu || 0)
  if (cpuData.value.length > maxDataPoints) {
    cpuData.value.shift()
  }

  memoryData.value.push(stats.memory || 0)
  if (memoryData.value.length > maxDataPoints) {
    memoryData.value.shift()
  }

  if (stats.network && stats.network.eth0) {
    networkRxData.value.push(stats.network.eth0.rx_bytes || 0)
    networkTxData.value.push(stats.network.eth0.tx_bytes || 0)
  } else if (stats.network) {
    const keys = Object.keys(stats.network as any)
    if (keys.length > 0) {
      const netKey = keys[0] as string
      const networkData = (stats.network as any)[netKey]
      networkRxData.value.push((networkData?.rx_bytes ?? 0) as number)
      networkTxData.value.push((networkData?.tx_bytes ?? 0) as number)
    }
  } else {
    networkRxData.value.push(0)
    networkTxData.value.push(0)
  }

  if (networkRxData.value.length > maxDataPoints) {
    networkRxData.value.shift()
  }
  if (networkTxData.value.length > maxDataPoints) {
    networkTxData.value.shift()
  }
}

const startStatsMonitoring = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws${token ? `?token=${encodeURIComponent(token)}` : ''}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    ws?.send(JSON.stringify({
      type: 'stats',
      data: {
        container_id: props.containerId
      }
    }))
  }

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)

      switch (message.type) {
        case 'stats_started':
          console.log('Stats monitoring started')
          break
        case 'stats_update':
          updateChartData(message.data)
          break
        case 'stats_error':
          console.error('Stats error:', message.data.error)
          break
      }
    } catch (error) {
      console.error('Error parsing stats message:', error)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  ws.onclose = () => {
    console.log('WebSocket disconnected')
  }
}

const stopStatsMonitoring = () => {
  if (ws) {
    ws.send(JSON.stringify({
      type: 'stats_stop'
    }))
    ws.close()
    ws = null
  }
}

onMounted(() => {
  for (let i = 0; i < maxDataPoints; i++) {
    const now = new Date(Date.now() - (maxDataPoints - i) * 1000)
    timeLabels.value.push(now.toLocaleTimeString('zh-CN'))
  }
  startStatsMonitoring()
})

onBeforeUnmount(() => {
  stopStatsMonitoring()
})
</script>

<style scoped>
.resource-monitor {
  width: 100%;
}
</style>
