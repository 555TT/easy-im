import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Chat, Conversation } from '@/types/im'

export const useChatStore = defineStore('chat', () => {
  const conversations = ref<Conversation[]>([])
  const currentConversation = ref<Conversation | null>(null)
  const messages = ref<Map<string, Chat[]>>(new Map())

  // 添加或更新会话
  function addOrUpdateConversation(conv: Conversation) {
    const index = conversations.value.findIndex((c) => c.conversationId === conv.conversationId)
    if (index > -1) {
      conversations.value[index] = { ...conversations.value[index], ...conv }
    } else {
      conversations.value.unshift(conv)
    }
  }

  // 添加消息
  function addMessage(conversationId: string, chat: Chat) {
    if (!messages.value.has(conversationId)) {
      messages.value.set(conversationId, [])
    }
    messages.value.get(conversationId)!.push(chat)
  }

  // 获取消息
  function getMessages(conversationId: string): Chat[] {
    return messages.value.get(conversationId) || []
  }

  // 设置当前会话
  function setCurrentConversation(conv: Conversation | null) {
    currentConversation.value = conv
  }

  return {
    conversations,
    currentConversation,
    messages,
    addOrUpdateConversation,
    addMessage,
    getMessages,
    setCurrentConversation,
  }
})
