import type { ApiResult } from '@/types/platform';
import type { BackendIslandVoyageRecord, IslandCruiseViewModel, IslandLockedOrder, IslandTicketResult } from '@/types/island-cruise';
import { getRuntimeConfig } from '@/utils/config';
import { getIslandCruiseMock } from '@/mock/island-cruise';
import { mapIslandVoyage, minFareLabel } from '@/utils/island-cruise-mappers';
import { requestJson } from './request';

interface BackendEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

interface BackendSmartSearchRoute {
  up_port_id: number;
  down_port_id: number;
  up_port_name: string;
  down_port_name: string;
  label: string;
  date: string;
  count: number;
  min_price: number;
  first_time: string;
  voyages: BackendIslandVoyageRecord[];
}

interface BackendSmartSearchPayload {
  recommended: BackendSmartSearchRoute | null;
  routes: BackendSmartSearchRoute[];
}

const DEFAULT_DATE = '2026-07-15';

function fallbackIslandCruise(scenarioId: string, errorMessage: string): ApiResult<IslandCruiseViewModel> {
  const mock = getIslandCruiseMock(scenarioId);
  return {
    data: {
      ...mock,
      dataSource: 'fallback',
      errorMessage
    },
    dataSource: 'fallback'
  };
}

function mapSmartSearch(payload: BackendSmartSearchPayload): IslandCruiseViewModel {
  const route = payload.recommended ?? payload.routes.find((item) => item.count > 0) ?? payload.routes[0];
  const fallbackDate = route?.date || DEFAULT_DATE;
  const voyages = (route?.voyages ?? []).map((voyage) => mapIslandVoyage(voyage, fallbackDate));
  const firstVoyage = voyages[0] ?? null;
  const minPriceLabel = firstVoyage ? minFareLabel(firstVoyage) : route?.min_price ? `¥${route.min_price} 起` : '暂无报价';

  return {
    dataSource: 'local',
    scenarioId: voyages.length > 0 ? 'phase3-success' : 'phase3-no-voyage',
    hero: {
      title: '澳门环岛游',
      summary: '本地 Go 后端 /island-cruise/smart-search 查询结果',
      imageUrl: '/static/phase2/macau-cruise-night-banner-web.jpg'
    },
    ports: {
      up: {
        id: route?.up_port_id ?? 2312,
        name: route?.up_port_name ?? '湾仔码头'
      },
      down: {
        id: route?.down_port_id ?? 2332,
        name: route?.down_port_name ?? '澳门环岛游A'
      }
    },
    certTypes: [
      { id: 1, name: '身份证' },
      { id: 2, name: '护照' }
    ],
    recommended: route
      ? {
          date: fallbackDate,
          label: route.label,
          firstTime: route.first_time || firstVoyage?.departureTime || '待供应商返回',
          count: route.count,
          minPriceLabel
        }
      : null,
    voyages,
    selectedVoyageId: firstVoyage?.id ?? null,
    draft: {
      contactName: '',
      contactMobile: '',
      passengers: [],
      quantity: 1,
      totalAmount: firstVoyage?.fares[0]?.price ?? 0
    },
    lockedOrder: null,
    ticket: null,
    serviceNotes: ['local 模式只读取公开查询接口；锁票、出票、退票和改签需供应商与交易联调。']
  };
}

export async function loadIslandCruiseFlow(scenarioOverride?: string): Promise<ApiResult<IslandCruiseViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getIslandCruiseMock(scenarioOverride ?? config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  if (config.apiMode !== 'local') {
    throw new Error(`暂未开放 ${config.apiMode} 环境。`);
  }

  try {
    const result = await requestJson<BackendEnvelope<BackendSmartSearchPayload>>({
      path: `/island-cruise/smart-search?start_date=${DEFAULT_DATE}&days=7&people_num=1`
    });
    return {
      data: mapSmartSearch(result.data.data),
      dataSource: 'local'
    };
  } catch (error) {
    if (config.allowFallback) {
      return fallbackIslandCruise(
        scenarioOverride ?? config.mockScenarioId,
        error instanceof Error ? error.message : '本地环岛游公开查询接口不可用，已回落 Mock。'
      );
    }
    throw error;
  }
}

export async function lockIslandCruiseOrder(viewModel: IslandCruiseViewModel): Promise<ApiResult<IslandLockedOrder | null>> {
  return {
    data: viewModel.lockedOrder,
    dataSource: viewModel.dataSource
  };
}

export async function saleIslandCruiseOrder(viewModel: IslandCruiseViewModel): Promise<ApiResult<IslandTicketResult | null>> {
  return {
    data: viewModel.ticket,
    dataSource: viewModel.dataSource
  };
}
