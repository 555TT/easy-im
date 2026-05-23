import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/im'

export const useContactStore = defineStore('contact', () => {
  const onlineUsers = ref<User[]>([])
  const friends = ref<User[]>([])

  // 更新在线用户列表
  function setOnlineUsers(users: string[]) {
    onlineUsers.value = users.map((uid) => ({
      userId: uid,
      online: true,
    }))
  }

  // 添加在线用户
  function addOnlineUser(uid: string) {
    if (!onlineUsers.value.find((u) => u.userId === uid)) {
      onlineUsers.value.push({ userId: uid, online: true })
    }
  }

  // 移除在线用户
  function removeOnlineUser(uid: string) {
    const index = onlineUsers.value.findIndex((u) => u.userId === uid)
    if (index > -1) {
      onlineUsers.value.splice(index, 1)
    }
  }

  // 设置好友列表
  function setFriends(users: User[]) {
    friends.value = users
  }

  return {
    onlineUsers,
    friends,
    setOnlineUsers,
    addOnlineUser,
    removeOnlineUser,
    setFriends,
  }
})
