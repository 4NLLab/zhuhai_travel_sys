import type { ApiResult } from '@/types/platform';
import type { DriverViewModel } from '@/types/driver';
import { getRuntimeConfig } from '@/utils/config';
import { getDriverMock } from '@/mock/driver';

export async function loadDriverViewModel(scenarioOverride?: string): Promise<ApiResult<DriverViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getDriverMock(scenarioOverride ?? config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  throw new Error('Phase 4 暂不接真实司机接口，请切换到 mock 模式。');
}
