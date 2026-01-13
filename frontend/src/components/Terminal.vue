<template>
  <div class="terminal-container">
    <div ref="terminalRef" class="terminal"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { containerApi } from '../plugins/api'

interface Props {
  containerId: string
  shell?: string
}

const props = withDefaults(defineProps<Props>(), {
  shell: '/bin/sh'
})

const terminalRef = ref<HTMLElement>()
let terminal: Terminal | null = null
let ws: WebSocket | null = null

const connectTerminal = async () => {
  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
    
    const termResponse = await containerApi.get(props.containerId)
    
    const wsHost = window.location.host.replace('3000', '12800')
    const wsUrl = `ws://${wsHost}/api/container/terminal/ws?exec_id=${termResponse.exec_id}${token ? `&token=${encodeURIComponent(token)}` : ''}`
    
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)

        switch (message.type) {
          case 'terminal_output':
            if (terminal && message.data) {
              terminal.write(atob(message.data))
            }
            break
          case 'terminal_error':
            console.error('Terminal error:', message.data)
            break
        }
      } catch (error) {
        console.error('Error parsing message:', error)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
    }
  } catch (error) {
    console.error('Failed to connect terminal:', error)
  }
}

const initTerminal = () => {
  if (!terminalRef.value) return

  terminal = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Fira Code, Consolas, monospace',
    theme: {
      background: '#0d1117',
      foreground: '#c9d1d9',
      cursor: '#58a6ff',
      black: '#0d1117',
      red: '#ff7b72',
      green: '#3fb950',
      yellow: '#d29922',
      blue: '#58a6ff',
      magenta: '#bc8cff',
      cyan: '#39c5cf',
      white: '#c9d1d9',
      brightBlack: '#484f58',
      brightRed: '#ff7b72',
      brightGreen: '#3fb950',
      brightYellow: '#d29922',
      brightBlue: '#58a6ff',
      brightMagenta: '#bc8cff',
      brightCyan: '#39c5cf',
      brightWhite: '#f0f6fc'
    }
  })

  const fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.open(terminalRef.value)
  fitAddon.fit()
  terminal.focus()

  terminal.onData((data: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      const encoded = btoa(data)
      ws.send(JSON.stringify({
        type: 'input',
        data: encoded
      }))
    }
  })

  terminal.onResize((size) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'resize',
        rows: size.rows,
        cols: size.cols
      }))
    }
  })

  connectTerminal()
}

onMounted(() => {
  initTerminal()
})

onBeforeUnmount(() => {
  if (ws) {
    ws.close()
  }
  if (terminal) {
    terminal.dispose()
  }
})
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
}

.terminal {
  height: 100%;
  padding: 10px;
}

:deep(.xterm) {
  padding: 10px;
}

:deep(.xterm-viewport) {
  overflow-y: auto;
}
</style>
