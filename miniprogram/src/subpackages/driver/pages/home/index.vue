<template>
  <view class="page-shell">
    <view class="card driver-card">
      <text class="driver-card__title">司机端分包占位</text>
      <text class="driver-card__body">Phase 4 会迁移登录、审核状态、推广二维码、钱包、佣金、提现和提现记录。</text>
      <text class="driver-card__meta">当前只验证分包路由与包体预算。</text>
    </view>
  </view>
</template>

<style scoped>
.driver-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
}

.driver-card__title {
  font-size: 20px;
  font-weight: 800;
}

.driver-card__body,
.driver-card__meta {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}
</style>
