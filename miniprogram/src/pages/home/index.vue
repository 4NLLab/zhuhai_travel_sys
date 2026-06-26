<script setup lang="ts">
import { computed, ref } from 'vue';
import ProductEntryCard from '@/components/ProductEntryCard.vue';
import StatusPill from '@/components/StatusPill.vue';
import EmptyState from '@/components/EmptyState.vue';
import { loadMainShellViewModel } from '@/api/main-shell';
import { navigationAdapter } from '@/adapters/navigation';
import type { MainShellViewModel, ProductCategory } from '@/types/main-shell';

const viewModel = ref<MainShellViewModel | null>(null);
const isLoading = ref(true);
const errorMessage = ref('');
const activeCategory = ref<ProductCategory | 'all'>('all');

const categoryTabs: Array<{ id: ProductCategory | 'all'; label: string }> = [
  { id: 'all', label: '推荐' },
  { id: 'ship', label: '船票' },
  { id: 'hotel', label: '酒店' },
  { id: 'tour', label: '港澳游' },
  { id: 'play', label: '玩乐' },
  { id: 'car', label: '接送' }
];

loadMainShellViewModel()
  .then((result) => {
    viewModel.value = result.data;
    errorMessage.value = result.data.errorMessage ?? '';
  })
  .catch((error: Error) => {
    errorMessage.value = error.message;
  })
  .finally(() => {
    isLoading.value = false;
  });

const filteredProducts = computed(() => {
  const products = viewModel.value?.products ?? [];
  if (activeCategory.value === 'all') return products;
  return products.filter((product) => product.category === activeCategory.value);
});

function openRoute(route?: string) {
  if (!route) return;
  navigationAdapter.navigateTo(route);
}
</script>

<template>
  <view class="page-shell home-page">
    <view class="home-hero">
      <view class="home-hero__media">
        <image class="home-hero__image" src="/static/phase2/macau-cruise-night-banner-web.jpg" mode="aspectFill" />
      </view>
      <view class="home-hero__content">
        <StatusPill text="小程序新增主包首页" tone="success" />
        <text class="home-hero__title">珠海湾游</text>
        <text class="home-hero__summary">海上岭南、珠澳假期和本地文旅服务统一入口。</text>
      </view>
    </view>

    <view class="search-shell">
      <text class="search-shell__icon">⌕</text>
      <text class="search-shell__placeholder">搜索服务名称、目的地或订单</text>
    </view>

    <view v-if="errorMessage" class="section">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="section">
      <view class="section-row">
        <text class="section-title">热门目的地</text>
        <text class="section-link">来源：index.html</text>
      </view>
      <view class="destination-grid">
        <button
          v-for="destination in viewModel?.destinations ?? []"
          :key="destination.id"
          class="destination-item"
          @click="openRoute(destination.route)"
        >
          <text class="destination-item__mark">{{ destination.title.slice(0, 1) }}</text>
          <view>
            <text class="destination-item__title">{{ destination.title }}</text>
            <text class="destination-item__subtitle">{{ destination.subtitle }}</text>
          </view>
        </button>
      </view>
    </view>

    <view class="section island-entry" @click="openRoute('/pages/island-cruise/index')">
      <view>
        <text class="island-entry__title">澳门环岛游</text>
        <text class="island-entry__body">进入环岛游订票切片，选择班次、实名购票并模拟出票。</text>
      </view>
      <text class="island-entry__action">立即预订</text>
    </view>

    <view class="section">
      <view class="section-row">
        <text class="section-title">精选商品</text>
        <text class="section-link">Mock：{{ viewModel?.scenarioId ?? '加载中' }}</text>
      </view>
      <scroll-view class="category-tabs" scroll-x>
        <button
          v-for="tab in categoryTabs"
          :key="tab.id"
          class="category-tab"
          :class="{ 'category-tab--active': activeCategory === tab.id }"
          @click="activeCategory = tab.id"
        >
          {{ tab.label }}
        </button>
      </scroll-view>

      <view v-if="isLoading" class="product-list">
        <EmptyState title="加载中" description="正在读取首页 Mock 场景。" />
      </view>
      <view v-else-if="filteredProducts.length === 0" class="product-list">
        <EmptyState title="暂无商品" description="当前筛选或 Mock 场景没有可展示商品。" />
      </view>
      <view v-else class="product-list">
        <ProductEntryCard
          v-for="product in filteredProducts"
          :key="product.id"
          :product="product"
          @click="openRoute(product.route)"
        />
      </view>
    </view>
  </view>
</template>

<style scoped>
.home-page {
  padding-top: 16px;
}

.home-hero {
  position: relative;
  min-height: 260px;
  overflow: hidden;
  border-radius: 8px;
  background: #126c91;
}

.home-hero__media,
.home-hero__image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.home-hero__content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 260px;
  box-sizing: border-box;
  justify-content: flex-end;
  padding: 24px;
  color: #ffffff;
  background: linear-gradient(180deg, rgba(9, 33, 42, 0.08), rgba(9, 33, 42, 0.72));
}

.home-hero__title {
  font-size: 30px;
  font-weight: 900;
}

.home-hero__summary {
  max-width: 520px;
  font-size: 15px;
  line-height: 24px;
}

.search-shell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 48px;
  margin-top: 14px;
  padding: 0 16px;
  border: 1px solid var(--zt-color-border);
  border-radius: 8px;
  background: var(--zt-color-surface);
}

.search-shell__icon {
  color: var(--zt-color-primary);
  font-size: 20px;
}

.search-shell__placeholder,
.section-link {
  color: var(--zt-color-text-muted);
  font-size: 13px;
}

.section-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.destination-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.destination-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-height: 72px;
  box-sizing: border-box;
  margin: 0;
  padding: 12px;
  border: 1px solid var(--zt-color-border);
  border-radius: 8px;
  background: var(--zt-color-surface);
  text-align: left;
}

.destination-item__mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: var(--zt-color-primary-weak);
  color: var(--zt-color-primary);
  font-weight: 900;
}

.destination-item__title,
.destination-item__subtitle {
  display: block;
}

.destination-item__title {
  font-size: 14px;
  font-weight: 800;
}

.destination-item__subtitle {
  margin-top: 2px;
  color: var(--zt-color-text-muted);
  font-size: 12px;
}

.island-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px;
  border-radius: 8px;
  background: linear-gradient(135deg, #0f766e, #b45309);
  color: #ffffff;
}

.island-entry__title,
.island-entry__body {
  display: block;
}

.island-entry__title {
  font-size: 18px;
  font-weight: 900;
}

.island-entry__body {
  margin-top: 6px;
  font-size: 13px;
  line-height: 20px;
}

.island-entry__action {
  flex: 0 0 auto;
  font-size: 13px;
  font-weight: 900;
}

.category-tabs {
  white-space: nowrap;
  margin-bottom: 12px;
}

.category-tab {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 68px;
  height: 36px;
  margin-right: 8px;
  margin-left: 0;
  padding: 0 14px;
  border-radius: 8px;
  color: var(--zt-color-text-muted);
  background: var(--zt-color-surface);
}

.category-tab--active {
  color: #ffffff;
  background: var(--zt-color-primary);
}

.product-list {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 760px) {
  .product-list {
    grid-template-columns: 1fr;
  }

  .destination-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
