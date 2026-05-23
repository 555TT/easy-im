import { ref, onUnmounted } from 'vue'
import { FrameType } from '@/types/im'
import type { WSMessage, Chat, MarkRead, Push } from '@/types/im'

export function useWebSocket() {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const messageHandlers: ((msg: WSMessage) => void)[] = []

  let userId = ''
  let heartbeatTimer: number | null = null
  let reconnectTimer: number | null = null

  // 连接
  function connect(host: string, uid: string) {
    userId = uid
    const url = `${host}/ws?userId=${uid}`

    ws.value = new WebSocket(url)

    ws.value.onopen = () => {
      isConnected.value = true
      startHeartbeat()
      // 连接成功后发送 user.online
      sendMessage('user.online', null)
    }

    ws.value.onclose = () => {
      isConnected.value = false
      stopHeartbeat()
      // 断线重连
      if (!reconnectTimer) {
        reconnectTimer = window.setTimeout(() => {
          reconnectTimer = null
          connect(host, userId)
        }, 3000)
      }
    }

    ws.value.onerror = (err) => {
      console.error('WebSocket error:', err)
    }

    ws.value.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        // ACK 确认不触发业务回调
        if (msg.frameType === 2) {
          return
        }
        messageHandlers.forEach((handler) => handler(msg))
      } catch (e) {
        console.error('Parse message error:', e)
      }
    }
  }

  // 断开
  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    stopHeartbeat()
    ws.value?.close()
    ws.value = null
    isConnected.value = false
  }

  // 发送消息
  function sendMessage(method: string, data: any, needAck = true) {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected')
      return
    }

    const msg: WSMessage = {
      id: generateId(),
      frameType: needAck ? FrameType.FrameData : FrameType.FrameNoAck,
      method,
      data,
    }

    ws.value.send(JSON.stringify(msg))
  }

  // 发送聊天消息
  function sendChat(chat: Chat) {
    sendMessage('conversation.chat', chat)
  }

  // 发送已读标记
  function sendMarkRead(markRead: MarkRead) {
    sendMessage('conversation.markChat', markRead)
  }

  // 发送推送
  function sendPush(push: Push) {
    sendMessage('push', push)
  }

  // 监听消息
  function onMessage(handler: (msg: WSMessage) => void) {
    messageHandlers.push(handler)
    return () => {
      const index = messageHandlers.indexOf(handler)
      if (index > -1) {
        messageHandlers.splice(index, 1)
      }
    }
  }

  // 心跳
  function startHeartbeat() {
    stopHeartbeat()
    heartbeatTimer = window.setInterval(() => {
      if (ws.value?.readyState === WebSocket.OPEN) {
        ws.value.send(JSON.stringify({ frameType: FrameType.FramePing }))
      }
    }, 30000)
  }

  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  // 生成唯一 ID
  function generateId(): string {
    return `${Date.now()}-${Math.random().toString(36).slice(2, 11)}`
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    ws,
    isConnected,
    connect,
    disconnect,
    sendMessage,
    sendChat,
    sendMarkRead,
    sendPush,
    onMessage,
  }
}
