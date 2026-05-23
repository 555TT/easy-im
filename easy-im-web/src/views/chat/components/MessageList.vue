<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import MessageBubble from './MessageBubble.vue'
import type { Message } from '@/types/domain'
import { useAutoScroll } from '@/composables/useAutoScroll'

const props = defineProps<{
  messages: Message[]
  selfId: string
  peerName: string
}>()

const containerRef = ref<HTMLElement | null>(null)
const { onScroll, scrollToBottom } = useAutoScroll(containerRef)

watch(() => props.messages.length, () => { scrollToBottom() })
onMounted(() => { scrollToBottom(true) })
</script>

<template>
  <div ref="containerRef" class="list" @scroll="onScroll">
    <MessageBubble
      v-for="m in messages"
      :key="m.msgId || m.clientMsgId"
      :message="m"
      :mine="m.sendId === selfId"
      :peer-name="peerName"
    />
  </div>
</template>

<style scoped>
.list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: var(--content-bg);
}
</style>
