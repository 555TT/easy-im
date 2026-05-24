import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Message, MessageStatus } from '@/types/domain'
import * as imApi from '@/api/im'
import { ChatType, MType } from '@/types/domain'
import { useAuthStore } from '@/stores/auth'

export const useChatStore = defineStore('chat', () => {
  const byConv = ref<Record<string, Message[]>>({})

  function list(conversationId: string): Message[] {
    return byConv.value[conversationId] ?? []
  }

  function append(conversationId: string, m: Message): void {
    const arr = (byConv.value[conversationId] ??= [])
    arr.push(m)
  }

  function replace(
    conversationId: string,
    predicate: (m: Message) => boolean,
    next: Message,
  ): boolean {
    const arr = byConv.value[conversationId]
    if (!arr) return false
    const idx = arr.findIndex(predicate)
    if (idx === -1) return false
    arr[idx] = next
    return true
  }

  function setStatus(
    conversationId: string,
    msgId: string,
    status: MessageStatus,
  ): void {
    const arr = byConv.value[conversationId]
    if (!arr) return
    const m = arr.find((x) => x.msgId === msgId || x.clientMsgId === msgId)
    if (m) m.status = status
  }

  async function fetchHistory(conversationId: string, count = 50): Promise<void> {
    const auth = useAuthStore()
    const resp = await imApi.listChatLog({
      userId: auth.userId,
      conversationId,
      endSendTime: Date.now(),
      count,
    })
    const items: Message[] = (resp.list ?? []).map((dto) => ({
      msgId: dto.id ?? '',
      conversationId: dto.conversationId ?? conversationId,
      chatType: (dto.chatType ?? ChatType.Single) as 1 | 2,
      sendId: dto.sendId ?? '',
      recvId: dto.recvId ?? '',
      content: dto.msgContent ?? '',
      mType: (dto.msgType ?? MType.Text) as 0,
      sendTime: dto.sendTime ?? 0,
      status: 'sent',
    }))
    items.sort((a, b) => a.sendTime - b.sendTime)
    byConv.value[conversationId] = items
  }

  return { byConv, list, append, replace, setStatus, fetchHistory }
})
