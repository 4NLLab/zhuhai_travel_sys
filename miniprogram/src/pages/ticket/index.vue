<script setup lang="ts">
import { ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import TicketSummaryCard from '@/components/TicketSummaryCard.vue';
import { loadTicketDetail } from '@/api/main-shell';
import type { TicketSummary } from '@/types/main-shell';

const ticket = ref<TicketSummary | null>(null);
const errorMessage = ref('');

function ticketStatusLabel(currentTicket: TicketSummary): string {
  const labels: Record<TicketSummary['status'], string> = {
    available: '待使用',
    unavailable: '不可用',
    refunding: '退款中'
  };
  return labels[currentTicket.status];
}

loadTicketDetail()
  .then((result) => {
    ticket.value = result.data;
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  });
</script>

<template>
  <view class="page-shell ticket-page">
    <view class="ticket-state">
      <text class="ticket-state__label">{{ ticket ? ticketStatusLabel(ticket) : '票券详情' }}</text>
      <text class="ticket-state__body">票券详情视觉参考自 ticket.html，二维码使用小程序兼容静态绘制 fallback。</text>
    </view>

    <EmptyState v-if="errorMessage" title="票券读取失败" :description="errorMessage" />
    <EmptyState v-else-if="!ticket" title="暂无票券" description="当前 Mock 场景没有可展示票券。" />
    <template v-else>
      <TicketSummaryCard :ticket="ticket" />

      <view class="card notice-card">
        <text class="notice-card__title">使用须知</text>
        <view v-for="notice in ticket.notice" :key="notice" class="notice-card__item">
          <text>{{ notice }}</text>
        </view>
      </view>

      <view class="card order-info-card">
        <text class="notice-card__title">订单信息</text>
        <view class="info-row">
          <text>实际付款</text>
          <text class="info-row__value">{{ ticket.amountLabel }}</text>
        </view>
        <view class="info-row">
          <text>订单编号</text>
          <text class="info-row__value">{{ ticket.orderNo }}</text>
        </view>
        <view class="info-row">
          <text>订单类型</text>
          <text class="info-row__value">{{ ticket.source === 'island_cruise' ? '环岛游订单' : '通用订单' }}</text>
        </view>
      </view>
    </template>
  </view>
</template>

<style scoped>
.ticket-page {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ticket-state {
  padding: 8px 2px;
}

.ticket-state__label,
.ticket-state__body {
  display: block;
}

.ticket-state__label {
  font-size: 32px;
  font-weight: 900;
}

.ticket-state__body {
  margin-top: 8px;
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.notice-card,
.order-info-card {
  padding: 16px;
}

.notice-card__title {
  display: block;
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 900;
}

.notice-card__item {
  padding: 8px 0;
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 0;
  border-top: 1px solid var(--zt-color-border);
  font-size: 14px;
}

.info-row text {
  color: var(--zt-color-text-muted);
}

.info-row .info-row__value {
  color: var(--zt-color-text);
  font-weight: 800;
  text-align: right;
}
</style>
