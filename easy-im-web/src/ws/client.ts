import { FrameType, type WSMessage } from './types'
import { clientMsgId } from '@/utils/id'

type EventMap = {
  open: void
  close: void
  online: string[]
  chat: WSMessage
  markRead: WSMessage
  error: string
}
type Listener<K extends keyof EventMap> = (payload: EventMap[K]) => void

class WSClient {
  private ws: WebSocket | null = null
  private token: string | null = null
  private heartbeatTimer: number | null = null
  private reconnectTimer: number | null = null
  private retries = 0
  private manualClose = false
  private listeners: { [K in keyof EventMap]?: Set<Listener<K>> } = {}

  connect(token: string): void {
    if (this.ws && this.ws.readyState <= 1 && this.token === token) return
    this.token = token
    this.manualClose = false
    this.openSocket()
  }

  disconnect(): void {
    this.manualClose = true
    this.clearTimers()
    if (this.ws) {
      try { this.ws.close() } catch { /* ignore */ }
    }
    this.ws = null
    this.token = null
  }

  send<T>(method: string, data: T, needAck = true): string {
    const id = clientMsgId()
    const msg: WSMessage<T> = {
      id,
      frameType: needAck ? FrameType.Data : FrameType.NoAck,
      method,
      data,
    }
    this.raw(JSON.stringify(msg))
    return id
  }

  on<K extends keyof EventMap>(event: K, fn: Listener<K>): () => void {
    const set = (this.listeners[event] ??= new Set()) as Set<Listener<K>>
    set.add(fn)
    return () => set.delete(fn)
  }

  private emit<K extends keyof EventMap>(event: K, payload: EventMap[K]): void {
    this.listeners[event]?.forEach((fn) => (fn as Listener<K>)(payload))
  }

  private raw(text: string): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return
    this.ws.send(text)
  }

  private openSocket(): void {
    if (!this.token) return
    const ws = new WebSocket('/ws', [this.token])
    this.ws = ws

    ws.onopen = () => {
      this.retries = 0
      this.startHeartbeat()
      this.send('user.online', null, false)
      this.emit('open', undefined)
    }

    ws.onclose = () => {
      this.clearTimers()
      this.emit('close', undefined)
      if (!this.manualClose) this.scheduleReconnect()
    }

    ws.onerror = () => {
      this.emit('error', 'socket error')
    }

    ws.onmessage = (evt) => {
      let msg: WSMessage
      try { msg = JSON.parse(evt.data) as WSMessage } catch { return }
      if (msg.frameType === FrameType.Ack) return
      if (msg.frameType === FrameType.Error) {
        this.emit('error', String(msg.data ?? 'error'))
        return
      }
      switch (msg.method) {
        case 'user.online':
          if (Array.isArray(msg.data)) this.emit('online', msg.data as string[])
          break
        case 'conversation.chat':
        case 'push':
          this.emit('chat', msg)
          break
        case 'conversation.markChat':
          this.emit('markRead', msg)
          break
      }
    }
  }

  private startHeartbeat(): void {
    this.heartbeatTimer = window.setInterval(() => {
      this.raw(JSON.stringify({ frameType: FrameType.Ping }))
    }, 30_000)
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer != null) return
    const delay = Math.min(30_000, 1_000 * 2 ** this.retries)
    this.retries += 1
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.openSocket()
    }, delay)
  }

  private clearTimers(): void {
    if (this.heartbeatTimer != null) {
      window.clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
    if (this.reconnectTimer != null) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }
}

export const wsClient = new WSClient()
