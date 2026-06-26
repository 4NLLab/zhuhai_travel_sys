<script setup lang="ts">
import { computed, ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import EmptyState from '@/components/EmptyState.vue';
import { navigationAdapter } from '@/adapters/navigation';
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
  { id: 'detail', label: '详情' },
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

const bottomPrice = computed(() => {
  if (activeStep.value === 'traveler') return money(viewModel.value?.draft.totalAmount ?? 0);
  if (activeStep.value === 'pay') return lockedOrder.value ? money(lockedOrder.value.amount) : '¥--';
  return selectedFare.value ? money(selectedFare.value.price) : '¥--';
});

onLoad((query) => {
  activeStep.value = normalizeStep(String(query?.step ?? 'detail'));
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

function openHome() {
  navigationAdapter.switchTab('/pages/home/index');
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
  <view class="island-page">
    <EmptyState v-if="errorMessage" title="环岛游服务暂不可用" :description="errorMessage" />
    <template v-else-if="viewModel">
      <view v-if="activeStep === 'detail'" class="detail-screen">
        <view class="detail-hero">
          <image class="media-image" src="/static/phase2/macau-cruise-night-banner-web.jpg" mode="aspectFill" />
          <view class="detail-top">
            <button class="nav-chip" aria-label="返回首页" @click="openHome">‹ 首页</button>
            <button class="nav-chip">▱ 海上夜游</button>
          </view>
          <view class="media-control"><text>◌</text><text>0:00</text><text>⋮</text></view>
        </view>

        <view class="detail-copy card-dark">
          <text class="eyebrow">珠海湾仔出发 · 海上看澳门</text>
          <text class="detail-title">{{ viewModel.hero.title }}</text>
          <text class="detail-body">{{ viewModel.hero.summary }}</text>
          <view class="hero-tags">
            <text>湾仔码头登船</text>
            <text>电子票验票</text>
            <text>夜景体验</text>
            <text>实时班次</text>
          </view>
        </view>

        <view class="action-card">
          <view class="action-head">
            <view>
              <text class="action-title">选择你的海上时段</text>
              <text class="action-body">可先咨询航线与天气，也可直接进入购票</text>
            </view>
            <text class="status-pill">在线可订</text>
          </view>
          <view class="route-box">
            <view class="route-node"><text>上船码头</text><strong>{{ viewModel.ports.up.name }}</strong></view>
            <text class="route-arrow">→</text>
            <view class="route-node"><text>航线体验</text><strong>{{ viewModel.ports.down.name }}</strong></view>
          </view>
          <view class="action-buttons">
            <button class="ghost-button">咨询客服</button>
            <button class="gold-button" @click="goStep('traveler')">立即购票</button>
          </view>
        </view>

        <view class="card-dark">
          <view class="section-head"><text>体验亮点</text><text>环岛游推荐</text></view>
          <view class="story-grid">
            <view><strong>海上看澳门</strong><span>沿湾仔水道近距离欣赏澳门天际线。</span></view>
            <view><strong>电子票核销</strong><span>购票后凭码检票，减少排队等待。</span></view>
          </view>
        </view>
      </view>

      <view v-else class="booking-screen">
        <view class="booking-top">
          <button class="icon-button" @click="goStep('detail')">‹</button>
          <text>澳门环岛游</text>
          <button class="icon-button">⌯</button>
        </view>

        <view class="video-card">
          <image class="media-image" src="/static/phase2/macau-cruise-night-banner-web.jpg" mode="aspectFill" />
          <view class="media-control"><text>◌</text><text>0:00</text><text>⋮</text></view>
        </view>

        <view class="hero-copy card-dark">
          <text class="detail-title">澳门环岛游</text>
          <text class="detail-body">湾仔码头登船，从海上看澳门天际线、城市灯影与珠澳湾区夜色。</text>
          <view class="hero-tags"><text>约 90 分钟</text><text>电子票</text><text>实时班次</text></view>
        </view>

        <view class="quick-metrics">
          <view><text>默认上船</text><strong>{{ viewModel.ports.up.name }}</strong></view>
          <view><text>推荐体验</text><strong>环岛/夜游</strong></view>
          <view><text>购票方式</text><strong>在线出票</strong></view>
        </view>

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

        <view v-if="activeStep === 'traveler'" class="card-dark form-card">
          <view class="section-head"><text>实名信息</text><text>{{ selectedVoyage ? selectedVoyage.departureTime : '待选班次' }}</text></view>
          <view class="field-grid">
            <view class="field-box"><text>联系人</text><input v-model="viewModel.draft.contactName" placeholder="请填写联系人姓名" /></view>
            <view class="field-box"><text>手机号</text><input v-model="viewModel.draft.contactMobile" type="number" placeholder="接收出票短信" /></view>
          </view>
          <view v-for="(passenger, index) in viewModel.draft.passengers" :key="index" class="passenger-card">
            <text class="passenger-title">乘客 {{ index + 1 }} · {{ passenger.certTypeName }}</text>
            <view class="field-grid">
              <view class="field-box"><text>姓名</text><input v-model="passenger.name" placeholder="真实姓名" /></view>
              <view class="field-box"><text>手机号</text><input v-model="passenger.mobile" type="number" placeholder="乘客手机号" /></view>
              <view class="field-box field-box--wide"><text>证件号码</text><input v-model="passenger.certNo" placeholder="用于实名核验" /></view>
            </view>
          </view>
        </view>

        <view v-if="activeStep === 'pay'" class="card-dark form-card">
          <view class="section-head"><text>支付订单</text><text>{{ lockedOrder?.status === 'expired' ? '已超时' : '待支付' }}</text></view>
          <text class="pay-amount">{{ lockedOrder ? money(lockedOrder.amount) : '¥--' }}</text>
          <view class="summary-row"><text>订单号</text><text>{{ lockedOrder?.localOrderNo ?? '请先确认订单' }}</text></view>
          <view class="summary-row"><text>票号</text><text>{{ lockedOrder?.ticketNo ?? '--' }}</text></view>
          <view class="summary-row"><text>支付截止</text><text>{{ lockedOrder?.expireTimeLabel ?? '--' }}</text></view>
        </view>

        <view v-if="activeStep === 'ticket'" class="card-dark island-ticket">
          <view class="section-head"><text>{{ ticket?.status === 'sale_failed' ? '出票失败' : '待使用' }}</text><text>电子票</text></view>
          <view class="ticket-qr"><text></text></view>
          <text class="ticket-code">核销码 {{ ticket?.maskedCode ?? '0105****3806' }}</text>
          <view class="summary-row"><text>航班</text><text>{{ ticket?.voyageLabel ?? selectedVoyage?.name }}</text></view>
          <view class="summary-row"><text>乘客</text><text>{{ ticket?.passengerLabel ?? '珠海游客' }}</text></view>
        </view>

        <view class="card-dark">
          <view class="section-head"><text>可售班次</text><text>{{ viewModel.voyages.length }} 个班次</text></view>
          <view v-if="viewModel.voyages.length === 0">
            <EmptyState title="暂无班次" description="当前日期暂无可售班次，请更换日期或航线。" />
          </view>
          <view v-else class="voyage-list">
            <button
              v-for="voyage in viewModel.voyages"
              :key="voyage.id"
              class="voyage-card"
              :class="{ 'voyage-card--active': selectedVoyage?.id === voyage.id }"
              @click="selectVoyage(voyage.id)"
            >
              <view>
                <strong>{{ voyage.departureTime }} · {{ voyage.name }}</strong>
                <span>{{ voyage.shipName }} · 航班号 {{ voyage.voyageNo }}</span>
              </view>
              <text>{{ minFareLabel(voyage) }} · 余 {{ totalStock(voyage) }}</text>
            </button>
          </view>
        </view>

        <EmptyState v-if="validationMessage" title="流程提示" :description="validationMessage" />
      </view>
    </template>

    <view class="bottom-bar" :class="{ 'bottom-bar--two': activeStep === 'pay' }">
      <view>
        <text>{{ activeStep === 'detail' ? '选择班次后显示票价' : '合计' }}</text>
        <strong>{{ bottomPrice }}</strong>
      </view>
      <button v-if="activeStep === 'detail'" class="cyan-button" @click="goStep('traveler')">立即预订</button>
      <button v-else-if="activeStep === 'traveler'" class="cyan-button" :disabled="isSubmitting" @click="confirmOrder">{{ isSubmitting ? '正在保留' : '确认订单' }}</button>
      <button v-else-if="activeStep === 'pay'" class="dark-button" @click="cancelLock">取消订单</button>
      <button v-if="activeStep === 'pay'" class="cyan-button" :disabled="isSubmitting" @click="payAndIssueTicket">{{ isSubmitting ? '正在出票' : '立即支付' }}</button>
      <button v-else-if="activeStep === 'ticket'" class="cyan-button">查看票券</button>
    </view>
  </view>
</template>

<style scoped>
.island-page {
  width: min(100vw, 430px);
  min-height: 100vh;
  margin: clamp(0px, calc((100vw - 430px) * 100), 28px) auto;
  padding-bottom: calc(100px + env(safe-area-inset-bottom));
  overflow-x: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px);
  background:
    radial-gradient(circle at 50% -10%, rgba(60, 198, 216, 0.22), transparent 32%),
    linear-gradient(180deg, #02070d 0%, #06131d 48%, #0a1821 100%);
  color: #f4fbff;
  box-shadow: 0 0 0 1px rgba(132, 189, 211, 0.1), 0 24px 80px rgba(0, 0, 0, 0.42);
}

button {
  margin: 0;
}

button::after {
  display: none;
}

.detail-hero,
.video-card {
  position: relative;
  overflow: hidden;
  background: #02070d;
}

.detail-hero {
  height: 282px;
  border-radius: 0 0 18px 18px;
}

.video-card {
  height: 264px;
  border-radius: 0 0 18px 18px;
}

.media-image {
  width: 100%;
  height: 100%;
  display: block;
}

.detail-top {
  position: absolute;
  left: 14px;
  right: 14px;
  top: 14px;
  z-index: 2;
  display: flex;
  justify-content: space-between;
}

.nav-chip,
.icon-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
  padding: 0 11px;
  border-radius: 999px;
  color: #ffffff;
  background: rgba(2, 8, 13, 0.54);
  backdrop-filter: blur(12px);
  font-size: 13px;
  font-weight: 850;
}

.media-control {
  position: absolute;
  left: 28px;
  right: 18px;
  bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #ffffff;
  font-weight: 900;
}

.card-dark,
.action-card {
  margin: 0 14px 10px;
  padding: 14px;
  border: 1px solid rgba(151, 190, 205, 0.2);
  border-radius: 8px;
  background: rgba(12, 31, 44, 0.9);
  box-shadow: 0 18px 42px rgba(0, 0, 0, 0.26);
}

.detail-copy {
  position: relative;
  z-index: 2;
  margin-top: -18px;
  padding: 18px 16px;
  background: linear-gradient(135deg, rgba(20, 48, 66, 0.98), rgba(9, 23, 35, 0.96));
}

.eyebrow {
  display: inline-block;
  margin-bottom: 10px;
  padding: 5px 9px;
  border-radius: 999px;
  color: #fff1c9;
  background: rgba(234, 192, 110, 0.16);
  font-size: 12px;
  font-weight: 850;
}

.detail-title,
.detail-body,
.action-title,
.action-body,
.route-node text,
.route-node strong,
.section-head text,
.quick-metrics text,
.quick-metrics strong,
.story-grid strong,
.story-grid span,
.field-box text,
.passenger-title,
.summary-row text,
.ticket-code {
  display: block;
}

.detail-title {
  font-size: 32px;
  font-weight: 900;
  line-height: 1.08;
}

.detail-body {
  max-width: 320px;
  margin-top: 10px;
  color: rgba(228, 244, 250, 0.86);
  font-size: 13px;
  line-height: 1.58;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
  margin-top: 14px;
}

.hero-tags text {
  padding: 6px 9px;
  border: 1px solid rgba(75, 208, 226, 0.22);
  border-radius: 999px;
  color: #d8f9ff;
  background: rgba(75, 208, 226, 0.12);
  font-size: 11px;
  font-weight: 850;
}

.action-card {
  padding: 0;
  overflow: hidden;
}

.action-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 14px;
  color: #07121a;
  background: linear-gradient(135deg, #eac06e, #f1d99d);
}

.action-title {
  font-size: 17px;
  font-weight: 900;
}

.action-body,
.route-node text,
.section-head text:last-child,
.quick-metrics text,
.story-grid span,
.voyage-card span {
  color: #9fb2bf;
  font-size: 12px;
  line-height: 1.42;
}

.status-pill {
  flex: 0 0 auto;
  padding: 6px 9px;
  border-radius: 999px;
  color: #082331;
  background: rgba(255, 255, 255, 0.74);
  font-size: 12px;
  font-weight: 900;
}

.route-box {
  display: grid;
  grid-template-columns: 1fr 38px 1fr;
  gap: 8px;
  align-items: center;
  padding: 14px;
}

.route-node {
  min-height: 64px;
  padding: 10px;
  border: 1px solid rgba(151, 190, 205, 0.2);
  border-radius: 8px;
  background: rgba(7, 21, 31, 0.68);
}

.route-node strong {
  margin-top: 5px;
  font-size: 17px;
}

.route-arrow {
  display: grid;
  place-items: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  color: #bff6ff;
  background: rgba(75, 208, 226, 0.18);
  font-size: 24px;
  font-weight: 900;
}

.action-buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  padding: 0 14px 14px;
}

.ghost-button,
.gold-button,
.cyan-button,
.dark-button {
  display: grid;
  place-items: center;
  min-height: 44px;
  border-radius: 999px;
  font-size: 16px;
  font-weight: 900;
}

.ghost-button {
  color: #dffbff;
  background: rgba(75, 208, 226, 0.14);
  border: 1px solid rgba(75, 208, 226, 0.24);
}

.gold-button {
  color: #07121a;
  background: linear-gradient(135deg, #f3d78c, #eac06e);
}

.cyan-button {
  color: #ffffff;
  background: #27afc4;
}

.dark-button {
  color: #e6f7fb;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(151, 190, 205, 0.2);
}

.section-head,
.summary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.section-head {
  margin-bottom: 12px;
  font-weight: 900;
}

.story-grid,
.quick-metrics,
.field-grid,
.voyage-list {
  display: grid;
  gap: 8px;
}

.story-grid {
  grid-template-columns: 1fr 1fr;
}

.story-grid view,
.quick-metrics view,
.recommend-card,
.field-box,
.passenger-card,
.voyage-card {
  padding: 12px;
  border: 1px solid rgba(151, 190, 205, 0.2);
  border-radius: 8px;
  background: rgba(7, 21, 31, 0.72);
}

.booking-top {
  position: sticky;
  top: 0;
  z-index: 20;
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) 42px;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: rgba(5, 15, 23, 0.82);
  backdrop-filter: blur(14px);
}

.booking-top text {
  text-align: center;
  font-size: 17px;
  font-weight: 900;
}

.icon-button {
  width: 38px;
  min-height: 38px;
  padding: 0;
}

.hero-copy {
  margin-top: 12px;
}

.quick-metrics {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin: 0 14px 10px;
}

.quick-metrics view {
  min-height: 78px;
}

.quick-metrics strong {
  margin-top: 7px;
  font-size: 16px;
}

.step-tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  margin: 0 14px 10px;
}

.step-tabs__item {
  min-height: 34px;
  border-radius: 8px;
  color: #9fb2bf;
  background: rgba(255, 255, 255, 0.08);
  font-size: 13px;
}

.step-tabs__item--active {
  color: #06131d;
  background: #3cc6d8;
  font-weight: 900;
}

.field-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.field-box input {
  width: 100%;
  height: 42px;
  margin-top: 8px;
  padding: 0 10px;
  box-sizing: border-box;
  border-radius: 8px;
  color: #f4fbff;
  background: rgba(255, 255, 255, 0.08);
  font-size: 14px;
}

.field-box--wide {
  grid-column: 1 / -1;
}

.passenger-card {
  margin-top: 10px;
}

.passenger-title {
  margin-bottom: 10px;
  font-weight: 900;
}

.pay-amount {
  display: block;
  margin-bottom: 12px;
  color: #f0c978;
  font-size: 34px;
  font-weight: 900;
}

.summary-row {
  padding: 10px 0;
  border-top: 1px solid rgba(151, 190, 205, 0.2);
  color: #9fb2bf;
  font-size: 13px;
}

.summary-row text:last-child {
  color: #f4fbff;
  font-weight: 900;
  text-align: right;
}

.ticket-qr {
  width: 150px;
  height: 150px;
  margin: 8px auto 14px;
  border: 10px solid #ffffff;
  background:
    linear-gradient(90deg, #000 14px, transparent 14px 26px, #000 26px 40px, transparent 40px) 0 0 / 58px 58px,
    linear-gradient(#000 14px, transparent 14px 26px, #000 26px 40px, transparent 40px) 0 0 / 58px 58px,
    #ffffff;
}

.ticket-code {
  text-align: center;
  font-size: 18px;
  font-weight: 900;
}

.voyage-card {
  width: 100%;
  margin: 0;
  text-align: left;
}

.voyage-card--active {
  border-color: #3cc6d8;
}

.voyage-card strong {
  display: block;
  font-size: 15px;
}

.voyage-card > text {
  display: block;
  margin-top: 8px;
  color: #f0c978;
  font-size: 13px;
  font-weight: 900;
}

.bottom-bar {
  position: fixed;
  left: 50%;
  bottom: 0;
  z-index: 30;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 150px;
  gap: 12px;
  align-items: center;
  width: min(100vw, 430px);
  transform: translateX(-50%);
  padding: 12px 14px calc(12px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  background: rgba(3, 12, 19, 0.96);
  border-top: 1px solid rgba(151, 190, 205, 0.18);
}

.bottom-bar--two {
  grid-template-columns: minmax(0, 1fr) 104px 130px;
}

.bottom-bar text,
.bottom-bar strong {
  display: block;
}

.bottom-bar text {
  color: #9fb2bf;
  font-size: 12px;
}

.bottom-bar strong {
  margin-top: 4px;
  color: #f0c978;
  font-size: 22px;
}
</style>
