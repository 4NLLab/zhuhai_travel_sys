<script setup lang="ts">
import type { TicketSummary } from '@/types/main-shell';
import StatusPill from './StatusPill.vue';

defineProps<{
  ticket: TicketSummary;
}>();

function ticketStatusText(status: TicketSummary['status']): string {
  const labels: Record<TicketSummary['status'], string> = {
    available: '可使用',
    unavailable: '不可用',
    refunding: '退款中'
  };
  return labels[status];
}
</script>

<template>
  <view class="card ticket-card">
    <view class="ticket-card__visual">
      <view class="ticket-card__qr">
        <text></text>
      </view>
      <view>
        <text class="ticket-card__label">出示核销码检票</text>
        <text class="ticket-card__copy">{{ ticket.verifyLocation }}</text>
      </view>
    </view>
    <view class="ticket-card__head">
      <text class="ticket-card__title">{{ ticket.productTitle }}</text>
      <StatusPill
        :text="ticketStatusText(ticket.status)"
        :tone="ticket.status === 'available' ? 'success' : ticket.status === 'refunding' ? 'warning' : 'muted'"
      />
    </view>
    <view class="ticket-card__grid">
      <view>
        <text>日期</text>
        <text class="ticket-card__value">{{ ticket.validDateLabel }}</text>
      </view>
      <view>
        <text>时段</text>
        <text class="ticket-card__value">{{ ticket.validTimeLabel }}</text>
      </view>
      <view>
        <text>数量</text>
        <text class="ticket-card__value">{{ ticket.quantityLabel }}</text>
      </view>
    </view>
    <view class="ticket-card__code">
      <text>券码</text>
      <text class="ticket-card__value">{{ ticket.maskedCode }}</text>
    </view>
  </view>
</template>

<style scoped>
.ticket-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
}

.ticket-card__visual {
  display: grid;
  grid-template-columns: 76px minmax(0, 1fr);
  align-items: center;
  gap: 12px;
  min-height: 92px;
  padding: 12px;
  border-radius: 8px;
  background: linear-gradient(135deg, #eefaf7, #fff3de);
}

.ticket-card__qr {
  width: 72px;
  height: 72px;
  border: 6px solid #152238;
  border-radius: 8px;
  background:
    linear-gradient(90deg, #152238 12px, transparent 12px 20px, #152238 20px 32px, transparent 32px),
    linear-gradient(#152238 12px, transparent 12px 20px, #152238 20px 32px, transparent 32px),
    #ffffff;
}

.ticket-card__label,
.ticket-card__title,
.ticket-card__value {
  font-weight: 800;
}

.ticket-card__value {
  display: block;
}

.ticket-card__label,
.ticket-card__title {
  display: block;
  font-size: 17px;
}

.ticket-card__copy,
.ticket-card__code text,
.ticket-card__grid text {
  display: block;
  margin-top: 4px;
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.ticket-card__grid .ticket-card__value,
.ticket-card__code .ticket-card__value {
  color: var(--zt-color-text);
}

.ticket-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.ticket-card__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.ticket-card__grid view {
  padding: 10px;
  border-radius: 8px;
  background: var(--zt-color-surface-muted);
}

.ticket-card__code {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--zt-color-border);
}
</style>
