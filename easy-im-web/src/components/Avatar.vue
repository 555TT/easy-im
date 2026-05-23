<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{ src?: string; name?: string; size?: number }>(),
  { size: 36, src: '', name: '' },
)

const initials = computed(() =>
  (props.name || '?').slice(0, 2).toUpperCase(),
)

const style = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  fontSize: `${Math.round(props.size * 0.4)}px`,
}))
</script>

<template>
  <div class="avatar" :style="style">
    <img v-if="src" :src="src" alt="" />
    <span v-else>{{ initials }}</span>
  </div>
</template>

<style scoped>
.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--brand-500);
  color: #fff;
  overflow: hidden;
  font-weight: 600;
  flex-shrink: 0;
}
.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
