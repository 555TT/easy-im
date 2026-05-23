<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  type FormRules,
} from 'element-plus'
import Avatar from '@/components/Avatar.vue'
import * as userApi from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import type { UserInfoDTO } from '@/types/api'
import { ApiError } from '@/api/http'

const auth = useAuthStore()
const info = ref<UserInfoDTO | null>(null)
const loading = ref(false)
const saving = ref(false)
const editing = ref(false)
const formRef = ref<InstanceType<typeof ElForm>>()

const form = reactive({
  nickname: '',
  sex: 0,
  email: '',
  avatar: '',
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const rules: FormRules<typeof form> = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback) => {
        if (!value?.trim()) callback(new Error('昵称不能为空'))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  sex: [
    {
      validator: (_rule, value: number, callback) => {
        if ([0, 1, 2].includes(value)) callback()
        else callback(new Error('请选择合法的性别'))
      },
      trigger: 'change',
    },
  ],
  email: [
    {
      validator: (_rule, value: string, callback) => {
        if (!value || emailRegex.test(value.trim())) callback()
        else callback(new Error('邮箱格式不正确'))
      },
      trigger: 'blur',
    },
  ],
  avatar: [
    {
      validator: (_rule, value: string, callback) => {
        if (!value) return callback()
        try {
          new URL(value)
          callback()
        } catch {
          callback(new Error('头像请填写合法的 URL'))
        }
      },
      trigger: 'blur',
    },
  ],
  oldPassword: [
    {
      validator: (_rule, value: string, callback) => {
        if (!form.newPassword && !value) return callback()
        if (!value) return callback(new Error('请输入当前密码'))
        callback()
      },
      trigger: 'blur',
    },
  ],
  newPassword: [
    {
      validator: (_rule, value: string, callback) => {
        if (!form.oldPassword && !value) return callback()
        if (!value) return callback(new Error('请输入新密码'))
        if (value.trim().length < 6) return callback(new Error('新密码至少 6 位'))
        if (value === form.oldPassword) return callback(new Error('新密码不能与旧密码相同'))
        callback()
      },
      trigger: 'blur',
    },
  ],
  confirmPassword: [
    {
      validator: (_rule, value: string, callback) => {
        if (!form.newPassword && !value) return callback()
        if (!value) return callback(new Error('请再次输入新密码'))
        if (value !== form.newPassword) return callback(new Error('两次输入的新密码不一致'))
        callback()
      },
      trigger: 'blur',
    },
  ],
}

const displayName = computed(() => info.value?.nickname || auth.userId)
const displaySex = computed(() => {
  if (!info.value) return '未知'
  return info.value.sex === 1 ? '男' : info.value.sex === 2 ? '女' : '未知'
})

function fillForm(user: UserInfoDTO): void {
  form.nickname = user.nickname
  form.sex = user.sex
  form.email = user.email
  form.avatar = user.avatar
  form.oldPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
}

async function loadProfile(): Promise<void> {
  if (!auth.userId) return
  loading.value = true
  try {
    const resp = await userApi.getUserInfo()
    info.value = resp.user
    fillForm(resp.user)
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '加载个人资料失败'
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}

onMounted(loadProfile)

function startEdit(): void {
  if (!info.value) return
  fillForm(info.value)
  editing.value = true
}

function cancelEdit(): void {
  if (info.value) fillForm(info.value)
  formRef.value?.clearValidate()
  editing.value = false
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const resp = await userApi.updateProfile({
      nickname: form.nickname.trim(),
      sex: form.sex,
      email: form.email.trim(),
      avatar: form.avatar.trim(),
      old_password: form.oldPassword,
      new_password: form.newPassword,
    })
    if (resp.success) {
      await loadProfile()
      editing.value = false
      formRef.value?.clearValidate()
      ElMessage.success('个人资料已更新')
    }
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '保存失败'
    ElMessage.error(msg)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="profile-page">
    <header class="head">个人信息</header>
    <div v-loading="loading" class="body">
      <div class="panel">
        <div class="profile-header">
          <Avatar :src="editing ? form.avatar : info?.avatar" :name="displayName" :size="80" />
          <div class="summary">
            <div class="name">{{ displayName }}</div>
            <div class="sub">{{ auth.userId }}</div>
          </div>
          <div class="actions">
            <template v-if="editing">
              <el-button @click="cancelEdit">取消</el-button>
              <el-button type="primary" :loading="saving" @click="submit">保存资料</el-button>
            </template>
            <el-button v-else type="primary" @click="startEdit">编辑资料</el-button>
          </div>
        </div>

        <template v-if="!editing">
          <ul class="meta">
            <li><span>用户 ID</span><span>{{ auth.userId }}</span></li>
            <li v-if="info"><span>手机号</span><span>{{ info.mobile }}</span></li>
            <li><span>昵称</span><span>{{ info?.nickname || '-' }}</span></li>
            <li><span>邮箱</span><span>{{ info?.email || '-' }}</span></li>
            <li><span>性别</span><span>{{ displaySex }}</span></li>
            <li><span>头像地址</span><span class="text-cell">{{ info?.avatar || '-' }}</span></li>
          </ul>
        </template>

        <el-form v-else ref="formRef" :model="form" :rules="rules" label-position="top" class="form">
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="form.nickname" maxlength="24" show-word-limit placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input :model-value="info?.mobile || ''" disabled />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="选填，示例：name@example.com" />
          </el-form-item>
          <el-form-item label="性别" prop="sex">
            <el-radio-group v-model="form.sex">
              <el-radio :value="0">未知</el-radio>
              <el-radio :value="1">男</el-radio>
              <el-radio :value="2">女</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="头像地址" prop="avatar">
            <el-input v-model="form.avatar" placeholder="选填，填写可访问的图片 URL" />
          </el-form-item>
          <el-form-item label="当前密码" prop="oldPassword">
            <el-input v-model="form.oldPassword" type="password" show-password placeholder="如需改密请填写当前密码" />
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="form.newPassword" type="password" show-password placeholder="不少于 6 位；不修改可留空" />
          </el-form-item>
          <el-form-item label="确认新密码" prop="confirmPassword">
            <el-input v-model="form.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
          </el-form-item>
        </el-form>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-page {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--content-bg);
}
.head {
  height: 56px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--divider);
  font-weight: 600;
  color: var(--text-primary);
}
.body {
  flex: 1;
  padding: 24px;
  overflow: auto;
}
.panel {
  background: var(--surface-bg);
  border-radius: 14px;
  box-shadow: var(--shadow-pop);
  padding: 24px;
}
.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}
.summary {
  min-width: 0;
}
.name {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}
.sub {
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 13px;
}
.actions {
  margin-left: auto;
  display: flex;
  gap: 12px;
}
.meta {
  list-style: none;
  padding: 0;
  margin: 0;
}
.meta li {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--divider);
  font-size: 14px;
}
.meta li:last-child {
  border-bottom: none;
}
.meta li span:first-child {
  color: var(--text-secondary);
}
.text-cell {
  word-break: break-all;
}
.form {
  max-width: 520px;
}
@media (max-width: 768px) {
  .body {
    padding: 16px;
  }
  .panel {
    padding: 16px;
  }
  .profile-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .actions {
    margin-left: 0;
    width: 100%;
    flex-wrap: wrap;
  }
  .meta li {
    grid-template-columns: 1fr;
    gap: 6px;
  }
}
</style>
