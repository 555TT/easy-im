<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElButton, ElMessage } from 'element-plus'
import Avatar from '@/components/Avatar.vue'
import { useContactStore } from '@/stores/contact'
import { useAuthStore } from '@/stores/auth'
import * as imApi from '@/api/im'
import { ChatType } from '@/types/domain'
import { ApiError } from '@/api/http'

const props = defineProps<{ friendId: string }>()
const router = useRouter()
const contact = useContactStore()
const auth = useAuthStore()

const friend = computed(() => contact.friends.find((f) => f.userId === props.friendId) ?? null)

function buildSingleConvId(a: string, b: string): string {
  const [x, y] = [a, b].sort((p, q) => (BigInt(p) < BigInt(q) ? -1 : 1))
  return `${x}_${y}`
}

async function startChat(): Promise<void> {
  if (!friend.value) return
  try {
    await imApi.setUpConversation({
      sendId: auth.userId,
      recvId: friend.value.userId,
      chatType: ChatType.Single,
    })
    const cid = buildSingleConvId(auth.userId, friend.value.userId)
    router.push(`/chat/${cid}`)
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '无法发起会话'
    ElMessage.error(msg)
  }
}
</script>

<template>
  <section class="detail">
    <div v-if="!friend" class="empty">请选择一位联系人</div>
    <div v-else class="body">
      <Avatar :src="friend.avatar" :name="friend.nickname || friend.userId" :size="72" />
      <div class="name">{{ friend.remark || friend.nickname || friend.userId }}</div>
      <div class="meta">
        <span>ID: {{ friend.userId }}</span>
        <span class="state">
          <span class="dot" :class="{ on: friend.online }" />
          {{ friend.online ? '在线' : '离线' }}
        </span>
      </div>
      <el-button type="primary" @click="startChat">发起会话</el-button>
    </div>
  </section>
</template>

<style scoped>
.detail { flex: 1; display: flex; align-items: center; justify-content: center; background: var(--content-bg); }
.empty { color: var(--text-secondary); }
.body { display: flex; flex-direction: column; align-items: center; gap: 12px; }
.name { font-size: 18px; font-weight: 600; color: var(--text-primary); }
.meta { display: flex; gap: 16px; color: var(--text-secondary); font-size: 13px; }
.state { display: inline-flex; gap: 6px; align-items: center; }
.dot { width: 8px; height: 8px; border-radius: 50%; background: #c9cdd4; }
.dot.on { background: var(--success); }
</style>
