<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import EmptyState from '@/components/EmptyState.vue';
import { loadMainShellViewModel } from '@/api/main-shell';
import { navigationAdapter } from '@/adapters/navigation';
import type { MainShellViewModel, ProductCategory } from '@/types/main-shell';

const viewModel = ref<MainShellViewModel | null>(null);
const isLoading = ref(true);
const errorMessage = ref('');
const activeCategory = ref<ProductCategory | 'all'>('all');
const activeHeroSlideIndex = ref(0);
let heroSlideTimer: ReturnType<typeof setInterval> | null = null;

const categoryTabs: Array<{ id: ProductCategory | 'all'; label: string }> = [
  { id: 'all', label: '推荐' },
  { id: 'ship', label: '船票' },
  { id: 'hotel', label: '酒店民宿' },
  { id: 'tour', label: '港澳游' },
  { id: 'play', label: '休闲娱乐' },
  { id: 'car', label: '接送租车' }
];

const destinationIconMap: Record<string, string> = {
  'island-cruise': '♨',
  'shekou-ferry': '△',
  'private-tour': '✦',
  bridge: '▥',
  chimelong: '⌘',
  macau: '$',
  hotel: '⌂',
  car: '⌁'
};

const heroSlides = [
  {
    imageUrl: '/static/phase2/home-hero-watercolor.png',
    title: '珠海湾区漫游',
    subtitle: '一屏看海岛、桥影与城市烟火'
  },
  {
    imageUrl: '/static/phase2/home-hero-moon-web.jpg',
    title: '日月贝夜色',
    subtitle: '城市灯火与海湾夜航同框'
  },
  {
    imageUrl: '/static/phase2/home-hero-sunset-web.jpg',
    title: '情侣路日落',
    subtitle: '椰影海风，慢逛海滨黄昏'
  }
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

const activeHeroTitle = computed(() => heroSlides[activeHeroSlideIndex.value]?.title ?? '珠海湾区漫游');
const activeHeroSubtitle = computed(() => heroSlides[activeHeroSlideIndex.value]?.subtitle ?? '一屏看海岛、桥影与城市烟火');

function openRoute(route?: string) {
  if (!route) return;
  navigationAdapter.navigateTo(route);
}

function switchTab(route: string) {
  navigationAdapter.switchTab(route);
}

function showHeroSlide(index: number) {
  activeHeroSlideIndex.value = (index + heroSlides.length) % heroSlides.length;
}

function restartHeroCarousel() {
  if (heroSlideTimer) {
    clearInterval(heroSlideTimer);
  }
  heroSlideTimer = setInterval(() => {
    showHeroSlide(activeHeroSlideIndex.value + 1);
  }, 4200);
}

function selectHeroSlide(index: number) {
  showHeroSlide(index);
  restartHeroCarousel();
}

onMounted(() => {
  restartHeroCarousel();
});

onBeforeUnmount(() => {
  if (heroSlideTimer) {
    clearInterval(heroSlideTimer);
  }
});
</script>

<template>
  <view class="home-page">
    <view class="home-hero" aria-label="九洲港至蛇口船票">
      <view class="hero-slider" aria-hidden="true">
        <image
          v-for="(slide, index) in heroSlides"
          :key="slide.title"
          class="home-hero__image"
          :class="{ 'home-hero__image--active': activeHeroSlideIndex === index }"
          :src="slide.imageUrl"
          mode="aspectFill"
        />
      </view>
      <view class="home-hero__shade" />
      <view class="home-hero__top">
        <view class="brand-lockup" aria-label="珠海湾游">
          <view class="brand-logo">
            <view class="brand-logo__sun" />
            <view class="brand-logo__pearl" />
            <view class="brand-logo__wave brand-logo__wave--one" />
            <view class="brand-logo__wave brand-logo__wave--two" />
          </view>
          <view class="brand-copy">
            <text class="brand-copy__name">珠海湾游</text>
            <text class="brand-copy__tagline">海上岭南 · 珠澳假期</text>
          </view>
        </view>
        <view class="hero-search">
          <text class="hero-search__icon">⌕</text>
          <text class="hero-search__placeholder">搜索服务名称</text>
        </view>
      </view>
      <view class="hero-greeting" aria-label="您好，珠海">
        <text class="hero-greeting__hello">您好～</text>
        <text class="hero-greeting__city">珠海！</text>
      </view>
      <view class="hero-carousel-copy" aria-label="首页轮播推荐">
        <text class="hero-carousel-copy__title">{{ activeHeroTitle }}</text>
        <text class="hero-carousel-copy__subtitle">{{ activeHeroSubtitle }}</text>
        <view class="hero-dots" aria-label="切换首页轮播">
          <button
            v-for="(slide, index) in heroSlides"
            :key="slide.title"
            class="hero-dots__item"
            :class="{ 'hero-dots__item--active': activeHeroSlideIndex === index }"
            :aria-label="`第 ${index + 1} 张轮播图`"
            @click="selectHeroSlide(index)"
          />
        </view>
      </view>
    </view>

    <view v-if="errorMessage" class="section">
      <EmptyState title="当前场景提示" :description="errorMessage" />
    </view>

    <view class="home-content">
      <view class="quick-sheet">
        <view class="section-title">
          <text class="section-title__text">热门目的地</text>
          <text class="section-title__link">全部分类</text>
        </view>
        <view class="destination-grid" aria-label="珠海热门目的地">
          <button
            v-for="destination in (viewModel?.destinations ?? []).slice(0, 6)"
            :key="destination.id"
            class="destination-item"
            @click="openRoute(destination.route)"
          >
            <text class="destination-item__icon">{{ destinationIconMap[destination.id] ?? destination.title.slice(0, 1) }}</text>
            <text class="destination-item__title">{{ destination.title }}</text>
            <text class="destination-item__subtitle">{{ destination.subtitle }}</text>
          </button>
        </view>
      </view>

      <view class="promo" @click="openRoute('/pages/island-cruise/index')">
        <image class="promo__image" src="/static/phase2/macau-cruise-night-banner-web.jpg" mode="aspectFill" />
        <view class="promo__shade" />
        <view class="promo__copy">
          <text class="promo__title">九洲港至蛇口船票</text>
          <text class="promo__body">九洲港出发，先选成人/儿童/长者人数，再查询接口真实可售船次</text>
        </view>
        <text class="promo__badge">立即购票</text>
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

      <view class="products">
        <view class="section-title">
          <text class="section-title__text">精选商品</text>
          <text class="section-title__link">上新提醒</text>
        </view>

        <view v-if="isLoading" class="product-list">
          <EmptyState title="加载中" description="正在读取首页 Mock 场景。" />
        </view>
        <view v-else-if="filteredProducts.length === 0" class="product-list">
          <EmptyState title="暂无商品" description="当前筛选或 Mock 场景没有可展示商品。" />
        </view>
        <view v-else class="product-list">
          <button
            v-for="product in filteredProducts"
            :key="product.id"
            class="product-card"
            @click="openRoute(product.route)"
          >
            <view class="product-card__media">
              <image
                v-if="product.imageUrl"
                class="product-card__image"
                :src="product.imageUrl"
                mode="aspectFill"
              />
              <text class="product-card__tag">{{ product.tag }}</text>
            </view>
            <view class="product-card__info">
              <text class="product-card__title">{{ product.title }}</text>
              <text class="product-card__subtitle">{{ product.subtitle }}</text>
              <view class="meta-row">
                <text>即时确认</text>
                <text>{{ product.category === 'ship' ? '可核销' : '库存同步' }}</text>
              </view>
              <view class="price-row">
                <text class="price-label">{{ product.priceLabel }}</text>
                <text class="reserve-button">{{ product.actionText }}</text>
              </view>
            </view>
          </button>
        </view>
      </view>
    </view>

    <view class="visual-bottom-nav" aria-label="底部导航">
      <view class="visual-bottom-nav__item visual-bottom-nav__item--active" @click="switchTab('/pages/home/index')">
        <text class="visual-bottom-nav__icon">⌂</text>
        <text>首页</text>
      </view>
      <view class="visual-bottom-nav__item" @click="switchTab('/pages/service/index')">
        <text class="visual-bottom-nav__icon">♬</text>
        <text>客服</text>
      </view>
      <view class="visual-bottom-nav__item" @click="switchTab('/pages/profile/index')">
        <text class="visual-bottom-nav__icon">♙</text>
        <text>我的</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.home-page {
  position: relative;
  width: min(100vw, 430px);
  min-height: 100vh;
  margin: clamp(0px, calc((100vw - 430px) * 100), 28px) auto;
  overflow-x: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px);
  background: #edf5f7;
  box-shadow:
    0 0 0 1px rgba(9, 50, 66, 0.04),
    0 24px 80px rgba(8, 45, 64, 0.18);
  color: #0d2b3a;
}

