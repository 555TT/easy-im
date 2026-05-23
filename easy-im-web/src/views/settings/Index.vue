<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage, ElRadio, ElRadioGroup } from 'element-plus'
import { applyStoredTheme, setStoredTheme, type ThemeMode } from '@/utils/theme'

const theme = ref<ThemeMode>(applyStoredTheme())

const themeLabel = computed(() => (theme.value === 'dark' ? '深色' : '浅色'))

function changeTheme(value: string | number | boolean | undefined): void {
  if (value !== 'light' && value !== 'dark') return
  theme.value = value
  setStoredTheme(value)
  ElMessage.success(`已切换为${value === 'dark' ? '深色' : '浅色'}主题`)
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

      <div class="panel preview-panel">
        <div class="title">预览</div>
        <div class="preview-shell">
          <div class="preview-rail"></div>
          <div class="preview-main">
            <div class="preview-toolbar"></div>
            <div class="preview-card">
              <div class="preview-line short"></div>
              <div class="preview-line"></div>
              <div class="preview-line"></div>
            </div>
          </div>
        </div>
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
.preview-panel .title {
  margin-bottom: 16px;
}
.preview-shell {
  display: grid;
  grid-template-columns: 64px 1fr;
  min-height: 220px;
  border: 1px solid var(--divider);
  border-radius: 16px;
  overflow: hidden;
  background: var(--content-bg);
}
.preview-rail {
  background: var(--rail-bg);
}
.preview-main {
  padding: 16px;
  background: var(--side-bg);
}
.preview-toolbar {
  height: 18px;
  width: 140px;
  border-radius: 999px;
  background: var(--side-active);
}
.preview-card {
  margin-top: 16px;
  padding: 20px;
  border-radius: 14px;
  background: var(--surface-bg);
  box-shadow: var(--shadow-pop);
}
.preview-line {
  height: 12px;
  border-radius: 999px;
  background: var(--divider);
  margin-top: 12px;
}
.preview-line.short {
  width: 40%;
  margin-top: 0;
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
