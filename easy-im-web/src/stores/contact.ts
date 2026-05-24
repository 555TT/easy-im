import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Friend, FriendRequest } from '@/types/domain'
import * as socialApi from '@/api/social'
import { mapFriendDTO } from '@/utils/contact'
import { buildOnlineMap } from '@/utils/online'

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
    friends.value = (fl.list ?? []).map((f) => mapFriendDTO(f, onlineMap.value))
    requests.value = (reqs.list ?? []).map((r) => ({
      id: r.id ?? '',
      fromUserId: r.req_uid ?? '',
      message: r.req_msg ?? '',
      reqTime: r.req_time ?? 0,
      handleResult: r.handle_result ?? 1,
    }))
  }

  async function refreshOnline(): Promise<void> {
    const ol = await socialApi.listOnlineFriends()
    setOnline(ol.onlineList ?? {})
  }

  function setOnlineUsers(userIds: string[]): void {
    setOnline(buildOnlineMap(friends.value.map((f) => f.userId), userIds))
  }

  function setOnline(map: Record<string, boolean>): void {
    onlineMap.value = map
    friends.value.forEach((f) => (f.online = !!map[f.userId]))
  }

  return { friends, requests, onlineMap, pendingRequestCount, fetchAll, refreshOnline, setOnlineUsers, setOnline }
})
