<script setup lang="ts">
import { computed, ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import EmptyState from '@/components/EmptyState.vue';
import StatusPill from '@/components/StatusPill.vue';
import { storageAdapter } from '@/adapters/storage';
import { loadDriverViewModel } from '@/api/driver';
import type { DriverViewModel } from '@/types/driver';

const viewModel = ref<DriverViewModel | null>(null);
const activeTab = ref<'login' | 'register' | 'console'>('login');
const withdrawMessage = ref('');

const canWithdraw = computed(() => {
  const wallet = viewModel.value?.wallet;
  const amount = viewModel.value?.withdrawDraft.amount ?? 0;
  return !!wallet && amount > 0 && amount <= wallet.availableAmount;
});

const sourceLabel = computed(() => {
  if (viewModel.value?.dataSource === 'local') return '本地接口';
  if (viewModel.value?.dataSource === 'fallback') return '本地回落';
  return 'Mock';
});

onLoad((query) => {
  loadPage(String(query?.scenario ?? ''));
});

function loadPage(scenarioOverride?: string) {
  loadDriverViewModel(scenarioOverride).then((result) => {
    viewModel.value = result.data;
    withdrawMessage.value = result.data.statusMessage ?? '';
  });
}

function loginAsDriver() {
  loadDriverViewModel('phase4-active').then((result) => {
    viewModel.value = result.data;
    withdrawMessage.value = result.data.statusMessage ?? '';
    activeTab.value = 'console';
  });
}

function submitRegister() {
  loadDriverViewModel('phase4-pending-review').then((result) => {
    viewModel.value = result.data;
    withdrawMessage.value = result.data.statusMessage ?? '注册申请已提交，请等待管理员审核。';
    activeTab.value = 'console';
  });
}

function submitWithdraw() {
  if (!viewModel.value?.wallet) {
    withdrawMessage.value = '请先登录 active 司机账号。';
    return;
  }
  if (!canWithdraw.value) {
    withdrawMessage.value = `余额不足，可提现 ${viewModel.value.wallet.availableLabel}。`;
    return;
  }
  withdrawMessage.value = '提现申请已提交，待管理员审核打款。';
}

function logout() {
  storageAdapter.clearSession();
  activeTab.value = 'login';
}
</script>

<template>
  <view class="driver-page">
    <view class="hero">
      <image class="hero-bg" src="/static/phase2/taxi-scan-illustration-web.jpg" mode="aspectFill" />
      <view class="top-line">
        <button class="back-link">‹ 返回首页</button>
        <button class="help-chip">◉ 司机入驻</button>
      </view>
      <view class="hero-copy">
        <text class="hero-title">司机登录与注册</text>
        <text class="hero-body">审核通过的司机可登录查看钱包、佣金和提现记录；新司机先提交身份与车辆信息等待后台审核。</text>
      </view>
      <view class="hero-card">
        <view class="metric"><strong>1</strong><text>提交身份信息</text></view>
        <view class="metric"><strong>2</strong><text>绑定车辆车牌</text></view>
        <view class="metric"><strong>3</strong><text>审核后生成二维码</text></view>
      </view>
    </view>

    <view class="content">
      <view class="mode-tabs">
        <button class="mode-tab" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">司机登录</button>
        <button class="mode-tab" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">司机注册</button>
      </view>

      <EmptyState v-if="viewModel?.errorMessage" title="司机端提示" :description="viewModel.errorMessage" />

      <view v-if="activeTab === 'login'" class="card form-card">
        <view class="section-head">
          <text>司机登录</text>
          <text>active drivers</text>
        </view>
        <view class="form-grid">
          <view class="field"><text>手机号</text><input value="13800138000" placeholder="请输入已审核司机手机号" /></view>
          <view class="field"><text>登录密码</text><input password value="******" placeholder="请输入登录密码" /></view>
        </view>
        <button class="primary-button" @click="loginAsDriver">登录司机端</button>
      </view>

      <view v-if="activeTab === 'register'" class="card form-card">
        <view class="section-head">
          <text>司机注册</text>
          <text>drivers / vehicles</text>
        </view>
        <view class="step-row">
          <view class="step active">身份信息</view>
          <view class="step">车辆信息</view>
          <view class="step">提交审核</view>
        </view>
        <view class="form-grid">
          <view class="field"><text>真实姓名</text><input value="陈师傅" placeholder="请输入与身份证一致的姓名" /></view>
          <view class="field"><text>手机号</text><input value="13800138000" placeholder="用于登录和接收审核通知" /></view>
          <view class="field"><text>身份证号</text><input value="440402********2345" placeholder="请输入18位身份证号码" /></view>
          <view class="field"><text>车牌号</text><input value="粤C·D4208" placeholder="请输入车辆号牌" /></view>
        </view>
        <button class="primary-button" @click="submitRegister">提交注册申请</button>
      </view>

      <view v-if="activeTab === 'console'" class="console">
        <view v-if="viewModel?.profile" class="card profile-card">
          <view class="driver-profile">
            <view>
              <text class="driver-name">{{ viewModel.profile.name }}</text>
              <text class="driver-meta">{{ viewModel.profile.driverNo }} · {{ viewModel.profile.commissionRateLabel }}</text>
              <text class="driver-meta">{{ viewModel.profile.maskedPhone }} · {{ viewModel.profile.vehicleLabel }}</text>
            </view>
            <StatusPill :text="viewModel.profile.status === 'active' ? '已审核' : '待审核'" :tone="viewModel.profile.status === 'active' ? 'success' : 'warning'" />
          </view>

          <view class="qr-panel">
            <view class="qr-box"><text></text></view>
            <view>
              <text class="qr-title">车座推广码</text>
              <text class="qr-body">{{ viewModel.profile.qrCodeText }}</text>
            </view>
          </view>
        </view>
        <EmptyState v-else title="未登录" description="请先登录司机端或提交注册申请。" />

        <view v-if="viewModel?.profile?.status === 'pending_review'" class="card panel-card">
          <text class="section-title">审核状态</text>
          <text class="panel-copy">{{ viewModel.statusMessage }}</text>
        </view>

        <view v-if="viewModel?.wallet" class="card panel-card">
          <view class="section-head">
            <text>钱包余额</text>
            <text>{{ sourceLabel }} 同步</text>
          </view>
          <view class="wallet-grid">
            <view><text>可提现</text><strong>{{ viewModel.wallet.availableLabel }}</strong></view>
            <view><text>待结算</text><strong>{{ viewModel.wallet.pendingLabel }}</strong></view>
            <view><text>已结算</text><strong>{{ viewModel.wallet.settledLabel }}</strong></view>
            <view><text>已提现</text><strong>{{ viewModel.wallet.withdrawnLabel }}</strong></view>
          </view>
        </view>

        <view v-if="viewModel?.wallet" class="card panel-card">
          <text class="section-title">最近佣金</text>
          <view v-if="viewModel.commissions.length === 0">
            <EmptyState title="暂无佣金记录" description="订单核销后自动计算佣金。" />
          </view>
          <view v-else class="record-list">
            <view v-for="item in viewModel.commissions" :key="item.id" class="record-item">
              <text>{{ item.orderNo }} · {{ item.amountLabel }}</text>
              <text>{{ item.statusLabel }} · {{ item.createdAtLabel }}</text>
            </view>
          </view>
        </view>

        <view v-if="viewModel?.wallet" class="card panel-card">
          <text class="section-title">申请提现</text>
          <view class="form-grid">
            <view class="field"><text>提现金额</text><input v-model.number="viewModel.withdrawDraft.amount" type="digit" /></view>
            <view class="field"><text>收款账号</text><input v-model="viewModel.withdrawDraft.account" /></view>
            <view class="field"><text>实名姓名</text><input v-model="viewModel.withdrawDraft.realName" /></view>
          </view>
          <button class="primary-button" @click="submitWithdraw">提交提现</button>
          <text v-if="withdrawMessage" class="panel-copy">{{ withdrawMessage }}</text>
        </view>

        <button class="logout-button" @click="logout">退出司机端</button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.driver-page {
  width: min(100vw, 430px);
  min-height: 100vh;
  margin: clamp(0px, calc((100vw - 430px) * 100), 28px) auto;
  overflow-x: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px);
  background:
    radial-gradient(circle at 20% 0%, rgba(15, 141, 173, 0.18), transparent 34%),
    linear-gradient(180deg, #f7fbfc, #dfeff2);
  color: #15242d;
  box-shadow: 0 0 0 1px rgba(20, 35, 45, 0.05), 0 24px 80px rgba(8, 45, 64, 0.18);
}

button {
  margin: 0;
}

button::after {
  display: none;
}

.hero {
  position: relative;
  min-height: 238px;
  padding: 18px 16px;
  overflow: hidden;
  color: #ffffff;
  background: #0f8dad;
}

.hero::before,
.hero::after {
  content: "";
  position: absolute;
}

.hero::before {
  inset: 0;
  z-index: 1;
  background: linear-gradient(135deg, rgba(8, 109, 137, 0.96), rgba(15, 141, 173, 0.82));
}

.hero::after {
  inset: auto -54px -86px 36%;
  z-index: 2;
  height: 190px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  transform: rotate(-12deg);
}

.hero-bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.top-line,
.hero-copy,
.hero-card {
  position: relative;
  z-index: 3;
}

.top-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.back-link,
.help-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: 34px;
  padding: 0 10px;
  border-radius: 999px;
  color: #ffffff;
  background: rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(8px);
  font-size: 13px;
  font-weight: 800;
}

