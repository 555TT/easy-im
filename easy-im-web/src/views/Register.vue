<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage, ElRadioGroup, ElRadio } from 'element-plus'
import * as userApi from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import { wsClient } from '@/ws/client'
import { ApiError } from '@/api/http'

const router = useRouter()
const auth = useAuthStore()

const phone = ref('')
const password = ref('')
const nickname = ref('')
const sex = ref<number>(1)
const submitting = ref(false)

async function submit(): Promise<void> {
  if (!phone.value || !password.value || !nickname.value) {
    ElMessage.warning('请填写完整')
    return
  }
  submitting.value = true
  try {
    const resp = await userApi.register({
      phone: phone.value,
      password: password.value,
      nickname: nickname.value,
      sex: sex.value,
      avatar: '',
    })
    auth.setSession(resp.token, resp.expire)
    wsClient.connect(resp.token)
    router.push('/chat')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '注册失败'
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
      <div class="subtitle">注册新账号</div>
      <el-form @submit.prevent="submit">
        <el-form-item>
          <el-input v-model="phone" placeholder="手机号" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="password" type="password" placeholder="密码" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-input v-model="nickname" placeholder="昵称" size="large" />
        </el-form-item>
        <el-form-item>
          <el-radio-group v-model="sex">
            <el-radio :value="1">男</el-radio>
            <el-radio :value="2">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-button type="primary" size="large" :loading="submitting" class="w-full" @click="submit">
          注册
        </el-button>
        <div class="footer">
          已有账号？
          <RouterLink to="/login" class="link">返回登录</RouterLink>
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
