<script setup lang="ts">
import { computed, ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import EmptyState from '@/components/EmptyState.vue';
import StatusPill from '@/components/StatusPill.vue';
import { loadIslandCruiseFlow, lockIslandCruiseOrder, saleIslandCruiseOrder } from '@/api/island-cruise';
import { minFareLabel, totalStock } from '@/utils/island-cruise-mappers';
import type { IslandCruiseStep, IslandCruiseViewModel, IslandLockedOrder, IslandTicketResult } from '@/types/island-cruise';

const viewModel = ref<IslandCruiseViewModel | null>(null);
const activeStep = ref<IslandCruiseStep>('detail');
const errorMessage = ref('');
const isSubmitting = ref(false);
const lockedOrder = ref<IslandLockedOrder | null>(null);
const ticket = ref<IslandTicketResult | null>(null);
const validationMessage = ref('');

const steps: Array<{ id: IslandCruiseStep; label: string }> = [
  { id: 'detail', label: '选班次' },
  { id: 'traveler', label: '实名' },
  { id: 'pay', label: '支付' },
  { id: 'ticket', label: '票券' }
];

const selectedVoyage = computed(() => {
  const data = viewModel.value;
  if (!data) return null;
  return data.voyages.find((voyage) => voyage.id === data.selectedVoyageId) ?? data.voyages[0] ?? null;
});

const selectedFare = computed(() => {
  const voyage = selectedVoyage.value;
  if (!voyage) return null;
  const draftFareId = viewModel.value?.draft.passengers[0]?.fareId;
  return voyage.fares.find((fare) => fare.id === draftFareId) ?? voyage.fares[0] ?? null;
});

onLoad((query) => {
  const step = normalizeStep(String(query?.step ?? 'detail'));
  activeStep.value = step;
  loadPage(String(query?.scenario ?? ''));
});

function normalizeStep(value: string): IslandCruiseStep {
  if (value === 'traveler' || value === 'pay' || value === 'ticket') return value;
  return 'detail';
}

function loadPage(scenarioOverride?: string) {
  loadIslandCruiseFlow(scenarioOverride)
    .then((result) => {
      viewModel.value = result.data;
      lockedOrder.value = result.data.lockedOrder;
      ticket.value = result.data.ticket;
      errorMessage.value = result.data.errorMessage ?? '';
      validationMessage.value = '';
    })
    .catch((error: Error) => {
      errorMessage.value = error.message;
    });
}

function selectVoyage(voyageId: number) {
  if (!viewModel.value) return;
  viewModel.value = {
    ...viewModel.value,
    selectedVoyageId: voyageId
  };
}

function goStep(step: IslandCruiseStep) {
  activeStep.value = step;
}

function validateDraft(): boolean {
  const draft = viewModel.value?.draft;
  if (!selectedVoyage.value || !selectedFare.value) {
    validationMessage.value = '请先选择可售班次。';
    return false;
  }
  if (!draft?.contactName || !draft.contactMobile) {
    validationMessage.value = '请填写联系人姓名和手机号。';
    return false;
  }
  if (!draft.passengers.length) {
    validationMessage.value = '请至少填写 1 位实名出行人。';
    return false;
  }
  const missingIndex = draft.passengers.findIndex((passenger) => !passenger.name || !passenger.mobile || !passenger.certNo);
  if (missingIndex >= 0) {
    validationMessage.value = `请补全乘客 ${missingIndex + 1} 的实名信息。`;
    return false;
  }
  if (selectedFare.value.stock < draft.quantity) {
    validationMessage.value = '当前班次余票不足，请更换班次或减少人数。';
    return false;
  }
  validationMessage.value = '';
  return true;
}

async function confirmOrder() {
  if (!viewModel.value || !validateDraft()) return;
  isSubmitting.value = true;
  const result = await lockIslandCruiseOrder(viewModel.value);
  lockedOrder.value = result.data;
  isSubmitting.value = false;
  activeStep.value = 'pay';
}

function cancelLock() {
  lockedOrder.value = null;
  activeStep.value = 'traveler';
}

async function payAndIssueTicket() {
  if (!viewModel.value || !lockedOrder.value) {
    validationMessage.value = '请先确认订单并保留座位。';
    return;
  }
  if (lockedOrder.value.status === 'expired') {
    validationMessage.value = '订单已超时，请重新选择班次。';
    return;
  }
  isSubmitting.value = true;
  const result = await saleIslandCruiseOrder(viewModel.value);
  ticket.value = result.data;
  isSubmitting.value = false;
  activeStep.value = 'ticket';
}

function money(value: number): string {
  return `¥${Number.isInteger(value) ? value : value.toFixed(2)}`;
}
</script>

<template>
  <view class="page-shell island-page">
    <view class="step-tabs">
      <button
        v-for="step in steps"
        :key="step.id"
        class="step-tabs__item"
        :class="{ 'step-tabs__item--active': activeStep === step.id }"
        @click="goStep(step.id)"
      >
        {{ step.label }}
      </button>
    </view>

    <EmptyState v-if="errorMessage" title="环岛游服务暂不可用" :description="errorMessage" />
    <template v-else-if="viewModel">
      <view v-if="activeStep === 'detail'" class="step-panel">
        <view class="island-hero">
          <image class="island-hero__image" :src="viewModel.hero.imageUrl" mode="aspectFill" />
          <view class="island-hero__content">
            <StatusPill text="实时班次 Mock" tone="success" />
            <text class="island-hero__title">{{ viewModel.hero.title }}</text>
            <text class="island-hero__summary">{{ viewModel.hero.summary }}</text>
          </view>
        </view>

        <view class="card route-card">
          <view>
            <text class="muted-label">上船码头</text>
            <text class="route-card__name">{{ viewModel.ports.up.name }}</text>
          </view>
          <text class="route-card__arrow">→</text>
          <view>
            <text class="muted-label">航线体验</text>
            <text class="route-card__name">{{ viewModel.ports.down.name }}</text>
          </view>
        </view>

        <view class="section card recommend-card">
          <view class="section-row">
            <text class="section-title">近期推荐</text>
            <text class="section-link">{{ viewModel.recommended ? viewModel.recommended.date : '暂无可订' }}</text>
          </view>
          <view v-if="viewModel.recommended">
            <text class="recommend-card__title">{{ viewModel.recommended.label }}</text>
            <text class="recommend-card__body">
              {{ viewModel.recommended.firstTime }} 开航 · {{ viewModel.recommended.count }} 个班次 · {{ viewModel.recommended.minPriceLabel }}
            </text>
          </view>
          <EmptyState v-else title="暂无推荐班次" description="当前日期或航线暂无可售班次，可切换日期后重试。" />
        </view>

        <view class="section">
          <view class="section-row">
            <text class="section-title">可售班次</text>
            <text class="section-link">{{ viewModel.voyages.length }} 个班次</text>
          </view>
          <view v-if="viewModel.voyages.length === 0" class="list-stack">
            <EmptyState title="暂无班次" description="当前日期暂无可售班次，请更换日期或航线。" />
          </view>
          <view v-else class="list-stack">
            <button
              v-for="voyage in viewModel.voyages"
              :key="voyage.id"
              class="card voyage-card"
              :class="{ 'voyage-card--active': selectedVoyage?.id === voyage.id }"
              @click="selectVoyage(voyage.id)"
            >
              <view class="voyage-card__top">
                <view>
                  <text class="voyage-card__time">{{ voyage.departureTime }} 开航</text>
                  <text class="voyage-card__name">{{ voyage.name }} · {{ voyage.shipName }}</text>
                </view>
                <view class="voyage-card__price">
                  <text>{{ minFareLabel(voyage) }}</text>
                  <text>余 {{ totalStock(voyage) }}</text>
                </view>
              </view>
              <view class="fare-grid">
                <view v-for="fare in voyage.fares" :key="fare.id" class="fare-chip">
                  <text>{{ fare.cabinClassName }} {{ fare.name }}</text>
                  <text>{{ money(fare.price) }} · 余{{ fare.stock }}</text>
                </view>
              </view>
            </button>
          </view>
        </view>

        <view class="section card note-card">
          <text class="section-title">预订须知</text>
          <text v-for="note in viewModel.serviceNotes" :key="note" class="note-card__item">{{ note }}</text>
        </view>

        <view class="bottom-action">
          <view>
            <text class="muted-label">所选票价</text>
            <text class="bottom-action__price">{{ selectedFare ? money(selectedFare.price) : '¥--' }}</text>
          </view>
          <button class="primary-button bottom-action__button" :disabled="!selectedVoyage" @click="goStep('traveler')">立即预订</button>
        </view>
      </view>

      <view v-if="activeStep === 'traveler'" class="step-panel">
        <view class="section page-head">
          <text class="page-title">填写出行人</text>
          <text class="page-subtitle">{{ selectedVoyage ? `${selectedVoyage.departureTime} · ${selectedVoyage.name}` : '未选择班次' }}</text>
        </view>

        <view class="card form-card">
          <text class="section-title">票种人数</text>
          <view class="summary-row">
            <text>{{ selectedFare ? `${selectedFare.cabinClassName} ${selectedFare.name}` : '未选择票种' }}</text>
            <text>{{ viewModel.draft.quantity }} 人 · {{ money(viewModel.draft.totalAmount) }}</text>
          </view>
        </view>

        <view class="card form-card">
          <text class="section-title">联系人信息</text>
          <view class="field-grid">
            <view class="field-box">
              <text>联系人</text>
              <input v-model="viewModel.draft.contactName" placeholder="请填写联系人姓名" />
            </view>
            <view class="field-box">
              <text>手机号</text>
              <input v-model="viewModel.draft.contactMobile" type="number" placeholder="接收出票短信" />
            </view>
          </view>
        </view>

        <view class="card form-card">
          <text class="section-title">出行人实名</text>
          <view v-if="viewModel.draft.passengers.length === 0">
            <EmptyState title="实名信息不完整" description="当前 Mock 场景模拟乘客校验失败，请补齐乘客后再提交。" />
          </view>
          <view v-for="(passenger, index) in viewModel.draft.passengers" :key="index" class="passenger-card">
            <text class="passenger-card__title">乘客 {{ index + 1 }} · {{ passenger.certTypeName }}</text>
            <view class="field-grid">
              <view class="field-box">
                <text>姓名</text>
                <input v-model="passenger.name" placeholder="真实姓名" />
              </view>
              <view class="field-box">
                <text>手机号</text>
                <input v-model="passenger.mobile" type="number" placeholder="乘客手机号" />
              </view>
              <view class="field-box field-box--wide">
                <text>证件号码</text>
                <input v-model="passenger.certNo" placeholder="用于实名核验" />
              </view>
            </view>
          </view>
        </view>

        <EmptyState v-if="validationMessage" title="请检查订单信息" :description="validationMessage" />

        <view class="bottom-action">
          <view>
            <text class="muted-label">合计</text>
            <text class="bottom-action__price">{{ money(viewModel.draft.totalAmount) }}</text>
          </view>
          <button class="primary-button bottom-action__button" :disabled="isSubmitting" @click="confirmOrder">
            {{ isSubmitting ? '正在保留' : '确认订单' }}
          </button>
        </view>
      </view>

      <view v-if="activeStep === 'pay'" class="step-panel">
        <view class="section page-head">
          <text class="page-title">支付订单</text>
          <text class="page-subtitle">开发期仅模拟支付和出票，不触发真实微信支付。</text>
        </view>

        <view v-if="lockedOrder" class="card pay-card">
          <view class="section-row">
            <text class="section-title">座位已保留</text>
            <StatusPill :text="lockedOrder.status === 'expired' ? '已超时' : '待支付'" :tone="lockedOrder.status === 'expired' ? 'warning' : 'success'" />
          </view>
          <text class="pay-card__amount">{{ money(lockedOrder.amount) }}</text>
          <view class="summary-row">
            <text>订单号</text>
            <text>{{ lockedOrder.localOrderNo }}</text>
          </view>
          <view class="summary-row">
            <text>票号</text>
            <text>{{ lockedOrder.ticketNo }}</text>
          </view>
          <view class="summary-row">
            <text>支付截止</text>
            <text>{{ lockedOrder.expireTimeLabel }}</text>
          </view>
        </view>
        <EmptyState v-else title="暂无锁座订单" description="请先返回实名页确认订单。" />

        <view class="card flow-card">
          <text class="section-title">购票进度</text>
          <view class="flow-step"><text>1</text><view><strong>选择班次</strong><span>已确认可售时间、舱位和票价。</span></view></view>
          <view class="flow-step"><text>2</text><view><strong>保留座位</strong><span>当前订单已生成，请在截止时间前支付。</span></view></view>
          <view class="flow-step"><text>3</text><view><strong>支付出票</strong><span>支付成功后生成电子票，可凭票检票登船。</span></view></view>
        </view>

        <EmptyState v-if="validationMessage" title="支付提示" :description="validationMessage" />

        <view class="bottom-action bottom-action--two">
          <button class="secondary-button bottom-action__button" @click="cancelLock">取消订单</button>
          <button class="primary-button bottom-action__button" :disabled="isSubmitting" @click="payAndIssueTicket">
            {{ isSubmitting ? '正在出票' : '立即支付并出票' }}
          </button>
        </view>
      </view>

      <view v-if="activeStep === 'ticket'" class="step-panel">
        <view class="ticket-state">
          <text class="ticket-state__title">{{ ticket?.status === 'sale_failed' ? '出票失败' : '待使用' }}</text>
          <text class="ticket-state__body">
            {{ ticket?.failureMessage ?? '请按所选开航时间到湾仔旅游码头检票登船。' }}
          </text>
        </view>

        <view v-if="ticket" class="card island-ticket">
          <view class="ticket-qr"><text></text></view>
          <text class="ticket-code">核销码 {{ ticket.maskedCode }}</text>
          <view class="summary">
            <view class="summary-row"><text>订单号</text><text>{{ ticket.localOrderNo }}</text></view>
            <view class="summary-row"><text>票号</text><text>{{ ticket.ticketNo }}</text></view>
            <view class="summary-row"><text>航班</text><text>{{ ticket.voyageLabel }}</text></view>
            <view class="summary-row"><text>乘客</text><text>{{ ticket.passengerLabel }}</text></view>
            <view class="summary-row"><text>付款时间</text><text>{{ ticket.paidAtLabel }}</text></view>
          </view>
        </view>
        <EmptyState v-else title="暂无电子票" description="出票后会在这里展示电子票核销信息。" />

        <view class="card service-card">
          <text class="section-title">售后服务</text>
          <text class="service-card__body">退票、改签入口已保留；真实费用试算、库存回滚和供应商退改签闭环留到后续接入。</text>
          <view class="service-actions">
            <button class="secondary-button">申请退票</button>
            <button class="secondary-button">申请改签</button>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<style scoped>
.island-page {
  padding-top: 12px;
}

.step-tabs {
  position: sticky;
  top: 0;
  z-index: 3;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  padding: 8px 0 12px;
  background: var(--zt-color-bg);
}

.step-tabs__item {
  height: 34px;
  margin: 0;
  border-radius: 8px;
  background: var(--zt-color-surface);
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.step-tabs__item--active {
  background: var(--zt-color-primary);
  color: #ffffff;
}

.step-panel {
  padding-bottom: 88px;
}

.island-hero {
  position: relative;
  min-height: 260px;
  overflow: hidden;
  border-radius: 8px;
  background: #0f766e;
}

.island-hero__image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.island-hero__content {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 12px;
  min-height: 260px;
  box-sizing: border-box;
  padding: 24px;
  color: #ffffff;
  background: linear-gradient(180deg, rgba(15, 28, 46, 0.12), rgba(15, 28, 46, 0.74));
}

.island-hero__title {
  font-size: 30px;
  font-weight: 900;
}

.island-hero__summary,
.page-subtitle,
.recommend-card__body,
.note-card__item,
.service-card__body {
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.route-card,
.recommend-card,
.note-card,
.form-card,
.pay-card,
.flow-card,
.island-ticket,
.service-card {
  padding: 16px;
}

.route-card {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 32px minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  margin-top: 14px;
}

.route-card__name,
.recommend-card__title,
.page-title,
.ticket-state__title {
  display: block;
  color: var(--zt-color-text);
  font-weight: 900;
}

.route-card__name {
  margin-top: 4px;
  font-size: 17px;
}

.route-card__arrow {
  color: var(--zt-color-primary);
  font-size: 22px;
  font-weight: 900;
  text-align: center;
}

.muted-label,
.section-link {
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.section-row,
.summary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.recommend-card__title {
  font-size: 18px;
}

.list-stack {
  display: grid;
  gap: 12px;
}

.voyage-card {
  width: 100%;
  margin: 0;
  padding: 16px;
  text-align: left;
}

.voyage-card--active {
  border-color: var(--zt-color-primary);
}

.voyage-card__top {
  display: flex;
  justify-content: space-between;
  gap: 14px;
}

.voyage-card__time,
.voyage-card__name,
.voyage-card__price text,
.fare-chip text,
.note-card__item,
.page-title,
.page-subtitle,
.passenger-card__title,
.ticket-code,
.service-card__body {
  display: block;
}

.voyage-card__time {
  font-size: 18px;
  font-weight: 900;
}

.voyage-card__name,
.voyage-card__price text:last-child,
.fare-chip text:last-child {
  margin-top: 4px;
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.voyage-card__price {
  flex: 0 0 auto;
  text-align: right;
}

.voyage-card__price text:first-child {
  color: var(--zt-color-accent);
  font-size: 16px;
  font-weight: 900;
}

.fare-grid,
.field-grid,
.service-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 12px;
}

.fare-chip,
.field-box {
  padding: 12px;
  border-radius: 8px;
  background: var(--zt-color-surface-muted);
}

.fare-chip text:first-child,
.field-box text,
.summary-row text:last-child {
  font-weight: 800;
}

.page-title {
  font-size: 28px;
}

.field-box input {
  margin-top: 8px;
  font-size: 15px;
}

.field-box--wide {
  grid-column: 1 / -1;
}

.passenger-card {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--zt-color-border);
}

.bottom-action {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 5;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(140px, 220px);
  gap: 12px;
  box-sizing: border-box;
  padding: 12px 16px calc(12px + env(safe-area-inset-bottom));
  border-top: 1px solid var(--zt-color-border);
  background: rgba(255, 255, 255, 0.96);
}

.bottom-action--two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.bottom-action__price,
.pay-card__amount {
  display: block;
  color: var(--zt-color-accent);
  font-size: 24px;
  font-weight: 900;
}

.bottom-action__button {
  width: 100%;
}

.pay-card__amount {
  margin: 14px 0;
  font-size: 34px;
}

.flow-step {
  display: grid;
  grid-template-columns: 32px minmax(0, 1fr);
  gap: 10px;
  padding: 10px 0;
}

.flow-step text {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: var(--zt-color-primary-weak);
  color: var(--zt-color-primary);
  font-weight: 900;
}

.flow-step strong,
.flow-step span {
  display: block;
}

.flow-step span {
  margin-top: 4px;
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.ticket-state {
  padding: 18px 2px 8px;
}

.ticket-state__title {
  font-size: 32px;
}

.ticket-state__body {
  display: block;
  margin-top: 8px;
  color: var(--zt-color-text-muted);
  font-size: 14px;
  line-height: 22px;
}

.ticket-qr {
  width: 132px;
  height: 132px;
  margin: 0 auto 12px;
  border: 8px solid #152238;
  border-radius: 8px;
  background:
    linear-gradient(90deg, #152238 18px, transparent 18px 30px, #152238 30px 48px, transparent 48px),
    linear-gradient(#152238 18px, transparent 18px 30px, #152238 30px 48px, transparent 48px),
    #ffffff;
}

.ticket-code {
  text-align: center;
  font-weight: 900;
}

.summary {
  margin-top: 14px;
}

.summary-row {
  padding: 10px 0;
  border-top: 1px solid var(--zt-color-border);
  color: var(--zt-color-text-muted);
  font-size: 14px;
}

.summary-row text:last-child {
  color: var(--zt-color-text);
  text-align: right;
}

.service-actions button {
  width: 100%;
}

@media (max-width: 640px) {
  .field-grid,
  .fare-grid,
  .service-actions,
  .bottom-action,
  .bottom-action--two {
    grid-template-columns: 1fr;
  }

  .route-card {
    grid-template-columns: 1fr;
  }

  .route-card__arrow {
    text-align: left;
  }
}

@media (min-width: 900px) {
  .bottom-action {
    position: static;
    margin-top: 16px;
  }
}
</style>
