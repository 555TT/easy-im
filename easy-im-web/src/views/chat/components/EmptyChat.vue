<template>
  <div class="empty flex-center">
    <div class="inner">
      <div class="icon i-ep-chat-dot-round" />
      <div class="hint">选择一个会话开始聊天</div>
    </div>
  </div>
</template>

<style scoped>
.empty {
  flex: 1;
  background: var(--content-bg);
  color: var(--text-secondary);
}
.inner {
  text-align: center;
}
.icon {
  font-size: 48px;
  color: var(--brand-400);
  margin-bottom: 12px;
}
.hint {
  font-size: 14px;
}
</style>
