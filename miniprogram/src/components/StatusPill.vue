<script setup lang="ts">
defineProps<{
  text: string;
  tone?: 'success' | 'warning' | 'muted';
}>();
</script>

<template>
  <text class="status-pill" :class="`status-pill--${tone ?? 'muted'}`">{{ text }}</text>
</template>

<style scoped>
.status-pill {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 0 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
}

.status-pill--success {
  background: var(--zt-color-primary-weak);
  color: var(--zt-color-primary);
}

.status-pill--warning {
  background: #fff2cc;
  color: var(--zt-color-accent);
}

.status-pill--muted {
  background: #eef2f6;
  color: var(--zt-color-text-muted);
}
</style>
