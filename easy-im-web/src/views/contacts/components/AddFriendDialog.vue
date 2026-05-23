<script setup lang="ts">
import { ref } from 'vue'
import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
import * as socialApi from '@/api/social'
import { ApiError } from '@/api/http'

const open = defineModel<boolean>('open', { default: false })
const userId = ref('')
const msg = ref('')
const submitting = ref(false)

async function submit(): Promise<void> {
  if (!userId.value) {
    ElMessage.warning('请输入用户 ID')
    return
  }
  submitting.value = true
  try {
    await socialApi.sendFriendRequest({
      user_id: userId.value,
      req_msg: msg.value,
      req_time: Date.now(),
    })
    ElMessage.success('申请已发送')
    open.value = false
    userId.value = ''
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
      <el-form-item label="对方用户 ID">
        <el-input v-model="userId" placeholder="对方 user_id" />
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
