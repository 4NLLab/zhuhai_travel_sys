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
  navigationAdapter.navigateTo('/pages/orders/index');
}

function openHome() {
  navigationAdapter.switchTab('/pages/home/index');
}

function openService() {
  navigationAdapter.switchTab('/pages/service/index');
}
</script>

<template>
  <view class="profile-page">
    <view class="page-heading">
      <text class="page-heading__title">我的</text>
      <text class="page-heading__body">查看我的票券、订单状态、出行人和售后服务。</text>
    </view>

    <view v-if="viewModel?.user" class="profile-card">
      <view class="profile-top">
        <text class="avatar">♙</text>
        <view class="profile-identity">
          <text class="profile-name">{{ viewModel.user.displayName }}</text>
          <text class="profile-mobile">手机号：{{ viewModel.user.maskedMobile }}</text>
        </view>
        <StatusPill text="已实名" tone="success" />
      </view>
      <view v-if="viewModel.user.nextTripLabel" class="relation-box">
        <text class="relation-box__title">{{ viewModel.user.nextTripLabel }}</text>
        <text class="relation-box__body">请按出行提醒提前到达码头或景区。</text>
      </view>
      <view class="rule-list">
        <view class="rule-item">
          <text class="rule-item__icon">▣</text>
          <text>可用票券：{{ viewModel.user.availableTicketCount }} 张</text>
        </view>
        <view class="rule-item">
          <text class="rule-item__icon">☷</text>
          <text>近30天订单：{{ viewModel.user.recentOrderCount }} 笔</text>
        </view>
      </view>
    </view>
    <EmptyState v-else title="未登录" description="请先登录后查看个人资料、订单和票券。" />

    <view v-if="errorMessage" class="admin-panel">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="admin-panel">
      <view class="section-title">
        <text class="section-title__text">我的票券</text>
        <text class="section-title__link">票券夹</text>
      </view>
      <view class="kpi-grid">
        <button class="kpi-card" @click="openOrders">
          <view class="kpi-card__head">
            <text>待使用</text>
            <text>▤</text>
          </view>
          <text class="kpi-card__value">{{ viewModel?.user?.availableTicketCount ?? 0 }}</text>
          <text class="kpi-card__trend">船票、门票可核销</text>
        </button>
        <button class="kpi-card" @click="openOrders">
          <view class="kpi-card__head">
            <text>退款中</text>
            <text>↺</text>
          </view>
          <text class="kpi-card__value">1</text>
          <text class="kpi-card__trend">售后处理中</text>
        </button>
      </view>
    </view>

    <view class="admin-panel">
      <button class="order-entry" @click="openOrders">
        <text class="order-entry__icon">☷</text>
        <view>
          <text class="order-entry__title">我的订单</text>
          <text class="order-entry__body">查看待使用、已预约、已完成和退款售后订单</text>
        </view>
        <text class="order-entry__arrow">›</text>
      </button>
    </view>

    <view class="admin-panel">
      <view class="section-title">
        <text class="section-title__text">常用服务</text>
        <text class="section-title__link">更多</text>
      </view>
      <view class="module-grid">
        <button class="module-tile">
          <text class="module-tile__icon">♙</text>
          <view><text class="module-tile__title">常用出行人</text><text class="module-tile__body">实名信息、证件管理</text></view>
        </button>
        <button class="module-tile">
          <text class="module-tile__icon">▤</text>
          <view><text class="module-tile__title">发票抬头</text><text class="module-tile__body">企业抬头、开票记录</text></view>
        </button>
        <button class="module-tile">
          <text class="module-tile__icon">♬</text>
          <view><text class="module-tile__title">售后进度</text><text class="module-tile__body">退款、改期、投诉建议</text></view>
        </button>
        <button class="module-tile">
          <text class="module-tile__icon">♡</text>
          <view><text class="module-tile__title">我的收藏</text><text class="module-tile__body">景点、酒店、路线</text></view>
        </button>
      </view>
    </view>

    <button class="logout-button" @click="clearSession">清理本地会话</button>

    <view class="visual-bottom-nav" aria-label="底部导航">
      <view class="visual-bottom-nav__item" @click="openHome">
        <text class="visual-bottom-nav__icon">⌂</text>
        <text>首页</text>
      </view>
      <view class="visual-bottom-nav__item" @click="openService">
        <text class="visual-bottom-nav__icon">♬</text>
        <text>客服</text>
      </view>
      <view class="visual-bottom-nav__item visual-bottom-nav__item--active">
        <text class="visual-bottom-nav__icon">♙</text>
        <text>我的</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.profile-page {
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

.page-heading__title,
.page-heading__body,
.profile-name,
.profile-mobile,
.relation-box__title,
.relation-box__body,
.kpi-card__value,
.kpi-card__trend,
.order-entry__title,
.order-entry__body,
.module-tile__title,
.module-tile__body {
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

.profile-card,
.admin-panel {
  margin: 0 14px 10px;
  padding: 14px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.profile-top {
  display: grid;
  grid-template-columns: 52px minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
}

.avatar {
  display: grid;
  place-items: center;
  width: 52px;
  height: 52px;
  border-radius: 8px;
  background: linear-gradient(135deg, #0f8b8d, #226bb8);
  color: #ffffff;
  font-size: 28px;
  font-weight: 900;
}

.profile-identity {
  min-width: 0;
}

.profile-name {
  font-size: 18px;
  font-weight: 900;
  line-height: 1.2;
}

.profile-mobile {
  margin-top: 6px;
  color: #71808a;
  font-size: 12px;
}

.relation-box {
  margin-top: 12px;
  padding: 11px;
  border-radius: 8px;
  background: #f7fbfc;
}

.relation-box__title {
  font-size: 14px;
  font-weight: 900;
}

.relation-box__body {
  margin-top: 6px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.4;
}

.rule-list {
  display: grid;
  gap: 8px;
  margin-top: 10px;
}

.rule-item {
  display: grid;
  grid-template-columns: 24px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  color: #71808a;
  font-size: 12px;
  line-height: 1.45;
}

.rule-item__icon {
  color: #087393;
  font-size: 18px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.section-title__text {
  font-size: 18px;
  font-weight: 900;
}

.section-title__link {
  color: #71808a;
  font-size: 13px;
}

.kpi-grid,
.module-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.kpi-card,
.order-entry,
.module-tile,
.logout-button {
  margin: 0;
  border: 0;
  border-radius: 8px;
  text-align: left;
}

.kpi-card::after,
.order-entry::after,
.module-tile::after,
.logout-button::after {
  display: none;
}

.kpi-card {
  padding: 12px;
  background: #f4fafb;
}

.kpi-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #71808a;
  font-size: 12px;
}

.kpi-card__value {
  margin-top: 8px;
  font-size: 22px;
  font-weight: 900;
  line-height: 1;
}

.kpi-card__trend {
  margin-top: 7px;
  color: #4b8a61;
  font-size: 12px;
}

.order-entry {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 12px;
  background: #ffffff;
  box-shadow: 0 10px 26px rgba(10, 49, 70, 0.075);
}

.order-entry__icon {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  border-radius: 8px;
  background: linear-gradient(135deg, #12a8c7, #087393);
  color: #ffffff;
  font-size: 22px;
}

.order-entry__title {
  color: #0d2b3a;
  font-size: 16px;
  font-weight: 900;
}

.order-entry__body {
  margin-top: 5px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.35;
}

.order-entry__arrow {
  color: #71808a;
  font-size: 22px;
}

.module-tile {
  display: grid;
  grid-template-columns: 34px minmax(0, 1fr);
  gap: 9px;
  align-items: center;
  min-height: 66px;
  padding: 10px;
  background: #f7fbfc;
}

.module-tile__icon {
  color: #087393;
  font-size: 22px;
}

.module-tile__title {
  color: #0d2b3a;
  font-size: 14px;
  font-weight: 900;
}

.module-tile__body {
  margin-top: 4px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.25;
}

.logout-button {
  width: calc(100% - 28px);
  min-height: 38px;
  margin: 0 14px 12px;
  border: 1px solid #dbe7eb;
  background: #ffffff;
  color: #71808a;
  font-size: 13px;
  text-align: center;
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
