<script setup lang="ts">
import { computed, onMounted, onUnmounted, watch } from 'vue'
import ChatHeader from './ChatHeader.vue'
import MessageList from './MessageList.vue'
import MessageInput from './MessageInput.vue'
import { useAuthStore } from '@/stores/auth'
import { useChatStore } from '@/stores/chat'
import { useConversationStore } from '@/stores/conversation'
import { useContactStore } from '@/stores/contact'
import { wsClient } from '@/ws/client'
import { ChatType, MType, type Message } from '@/types/domain'
import type { WSMessage, WSPushPayload } from '@/ws/types'
import { clientMsgId } from '@/utils/id'

const props = defineProps<{ conversationId: string }>()

const auth = useAuthStore()
const chat = useChatStore()
const convo = useConversationStore()
const contact = useContactStore()

const messages = computed<Message[]>(() => chat.list(props.conversationId))

const peer = computed(() => {
  const c = convo.list.find((x) => x.conversationId === props.conversationId)
  if (!c) return { id: '', name: props.conversationId }
  const f = contact.friends.find((x) => x.userId === c.peerUserId)
  return {
    id: c.peerUserId,
    name: f?.remark || f?.nickname || c.peerNickname || c.peerUserId,
  }
})

let off: (() => void) | null = null

async function load(): Promise<void> {
  if (!props.conversationId) return
  await chat.fetchHistory(props.conversationId).catch(() => {})
  convo.clearUnread(props.conversationId)
}

function onChatPush(msg: WSMessage): void {
  if (msg.method !== 'push') return
  const p = msg.data as WSPushPayload | undefined
  if (!p || p.conversationId !== props.conversationId) return

  if (p.sendId === auth.userId) {
    const matched = chat.replace(
      props.conversationId,
      (m) =>
        m.status === 'sending' &&
        m.content === p.content &&
        Math.abs(m.sendTime - p.sendTime) < 10_000,
      {
        msgId: p.msgId,
        conversationId: p.conversationId,
        chatType: (p.chatType ?? ChatType.Single) as 1 | 2,
        sendId: p.sendId,
        recvId: p.recvId,
        content: p.content,
        mType: (p.mType ?? MType.Text) as 0,
        sendTime: p.sendTime,
        status: 'sent',
      },
    )
    if (matched) return
  }

  chat.append(props.conversationId, {
    msgId: p.msgId,
    conversationId: p.conversationId,
    chatType: (p.chatType ?? ChatType.Single) as 1 | 2,
    sendId: p.sendId,
    recvId: p.recvId,
    content: p.content,
    mType: (p.mType ?? MType.Text) as 0,
    sendTime: p.sendTime,
    status: 'sent',
  })
  convo.touch(p.conversationId, p.content, p.sendTime, p.sendId !== auth.userId)
}

function onSend(text: string): void {
  if (!props.conversationId || !peer.value.id) return
  const cmid = clientMsgId()
  const sendTime = Date.now()
  chat.append(props.conversationId, {
    msgId: '',
    clientMsgId: cmid,
    conversationId: props.conversationId,
    chatType: ChatType.Single,
    sendId: auth.userId,
    recvId: peer.value.id,
    content: text,
    mType: MType.Text,
    sendTime,
    status: 'sending',
  })
  convo.touch(props.conversationId, text, sendTime, false)
  wsClient.send('conversation.chat', {
    conversationId: props.conversationId,
    chatType: ChatType.Single,
    sendId: auth.userId,
    recvId: peer.value.id,
    sendTime,
    msg: {
      msgId: '',
      mType: MType.Text,
      content: text,
      readRecords: {},
    },
  })

  window.setTimeout(() => {
    chat.setStatus(props.conversationId, cmid, 'failed')
  }, 10_000)
}

onMounted(() => {
  off = wsClient.on('chat', onChatPush)
  load()
})
onUnmounted(() => { off?.() })
watch(() => props.conversationId, () => { load() })
</script>

<template>
  <section class="panel">
    <ChatHeader :name="peer.name" />
    <MessageList :messages="messages" :self-id="auth.userId" :peer-name="peer.name" />
    <MessageInput @send="onSend" />
  </section>
</template>

<style scoped>
.panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--content-bg);
  min-width: 0;
}
</style>
