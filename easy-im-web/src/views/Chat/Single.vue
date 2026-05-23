<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'
import { ChatType, type Chat } from '@/types/im'

const props = defineProps<{
  userId: string
}>()

const userStore = useUserStore()
const chatStore = useChatStore()
const { sendChat } = useWebSocket()

const messageText = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

const conversationId = computed(() => {
  return [userStore.userId, props.userId].sort().join('-')
})

const messages = computed(() => chatStore.getMessages(conversationId.value))

watch(
  () => messages.value.length,
  async () => {
    await nextTick()
    scrollToBottom()
  }
)

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function handleSend() {
  if (!messageText.value.trim()) {
    return
  }

  const chat: Chat = {
    conversationId: conversationId.value,
    chatType: ChatType.SingleChatType,
    sendId: userStore.userId,
    recvId: props.userId,
    sendTime: Date.now(),
    msg: {
      msgId: '',
      mType: 0,
      content: messageText.value,
      readRecords: {},
    },
  }

  sendChat(chat)
  chatStore.addMessage(conversationId.value, chat)
  messageText.value = ''
}

function formatTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div class="single-chat">
    <!-- 聊天头部 -->
    <header class="chat-header flex-between p-4">
      <span class="font-bold">{{ userId }}</span>
    </header>

    <!-- 消息列表 -->
    <div ref="messagesContainer" class="messages-container">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        class="message-item"
        :class="{ mine: msg.sendId === userStore.userId, theirs: msg.sendId !== userStore.userId }"
      >
        <div class="message-content">
          <p class="message-text">{{ msg.msg.content }}</p>
          <span class="message-time">{{ formatTime(msg.sendTime) }}</span>
        </div>
      </div>
    </div>

    <!-- 输入框 -->
    <footer class="chat-footer p-4">
      <el-input
        v-model="messageText"
        placeholder="输入消息..."
        @keyup.enter="handleSend"
      >
        <template #append>
          <el-button @click="handleSend">发送</el-button>
        </template>
      </el-input>
    </footer>
  </div>
</template>

<style scoped>
.single-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  border-bottom: 1px solid #e0e0e0;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
}

.message-item.mine {
  justify-content: flex-end;
}

.message-item.theirs {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: 8px;
  position: relative;
}

.mine .message-content {
  background: #1890ff;
  color: white;
}

.theirs .message-content {
  background: #f0f0f0;
}

.message-text {
  word-break: break-word;
}

.message-time {
  font-size: 10px;
  opacity: 0.7;
  display: block;
  text-align: right;
  margin-top: 4px;
}

.chat-footer {
  border-top: 1px solid #e0e0e0;
}
</style>
