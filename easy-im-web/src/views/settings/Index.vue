<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElButton, ElMessage } from 'element-plus'
import Avatar from '@/components/Avatar.vue'
import { useAuthStore } from '@/stores/auth'
import { wsClient } from '@/ws/client'
import * as userApi from '@/api/user'
import type { UserInfoDTO } from '@/types/api'
import { ApiError } from '@/api/http'

const router = useRouter()
const auth = useAuthStore()
const info = ref<UserInfoDTO | null>(null)

onMounted(async () => {
  if (!auth.userId) return
  try {
    const resp = await userApi.getUserInfo(auth.userId)
    info.value = resp.user
  } catch (err) {
    if (err instanceof ApiError) ElMessage.error(err.message)
  }
})

function logout(): void {
  auth.clear()
  wsClient.disconnect()
  router.push({ name: 'login' })
}
</script>

<template>
  <section class="settings">
    <header class="head">设置</header>
    <div class="body">
      <Avatar :src="info?.avatar" :name="info?.nickname || auth.userId" :size="80" />
      <div class="name">{{ info?.nickname || auth.userId }}</div>
      <ul class="meta">
        <li><span>用户 ID</span><span>{{ auth.userId }}</span></li>
        <li v-if="info"><span>手机号</span><span>{{ info.mobile }}</span></li>
        <li v-if="info"><span>性别</span><span>{{ info.sex === 1 ? '男' : info.sex === 2 ? '女' : '未知' }}</span></li>
      </ul>
      <el-button type="danger" plain @click="logout">退出登录</el-button>
    </div>
  </section>
</template>

<style scoped>
.settings { flex: 1; display: flex; flex-direction: column; background: var(--content-bg); }
.head {
  height: 56px; padding: 0 20px;
  display: flex; align-items: center;
  border-bottom: 1px solid var(--divider);
  font-weight: 600; color: var(--text-primary);
}
.body {
  flex: 1;
  display: flex; flex-direction: column; align-items: center;
  gap: 16px; padding: 32px;
}
.name { font-size: 18px; font-weight: 600; }
.meta { list-style: none; padding: 0; margin: 0; width: 320px; }
.meta li {
  display: flex; justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid var(--divider);
  font-size: 13px;
}
.meta li span:first-child { color: var(--text-secondary); }
</style>