.hero-copy {
  margin-top: 30px;
}

.hero-title,
.hero-body,
.metric strong,
.metric text,
.section-head text,
.field text,
.driver-name,
.driver-meta,
.qr-title,
.qr-body,
.panel-copy,
.record-item text,
.section-title {
  display: block;
}

.hero-title {
  font-size: 30px;
  font-weight: 900;
  line-height: 1.12;
}

.hero-body {
  max-width: 300px;
  margin-top: 10px;
  color: rgba(255, 255, 255, 0.88);
  font-size: 13px;
  line-height: 1.5;
}

.hero-card {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin-top: 20px;
}

.metric {
  min-height: 72px;
  padding: 11px 9px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.metric strong {
  font-size: 18px;
}

.metric text {
  margin-top: 5px;
  color: rgba(255, 255, 255, 0.78);
  font-size: 11px;
  line-height: 1.35;
}

.content {
  padding: 12px 14px calc(36px + env(safe-area-inset-bottom));
}

.mode-tabs {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 118px;
}

.console .mode-tabs {
  margin-bottom: 12px;
}

.mode-tab {
  height: 42px;
  border: 1px solid #d8e6eb;
  border-radius: 8px;
  color: #465a64;
  background: #ffffff;
  font-size: 17px;
  font-weight: 900;
}

.mode-tab.active {
  color: #ffffff;
  background: #0f8dad;
  border-color: #0f8dad;
}

.card {
  margin-bottom: 10px;
  padding: 14px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 14px 34px rgba(8, 45, 64, 0.12);
}

.section-head,
.driver-profile,
.qr-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-head {
  margin-bottom: 12px;
}

.section-head text:first-child,
.section-title {
  font-size: 18px;
  font-weight: 900;
}

.section-head text:last-child,
.driver-meta,
.qr-body,
.panel-copy,
.record-item text:last-child {
  color: #6c7a84;
  font-size: 12px;
  line-height: 1.45;
}

.step-row,
.form-grid,
.wallet-grid,
.record-list {
  display: grid;
  gap: 10px;
}

.step-row {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin-bottom: 12px;
}

.step {
  padding: 9px 8px;
  border-radius: 8px;
  color: #56707a;
  background: #f2f8fa;
  font-size: 12px;
  font-weight: 800;
  text-align: center;
}

.step.active {
  color: #086d89;
  background: #e8f7fb;
}

.field {
  display: grid;
  gap: 7px;
}

.field text {
  color: #4c5f69;
  font-size: 13px;
  font-weight: 800;
}

.field input {
  width: 100%;
  height: 46px;
  padding: 0 12px;
  box-sizing: border-box;
  border: 1px solid #e4edf1;
  border-radius: 8px;
  color: #15242d;
  background: #fbfdfe;
  font-size: 15px;
}

.primary-button,
.logout-button {
  width: 100%;
  min-height: 46px;
  margin-top: 14px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 900;
}

.primary-button {
  color: #ffffff;
  background: #0f8dad;
}

.logout-button {
  color: #6c7a84;
  background: #ffffff;
}

.console {
  margin-top: -96px;
}

.driver-name {
  font-size: 20px;
  font-weight: 900;
}

.qr-panel {
  justify-content: flex-start;
  margin-top: 14px;
  padding: 12px;
  border-radius: 8px;
  background: #edf5f7;
}

.qr-box {
  width: 82px;
  height: 82px;
  border: 6px solid #152238;
  border-radius: 8px;
  background:
    linear-gradient(90deg, #152238 12px, transparent 12px 20px, #152238 20px 32px, transparent 32px),
    linear-gradient(#152238 12px, transparent 12px 20px, #152238 20px 32px, transparent 32px),
    #ffffff;
}

.qr-title {
  font-size: 16px;
  font-weight: 900;
}

.wallet-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.wallet-grid view,
.record-item {
  padding: 12px;
  border-radius: 8px;
  background: #edf5f7;
}

.wallet-grid text,
.wallet-grid strong,
.record-item text:first-child {
  display: block;
}

.wallet-grid text {
  color: #6c7a84;
  font-size: 12px;
}

.wallet-grid strong {
  margin-top: 6px;
  color: #f15a2b;
  font-size: 20px;
}

.panel-copy {
  margin-top: 10px;
}
</style>
