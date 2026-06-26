<script setup lang="ts">
import { computed, ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import { loadMainShellViewModel } from '@/api/main-shell';
import { navigationAdapter } from '@/adapters/navigation';
import type { MainShellViewModel, OrderStatus } from '@/types/main-shell';

const viewModel = ref<MainShellViewModel | null>(null);
const errorMessage = ref('');
const activeStatus = ref<OrderStatus | 'all'>('all');

const filters: Array<{ id: OrderStatus | 'all'; label: string }> = [
  { id: 'all', label: '全部' },
  { id: 'pending_use', label: '待使用' },
  { id: 'reserved', label: '已预约' },
  { id: 'completed', label: '已完成' },
  { id: 'refunding', label: '退款售后' }
];

const statusLabels: Record<OrderStatus, string> = {
  pending_use: '待使用',
  pending_pay: '待支付',
  reserved: '已预约',
  completed: '已完成',
  refunded: '已退款',
  refunding: '退款中'
};

loadMainShellViewModel()
  .then((result) => {
    viewModel.value = result.data;
    errorMessage.value = result.data.errorMessage ?? '';
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  });

const filteredOrders = computed(() => {
  const orders = (viewModel.value?.orders ?? []).filter((order) => order.status !== 'pending_pay');
  if (activeStatus.value === 'all') return orders;
  if (activeStatus.value === 'refunding') {
    return orders.filter((order) => order.status === 'refunding' || order.status === 'refunded');
  }
  return orders.filter((order) => order.status === activeStatus.value);
});

function openTicket() {
  navigationAdapter.navigateTo('/pages/ticket/index');
}

function openProfile() {
  navigationAdapter.switchTab('/pages/profile/index');
}

function openHome() {
  navigationAdapter.switchTab('/pages/home/index');
}
</script>

<template>
  <view class="orders-page">
    <view class="page-heading">
      <view class="orders-page-head">
        <button class="back-to-mine" aria-label="返回我的" @click="openProfile"><text>‹</text></button>
        <view>
          <text class="page-heading__title">我的订单</text>
          <text class="page-heading__body">查看待使用、已预约、已完成和退款售后订单。</text>
        </view>
      </view>
    </view>

    <view v-if="errorMessage" class="orders-panel">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="orders-panel">
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

      <view v-if="filteredOrders.length === 0" class="order-list">
        <EmptyState title="暂无订单" description="当前状态下没有可展示订单。" />
      </view>
      <view v-else class="order-list">
        <view v-for="order in filteredOrders" :key="order.id" class="order-card">
          <view class="order-card__head">
            <text class="order-card__title">{{ order.title }}</text>
            <text
              class="order-card__status"
              :class="{ 'order-card__status--pending': order.status === 'pending_use' || order.status === 'reserved', 'order-card__status--refund': order.status === 'refunded' || order.status === 'refunding' }"
            >
              {{ statusLabels[order.status] }}
            </text>
          </view>
          <view class="order-card__meta">
            <text>订单号：{{ order.orderNo }}</text>
            <text>出行日：{{ order.travelDateLabel }}</text>
            <text>数量：{{ order.quantityLabel }}</text>
            <text>金额：{{ order.amountLabel }}</text>
          </view>
          <view class="order-card__foot">
            <text>{{ order.hint }}</text>
            <button class="order-detail" @click="openTicket">详情</button>
          </view>
        </view>
      </view>
    </view>

    <view class="visual-bottom-nav" aria-label="底部导航">
      <view class="visual-bottom-nav__item" @click="openHome">
        <text class="visual-bottom-nav__icon">⌂</text>
        <text>首页</text>
      </view>
      <view class="visual-bottom-nav__item">
        <text class="visual-bottom-nav__icon">♬</text>
        <text>客服</text>
      </view>
      <view class="visual-bottom-nav__item visual-bottom-nav__item--active" @click="openProfile">
        <text class="visual-bottom-nav__icon">♙</text>
        <text>我的</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.orders-page {
  width: min(100vw, 430px);
  min-height: 100vh;
  margin: clamp(0px, calc((100vw - 430px) * 100), 28px) auto;
  padding: 14px 0 calc(86px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  overflow-x: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px);
  background: #edf5f7;
  color: #0d2b3a;
  box-shadow:
    0 0 0 1px rgba(9, 50, 66, 0.04),
    0 24px 80px rgba(8, 45, 64, 0.18);
}

.page-heading {
  margin: 0 14px 12px;
  padding: 16px 14px;
  border-radius: 8px;
  color: #ffffff;
  background: linear-gradient(135deg, #0f8b8d, #226bb8);
  box-shadow: 0 12px 30px rgba(9, 50, 66, 0.12);
}

.orders-page-head {
  display: flex;
  align-items: center;
  gap: 10px;
}

.back-to-mine {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  margin: 0;
  border: 0;
  border-radius: 8px;
  color: #ffffff;
  background: rgba(255, 255, 255, 0.18);
  font-size: 28px;
  line-height: 1;
}

.back-to-mine text {
  display: block;
  margin-top: -2px;
}

.back-to-mine::after,
.order-filter__item::after,
.order-detail::after {
  display: none;
}

.page-heading__title,
.page-heading__body {
  display: block;
}

.page-heading__title {
  font-size: 22px;
  font-weight: 900;
  line-height: 1.2;
}

.page-heading__body {
  margin-top: 7px;
  color: rgba(255, 255, 255, 0.86);
  font-size: 13px;
  line-height: 1.45;
}

.orders-panel {
  margin: 0 14px 10px;
  padding: 14px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.order-filter {
  white-space: nowrap;
  margin-bottom: 10px;
}

.order-filter__item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 30px;
  margin: 0 7px 0 0;
  padding: 0 10px;
  border-radius: 8px;
  color: #71808a;
  background: #f2f7f8;
  font-size: 12px;
}

.order-filter__item--active {
  color: #ffffff;
  background: #087393;
  font-weight: 900;
}

.order-list {
  display: grid;
  gap: 9px;
}

.order-card {
  padding: 11px;
  border: 1px solid #dbe7eb;
  border-radius: 8px;
  background: #fbfdfe;
}

.order-card__head,
.order-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.order-card__title {
  min-width: 0;
  font-size: 14px;
  font-weight: 900;
  line-height: 1.25;
}

.order-card__status {
  flex: 0 0 auto;
  padding: 4px 7px;
  border-radius: 8px;
  color: #11623a;
  background: #def8eb;
  font-size: 11px;
  font-weight: 900;
}

.order-card__status--pending {
  color: #7c4b00;
  background: #fff1d1;
}

.order-card__status--refund {
  color: #9b3b47;
  background: #ffe8ec;
}

.order-card__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px 8px;
  margin-top: 9px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.35;
}

.order-card__foot {
  margin-top: 10px;
  color: #71808a;
  font-size: 12px;
}

.order-detail {
  flex: 0 0 auto;
  height: 28px;
  margin: 0;
  padding: 0 10px;
  border-radius: 8px;
  color: #087393;
  background: #e9f7fa;
  font-size: 12px;
  font-weight: 900;
}

.visual-bottom-nav {
  position: fixed;
  z-index: 30;
  left: 50%;
  bottom: 0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  width: min(100vw, 430px);
  transform: translateX(-50%);
  padding: 8px 8px calc(8px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  border-top: 1px solid rgba(9, 50, 66, 0.08);
  border-radius: 0 0 clamp(0px, calc((100vw - 430px) * 100), 26px) clamp(0px, calc((100vw - 430px) * 100), 26px);
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 -10px 26px rgba(14, 48, 64, 0.08);
  backdrop-filter: blur(12px);
}

.visual-bottom-nav__item {
  display: grid;
  place-items: center;
  gap: 3px;
  min-height: 50px;
  color: #66727a;
  font-size: 12px;
}

.visual-bottom-nav__item--active {
  color: #087393;
  font-weight: 900;
}

.visual-bottom-nav__icon {
  font-size: 22px;
  line-height: 1;
}
</style>
