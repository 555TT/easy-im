<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import ConversationList from './components/ConversationList.vue'
import ChatPanel from './components/ChatPanel.vue'
import EmptyChat from './components/EmptyChat.vue'
import { useConversationStore } from '@/stores/conversation'

const route = useRoute()
const convo = useConversationStore()

const cid = computed(() => {
  const fromRoute = (route.params.conversationId as string) || ''
  return fromRoute || convo.currentId
})

onMounted(async () => {
  await convo.fetchAll().catch(() => {})
  if (cid.value) convo.setCurrent(cid.value)
})

watch(
  () => route.params.conversationId,
  (v) => {
    const id = (v as string) || ''
    if (id) convo.setCurrent(id)
  },
)
</script>

<template>
  <ConversationList />
  <ChatPanel v-if="cid" :conversation-id="cid" :key="cid" />
  <EmptyChat v-else />
</template>
