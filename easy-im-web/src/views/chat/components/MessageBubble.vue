<script setup lang="ts">
import type { Message } from '@/types/domain'
import Avatar from '@/components/Avatar.vue'
import { formatHM } from '@/utils/time'

defineProps<{ message: Message; mine: boolean; peerName: string }>()
</script>

<template>
  <div class="row" :class="{ mine }">
    <Avatar v-if="!mine" :name="peerName" :size="32" />
    <div class="col">
      <div class="bubble" :class="{ mine }">
        <span class="text">{{ message.content }}</span>
      </div>
      <div class="meta">
        <span class="time">{{ formatHM(message.sendTime) }}</span>
        <span v-if="mine && message.status === 'sending'" class="status">发送中</span>
        <span v-if="mine && message.status === 'failed'" class="status fail">失败</span>
      </div>
    </div>
    <Avatar v-if="mine" :name="'我'" :size="32" />
  </div>
</template>

<style scoped>
.row {
  display: flex;
  gap: 8px;
  margin-bottom: 14px;
  align-items: flex-end;
}
.row.mine { flex-direction: row-reverse; }
.col { display: flex; flex-direction: column; max-width: 60%; }
.bubble {
  padding: 8px 12px;
  border-radius: 12px;
  background: var(--bubble-theirs);
  color: var(--text-primary);
  line-height: 1.5;
  word-break: break-word;
}
.bubble.mine {
  background: var(--bubble-mine);
  color: #fff;
}
.text { white-space: pre-wrap; }
.meta {
  display: flex;
  gap: 6px;
  margin-top: 4px;
  font-size: 11px;
  color: var(--text-secondary);
}
.row.mine .meta { justify-content: flex-end; }
.status.fail { color: var(--danger); }
</style>