.home-hero {
  position: relative;
  height: 430px;
  overflow: hidden;
  border-radius: clamp(0px, calc((100vw - 430px) * 100), 26px) clamp(0px, calc((100vw - 430px) * 100), 26px) 0 0;
  color: #ffffff;
  background: #126c91;
}

.hero-slider,
.home-hero__image,
.home-hero__shade {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.home-hero__image {
  opacity: 0;
  transition: opacity 580ms ease;
}

.home-hero__image--active {
  opacity: 1;
}

.home-hero__shade {
  background:
    linear-gradient(180deg, rgba(5, 23, 32, 0.18) 0%, rgba(5, 23, 32, 0.02) 42%, rgba(5, 23, 32, 0.22) 100%),
    linear-gradient(90deg, rgba(6, 40, 53, 0.2), transparent 55%);
  pointer-events: none;
}

.home-hero__top {
  position: absolute;
  z-index: 2;
  top: env(safe-area-inset-top);
  left: 0;
  right: 0;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 12px;
  align-items: start;
  min-height: 76px;
  padding: 18px 14px 0;
  box-sizing: border-box;
}

.brand-lockup {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr);
  align-items: center;
  gap: 8px;
  min-width: 0;
  text-shadow: 0 2px 12px rgba(8, 37, 48, 0.28);
}

