import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Conversation } from '@/types/domain'
import * as imApi from '@/api/im'
import { ChatType } from '@/types/domain'

export const useConversationStore = defineStore('conversation', () => {
  const list = ref<Conversation[]>([])
  const currentId = ref<string>('')

  const sorted = computed(() =>
    [...list.value].sort((a, b) => b.lastTime - a.lastTime),
  )
  const current = computed(() =>
    list.value.find((c) => c.conversationId === currentId.value) ?? null,
  )

  async function fetchAll(): Promise<void> {
    const resp = await imApi.listConversations()
    const map = resp.conversationList ?? {}
    list.value = Object.entries(map).map(([cid, c]) => ({
      conversationId: cid,
      chatType: (c.chatType ?? ChatType.Single) as 1 | 2,
      peerUserId: '',
      peerNickname: '',
      peerAvatar: '',
      lastContent: '',
      lastTime: 0,
      unread: 0,
    }))
  }

  function setCurrent(id: string): void {
    currentId.value = id
  }

  function upsert(c: Conversation): void {
    const idx = list.value.findIndex((x) => x.conversationId === c.conversationId)
    if (idx === -1) list.value.push(c)
    else list.value[idx] = { ...list.value[idx], ...c }
  }

  function touch(
    conversationId: string,
    lastContent: string,
    lastTime: number,
    incrementUnread: boolean,
  ): void {
    const c = list.value.find((x) => x.conversationId === conversationId)
    if (!c) return
    c.lastContent = lastContent
    c.lastTime = lastTime
    if (incrementUnread) c.unread += 1
  }

  function clearUnread(conversationId: string): void {
    const c = list.value.find((x) => x.conversationId === conversationId)
    if (c) c.unread = 0
  }

  return { list, currentId, sorted, current, fetchAll, setCurrent, upsert, touch, clearUnread }
})
