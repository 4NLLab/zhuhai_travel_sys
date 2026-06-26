import type { ApiResult } from '@/types/platform';
import type { IslandCruiseViewModel, IslandLockedOrder, IslandTicketResult } from '@/types/island-cruise';
import { getRuntimeConfig } from '@/utils/config';
import { getIslandCruiseMock } from '@/mock/island-cruise';

export async function loadIslandCruiseFlow(scenarioOverride?: string): Promise<ApiResult<IslandCruiseViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getIslandCruiseMock(scenarioOverride ?? config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  throw new Error('Phase 3 暂不接真实环岛游接口，请切换到 mock 模式。');
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
