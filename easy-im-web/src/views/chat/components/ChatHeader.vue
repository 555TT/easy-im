<script setup lang="ts">
import Avatar from '@/components/Avatar.vue'

defineProps<{ name: string; subtitle?: string }>()
</script>

<template>
  <header class="head">
    <Avatar :name="name" :size="36" />
    <div class="meta">
      <div class="name">{{ name }}</div>
      <div v-if="subtitle" class="sub">{{ subtitle }}</div>
    </div>
  </header>
</template>

<style scoped>
.head {
  height: 56px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  border-bottom: 1px solid var(--divider);
  background: var(--content-bg);
}
.meta { display: flex; flex-direction: column; }
.name { font-weight: 600; color: var(--text-primary); }
.sub { font-size: 12px; color: var(--text-secondary); }
</style>
