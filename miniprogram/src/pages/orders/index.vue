<script setup lang="ts">
import { computed, ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import OrderSummaryCard from '@/components/OrderSummaryCard.vue';
import TicketSummaryCard from '@/components/TicketSummaryCard.vue';
import { loadMainShellViewModel } from '@/api/main-shell';
import { navigationAdapter } from '@/adapters/navigation';
import type { MainShellViewModel, OrderStatus } from '@/types/main-shell';

const viewModel = ref<MainShellViewModel | null>(null);
const errorMessage = ref('');
const activeStatus = ref<OrderStatus | 'all'>('all');

const filters: Array<{ id: OrderStatus | 'all'; label: string }> = [
  { id: 'all', label: '全部' },
  { id: 'pending_use', label: '待使用' },
  { id: 'pending_pay', label: '待支付' },
  { id: 'reserved', label: '已预约' },
  { id: 'completed', label: '已完成' },
  { id: 'refunded', label: '已退款' },
  { id: 'refunding', label: '退款售后' }
];

loadMainShellViewModel()
  .then((result) => {
    viewModel.value = result.data;
    errorMessage.value = result.data.errorMessage ?? '';
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  });

const filteredOrders = computed(() => {
  const orders = viewModel.value?.orders ?? [];
  if (activeStatus.value === 'all') return orders;
  return orders.filter((order) => order.status === activeStatus.value);
});

function openTicket() {
  navigationAdapter.navigateTo('/pages/ticket/index');
}
</script>

<template>
  <view class="page-shell orders-page">
    <view class="section page-head">
      <view>
        <text class="page-title">我的订单</text>
        <text class="page-subtitle">通用订单与环岛游订单分字段展示，真实用户接口留到 Phase 5 联调。</text>
      </view>
    </view>

    <view v-if="errorMessage" class="section">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="section">
      <scroll-view class="order-filter" scroll-x>
        <button
          v-for="filter in filters"
          :key="filter.id"
          class="order-filter__item"
          :class="{ 'order-filter__item--active': activeStatus === filter.id }"
          @click="activeStatus = filter.id"
        >
          {{ filter.label }}
        </button>
      </scroll-view>
    </view>

    <view class="section">
      <text class="section-title">订单列表</text>
      <view v-if="filteredOrders.length === 0" class="order-list">
        <EmptyState title="暂无订单" description="当前 Mock 场景没有匹配订单，可切换 phase2-success 查看完整列表。" />
      </view>
      <view v-else class="order-list">
        <OrderSummaryCard v-for="order in filteredOrders" :key="order.id" :order="order">
          <template v-if="order.ticketId" #action>
            <button class="mini-button" @click="openTicket">票券</button>
          </template>
        </OrderSummaryCard>
      </view>
    </view>

    <view class="section">
      <text class="section-title">票券基础壳</text>
      <view v-if="(viewModel?.tickets ?? []).length === 0" class="ticket-list">
        <EmptyState title="暂无可用票券" description="未登录、空态或接口失败时票券列表保持可解释状态。" />
      </view>
      <view v-else class="ticket-list">
        <TicketSummaryCard
          v-for="ticket in viewModel?.tickets ?? []"
          :key="ticket.id"
          :ticket="ticket"
          @click="openTicket"
        />
      </view>
    </view>
  </view>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.page-title,
.page-subtitle {
  display: block;
}

.page-title {
  font-size: 28px;
  font-weight: 900;
}

.page-subtitle {
  margin-top: 8px;
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.order-filter {
  white-space: nowrap;
}

.order-filter__item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 76px;
  height: 36px;
  margin-right: 8px;
  padding: 0 14px;
  border-radius: 8px;
  color: var(--zt-color-text-muted);
  background: var(--zt-color-surface);
}

.order-filter__item--active {
  color: #ffffff;
  background: var(--zt-color-primary);
}

.order-list,
.ticket-list {
  display: grid;
  gap: 12px;
  margin-top: 12px;
}

.mini-button {
  min-width: 56px;
  height: 32px;
  border-radius: 8px;
  color: #ffffff;
  background: var(--zt-color-primary);
  font-size: 13px;
  font-weight: 800;
}
</style>
