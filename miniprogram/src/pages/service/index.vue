<script setup lang="ts">
import { navigationAdapter } from '@/adapters/navigation';

const supportItems = [
  {
    icon: '♬',
    title: '在线客服',
    body: '船票、景区、酒店与行程咨询'
  },
  {
    icon: '▤',
    title: '订单售后',
    body: '退款、改签、发票和异常订单'
  },
  {
    icon: '⌖',
    title: '团队定制',
    body: '企业出行、研学和包车路线'
  }
];

function openHome() {
  navigationAdapter.switchTab('/pages/home/index');
}

function openProfile() {
  navigationAdapter.switchTab('/pages/profile/index');
}

function openService() {
  navigationAdapter.switchTab('/pages/service/index');
}
</script>

<template>
  <view class="service-page">
    <view class="page-heading">
      <text class="page-heading__title">客服中心</text>
      <text class="page-heading__body">售前咨询、订单售后、团队定制与紧急协助统一从这里进入。</text>
    </view>

    <view class="support-panel">
      <view v-for="item in supportItems" :key="item.title" class="support-action">
        <text class="support-icon">{{ item.icon }}</text>
        <view>
          <text class="support-title">{{ item.title }}</text>
          <text class="support-body">{{ item.body }}</text>
        </view>
        <text class="support-arrow">›</text>
      </view>
    </view>

    <view class="support-card">
      <view class="section-title">
        <text class="section-title__text">服务时间</text>
        <text class="section-title__link">09:00-22:00</text>
      </view>
      <view class="service-grid">
        <view><text>热线电话</text><strong>400-888-0756</strong></view>
        <view><text>响应承诺</text><strong>3 分钟内</strong></view>
      </view>
    </view>

    <view class="support-card">
      <view class="section-title">
        <text class="section-title__text">常见问题</text>
        <text class="section-title__link">更多</text>
      </view>
      <view class="faq-list">
        <text>· 出票后请在票券页查看核销码。</text>
        <text>· 天气或航班调整以码头现场通知为准。</text>
        <text>· 退款和改期请保留订单编号便于客服查询。</text>
      </view>
    </view>

    <view class="visual-bottom-nav" aria-label="底部导航">
      <view class="visual-bottom-nav__item" @click="openHome">
        <text class="visual-bottom-nav__icon">⌂</text>
        <text>首页</text>
      </view>
      <view class="visual-bottom-nav__item visual-bottom-nav__item--active" @click="openService">
        <text class="visual-bottom-nav__icon">♬</text>
        <text>客服</text>
      </view>
      <view class="visual-bottom-nav__item" @click="openProfile">
        <text class="visual-bottom-nav__icon">♙</text>
        <text>我的</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.service-page {
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
  padding: 18px 14px;
  border-radius: 8px;
  color: #ffffff;
  background: linear-gradient(135deg, #0f8b8d, #226bb8);
  box-shadow: 0 12px 30px rgba(9, 50, 66, 0.12);
}

.page-heading__title,
.page-heading__body,
.support-title,
.support-body,
.section-title__text,
.section-title__link,
.faq-list text,
.service-grid text,
.service-grid strong {
  display: block;
}

.page-heading__title {
  font-size: 26px;
  font-weight: 900;
}

.page-heading__body {
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.86);
  font-size: 13px;
  line-height: 1.45;
}

.support-panel,
.support-card {
  margin: 0 14px 12px;
  padding: 14px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.support-action {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr) auto;
  align-items: center;
  gap: 11px;
  min-height: 72px;
  padding: 10px 0;
  border-bottom: 1px solid #e3eef2;
}

.support-action:last-child {
  border-bottom: 0;
}

.support-icon {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  border-radius: 8px;
  color: #ffffff;
  background: linear-gradient(135deg, #12a8c7, #087393);
  font-size: 22px;
}

.support-title {
  font-size: 16px;
  font-weight: 900;
}

.support-body,
.faq-list text,
.service-grid text {
  margin-top: 5px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.45;
}

.support-arrow {
  color: #71808a;
  font-size: 22px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 12px;
}

.section-title__text {
  font-size: 18px;
  font-weight: 900;
}

.section-title__link {
  color: #71808a;
  font-size: 13px;
  font-weight: 800;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.service-grid view {
  padding: 12px;
  border-radius: 8px;
  background: #f4fafb;
}

.service-grid strong {
  margin-top: 7px;
  font-size: 18px;
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
