import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Conversation } from '@/types/domain'
import type { Friend } from '@/types/domain'
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

  /**
   * Derive peerUserId / peerNickname from conversationId + friends list.
   * conversationId format: "userId1_userId2" (sorted). One of the two parts
   * should match a friend of the current user.
   * Call this after both fetchAll (conversations) and contact.fetchAll complete.
   */
  function populatePeerFromFriends(friends: Friend[]): void {
    const friendMap = new Map<string, Friend>()
    for (const f of friends) {
      friendMap.set(f.userId, f)
    }
    list.value.forEach((c) => {
      if (c.peerUserId) return
      const parts = c.conversationId.split('_')
      if (parts.length !== 2) return
      // Try both parts — one of them is the peer (friend)
      const peer = friendMap.get(parts[0]) || friendMap.get(parts[1])
      if (peer) {
        c.peerUserId = peer.userId
        c.peerNickname = peer.nickname
        c.peerAvatar = peer.avatar
      }
    })
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

  return { list, currentId, sorted, current, fetchAll, populatePeerFromFriends, setCurrent, upsert, touch, clearUnread }
})
