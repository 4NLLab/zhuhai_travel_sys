import type { ApiResult } from '@/types/platform';
import type { MainShellViewModel, TicketSummary } from '@/types/main-shell';
import { getRuntimeConfig } from '@/utils/config';
import { getMainShellMock } from '@/mock/main-shell';

export async function loadMainShellViewModel(): Promise<ApiResult<MainShellViewModel>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    return {
      data: getMainShellMock(config.mockScenarioId),
      dataSource: 'mock'
    };
  }

  throw new Error('Phase 2 暂不接真实用户订单接口，请切换到 mock 模式。');
}

export async function loadTicketDetail(ticketId = 'ticket-island-1'): Promise<ApiResult<TicketSummary | null>> {
  const result = await loadMainShellViewModel();
  return {
    data: result.data.tickets.find((ticket) => ticket.id === ticketId) ?? null,
    dataSource: result.dataSource
  };
}
