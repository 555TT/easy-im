<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useContactStore } from '@/stores/contact'
import { wsClient } from '@/ws/client'
import Avatar from '@/components/Avatar.vue'
import { ElBadge, ElTooltip } from 'element-plus'
import { ChatDotRound, User, Setting, SwitchButton } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const contact = useContactStore()

const tabs = [
  { key: 'chat', label: '消息', icon: ChatDotRound, route: '/chat' },
  { key: 'contacts', label: '联系人', icon: User, route: '/contacts' },
  { key: 'settings', label: '设置', icon: Setting, route: '/settings' },
]

const activeKey = computed(() => {
  if (route.path.startsWith('/contacts')) return 'contacts'
  if (route.path.startsWith('/settings')) return 'settings'
  return 'chat'
})

const isProfileRoute = computed(() => route.path.startsWith('/profile'))

function go(path: string): void {
  router.push(path)
}

function logout(): void {
  auth.clear()
  wsClient.disconnect()
  router.push({ name: 'login' })
}
</script>

<template>
  <aside class="rail">
    <div class="rail-top">
      <button
        v-for="t in tabs"
        :key="t.key"
        class="rail-item"
        :class="{ active: activeKey === t.key }"
        :title="t.label"
        @click="go(t.route)"
      >
        <el-badge
          v-if="t.key === 'contacts' && contact.pendingRequestCount > 0"
          :value="contact.pendingRequestCount"
        >
          <component :is="t.icon" class="icon" />
        </el-badge>
        <component v-else :is="t.icon" class="icon" />
      </button>
    </div>
    <div class="rail-bottom">
      <button
        class="profile-entry"
        :class="{ active: isProfileRoute }"
        title="个人信息"
        @click="go('/profile')"
      >
        <Avatar :name="auth.userId" :size="36" />
      </button>
      <el-tooltip content="退出登录" placement="right">
        <button class="logout-btn" title="退出登录" @click="logout">
          <SwitchButton class="logout-icon" />
        </button>
      </el-tooltip>
    </div>
  </aside>
</template>

<style scoped>
.rail {
  width: 64px;
  background: var(--rail-bg);
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}
.rail-top {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.rail-item {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.72);
  cursor: pointer;
  transition: background 120ms, color 120ms;
}
.rail-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}
.rail-item.active {
  background: var(--rail-active);
  color: #fff;
}
.icon {
  width: 20px;
  height: 20px;
}
.rail-bottom {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
.profile-entry,
.logout-btn {
  border: none;
  background: transparent;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.profile-entry {
  border-radius: 999px;
  outline: 2px solid transparent;
  transition: outline-color 120ms, transform 120ms;
}
.profile-entry:hover {
  transform: translateY(-1px);
}
.profile-entry.active {
  outline-color: var(--rail-active);
}
.logout-btn {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  color: rgba(255, 255, 255, 0.72);
  transition: background 120ms, color 120ms;
}
.logout-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}
.logout-icon {
  width: 16px;
  height: 16px;
}
</style>
