<script setup lang="ts">
import { useContactStore } from '@/stores/contact'

const contactStore = useContactStore()
</script>

<template>
  <div class="contact-page">
    <header class="header p-4">
      <h2>联系人</h2>
    </header>

    <el-tabs>
      <el-tab-pane label="在线用户">
        <div class="user-list">
          <div v-for="user in contactStore.onlineUsers" :key="user.userId" class="user-item">
            <el-avatar :size="48">{{ user.userId.slice(0, 2) }}</el-avatar>
            <div class="user-info">
              <span class="user-id">{{ user.userId }}</span>
              <span class="online-status">{{ user.online ? '在线' : '离线' }}</span>
            </div>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="好友列表">
        <el-empty description="暂无好友" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.contact-page {
  height: 100%;
  background: white;
}

.header {
  border-bottom: 1px solid #e0e0e0;
}

.user-list {
  padding: 16px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.user-item:hover {
  background: #f5f5f5;
}

.user-info {
  margin-left: 12px;
  display: flex;
  flex-direction: column;
}

.user-id {
  font-weight: 500;
}

.online-status {
  font-size: 12px;
  color: #52c41a;
}
</style>
