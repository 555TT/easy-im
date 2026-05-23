<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useWebSocket } from '@/composables/useWebSocket'
import { useContactStore } from '@/stores/contact'

const router = useRouter()
const userStore = useUserStore()
const contactStore = useContactStore()
const { connect, isConnected } = useWebSocket()

const userId = ref('')
const wsHost = ref('ws://localhost:8080')

function handleLogin() {
  if (!userId.value.trim()) {
    return
  }
  userStore.login(userId.value)
  connect(wsHost.value, userId.value)
  router.push('/chat')
}
</script>

<template>
  <div class="login-container flex-center">
    <el-card class="login-card">
      <template #header>
        <h2 class="text-center">IM 登录</h2>
      </template>

      <el-form label-width="80px">
        <el-form-item label="WebSocket">
          <el-input v-model="wsHost" placeholder="WebSocket 地址" />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="userId" placeholder="请输入用户ID" @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="w-full" @click="handleLogin">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-container {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
}

.text-center {
  text-align: center;
}

.w-full {
  width: 100%;
}
</style>
