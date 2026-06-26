<script setup lang="ts">
import StatusPill from './StatusPill.vue';
import type { OrderStatus, OrderSummary } from '@/types/main-shell';

defineProps<{
  order: OrderSummary;
}>();

function statusText(status: OrderStatus): string {
  const labels: Record<OrderStatus, string> = {
    pending_use: '待使用',
    pending_pay: '待支付',
    reserved: '已预约',
    completed: '已完成',
    refunded: '已退款',
    refunding: '退款中'
  };
  return labels[status];
}
</script>

<template>
  <view class="card order-card">
    <view class="order-card__head">
      <text class="order-card__title">{{ order.title }}</text>
      <StatusPill
        :text="statusText(order.status)"
        :tone="order.status === 'pending_pay' || order.status === 'refunding' ? 'warning' : 'success'"
      />
    </view>
    <view class="order-card__meta">
      <text>订单号：{{ order.orderNo }}</text>
      <text>日期：{{ order.travelDateLabel }}</text>
      <text>数量：{{ order.quantityLabel }}</text>
      <text>金额：{{ order.amountLabel }}</text>
    </view>
    <view class="order-card__foot">
      <text>{{ order.hint }}</text>
      <slot name="action" />
    </view>
  </view>
</template>

<style scoped>
.order-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.order-card__head,
.order-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.order-card__title {
  min-width: 0;
  font-size: 17px;
  font-weight: 800;
}

.order-card__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px 12px;
  color: var(--zt-color-text-muted);
  font-size: 13px;
  line-height: 20px;
}

.order-card__foot {
  color: var(--zt-color-text-muted);
  font-size: 13px;
}
</style>
