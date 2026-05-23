<script setup lang="ts">
import { ref } from 'vue'
import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
import * as socialApi from '@/api/social'
import { ApiError } from '@/api/http'

const open = defineModel<boolean>('open', { default: false })
const phone = ref('')
const msg = ref('')
const submitting = ref(false)

// 中国大陆手机号：11 位，1 开头，第二位 3-9
const phoneRegex = /^1[3-9]\d{9}$/

async function submit(): Promise<void> {
  if (!phone.value) {
    ElMessage.warning('请输入对方手机号')
    return
  }
  if (!phoneRegex.test(phone.value)) {
    ElMessage.warning('手机号格式不正确')
    return
  }
  submitting.value = true
  try {
    await socialApi.sendFriendRequest({
      phone: phone.value,
      req_msg: msg.value,
      req_time: Date.now(),
    })
    ElMessage.success('申请已发送')
    open.value = false
    phone.value = ''
    msg.value = ''
  } catch (err) {
    const m = err instanceof ApiError ? err.message : '发送失败'
    ElMessage.error(m)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <el-dialog v-model="open" title="添加好友" width="380px">
    <el-form>
      <el-form-item label="对方手机号">
        <el-input v-model="phone" placeholder="请输入对方手机号" maxlength="11" />
      </el-form-item>
      <el-form-item label="附言">
        <el-input v-model="msg" type="textarea" :rows="2" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="submit">发送</el-button>
    </template>
  </el-dialog>
</template>
