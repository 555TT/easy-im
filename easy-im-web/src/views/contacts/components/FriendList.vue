<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElButton } from 'element-plus'
import Avatar from '@/components/Avatar.vue'
import { useContactStore } from '@/stores/contact'

const router = useRouter()
const contact = useContactStore()
const open = defineModel<boolean>('addOpen', { default: false })

const items = computed(() => contact.friends)
const pendingCount = computed(() => contact.pendingRequestCount)

function go(id: string): void {
  router.push(`/contacts/${id}`)
}
function showRequests(): void {
  router.push('/contacts/requests')
}
</script>

<template>
  <aside class="side">
    <header class="head">
      <span class="title">联系人</span>
      <el-button size="small" type="primary" @click="open = true">添加</el-button>
    </header>

    <button class="item special" @click="showRequests">
      <div class="icon-box">🛎️</div>
      <div class="body">
        <div class="name">好友申请</div>
        <div class="sub">{{ pendingCount > 0 ? `${pendingCount} 个待处理` : '暂无待处理' }}</div>
      </div>
      <span v-if="pendingCount > 0" class="badge">{{ pendingCount }}</span>
    </button>

    <div class="section">我的好友</div>
    <div class="list">
      <button
        v-for="f in items"
        :key="f.id"
        class="item"
        @click="go(f.userId)"
      >
        <Avatar :src="f.avatar" :name="f.nickname || f.userId" :size="36" />
        <div class="body">
          <div class="name">{{ f.remark || f.nickname || f.userId }}</div>
          <div class="sub">
            <span class="dot" :class="{ on: f.online }" />
            {{ f.online ? '在线' : '离线' }}
          </div>
        </div>
      </button>
      <div v-if="items.length === 0" class="empty">暂无好友</div>
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title { font-size: 18px; font-weight: 600; color: var(--text-primary); }
.section {
  padding: 8px 16px;
  font-size: 12px;
  color: var(--text-secondary);
}
.list { flex: 1; overflow-y: auto; padding: 0 8px 8px; }
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
  position: relative;
}
.item:hover { background: var(--side-hover); }
.item.special { margin: 0 8px; }
.icon-box {
  width: 36px; height: 36px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  background: #fff7e6; font-size: 18px;
}
.body { flex: 1; min-width: 0; }
.name { font-weight: 600; color: var(--text-primary); }
.sub { font-size: 12px; color: var(--text-secondary); display: flex; gap: 6px; align-items: center; }
.dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: #c9cdd4;
}
.dot.on { background: var(--success); }
.badge {
  position: absolute; right: 12px; top: 50%; transform: translateY(-50%);
  background: var(--danger); color: #fff;
  min-width: 18px; height: 18px; border-radius: 9px;
  font-size: 11px; padding: 0 6px;
  display: inline-flex; align-items: center; justify-content: center;
}
.empty { text-align: center; color: var(--text-secondary); padding: 24px 0; font-size: 13px; }
</style>
