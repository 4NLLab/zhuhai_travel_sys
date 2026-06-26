import type { ApiResult } from '@/types/platform';

export interface PhaseOneStatus {
  title: string;
  summary: string;
  enabledAdapters: string[];
}

const scenarios: Record<string, PhaseOneStatus> = {
  'phase1-success': {
    title: '小程序工程已接入 Mock',
    summary: '当前为 Phase 1 占位数据，后续业务切片会按页面补齐场景。',
    enabledAdapters: ['request', 'storage', 'logger', 'modal', 'navigation']
  },
  'phase1-empty': {
    title: '暂无可展示业务',
    summary: '平台护栏已就绪，等待 Phase 2 开始迁移首页、订单和我的。',
    enabledAdapters: []
  },
  'phase1-failure': {
    title: 'Mock 场景模拟异常',
    summary: '用于验证错误态、日志脱敏和 fallback 边界。',
    enabledAdapters: ['logger']
  }
};

export function getPhaseOneStatus(scenarioId: string): ApiResult<PhaseOneStatus> {
  const fallbackScenario = scenarios['phase1-success'];
  if (!fallbackScenario) {
    throw new Error('phase1-success mock scenario is missing');
  }

  return {
    data: scenarios[scenarioId] ?? fallbackScenario,
    dataSource: 'mock'
  };
}
