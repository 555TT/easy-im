<script setup lang="ts">
import { onMounted } from 'vue'
import NavRail from '@/components/NavRail.vue'
import { wsClient } from '@/ws/client'
import { useAuthStore } from '@/stores/auth'
import { useContactStore } from '@/stores/contact'
import { useConversationStore } from '@/stores/conversation'

const auth = useAuthStore()
const contact = useContactStore()
const convo = useConversationStore()

onMounted(() => {
  if (auth.token) wsClient.connect(auth.token)
  contact.fetchAll().then(() => {
    // Derive peer info from friends list for existing conversations
    convo.populatePeerFromFriends(contact.friends)
  }).catch(() => { /* surfaced in views */ })

  wsClient.on('online', () => {
    contact.refreshOnline().catch(() => {})
  })
})
</script>

<template>
  <div class="layout">
    <NavRail />
    <RouterView />
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  height: 100%;
}
</style>
