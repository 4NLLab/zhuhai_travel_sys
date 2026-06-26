<script setup lang="ts">
import { ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import StatusPill from '@/components/StatusPill.vue';
import { loadMainShellViewModel } from '@/api/main-shell';
import { storageAdapter } from '@/adapters/storage';
import { navigationAdapter } from '@/adapters/navigation';
import type { MainShellViewModel } from '@/types/main-shell';

const viewModel = ref<MainShellViewModel | null>(null);
const errorMessage = ref('');

loadMainShellViewModel()
  .then((result) => {
    viewModel.value = result.data;
    errorMessage.value = result.data.errorMessage ?? '';
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  });

function clearSession() {
  storageAdapter.clearSession();
}

function openOrders() {
  navigationAdapter.switchTab('/pages/orders/index');
}

function openDriver() {
  navigationAdapter.navigateTo('/subpackages/driver/pages/home/index');
}
</script>

<template>
  <view class="page-shell profile-page">
    <view class="section">
      <text class="profile-title">我的</text>
      <text class="profile-subtitle">查看票券、订单状态、常用出行人与售后服务。</text>
    </view>

    <view v-if="viewModel?.user" class="card profile-card">
      <view class="profile-card__top">
        <text class="avatar">{{ viewModel.user.displayName.slice(0, 1) }}</text>
        <view class="profile-card__identity">
          <text class="profile-card__name">{{ viewModel.user.displayName }}</text>
          <text class="profile-card__mobile">手机号：{{ viewModel.user.maskedMobile }}</text>
        </view>
        <StatusPill text="已实名" tone="success" />
      </view>
      <view v-if="viewModel.user.nextTripLabel" class="next-trip">
        <text>{{ viewModel.user.nextTripLabel }}</text>
      </view>
      <view class="profile-metrics">
        <view>
          <text>{{ viewModel.user.availableTicketCount }}</text>
          <text class="profile-metrics__label">可用票券</text>
        </view>
        <view>
          <text>{{ viewModel.user.recentOrderCount }}</text>
          <text class="profile-metrics__label">近30天订单</text>
        </view>
      </view>
    </view>
    <EmptyState v-else title="未登录" description="请先登录后查看个人资料、订单和票券。" />

    <view v-if="errorMessage" class="section">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="section card ticket-entry" @click="openOrders">
      <view>
        <text class="entry-title">我的票券</text>
        <text class="entry-body">待使用、退款中和不可用票券统一从订单/票券页进入。</text>
      </view>
      <text class="entry-action">查看</text>
    </view>

    <view class="section card order-entry" @click="openOrders">
      <view>
        <text class="entry-title">我的订单</text>
        <text class="entry-body">查看待使用、已预约、已完成和退款售后订单。</text>
      </view>
      <text class="entry-action">进入</text>
    </view>

    <view class="section">
      <text class="section-title">常用服务</text>
      <view class="service-grid">
        <view class="card service-item">
          <text>常用出行人</text>
          <text class="service-item__body">实名信息、证件管理</text>
        </view>
        <view class="card service-item">
          <text>发票抬头</text>
          <text class="service-item__body">企业抬头、开票记录</text>
        </view>
        <view class="card service-item">
          <text>客服与协议</text>
          <text class="service-item__body">售后进度、服务条款</text>
        </view>
        <view class="card service-item" @click="openDriver">
          <text>司机入口</text>
          <text class="service-item__body">登录、钱包、佣金、提现和推广码</text>
        </view>
      </view>
    </view>

    <button class="secondary-button logout-button" @click="clearSession">清理本地会话</button>
  </view>
</template>

<style scoped>
.profile-title,
.profile-subtitle,
.entry-title,
.entry-body,
.service-item text {
  display: block;
}

.profile-title {
  font-size: 28px;
  font-weight: 900;
}

.profile-subtitle,
.entry-body,
.service-item__body,
.profile-card__mobile {
  margin-top: 6px;
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.profile-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 18px;
}

.profile-card__top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: var(--zt-color-primary);
  color: #ffffff;
  font-size: 22px;
  font-weight: 900;
}

.profile-card__identity {
  flex: 1;
  min-width: 0;
}

.profile-card__name {
  display: block;
  font-size: 18px;
  font-weight: 900;
}

.next-trip {
  padding: 12px;
  border-radius: 8px;
  background: var(--zt-color-surface-muted);
  color: var(--zt-color-primary);
  font-size: 14px;
  line-height: 22px;
}

.profile-metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.profile-metrics view {
  padding: 12px;
  border-radius: 8px;
  background: #f8fafc;
}

.profile-metrics text {
  display: block;
  font-size: 24px;
  font-weight: 900;
}

.profile-metrics .profile-metrics__label {
  display: block;
  color: var(--zt-color-text-muted);
  font-size: 13px;
  font-weight: 400;
}

.ticket-entry,
.order-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
}

.entry-title,
.service-item text {
  font-size: 16px;
  font-weight: 900;
}

.entry-action {
  color: var(--zt-color-primary);
  font-size: 14px;
  font-weight: 900;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.service-item {
  padding: 14px;
}

.logout-button {
  width: 100%;
  margin-top: 20px;
}

@media (max-width: 520px) {
  .service-grid {
    grid-template-columns: 1fr;
  }
}
</style>
