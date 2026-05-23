<script setup lang="ts">
import { computed } from 'vue'
import { ElButton, ElMessage } from 'element-plus'
import Avatar from '@/components/Avatar.vue'
import { useContactStore } from '@/stores/contact'
import * as socialApi from '@/api/social'
import { ApiError } from '@/api/http'
import { formatRelative } from '@/utils/time'

const contact = useContactStore()
const items = computed(() => contact.requests)

async function handle(id: string, result: number): Promise<void> {
  try {
    await socialApi.handleFriendRequest({ friend_req_id: id, handle_result: result })
    ElMessage.success('已处理')
    contact.fetchAll()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '处理失败'
    ElMessage.error(msg)
  }
}
</script>

<template>
  <section class="detail">
    <header class="head">好友申请</header>
    <div v-if="items.length === 0" class="empty">暂无好友申请</div>
    <ul v-else class="list">
      <li v-for="r in items" :key="r.id" class="row">
        <Avatar :name="r.fromUserId" :size="40" />
        <div class="body">
          <div class="name">{{ r.fromUserId }}</div>
          <div class="msg">{{ r.message || '请求添加好友' }}</div>
          <div class="time">{{ formatRelative(r.reqTime) }}</div>
        </div>
        <div v-if="r.handleResult === 1" class="actions">
          <el-button size="small" type="primary" @click="handle(r.id, 2)">通过</el-button>
          <el-button size="small" @click="handle(r.id, 3)">拒绝</el-button>
        </div>
        <div v-else class="result">
          {{ r.handleResult === 2 ? '已通过' : r.handleResult === 3 ? '已拒绝' : '已取消' }}
        </div>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.detail { flex: 1; display: flex; flex-direction: column; background: var(--content-bg); }
.head {
  height: 56px; display: flex; align-items: center;
  padding: 0 20px; border-bottom: 1px solid var(--divider);
  font-weight: 600; color: var(--text-primary);
}
.empty { padding: 48px; text-align: center; color: var(--text-secondary); }
.list { list-style: none; padding: 8px; margin: 0; flex: 1; overflow-y: auto; }
.row {
  display: flex; gap: 12px; align-items: center;
  padding: 12px; border-radius: 8px;
}
.row:hover { background: var(--side-hover); }
.body { flex: 1; min-width: 0; }
.name { font-weight: 600; color: var(--text-primary); }
.msg { color: var(--text-secondary); font-size: 13px; margin-top: 2px; }
.time { color: var(--text-secondary); font-size: 12px; margin-top: 2px; }
.actions { display: flex; gap: 8px; }
.result { color: var(--text-secondary); font-size: 13px; }
</style>