.brand-logo {
  position: relative;
  width: 42px;
  height: 42px;
  overflow: hidden;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.94);
  filter: drop-shadow(0 8px 14px rgba(7, 42, 55, 0.18));
}

.brand-logo__sun,
.brand-logo__pearl,
.brand-logo__wave {
  position: absolute;
}

.brand-logo__sun {
  right: 8px;
  top: 7px;
  width: 15px;
  height: 15px;
  border-radius: 50%;
  background: #f2b544;
}

.brand-logo__pearl {
  left: 10px;
  top: 10px;
  width: 13px;
  height: 13px;
  border: 2px solid #f2b544;
  border-radius: 50%;
  background: #fff7d7;
}

.brand-logo__wave {
  left: 7px;
  right: 7px;
  height: 10px;
  border-radius: 999px 999px 0 0;
}

.brand-logo__wave--one {
  bottom: 10px;
  background: #53c4ce;
}

.brand-logo__wave--two {
  bottom: 5px;
  background: #0f9cbd;
}

.brand-copy__name,
.brand-copy__tagline {
  display: block;
}

.brand-copy__name {
  font-size: 18px;
  font-weight: 900;
  line-height: 1.05;
}

.brand-copy__tagline {
  margin-top: 4px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 11px;
}

.hero-search {
  justify-self: end;
  display: flex;
  align-items: center;
  gap: 8px;
  width: min(218px, 58vw);
  height: 36px;
  box-sizing: border-box;
  padding: 0 13px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.68);
  box-shadow: 0 10px 24px rgba(11, 45, 65, 0.12);
  backdrop-filter: blur(16px);
}

.hero-search__icon,
.hero-search__placeholder {
  color: #71808a;
  font-size: 13px;
}

.hero-search__icon {
  font-size: 17px;
}

.hero-greeting {
  position: absolute;
  z-index: 2;
  top: 106px;
  left: 18px;
  display: flex;
  flex-direction: column;
  max-width: min(62%, 240px);
  text-shadow: 0 4px 16px rgba(6, 31, 40, 0.34);
  pointer-events: none;
}

.hero-greeting::before {
  content: "";
  position: absolute;
  left: 4px;
  right: 8px;
  bottom: 7px;
  height: 16px;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.36), rgba(255, 255, 255, 0.06));
  filter: blur(8px);
  transform: rotate(-6deg);
}

.hero-greeting__hello,
.hero-greeting__city {
  position: relative;
  display: block;
  font-family: "STKaiti", "Kaiti SC", "KaiTi", serif;
  font-weight: 900;
  line-height: 1;
  letter-spacing: 2px;
}

