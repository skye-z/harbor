import type { WSMessage } from '../../types'

type WSMessageHandler = (message: WSMessage) => void

class WSClient {
  private ws: WebSocket | null = null
  private reconnectTimer: number | null = null
  private messageHandlers: Set<WSMessageHandler> = new Set()
  private url: string

  constructor(url: string = '/api/ws') {
    this.url = url
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
        try {
          const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
          const host = window.location.host
          // 通过查询参数传递 token（如果有的话），以便后端 AuthMiddleware 能够校验
          const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
          const wsUrl = `${protocol}//${host}${this.url}${token ? `?token=${encodeURIComponent(token)}` : ''}`

          this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log('WebSocket connected')
          this.clearReconnectTimer()
          resolve()
        }

        this.ws.onclose = () => {
          console.log('WebSocket closed')
          this.scheduleReconnect()
        }

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error)
          reject(error)
        }

        this.ws.onmessage = (event) => {
          try {
            const message: WSMessage = JSON.parse(event.data)
            this.messageHandlers.forEach(handler => handler(message))
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error)
          }
        }
      } catch (error) {
        reject(error)
      }
    })
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.clearReconnectTimer()
  }

  send(message: WSMessage) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
    }
  }

  onMessage(handler: WSMessageHandler) {
    this.messageHandlers.add(handler)
    return () => {
      this.messageHandlers.delete(handler)
    }
  }

  private scheduleReconnect() {
    if (!this.reconnectTimer) {
      this.reconnectTimer = window.setTimeout(() => {
        console.log('Attempting to reconnect WebSocket...')
        this.connect().catch(error => {
          console.error('Reconnection failed:', error)
          this.scheduleReconnect()
        })
      }, 3000)
    }
  }

  private clearReconnectTimer() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }
}

export const wsClient = new WSClient()

export function useWebSocket() {
  const connect = () => wsClient.connect()
  const disconnect = () => wsClient.disconnect()
  const send = (message: WSMessage) => wsClient.send(message)
  const onMessage = (handler: WSMessageHandler) => wsClient.onMessage(handler)

  return {
    connect,
    disconnect,
    send,
    onMessage
  }
}
