<script setup lang="ts">
import { storageAdapter } from '@/adapters/storage';

function clearSession() {
  storageAdapter.clearSession();
}
</script>

<template>
  <view class="page-shell">
    <view class="card profile-card">
      <text class="profile-card__title">我的</text>
      <text class="profile-card__body">Phase 1 只保留最小用户摘要和退出登录存储清理能力。</text>
      <button class="secondary-button" @click="clearSession">清理本地会话</button>
    </view>
  </view>
</template>

<style scoped>
.profile-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
}

.profile-card__title {
  font-size: 20px;
  font-weight: 800;
}

.profile-card__body {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}
</style>
