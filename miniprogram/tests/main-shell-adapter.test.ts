import { describe, expect, it } from 'vitest';
import type { BackendOrderRecord, BackendTicketRecord } from '@/types/main-shell';
import { getMainShellMock, normalizePhaseTwoScenario } from '@/mock/main-shell';
import { mapBackendOrder, mapBackendTicket } from '@/utils/main-shell-mappers';

describe('phase two main shell adapters', () => {
  it('maps backend order list records into order summaries', () => {
    const record: BackendOrderRecord = {
      id: 'order-001',
      product_name: '澳门环岛游夜景船票',
      order_no: 'ZH20260715001028',
      status: 'pending_use',
      source_type: 'island_cruise',
      travel_date: '2026-07-15',
      quantity_text: '成人票 x1',
      paid_amount_text: '¥88',
      usage_hint: '凭电子票码到湾仔旅游码头核销',
      ticket_id: 'ticket-001'
    };

    expect(mapBackendOrder(record)).toEqual({
      id: 'order-001',
      title: '澳门环岛游夜景船票',
      orderNo: 'ZH20260715001028',
      status: 'pending_use',
      source: 'island_cruise',
      travelDateLabel: '07-15',
      quantityLabel: '成人票 x1',
      amountLabel: '¥88',
      hint: '凭电子票码到湾仔旅游码头核销',
      ticketId: 'ticket-001'
    });
  });

  it('maps backend ticket detail records and masks ticket code', () => {
    const record: BackendTicketRecord = {
      id: 'ticket-001',
      status: 'available',
      title: '待使用券码',
      product_name: '澳门环岛游夜景船票',
      valid_date: '2026-07-15',
      valid_time: '19:30-21:00',
      quantity_text: '成人票 x1',
      ticket_code: '0105049143806',
      verify_location: '湾仔旅游码头检票口',
      notice: ['核销前请勿将券码提供给无关人员。'],
      order_no: 'ZH20260715001028',
      paid_amount_text: '¥88',
      source_type: 'island_cruise'
    };

    const ticket = mapBackendTicket(record);

    expect(ticket.productTitle).toBe('澳门环岛游夜景船票');
    expect(ticket.validDateLabel).toBe('07-15');
    expect(ticket.maskedCode).toBe('0105********3806');
    expect(ticket.notice).toHaveLength(1);
  });

  it('exposes required phase two mock scenarios', () => {
    const success = getMainShellMock('phase2-success');
    const empty = getMainShellMock('phase2-empty');
    const unauthorized = getMainShellMock('phase2-unauthorized');
    const failure = getMainShellMock('phase2-failure');

    expect(success.orders.map((order) => order.status)).toEqual(
      expect.arrayContaining(['pending_use', 'pending_pay', 'reserved', 'refunded'])
    );
    expect(success.tickets.map((ticket) => ticket.status)).toEqual(expect.arrayContaining(['available', 'unavailable']));
    expect(empty.orders).toEqual([]);
    expect(unauthorized.user).toBeNull();
    expect(failure.errorMessage).toContain('订单服务');
  });

  it('falls back to the success scenario for unknown phase two ids', () => {
    expect(normalizePhaseTwoScenario('unknown')).toBe('phase2-success');
  });
});
