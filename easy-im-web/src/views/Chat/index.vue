<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useContactStore } from '@/stores/contact'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'
import { ChatType, type Chat, type WSMessage } from '@/types/im'
import ChatSingle from './Single.vue'

const router = useRouter()
const userStore = useUserStore()
const contactStore = useContactStore()
const chatStore = useChatStore()

const { connect, disconnect, sendChat, onMessage } = useWebSocket()

const activeTab = ref('chat')
const currentChatUser = ref('')
const wsHost = ref('ws://localhost:8080')

const unsubMessage = ref<(() => void) | null>(null)

onMounted(() => {
  // 如果未登录，重定向到登录页
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  // 建立连接
  connect(wsHost.value, userStore.userId)

  // 监听消息
  unsubMessage.value = onMessage(handleMessage)
})

onUnmounted(() => {
  unsubMessage.value?.()
  disconnect()
})

function handleMessage(msg: WSMessage) {
  // user.online 返回在线用户列表
  if (msg.method === 'user.online' && Array.isArray(msg.data)) {
    contactStore.setOnlineUsers(msg.data as string[])
  }

  // conversation.chat 收到聊天消息
  if (msg.method === 'conversation.chat' || msg.method === 'push') {
    const chat = msg.data as Chat
    chatStore.addMessage(chat.conversationId, chat)
    chatStore.addOrUpdateConversation({
      conversationId: chat.conversationId,
      chatType: chat.chatType,
      lastMsg: chat.msg,
      lastMsgTime: chat.sendTime,
      targetUser: { userId: chat.sendId },
    })
  }
}

function selectChat(userId: string) {
  currentChatUser.value = userId
  const conversationId = generateConversationId(userStore.userId, userId)
  chatStore.setCurrentConversation({
    conversationId,
    chatType: ChatType.SingleChatType,
    targetUser: { userId },
  })
}

function generateConversationId(uid1: string, uid2: string): string {
  return [uid1, uid2].sort().join('-')
}

function handleLogout() {
  disconnect()
  userStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="chat-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="sidebar-header flex-between p-4">
        <span class="font-bold">{{ userStore.userId }}</span>
        <el-button link @click="handleLogout">退出</el-button>
      </div>

      <el-tabs v-model="activeTab" class="sidebar-tabs">
        <el-tab-pane label="聊天" name="chat">
          <div class="user-list">
            <div
              v-for="user in contactStore.onlineUsers"
              :key="user.userId"
              class="user-item"
              :class="{ active: currentChatUser === user.userId }"
              @click="selectChat(user.userId)"
            >
              <el-avatar :size="40">{{ user.userId.slice(0, 2) }}</el-avatar>
              <span class="ml-3">{{ user.userId }}</span>
              <span v-if="user.online" class="online-dot" />
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="联系人" name="contact">
          <div class="user-list">
            <div
              v-for="user in contactStore.onlineUsers"
              :key="user.userId"
              class="user-item"
              @click="selectChat(user.userId)"
            >
              <el-avatar :size="40">{{ user.userId.slice(0, 2) }}</el-avatar>
              <span class="ml-3">{{ user.userId }}</span>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </aside>

    <!-- 聊天区域 -->
    <main class="chat-main">
      <ChatSingle v-if="currentChatUser" :user-id="currentChatUser" />
      <div v-else class="empty-state flex-center">
        <el-empty description="选择联系人开始聊天" />
      </div>
    </main>
  </div>
</template>

<style scoped>
.chat-layout {
  display: flex;
  height: 100%;
}

.sidebar {
  width: 280px;
  background: #f5f5f5;
  border-right: 1px solid #e0e0e0;
}

.sidebar-header {
  border-bottom: 1px solid #e0e0e0;
}

.sidebar-tabs {
  height: calc(100% - 60px);
}

.sidebar-tabs :deep(.el-tabs__content) {
  height: 100%;
  overflow-y: auto;
}

.user-list {
  padding: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  transition: background 0.2s;
}

.user-item:hover,
.user-item.active {
  background: #e0e0e0;
}

.online-dot {
  position: absolute;
  right: 10px;
  width: 8px;
  height: 8px;
  background: #52c41a;
  border-radius: 50%;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.empty-state {
  flex: 1;
}
</style>