.hero-greeting__hello {
  font-size: 30px;
  transform: rotate(-7deg);
}

.hero-greeting__city {
  margin-top: -2px;
  font-size: 52px;
  transform: rotate(-5deg);
}

.hero-carousel-copy {
  position: absolute;
  z-index: 2;
  left: 18px;
  right: 18px;
  bottom: 50px;
  display: grid;
  gap: 6px;
  text-shadow: 0 4px 18px rgba(5, 27, 38, 0.38);
}

.hero-carousel-copy__title {
  font-size: 25px;
  font-weight: 900;
  line-height: 1.05;
}

.hero-carousel-copy__subtitle {
  color: rgba(255, 255, 255, 0.9);
  font-size: 13px;
}

.hero-dots {
  display: flex;
  gap: 5px;
  margin-top: 4px;
}

.hero-dots__item {
  width: 5px;
  height: 5px;
  margin: 0;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.56);
}

.hero-dots__item::after {
  display: none;
}

.hero-dots__item--active {
  width: 18px;
  background: #ffffff;
}

.section {
  margin: 12px 14px;
}

.home-content {
  position: relative;
  z-index: 3;
  margin-top: -24px;
  padding-bottom: calc(86px + env(safe-area-inset-bottom));
}

.quick-sheet {
  padding: 14px;
  border-radius: 18px 18px 8px 8px;
  background: #ffffff;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.section-title__text {
  position: relative;
  padding-left: 10px;
  font-size: 18px;
  font-weight: 900;
  line-height: 1.2;
}

.section-title__text::before {
  content: "";
  position: absolute;
  left: 0;
  top: 2px;
  bottom: 2px;
  width: 4px;
  border-radius: 999px;
  background: #ff8f2f;
}

.section-title__link {
  color: #6d7b84;
  font-size: 13px;
}

.destination-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-auto-rows: minmax(88px, auto);
  gap: 8px;
}

.destination-item,
.product-card,
.category-tab {
  margin: 0;
  border: 0;
  text-align: inherit;
}

.destination-item::after,
.product-card::after,
.category-tab::after {
  display: none;
}

.destination-item {
  display: grid;
  justify-items: center;
  align-content: center;
  gap: 7px;
  min-height: 88px;
  padding: 10px 6px;
  border-radius: 8px;
  color: #875a32;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.82), rgba(255, 246, 236, 0.42)),
    #fff0df;
  box-shadow:
    inset 0 0 0 1px rgba(177, 104, 45, 0.08),
    0 8px 16px rgba(125, 83, 45, 0.06);
}

.destination-item:nth-child(2n) {
  color: #4f7c51;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.82), rgba(236, 248, 234, 0.5)),
    #edf8eb;
}

.destination-item:nth-child(3n) {
  color: #3f7281;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.82), rgba(231, 246, 250, 0.5)),
    #eaf7fa;
}

.destination-item__icon {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.66);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.58);
  font-size: 19px;
  font-weight: 900;
}

.destination-item__title,
.destination-item__subtitle {
  display: block;
  min-width: 0;
  text-align: center;
  line-height: 1.18;
}

.destination-item__title {
  color: #073142;
  font-size: 14px;
  font-weight: 900;
}

.destination-item__subtitle {
  color: rgba(84, 74, 67, 0.68);
  font-size: 11px;
}

.promo {
  position: relative;
  min-height: 116px;
  margin: 12px 14px 0;
  overflow: hidden;
  border-radius: 8px;
  color: #ffffff;
  box-shadow: 0 12px 22px rgba(10, 49, 70, 0.14);
}

