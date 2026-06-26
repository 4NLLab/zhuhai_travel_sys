<script setup lang="ts">
import { computed, ref } from 'vue';
import StatusPill from '@/components/StatusPill.vue';
import { loadPhaseOneStatus } from '@/api/phase-one';
import { getRuntimeConfig } from '@/utils/config';
import { navigationAdapter } from '@/adapters/navigation';

const config = getRuntimeConfig();
const isLoading = ref(true);
const title = ref('正在加载平台护栏');
const summary = ref('请稍候');
const adapters = ref<string[]>([]);
const errorMessage = ref('');

loadPhaseOneStatus()
  .then((result) => {
    title.value = result.data.title;
    summary.value = result.data.summary;
    adapters.value = result.data.enabledAdapters;
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
    title.value = '平台状态读取失败';
    summary.value = '请检查 API 模式或 Mock 场景配置。';
  })
  .finally(() => {
    isLoading.value = false;
  });

const adapterText = computed(() => (adapters.value.length ? adapters.value.join(' / ') : '待业务切片接入'));

function openDriverPlaceholder() {
  navigationAdapter.navigateTo('/src/subpackages/driver/pages/home/index');
}
</script>

<template>
  <view class="page-shell home-page">
    <view class="hero">
      <StatusPill text="Phase 1" tone="success" />
      <text class="hero__title">珠海文旅小程序</text>
      <text class="hero__summary">{{ summary }}</text>
    </view>

    <view class="section">
      <text class="section-title">平台状态</text>
      <view class="card status-card">
        <text class="status-card__title">{{ isLoading ? '加载中' : title }}</text>
        <text class="status-card__line">API 模式：{{ config.apiMode }}</text>
        <text class="status-card__line">Mock 场景：{{ config.mockScenarioId }}</text>
        <text class="status-card__line">当前平台：{{ config.platform }}</text>
        <text class="status-card__line">已接入 adapter：{{ adapterText }}</text>
        <text v-if="errorMessage" class="status-card__error">{{ errorMessage }}</text>
      </view>
    </view>

    <view class="section grid">
      <view class="card feature-card">
        <text class="feature-card__title">首页主包</text>
        <text class="feature-card__body">Phase 2 迁移品牌、分类、产品与环岛游入口。</text>
      </view>
      <view class="card feature-card">
        <text class="feature-card__title">分包预留</text>
        <text class="feature-card__body">司机端按低频入口放入独立分包。</text>
      </view>
    </view>

    <button class="primary-button driver-button" @click="openDriverPlaceholder">司机端占位</button>
  </view>
</template>

<style scoped>
.hero {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 28px 20px;
  border-radius: 8px;
  background: linear-gradient(135deg, #e9f7f2 0%, #fff7ed 100%);
}

.hero__title {
  font-size: 26px;
  line-height: 34px;
  font-weight: 800;
}

.hero__summary,
.feature-card__body,
.status-card__line {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}

.status-card,
.feature-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
}

.status-card__title,
.feature-card__title {
  font-size: 16px;
  font-weight: 700;
}

.status-card__error {
  color: var(--zt-color-danger);
  font-size: 14px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.driver-button {
  width: 100%;
  margin-top: 20px;
}

@media (max-width: 520px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
