import type { ApiResult } from '@/types/platform';
import type { DriverViewModel } from '@/types/driver';
import { getRuntimeConfig } from '@/utils/config';
import { getDriverMock } from '@/mock/driver';
import { storageAdapter } from '@/adapters/storage';
import { requestJson } from './request';

function fallbackDriver(errorMessage: string): ApiResult<DriverViewModel> {
  const mock = getDriverMock('phase4-active');
  return {
    data: {
      ...mock,
      dataSource: 'fallback',
      errorMessage
    },
    dataSource: 'fallback'
  };
}

export async function loadDriverViewModel(scenarioOverride?: string): Promise<ApiResult<DriverViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getDriverMock(scenarioOverride ?? config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  if (config.apiMode !== 'local') {
    throw new Error(`暂未开放 ${config.apiMode} 环境。`);
  }

  try {
    let session: ReturnType<typeof storageAdapter.getSession> = null;
    try {
      session = storageAdapter.getSession();
    } catch {
      session = null;
    }
    const headers = session?.profile.role === 'driver' ? { Authorization: `Bearer ${session.accessToken}` } : undefined;
    await requestJson({ path: '/driver/me', headers });
  } catch (error) {
    if (config.allowFallback) {
      return fallbackDriver(error instanceof Error ? error.message : '本地司机接口需要 active driver token，已回落 Mock。');
    }
    throw error;
  }

  if (config.allowFallback) {
    return fallbackDriver('本地司机鉴权已通过，但钱包/佣金聚合仍需真实 active driver seed，当前保留 Mock 展示。');
  }
  throw new Error('本地司机鉴权已通过，但 Phase 5 未完成司机钱包真实数据聚合映射。');
}
