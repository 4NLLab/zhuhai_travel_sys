<script setup lang="ts">
defineProps<{
  title: string;
  description: string;
}>();
</script>

<template>
  <view class="empty-state">
    <text class="empty-state__title">{{ title }}</text>
    <text class="empty-state__description">{{ description }}</text>
  </view>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 20px;
  border: 1px dashed var(--zt-color-border);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.66);
}

.empty-state__title {
  font-size: 16px;
  font-weight: 800;
}

.empty-state__description {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}
</style>
