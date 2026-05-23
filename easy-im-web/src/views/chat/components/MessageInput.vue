<script setup lang="ts">
import { ref } from 'vue'
import { ElInput, ElButton } from 'element-plus'

const emit = defineEmits<{ send: [text: string] }>()
const text = ref('')

function send(): void {
  const v = text.value.trim()
  if (!v) return
  emit('send', v)
  text.value = ''
}

function onKeydown(e: Event): void {
  const ke = e as KeyboardEvent
  if (ke.key === 'Enter' && !ke.shiftKey) {
    ke.preventDefault()
    send()
  }
}
</script>

<template>
  <div class="input">
    <el-input
      v-model="text"
      type="textarea"
      :autosize="{ minRows: 2, maxRows: 5 }"
      resize="none"
      placeholder="输入消息，Enter 发送，Shift+Enter 换行"
      @keydown="onKeydown"
    />
    <div class="actions">
      <el-button type="primary" :disabled="!text.trim()" @click="send">发送</el-button>
    </div>
  </div>
</template>

<style scoped>
.input {
  border-top: 1px solid var(--divider);
  padding: 12px 16px;
  background: var(--content-bg);
}
.actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}
</style>
