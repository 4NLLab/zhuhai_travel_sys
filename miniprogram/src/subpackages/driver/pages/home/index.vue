<script setup lang="ts">
import { computed, ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import EmptyState from '@/components/EmptyState.vue';
import StatusPill from '@/components/StatusPill.vue';
import { storageAdapter } from '@/adapters/storage';
import { loadDriverViewModel } from '@/api/driver';
import type { DriverViewModel } from '@/types/driver';

const viewModel = ref<DriverViewModel | null>(null);
const activeTab = ref<'login' | 'register' | 'console'>('console');
const withdrawMessage = ref('');

const canWithdraw = computed(() => {
  const wallet = viewModel.value?.wallet;
  const amount = viewModel.value?.withdrawDraft.amount ?? 0;
  return !!wallet && amount > 0 && amount <= wallet.availableAmount;
});

onLoad((query) => {
  loadPage(String(query?.scenario ?? ''));
});

function loadPage(scenarioOverride?: string) {
  loadDriverViewModel(scenarioOverride).then((result) => {
    viewModel.value = result.data;
    activeTab.value = result.data.profile ? 'console' : 'login';
    withdrawMessage.value = result.data.statusMessage ?? '';
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
  <view class="page-shell driver-page">
    <view class="driver-hero">
      <text class="driver-hero__title">司机端</text>
      <text class="driver-hero__body">审核通过的司机可查看钱包、佣金、提现记录和车座推广码。</text>
      <view class="driver-hero__metrics">
        <view><strong>3</strong><text>步骤入驻</text></view>
        <view><strong>12%</strong><text>示例提成</text></view>
        <view><strong>分包</strong><text>低频入口</text></view>
      </view>
    </view>

    <view class="mode-tabs">
      <button class="mode-tabs__item" :class="{ 'mode-tabs__item--active': activeTab === 'login' }" @click="activeTab = 'login'">登录</button>
      <button class="mode-tabs__item" :class="{ 'mode-tabs__item--active': activeTab === 'register' }" @click="activeTab = 'register'">注册</button>
      <button class="mode-tabs__item" :class="{ 'mode-tabs__item--active': activeTab === 'console' }" @click="activeTab = 'console'">工作台</button>
    </view>

    <EmptyState v-if="viewModel?.errorMessage" title="司机端提示" :description="viewModel.errorMessage" />

    <view v-if="activeTab === 'login'" class="card panel-card">
      <view class="section-row">
        <text class="section-title">司机登录</text>
        <StatusPill text="active drivers" tone="success" />
      </view>
      <view class="field-box"><text>手机号</text><input value="13800138000" placeholder="请输入已审核司机手机号" /></view>
      <view class="field-box"><text>登录密码</text><input password value="******" placeholder="请输入登录密码" /></view>
      <button class="primary-button panel-button" @click="loadPage('phase4-active')">登录司机端</button>
    </view>

    <view v-if="activeTab === 'register'" class="card panel-card">
      <view class="section-row">
        <text class="section-title">司机注册</text>
        <StatusPill text="等待审核" tone="warning" />
      </view>
      <view class="field-grid">
        <view class="field-box"><text>真实姓名</text><input value="陈师傅" /></view>
        <view class="field-box"><text>手机号</text><input value="13800138000" /></view>
        <view class="field-box"><text>身份证号</text><input value="440402********2345" /></view>
        <view class="field-box"><text>车牌号</text><input value="粤C·D4208" /></view>
      </view>
      <button class="primary-button panel-button" @click="loadPage('phase4-pending-review')">提交注册申请</button>
    </view>

    <view v-if="activeTab === 'console'">
      <view v-if="viewModel?.profile" class="card panel-card">
        <view class="driver-profile">
          <view>
            <text class="driver-profile__name">{{ viewModel.profile.name }}</text>
            <text class="driver-profile__meta">{{ viewModel.profile.driverNo }} · {{ viewModel.profile.commissionRateLabel }}</text>
            <text class="driver-profile__meta">{{ viewModel.profile.maskedPhone }} · {{ viewModel.profile.vehicleLabel }}</text>
          </view>
          <StatusPill :text="viewModel.profile.status === 'active' ? '已审核' : '待审核'" :tone="viewModel.profile.status === 'active' ? 'success' : 'warning'" />
        </view>

        <view class="qr-panel">
          <view class="qr-box"><text></text></view>
          <view>
            <text class="qr-panel__title">车座推广码</text>
            <text class="qr-panel__body">{{ viewModel.profile.qrCodeText }}</text>
          </view>
        </view>
      </view>
      <EmptyState v-else title="未登录" description="请先登录司机端或提交注册申请。" />

      <view v-if="viewModel?.profile?.status === 'pending_review'" class="card panel-card">
        <text class="section-title">审核状态</text>
        <text class="panel-copy">{{ viewModel.statusMessage }}</text>
      </view>

      <view v-if="viewModel?.wallet" class="card panel-card">
        <view class="section-row">
          <text class="section-title">钱包余额</text>
          <text class="section-link">Mock 同步</text>
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
        <view class="field-grid">
          <view class="field-box"><text>提现金额</text><input v-model.number="viewModel.withdrawDraft.amount" type="digit" /></view>
          <view class="field-box"><text>收款账号</text><input v-model="viewModel.withdrawDraft.account" /></view>
          <view class="field-box"><text>实名姓名</text><input v-model="viewModel.withdrawDraft.realName" /></view>
        </view>
        <button class="primary-button panel-button" @click="submitWithdraw">提交提现</button>
        <text v-if="withdrawMessage" class="panel-copy">{{ withdrawMessage }}</text>
      </view>

      <view v-if="viewModel?.wallet" class="card panel-card">
        <text class="section-title">提现记录</text>
        <view v-if="viewModel.withdrawals.length === 0">
          <EmptyState title="暂无提现记录" description="提交提现申请后会在这里展示审核状态。" />
        </view>
        <view v-else class="record-list">
          <view v-for="item in viewModel.withdrawals" :key="item.id" class="record-item">
            <text>{{ item.withdrawalNo }} · {{ item.amountLabel }}</text>
            <text>{{ item.statusLabel }} · {{ item.accountLabel }} · {{ item.createdAtLabel }}</text>
          </view>
        </view>
      </view>

      <button class="secondary-button panel-button" @click="logout">退出司机端</button>
    </view>
  </view>
</template>

<style scoped>
.driver-page {
  padding-top: 16px;
}

.driver-hero {
  padding: 22px;
  border-radius: 8px;
  color: #ffffff;
  background: linear-gradient(135deg, #0f766e, #264653);
}

.driver-hero__title,
.driver-hero__body,
.driver-profile__name,
.driver-profile__meta,
.qr-panel__title,
.qr-panel__body,
.panel-copy,
.record-item text {
  display: block;
}

.driver-hero__title {
  font-size: 30px;
  font-weight: 900;
}

.driver-hero__body {
  margin-top: 8px;
  font-size: 14px;
  line-height: 22px;
}

.driver-hero__metrics,
.wallet-grid,
.field-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.driver-hero__metrics view,
.wallet-grid view,
.field-box {
  padding: 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.14);
}

.driver-hero__metrics strong,
.driver-hero__metrics text {
  display: block;
}

.driver-hero__metrics strong {
  font-size: 20px;
}

.driver-hero__metrics text {
  margin-top: 4px;
  font-size: 12px;
}

.mode-tabs {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin: 14px 0;
}

.mode-tabs__item {
  height: 36px;
  margin: 0;
  border-radius: 8px;
  background: var(--zt-color-surface);
  color: var(--zt-color-text-muted);
}

.mode-tabs__item--active {
  background: var(--zt-color-primary);
  color: #ffffff;
}

.panel-card {
  margin-top: 14px;
  padding: 16px;
}

.section-row,
.driver-profile,
.qr-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-link,
.driver-profile__meta,
.qr-panel__body,
.panel-copy,
.record-item text:last-child {
  color: var(--zt-color-text-muted);
  font-size: 13px;
  line-height: 20px;
}

.field-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.field-box {
  background: var(--zt-color-surface-muted);
}

.field-box text,
.wallet-grid text,
.wallet-grid strong,
.record-item text:first-child {
  display: block;
  font-weight: 800;
}

.field-box input {
  margin-top: 8px;
}

.panel-button {
  width: 100%;
  margin-top: 14px;
}

.driver-profile__name {
  font-size: 20px;
  font-weight: 900;
}

.qr-panel {
  justify-content: flex-start;
  margin-top: 14px;
  padding: 12px;
  border-radius: 8px;
  background: var(--zt-color-surface-muted);
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

.qr-panel__title {
  font-size: 16px;
  font-weight: 900;
}

.wallet-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.wallet-grid view {
  background: var(--zt-color-surface-muted);
}

.wallet-grid strong {
  margin-top: 6px;
  color: var(--zt-color-accent);
  font-size: 20px;
}

.record-list {
  display: grid;
  gap: 10px;
  margin-top: 10px;
}

.record-item {
  padding: 12px;
  border-radius: 8px;
  background: var(--zt-color-surface-muted);
}

@media (max-width: 640px) {
  .driver-hero__metrics,
  .field-grid,
  .wallet-grid {
    grid-template-columns: 1fr;
  }

  .driver-profile,
  .qr-panel,
  .section-row {
    align-items: flex-start;
  }
}
</style>
