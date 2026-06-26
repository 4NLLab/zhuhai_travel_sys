<template>
  <view class="page-shell">
    <view class="card empty-card">
      <text class="empty-card__title">订单基础页</text>
      <text class="empty-card__body">Phase 1 仅验证 tabBar 与空态容器。Phase 2 会迁移订单列表和票券基础壳。</text>
      <text class="empty-card__tag">Mock 场景将在业务切片中补齐</text>
    </view>
  </view>
</template>

<style scoped>
.empty-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
}

.empty-card__title {
  font-size: 20px;
  font-weight: 800;
}

.empty-card__body,
.empty-card__tag {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}
</style>
