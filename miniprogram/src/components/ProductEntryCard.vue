<script setup lang="ts">
import type { ProductEntry } from '@/types/main-shell';

defineProps<{
  product: ProductEntry;
}>();
</script>

<template>
  <view class="card product-card">
    <image v-if="product.imageUrl" class="product-card__image" :src="product.imageUrl" mode="aspectFill" />
    <view class="product-card__body">
      <view class="product-card__head">
        <text class="product-card__title">{{ product.title }}</text>
        <text class="product-card__tag">{{ product.tag }}</text>
      </view>
      <text class="product-card__subtitle">{{ product.subtitle }}</text>
      <view class="product-card__foot">
        <text class="product-card__price">{{ product.priceLabel }}</text>
        <text class="product-card__action">{{ product.actionText }}</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.product-card {
  overflow: hidden;
}

.product-card__image {
  display: block;
  width: 100%;
  height: 128px;
  background: var(--zt-color-surface-muted);
}

.product-card__body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px;
}

.product-card__head,
.product-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.product-card__title {
  min-width: 0;
  font-size: 17px;
  font-weight: 800;
}

.product-card__tag,
.product-card__action {
  flex: 0 0 auto;
  font-size: 12px;
  font-weight: 800;
}

.product-card__tag {
  padding: 4px 8px;
  border-radius: 8px;
  color: var(--zt-color-primary);
  background: var(--zt-color-primary-weak);
}

.product-card__subtitle {
  font-size: 14px;
  line-height: 22px;
  color: var(--zt-color-text-muted);
}

.product-card__price {
  color: var(--zt-color-accent);
  font-size: 18px;
  font-weight: 900;
}

.product-card__action {
  color: var(--zt-color-primary);
}
</style>
