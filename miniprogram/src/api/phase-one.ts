import type { ApiResult } from '@/types/platform';
import { getPhaseOneStatus, type PhaseOneStatus } from '@/mock/scenario';
import { getRuntimeConfig } from '@/utils/config';
import { requestJson } from './request';

export async function loadPhaseOneStatus(): Promise<ApiResult<PhaseOneStatus>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return getPhaseOneStatus(config.mockScenarioId);
  }

  if (config.apiMode === 'local') {
    return requestJson<PhaseOneStatus>({
      path: '/health',
      allowFallback: false
    });
  }

  throw new Error(`暂未开放 ${config.apiMode} 环境。`);
}
