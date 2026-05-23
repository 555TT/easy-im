<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useConversationStore } from '@/stores/conversation'
import Avatar from '@/components/Avatar.vue'
import { formatRelative } from '@/utils/time'

const router = useRouter()
const store = useConversationStore()

const items = computed(() => store.sorted)

function open(id: string): void {
  store.setCurrent(id)
  store.clearUnread(id)
  router.push(`/chat/${id}`)
}
</script>

<template>
  <aside class="side">
    <header class="head">消息</header>
    <div class="list">
      <button
        v-for="c in items"
        :key="c.conversationId"
        class="item"
        :class="{ active: store.currentId === c.conversationId }"
        @click="open(c.conversationId)"
      >
        <Avatar :name="c.peerNickname || c.peerUserId" :size="40" />
        <div class="body">
          <div class="row1">
            <span class="name">{{ c.peerNickname || c.peerUserId || c.conversationId }}</span>
            <span class="time">{{ formatRelative(c.lastTime) }}</span>
          </div>
          <div class="row2">
            <span class="preview">{{ c.lastContent }}</span>
            <span v-if="c.unread > 0" class="badge">{{ c.unread > 99 ? '99+' : c.unread }}</span>
          </div>
        </div>
      </button>
      <div v-if="items.length === 0" class="empty">暂无会话</div>
    </div>
  </aside>
</template>

<style scoped>
.side {
  width: 280px;
  background: var(--side-bg);
  border-right: 1px solid var(--divider);
  display: flex;
  flex-direction: column;
}
.head {
  padding: 20px 16px 12px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}
.list {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px 8px;
}
.item {
  width: 100%;
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 10px 12px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  text-align: left;
}
.item:hover { background: var(--side-hover); }
.item.active { background: var(--side-active); }
.body { flex: 1; min-width: 0; }
.row1 {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.name {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.time {
  font-size: 12px;
  color: var(--text-secondary);
  flex-shrink: 0;
}
.row2 {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 2px;
}
.preview {
  font-size: 13px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.badge {
  min-width: 18px;
  height: 18px;
  padding: 0 6px;
  background: var(--danger);
  color: #fff;
  border-radius: 9px;
  font-size: 11px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.empty {
  text-align: center;
  color: var(--text-secondary);
  padding: 24px 0;
  font-size: 13px;
}
</style>
