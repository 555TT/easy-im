<script setup lang="ts">
import { onMounted } from 'vue'
import NavRail from '@/components/NavRail.vue'
import { wsClient } from '@/ws/client'
import { useAuthStore } from '@/stores/auth'
import { useContactStore } from '@/stores/contact'

const auth = useAuthStore()
const contact = useContactStore()

onMounted(() => {
  if (auth.token) wsClient.connect(auth.token)
  contact.fetchAll().catch(() => { /* surfaced in views */ })

  wsClient.on('online', (ids) => {
    const map: Record<string, boolean> = {}
    ids.forEach((id) => (map[id] = true))
    contact.setOnline(map)
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
