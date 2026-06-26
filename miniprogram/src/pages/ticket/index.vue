<script setup lang="ts">
import { computed, ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import { loadTicketDetail } from '@/api/main-shell';
import type { TicketSummary } from '@/types/main-shell';

const ticket = ref<TicketSummary | null>(null);
const errorMessage = ref('');

const ticketDateCopy = computed(() => {
  if (!ticket.value) return '游玩日期 --，有效期至当日 23:59';
  return `游玩日期 2026-${ticket.value.validDateLabel.replace('-', '-')}，有效期至当日 23:59`;
});

loadTicketDetail()
  .then((result) => {
    ticket.value = result.data;
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  });
</script>

<template>
  <view class="ticket-page">
    <view class="status-bar" aria-hidden="true">
      <text>21:47</text>
      <view class="status-bar__icons"><text>⋯</text><text>⌁</text><text class="battery" /></view>
    </view>
    <view class="ticket-nav" aria-label="页面导航">
      <button class="ticket-nav__button ticket-nav__button--back" aria-label="返回" />
      <view class="ticket-nav__actions">
        <button class="ticket-nav__button ticket-nav__button--service" aria-label="客服" />
        <button class="ticket-nav__button ticket-nav__button--share" aria-label="分享" />
        <button class="ticket-nav__button ticket-nav__button--home" aria-label="首页" />
      </view>
    </view>

    <EmptyState v-if="errorMessage" title="票券读取失败" :description="errorMessage" />
    <EmptyState v-else-if="!ticket" title="暂无票券" description="当前 Mock 场景没有可展示票券。" />
    <template v-else>
      <view class="ticket-state">
        <view class="state-title"><text>◷</text><text>待使用</text></view>
        <text class="state-subtitle">{{ ticketDateCopy }}</text>
      </view>

      <view class="ticket-visual" aria-label="电子票核销">
        <image class="ticket-visual__image" src="/static/phase3/verify-hero-wide-clean-web.png" mode="aspectFill" />
        <view class="ticket-visual__copy">
          <text class="ticket-visual__title">出示核销码检票</text>
          <text class="ticket-visual__body">{{ ticket.verifyLocation }}核验后登船。</text>
        </view>
      </view>

      <view class="card ticket-main">
        <view class="qr-wrap">
          <view class="fake-qr-large">
            <text class="qr-eye qr-eye--tl" />
            <text class="qr-eye qr-eye--tr" />
            <text class="qr-eye qr-eye--bl" />
          </view>
        </view>
        <view class="coupon">
          <text class="coupon__title">待使用券码 <text>1</text> 张  ◉</text>
          <view class="code-line">
            <text>· {{ ticket.maskedCode.replace('********', '0491') }}</text>
            <text class="copy-button">复制</text>
          </view>
        </view>
      </view>

      <view class="card visit-card">
        <view class="visit-head">
          <text class="visit-head__title">{{ ticket.productTitle }}</text>
          <text class="tag">电子票</text>
        </view>
        <view class="visit-grid">
          <view><text>游玩日期</text><text class="visit-value">{{ ticket.validDateLabel }} 周三</text></view>
          <view><text>入园时段</text><text class="visit-value">{{ ticket.validTimeLabel }}</text></view>
          <view><text>数量</text><text class="visit-value">{{ ticket.quantityLabel }}</text></view>
        </view>
      </view>

      <view class="card product-card">
        <image class="product-thumb" src="/static/phase2/macau-cruise-night-banner-web.jpg" mode="aspectFill" />
        <view>
          <text class="product-title">{{ ticket.productTitle }}</text>
          <text class="product-body">湾仔旅游码头出发，凭核销码检票登船</text>
        </view>
        <view class="product-price">
          <text>¥128</text>
          <text>x1</text>
          <text>实付 {{ ticket.amountLabel }}</text>
        </view>
      </view>

      <view class="card spot-card">
        <view>
          <text class="spot-title">{{ ticket.verifyLocation }}</text>
          <text class="spot-line">◷ 营业中 周一至周日 09:30-21:30</text>
          <text class="spot-line">⌖ 珠海市香洲区湾仔旅游码头售票大厅旁</text>
        </view>
        <button class="call-button">☎</button>
      </view>

      <view class="card notice-card">
        <view class="visit-head">
          <text class="visit-head__title">使用须知</text>
          <text class="notice-link">查看全部 ›</text>
        </view>
        <text class="notice-copy">请在游玩日期当天到检票口出示二维码，工作人员扫码后入园/登船。</text>
        <text v-for="notice in ticket.notice" :key="notice" class="notice-item">· {{ notice }}</text>
      </view>

      <view class="card order-info-card">
        <text class="visit-head__title">订单信息</text>
        <view class="info-row"><text>实际付款</text><text>{{ ticket.amountLabel }}</text></view>
        <view class="info-row"><text>手机号码</text><text>189****6928</text></view>
        <view class="info-row"><text>订单编号</text><text>{{ ticket.orderNo }} 复制</text></view>
        <view class="info-row"><text>下单来源</text><text>车座购票码</text></view>
      </view>
    </template>

    <view class="bottom-actions">
      <button>☎ 景区电话</button>
      <button>申请退款</button>
    </view>
  </view>
</template>

<style scoped>
.ticket-page {
  width: min(100vw, 430px);
  min-height: 100vh;
  margin: clamp(0px, calc((100vw - 430px) * 100), 28px) auto;
  padding: 8px 12px calc(104px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  overflow-x: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px);
  background: linear-gradient(180deg, #f6f8fb 0%, #eef4f5 100%);
  color: #071627;
  box-shadow:
    0 0 0 1px rgba(9, 50, 66, 0.04),
    0 24px 80px rgba(8, 45, 64, 0.18);
}

.ticket-nav {
  height: 54px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.status-bar {
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 6px;
  box-sizing: border-box;
  color: #000000;
  font-size: 24px;
  font-weight: 900;
}

.status-bar__icons {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
}

.battery {
  position: relative;
  display: block;
  width: 30px;
  height: 16px;
  border: 3px solid currentColor;
  border-radius: 4px;
}

.battery::before {
  content: "";
  position: absolute;
  right: -6px;
  top: 4px;
  width: 3px;
  height: 6px;
  border-radius: 0 2px 2px 0;
  background: currentColor;
}

.battery::after {
  content: "";
  position: absolute;
  inset: 3px;
  border-radius: 2px;
  background: currentColor;
}

.ticket-nav__button,
.bottom-actions button,
.call-button {
  margin: 0;
  border: 0;
}

.ticket-nav__button::after,
.bottom-actions button::after,
.call-button::after {
  display: none;
}

.ticket-nav__button {
  position: relative;
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  padding: 0;
  border-radius: 50%;
  color: #111111;
  background: transparent;
  background-repeat: no-repeat;
  background-position: center;
  background-size: 28px 28px;
  overflow: hidden;
}

.ticket-nav__actions {
  display: flex;
  gap: 16px;
}

.ticket-nav__button--back {
  background-image: url("data:image/svg+xml,%3Csvg width='28' height='28' viewBox='0 0 28 28' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M17.5 6.5L10 14l7.5 7.5' stroke='%23000' stroke-width='3.2' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
}

.ticket-nav__button--service {
  background-image: url("data:image/svg+xml,%3Csvg width='28' height='28' viewBox='0 0 28 28' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M5 15v-2.5C5 7.8 8.8 4 14 4s9 3.8 9 8.5V15' stroke='%23000' stroke-width='3.1' stroke-linecap='round'/%3E%3Cpath d='M5 15h4v7H5zM19 15h4v7h-4z' fill='%23000'/%3E%3C/svg%3E");
}

.ticket-nav__button--share {
  background-image: url("data:image/svg+xml,%3Csvg width='28' height='28' viewBox='0 0 28 28' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='8' cy='14' r='3.4' fill='%23000'/%3E%3Ccircle cx='20' cy='7' r='3.4' fill='%23000'/%3E%3Ccircle cx='20' cy='21' r='3.4' fill='%23000'/%3E%3Cpath d='M10.9 12.3l6.2-3.6M10.9 15.7l6.2 3.6' stroke='%23000' stroke-width='2.7' stroke-linecap='round'/%3E%3C/svg%3E");
}

.ticket-nav__button--home {
  background-image: url("data:image/svg+xml,%3Csvg width='28' height='28' viewBox='0 0 28 28' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M5 13.2L14 5l9 8.2' stroke='%23000' stroke-width='3' stroke-linecap='round' stroke-linejoin='round'/%3E%3Cpath d='M8 12.8V23h12V12.8' stroke='%23000' stroke-width='3' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
}

.ticket-state {
  padding: 6px 8px 18px;
}

.state-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 34px;
  font-weight: 900;
  line-height: 1.15;
}

.state-subtitle {
  display: block;
  margin-top: 10px;
  color: #61666b;
  font-size: 17px;
  line-height: 1.4;
}

.ticket-visual {
  position: relative;
  min-height: 116px;
  margin-bottom: 12px;
  overflow: hidden;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.ticket-visual__image,
.ticket-visual::after {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.ticket-visual::after {
  content: "";
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.74), rgba(255, 255, 255, 0.05) 58%);
}

.ticket-visual__copy {
  position: relative;
  z-index: 1;
  max-width: 210px;
  padding: 18px;
}

.ticket-visual__title,
.ticket-visual__body,
.coupon__title,
.product-title,
.product-body,
.spot-title,
.spot-line,
.notice-copy,
.notice-item {
  display: block;
}

.ticket-visual__title {
  font-size: 20px;
  font-weight: 900;
  line-height: 1.2;
}

.ticket-visual__body {
  margin-top: 8px;
  color: #61666b;
  font-size: 13px;
  line-height: 1.45;
}

.card {
  margin-bottom: 12px;
  padding: 14px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.ticket-main {
  padding: 0;
  overflow: hidden;
}

.qr-wrap {
  display: grid;
  place-items: center;
  padding: 40px 0 34px;
  background: #ffffff;
}

.fake-qr-large {
  position: relative;
  width: 246px;
  height: 246px;
  background:
    linear-gradient(90deg, transparent 0 16px, #000 16px 28px, transparent 28px 44px, #000 44px 56px, transparent 56px) 0 0 / 74px 74px,
    linear-gradient(0deg, transparent 0 18px, #000 18px 30px, transparent 30px 45px, #000 45px 57px, transparent 57px) 0 0 / 74px 74px,
    linear-gradient(135deg, transparent 0 40%, #000 41% 57%, transparent 58%) 0 0 / 34px 34px,
    #ffffff;
}

.qr-eye {
  position: absolute;
  width: 74px;
  height: 74px;
  border: 12px solid #000000;
  box-sizing: border-box;
  background: #ffffff;
}

.qr-eye::after {
  content: "";
  position: absolute;
  inset: 13px;
  background: #000000;
}

.qr-eye--tl {
  left: 0;
  top: 0;
}

.qr-eye--tr {
  right: 0;
  top: 0;
}

.qr-eye--bl {
  left: 0;
  bottom: 0;
}

.coupon {
  padding: 17px 18px 18px;
  border-top: 1px dashed #dbe2e6;
}

.coupon__title {
  font-size: 21px;
  font-weight: 900;
}

.coupon__title text {
  color: #f05b23;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 9px;
  font-size: 20px;
  font-weight: 900;
}

.copy-button {
  padding: 5px 8px;
  border-radius: 7px;
  color: #73808a;
  background: #eef2f5;
  font-size: 12px;
  font-weight: 700;
}

.visit-head,
.spot-card,
.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.visit-head__title,
.spot-title {
  font-size: 20px;
  font-weight: 900;
}

.tag {
  padding: 5px 8px;
  border-radius: 8px;
  color: #087393;
  background: #e9f7fa;
  font-size: 12px;
  font-weight: 900;
}

.visit-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin-top: 12px;
}

.visit-grid view {
  padding: 11px;
  border-radius: 8px;
  background: #f1fafb;
  color: #71808a;
  font-size: 12px;
}

.visit-value {
  display: block;
  margin-top: 6px;
  color: #071627;
  font-weight: 900;
}

.product-card {
  display: grid;
  grid-template-columns: 82px minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
}

.product-thumb {
  width: 82px;
  height: 64px;
  border-radius: 8px;
}

.product-title {
  font-size: 15px;
  font-weight: 900;
}

.product-body,
.spot-line,
.notice-copy,
.notice-item {
  margin-top: 6px;
  color: #71808a;
  font-size: 12px;
  line-height: 1.4;
}

.product-price {
  display: grid;
  gap: 4px;
  text-align: right;
  color: #71808a;
  font-size: 12px;
}

.product-price text:first-child,
.product-price text:last-child {
  color: #e65337;
  font-weight: 900;
}

.call-button {
  flex: 0 0 auto;
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  border-radius: 50%;
  color: #ffffff;
  background: #087393;
  font-size: 18px;
}

.notice-link {
  color: #71808a;
  font-size: 14px;
}

.info-row {
  padding: 12px 0;
  border-top: 1px solid #dbe7eb;
  color: #71808a;
  font-size: 14px;
}

.info-row text:last-child {
  color: #071627;
  font-weight: 900;
  text-align: right;
}

.bottom-actions {
  position: fixed;
  z-index: 30;
  left: 50%;
  bottom: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  width: min(100vw, 430px);
  transform: translateX(-50%);
  padding: 12px 14px calc(12px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  border-top: 1px solid #dbe7eb;
  border-radius: 0 0 clamp(0px, calc((100vw - 430px) * 100), 26px) clamp(0px, calc((100vw - 430px) * 100), 26px);
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
}

.bottom-actions button {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 52px;
  border: 1px solid #dbe7eb;
  border-radius: 999px;
  color: #1d2733;
  background: #ffffff;
  font-size: 16px;
  font-weight: 900;
}
</style>
