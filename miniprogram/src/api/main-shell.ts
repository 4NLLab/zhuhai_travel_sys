import type { ApiResult } from '@/types/platform';
import type { DestinationEntry, MainShellViewModel, ProductCategory, ProductEntry, TicketSummary } from '@/types/main-shell';
import { getRuntimeConfig } from '@/utils/config';
import { getMainShellMock } from '@/mock/main-shell';
import { requestJson } from './request';

interface BackendPage<T> {
  code: number;
  message: string;
  data: T[];
  total: number;
}

interface BackendEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

interface BackendProductRecord {
  id: number | string;
  title: string;
  subtitle?: string | null;
  product_type?: string;
  cover_url?: string | null;
  category?: BackendCategoryRecord;
  skus?: Array<{ sale_price?: number | string; status?: string }>;
  images?: Array<{ image_url?: string; sort_order?: number }>;
}

interface BackendCategoryRecord {
  id: number | string;
  name: string;
  slug?: string;
}

function categoryFromProductType(productType?: string): ProductCategory {
  if (productType?.includes('hotel')) return 'hotel';
  if (productType?.includes('car')) return 'car';
  if (productType?.includes('tour')) return 'tour';
  if (productType?.includes('ship') || productType?.includes('cruise')) return 'ship';
  if (productType?.includes('play') || productType?.includes('ticket')) return 'play';
  return 'service';
}

function priceLabel(record: BackendProductRecord): string {
  const prices = (record.skus ?? [])
    .filter((sku) => !sku.status || sku.status === 'active')
    .map((sku) => Number(sku.sale_price))
    .filter((value) => Number.isFinite(value) && value > 0)
    .sort((left, right) => left - right);
  return prices[0] ? `¥${prices[0]} 起` : '查价';
}

function productRoute(record: BackendProductRecord): string | undefined {
  const text = `${record.title} ${record.product_type ?? ''}`.toLowerCase();
  if (text.includes('环岛') || text.includes('cruise') || text.includes('ship')) {
    return '/pages/island-cruise/index';
  }
  return undefined;
}

function mapProduct(record: BackendProductRecord): ProductEntry {
  const image = record.cover_url ?? record.images?.sort((left, right) => (left.sort_order ?? 0) - (right.sort_order ?? 0))[0]?.image_url;
  return {
    id: String(record.id),
    title: record.title,
    subtitle: record.subtitle ?? record.category?.name ?? '珠海本地旅游产品',
    category: categoryFromProductType(record.product_type),
    priceLabel: priceLabel(record),
    tag: record.category?.name ?? '本地接口',
    imageUrl: image ?? undefined,
    actionText: '查看',
    route: productRoute(record),
    sourceNote: '来源：本地 Go 后端 /products'
  };
}

function mapCategory(record: BackendCategoryRecord): DestinationEntry {
  const category = categoryFromProductType(record.slug);
  return {
    id: String(record.id),
    title: record.name,
    subtitle: record.slug ?? '本地分类',
    category,
    route: record.slug?.includes('cruise') || record.name.includes('环岛') ? '/pages/island-cruise/index' : undefined
  };
}

function fallbackMainShell(errorMessage: string): ApiResult<MainShellViewModel> {
  const mock = getMainShellMock('phase2-success');
  return {
    data: {
      ...mock,
      dataSource: 'fallback',
      errorMessage
    },
    dataSource: 'fallback'
  };
}

export async function loadMainShellViewModel(): Promise<ApiResult<MainShellViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getMainShellMock(config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  if (config.apiMode !== 'local') {
    throw new Error(`暂未开放 ${config.apiMode} 环境。`);
  }

  try {
    const [productsResult, categoriesResult] = await Promise.all([
      requestJson<BackendPage<BackendProductRecord>>({ path: '/products?size=6' }),
      requestJson<BackendEnvelope<BackendCategoryRecord[]>>({ path: '/categories' })
    ]);
    const products = productsResult.data.data.map(mapProduct);
    const destinations = categoriesResult.data.data.map(mapCategory);

    return {
      data: {
        dataSource: 'local',
        scenarioId: 'phase2-success',
        user: null,
        destinations,
        products,
        orders: [],
        tickets: [],
        errorMessage: '本地后端公开产品/分类已读取；用户订单与票券需要 user token 和 seed，当前保留为空态。'
      },
      dataSource: 'local'
    };
  } catch (error) {
    if (config.allowFallback) {
      return fallbackMainShell(error instanceof Error ? error.message : '本地产品/分类接口不可用，已回落 Mock。');
    }
    throw error;
  }
}

export async function loadTicketDetail(ticketId = 'ticket-island-1'): Promise<ApiResult<TicketSummary | null>> {
  const result = await loadMainShellViewModel();
  return {
    data: result.data.tickets.find((ticket) => ticket.id === ticketId) ?? null,
    dataSource: result.dataSource
  };
}