.promo__image,
.promo__shade {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.promo__shade {
  background: linear-gradient(90deg, rgba(6, 31, 47, 0.82), rgba(6, 31, 47, 0.22));
}

.promo__copy {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 8px;
  max-width: 245px;
  padding: 18px 88px 16px 16px;
}

.promo__title {
  font-size: 23px;
  font-weight: 900;
  line-height: 1.08;
}

.promo__body {
  font-size: 12px;
  line-height: 1.45;
}

.promo__badge {
  position: absolute;
  z-index: 1;
  right: 14px;
  top: 18px;
  padding: 7px 10px;
  border-radius: 999px;
  color: #5a3b00;
  background: #ffd979;
  font-size: 12px;
  font-weight: 900;
}

.category-tabs {
  position: sticky;
  top: 0;
  z-index: 10;
  white-space: nowrap;
  margin-top: 10px;
  padding: 12px 14px;
  background: rgba(237, 245, 247, 0.94);
  backdrop-filter: blur(8px);
}

.category-tab {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 34px;
  margin-right: 8px;
  padding: 0 13px;
  border: 1px solid #dbe7eb;
  border-radius: 8px;
  color: #6d7b84;
  background: #ffffff;
  font-size: 14px;
}

.category-tab--active {
  color: #ffffff;
  border-color: #087393;
  background: #087393;
  font-weight: 900;
}

.products {
  padding: 0 14px;
}

.product-list {
  display: grid;
  gap: 10px;
}

.product-card {
  display: grid;
  grid-template-columns: 112px minmax(0, 1fr);
  align-items: start;
  gap: 12px;
  min-height: 118px;
  padding: 10px;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(10, 49, 70, 0.08);
}

.product-card__media {
  position: relative;
  width: 100%;
  aspect-ratio: 8 / 7;
  overflow: hidden;
  border-radius: 8px;
  background: #d7e9ec;
}

.product-card__image {
  width: 100%;
  height: 100%;
}

.product-card__tag {
  position: absolute;
  left: 7px;
  top: 7px;
  padding: 3px 6px;
  border-radius: 6px;
  color: #ffffff;
  background: rgba(8, 115, 147, 0.9);
  font-size: 11px;
  font-weight: 900;
}

.product-card__info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 100%;
}

.product-card__title {
  font-size: 16px;
  font-weight: 900;
  line-height: 1.28;
}

.product-card__subtitle {
  margin-top: 6px;
  color: #6d7b84;
  font-size: 12px;
  line-height: 1.45;
}

.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin: 8px 0;
}

.meta-row text {
  padding: 3px 6px;
  border-radius: 6px;
  color: #087393;
  background: #e9f7fa;
  font-size: 11px;
}

.price-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 8px;
  margin-top: auto;
}

.price-label {
  color: #e65337;
  font-size: 18px;
  font-weight: 900;
}

.reserve-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 58px;
  height: 30px;
  padding: 0 10px;
  border-radius: 8px;
  color: #ffffff;
  background: #f35d68;
  font-size: 13px;
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
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 -10px 26px rgba(14, 48, 64, 0.08);
  border-radius: 0 0 clamp(0px, calc((100vw - 430px) * 100), 26px) clamp(0px, calc((100vw - 430px) * 100), 26px);
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

@media (max-width: 360px) {
  .home-hero {
    height: 400px;
  }

  .brand-lockup {
    grid-template-columns: 34px minmax(0, 72px);
  }

  .brand-logo {
    width: 34px;
    height: 34px;
  }

  .brand-copy__name {
    font-size: 16px;
  }

  .brand-copy__tagline {
    display: none;
  }

  .hero-search {
    width: min(180px, 50vw);
  }

  .hero-greeting {
    top: 96px;
    left: 16px;
    max-width: 58%;
  }

  .hero-greeting__hello {
    font-size: 26px;
  }

  .hero-greeting__city {
    font-size: 44px;
  }

  .destination-grid {
    grid-auto-rows: minmax(84px, auto);
    gap: 7px;
  }

  .destination-item {
    gap: 6px;
    min-height: 84px;
    padding: 9px 5px;
  }

  .destination-item__icon {
    width: 30px;
    height: 30px;
  }

  .destination-item__title {
    font-size: 14px;
  }

  .destination-item__subtitle {
    font-size: 11px;
  }

  .product-card {
    grid-template-columns: 98px minmax(0, 1fr);
  }
}

@media (min-width: 431px) {
  .home-page {
    min-height: 860px;
    margin-top: 28px;
    margin-bottom: 28px;
    border-radius: 26px;
    box-shadow:
      0 0 0 1px rgba(9, 50, 66, 0.04),
      0 24px 80px rgba(8, 45, 64, 0.18);
  }

  .home-hero {
    border-radius: 26px 26px 0 0;
  }

  .visual-bottom-nav {
    border-radius: 0 0 26px 26px;
  }
}
</style>
