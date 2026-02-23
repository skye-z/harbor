<template>
  <div ref="terminalRef" class="terminal-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
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
let fitAddon: FitAddon | null = null

const connectTerminal = async () => {
  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
    
    const termResponse = await containerApi.terminal(props.containerId)
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/api/container/terminal/ws?exec_id=${termResponse.exec_id}${token ? `&token=${encodeURIComponent(token)}` : ''}`
    
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      if (terminal) {
        terminal.write(event.data)
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

const fitTerminal = () => {
  if (fitAddon && terminal) {
    fitAddon.fit()
  }
}

const handleResize = () => {
  fitTerminal()
}

const initTerminal = () => {
  if (!terminalRef.value) return

  terminal = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Fira Code, Consolas, "Courier New", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#ffffff',
      cursorAccent: '#1e1e1e',
      selectionBackground: '#264f78',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#ffffff'
    },
    allowProposedApi: true
  })

  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.loadAddon(new WebLinksAddon())
  terminal.open(terminalRef.value)
  
  setTimeout(() => {
    fitTerminal()
    terminal?.focus()
  }, 0)

  terminal.onData((data: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })

  terminalRef.value.addEventListener('contextmenu', (e) => {
    e.preventDefault()
    if (terminal && terminal.hasSelection()) {
      navigator.clipboard.writeText(terminal.getSelection())
      terminal.clearSelection()
    } else {
      navigator.clipboard.readText().then(text => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          ws.send(text)
        }
      })
    }
  })

  window.addEventListener('resize', handleResize)

  connectTerminal()
}

onMounted(() => {
  initTerminal()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (ws) {
    ws.close()
  }
  if (terminal) {
    terminal.dispose()
  }
})

defineExpose({
  fit: fitTerminal
})
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  background: #1e1e1e;
  overflow: hidden;
}

:deep(.xterm) {
  height: 100%;
  padding: 8px;
}

:deep(.xterm-viewport) {
  overflow-y: auto !important;
}
</style>
