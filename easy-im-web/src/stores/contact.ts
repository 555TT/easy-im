import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Friend, FriendRequest } from '@/types/domain'
import * as socialApi from '@/api/social'

export const useContactStore = defineStore('contact', () => {
  const friends = ref<Friend[]>([])
  const requests = ref<FriendRequest[]>([])
  const onlineMap = ref<Record<string, boolean>>({})

  const pendingRequestCount = computed(
    () => requests.value.filter((r) => r.handleResult === 1).length,
  )

  async function fetchAll(): Promise<void> {
    const [fl, ol, reqs] = await Promise.all([
      socialApi.listFriends(),
      socialApi.listOnlineFriends(),
      socialApi.listFriendRequests(),
    ])
    onlineMap.value = ol.onlineList ?? {}
    friends.value = (fl.list ?? []).map((f) => ({
      id: f.id ?? '',
      userId: f.friend_uid ?? '',
      nickname: f.nickname ?? '',
      avatar: f.avatar ?? '',
      remark: f.remark ?? '',
      online: !!(onlineMap.value[f.friend_uid ?? ''] ?? false),
    }))
    requests.value = (reqs.list ?? []).map((r) => ({
      id: r.id ?? '',
      fromUserId: r.req_uid ?? '',
      message: r.req_msg ?? '',
      reqTime: r.req_time ?? 0,
      handleResult: r.handle_result ?? 1,
    }))
  }

  function setOnline(map: Record<string, boolean>): void {
    onlineMap.value = map
    friends.value.forEach((f) => (f.online = !!map[f.userId]))
  }

  return { friends, requests, onlineMap, pendingRequestCount, fetchAll, setOnline }
})
