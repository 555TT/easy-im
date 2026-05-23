<script setup lang="ts">
import { computed, ref, reactive } from 'vue'
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
import { applyStoredTheme, setStoredTheme, type ThemeMode } from '@/utils/theme'
import * as userApi from '@/api/user'
import { ApiError } from '@/api/http'

const theme = ref<ThemeMode>(applyStoredTheme())

const themeLabel = computed(() => (theme.value === 'dark' ? '深色' : '浅色'))

function changeTheme(value: string | number | boolean | undefined): void {
  if (value !== 'light' && value !== 'dark') return
  theme.value = value
  setStoredTheme(value)
  ElMessage.success(`已切换为${value === 'dark' ? '深色' : '浅色'}主题`)
}

// Password change
const passwordFormRef = ref<InstanceType<typeof ElForm>>()
const passwordSaving = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const passwordRules: FormRules<typeof passwordForm> = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '新密码至少 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback) => {
        if (value !== passwordForm.newPassword) callback(new Error('两次输入的新密码不一致'))
        else callback()
      },
      trigger: 'blur',
    },
  ],
}

async function submitPassword(): Promise<void> {
  const valid = await passwordFormRef.value?.validate().catch(() => false)
  if (!valid) return

  passwordSaving.value = true
  try {
    const resp = await userApi.updatePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword,
    })
    if (resp.success) {
      ElMessage.success('密码修改成功')
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
      passwordFormRef.value?.clearValidate()
    }
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '修改失败'
    ElMessage.error(msg)
  } finally {
    passwordSaving.value = false
  }
}
</script>

<template>
  <section class="settings-page">
    <header class="head">界面设置</header>
    <div class="body">
      <div class="panel">
        <div class="panel-head">
          <div>
            <div class="title">主题模式</div>
            <div class="desc">这里仅控制应用界面外观，不包含个人资料信息。</div>
          </div>
          <div class="current">当前：{{ themeLabel }}</div>
        </div>

        <el-radio-group :model-value="theme" class="theme-group" @update:model-value="changeTheme">
          <el-radio value="light" class="theme-card">
            <div>
              <div class="card-title">浅色</div>
              <div class="card-desc">明亮清爽，适合白天使用</div>
            </div>
          </el-radio>
          <el-radio value="dark" class="theme-card">
            <div>
              <div class="card-title">深色</div>
              <div class="card-desc">降低刺眼感，适合夜间使用</div>
            </div>
          </el-radio>
        </el-radio-group>
      </div>

      <div class="panel">
        <div class="title">修改密码</div>
        <el-form
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          label-position="top"
          class="password-form"
        >
          <el-form-item label="当前密码" prop="oldPassword">
            <el-input v-model="passwordForm.oldPassword" type="password" show-password placeholder="请输入当前密码" />
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="不少于 6 位" />
          </el-form-item>
          <el-form-item label="确认新密码" prop="confirmPassword">
            <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
          </el-form-item>
          <el-button type="primary" :loading="passwordSaving" @click="submitPassword">修改密码</el-button>
        </el-form>
      </div>
    </div>
  </section>
</template>

<style scoped>
.settings-page {
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
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow: auto;
}
.panel {
  background: var(--surface-bg);
  border-radius: 14px;
  box-shadow: var(--shadow-pop);
  padding: 24px;
}
.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}
.title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}
.desc,
.current {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 13px;
}
.theme-group {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, minmax(220px, 1fr));
  gap: 16px;
}
:deep(.theme-card) {
  margin-right: 0;
  border: 1px solid var(--divider);
  border-radius: 12px;
  padding: 16px;
  background: var(--content-bg);
}
:deep(.theme-card .el-radio__label) {
  display: block;
  width: 100%;
  color: inherit;
}
.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}
.card-desc {
  margin-top: 8px;
  color: var(--text-secondary);
  font-size: 13px;
}
.password-form {
  max-width: 400px;
  margin-top: 16px;
}
@media (max-width: 768px) {
  .body {
    padding: 16px;
  }
  .panel {
    padding: 16px;
  }
  .panel-head {
    flex-direction: column;
  }
  .theme-group {
    grid-template-columns: 1fr;
  }
}
</style>
