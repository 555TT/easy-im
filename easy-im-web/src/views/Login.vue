<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus'
import * as userApi from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import { wsClient } from '@/ws/client'
import { ApiError } from '@/api/http'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const phone = ref('')
const password = ref('')
const submitting = ref(false)

async function submit(): Promise<void> {
  if (!phone.value || !password.value) {
    ElMessage.warning('请输入手机号和密码')
    return
  }
  submitting.value = true
  try {
    const resp = await userApi.login({ phone: phone.value, password: password.value })
    auth.setSession(resp.token, resp.expire)
    wsClient.connect(resp.token)
    const redirect = (route.query.redirect as string) || '/chat'
    router.push(redirect)
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '登录失败'
    ElMessage.error(msg)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="auth">
    <div class="card">
      <div class="brand">easy-im</div>
      <div class="subtitle">即时通讯 · 登录</div>
      <el-form @submit.prevent="submit">
        <el-form-item>
          <el-input v-model="phone" placeholder="手机号" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="password" type="password" placeholder="密码" size="large" show-password />
        </el-form-item>
        <el-button type="primary" size="large" :loading="submitting" class="w-full" @click="submit">
          登录
        </el-button>
        <div class="footer">
          没有账号？
          <RouterLink to="/register" class="link">立即注册</RouterLink>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.auth {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e6f4ff 0%, #f5f7fa 100%);
}
.card {
  width: 360px;
  padding: 32px;
  background: #fff;
  border-radius: 12px;
  box-shadow: var(--shadow-pop);
}
.brand {
  font-size: 26px;
  font-weight: 700;
  color: var(--brand-500);
}
.subtitle {
  color: var(--text-secondary);
  margin: 4px 0 24px;
}
.footer {
  text-align: center;
  margin-top: 16px;
  color: var(--text-secondary);
  font-size: 13px;
}
.link {
  color: var(--brand-500);
}
.w-full {
  width: 100%;
}
</style>
