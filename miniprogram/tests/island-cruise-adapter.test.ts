import { describe, expect, it } from 'vitest';
import type { BackendIslandLockRecord, BackendIslandSaleRecord, BackendIslandVoyageRecord } from '@/types/island-cruise';
import { getIslandCruiseMock, normalizeIslandCruiseScenario } from '@/mock/island-cruise';
import { mapIslandLock, mapIslandSale, mapIslandVoyage, minFareLabel, totalStock } from '@/utils/island-cruise-mappers';

describe('phase three island cruise adapters', () => {
  it('maps backend voyage records into selectable voyages', () => {
    const record: BackendIslandVoyageRecord = {
      voyageId: 1001,
      voyageName: '澳门环岛游夜景航班',
      voyageNo: 'MACAU-NIGHT-1930',
      shipName: '珠澳湾游 1 号',
      departureDate: '2026-07-15',
      departureTime: '19:30',
      cabinPriceList: [
        {
          cabinClassId: 'S',
          cabinClassName: '普通舱',
          count: 12,
          fareTypeList: [{ fareTypeId: 'adult', fareTypeName: '成人票', price: '88', originalPrice: 128 }]
        }
      ]
    };

    const voyage = mapIslandVoyage(record, '2026-07-15');

    expect(voyage.id).toBe(1001);
    expect(voyage.fares[0]).toMatchObject({ id: 'adult', name: '成人票', price: 88, stock: 12 });
    expect(minFareLabel(voyage)).toBe('¥88 起');
    expect(totalStock(voyage)).toBe(12);
  });

  it('maps lock and sale records into payment and ticket view models', () => {
    const voyage = mapIslandVoyage(
      {
        voyageId: 1001,
        voyageName: '澳门环岛游',
        departureDate: '2026-07-15',
        departureTime: '19:30',
        cabinPriceList: []
      },
      '2026-07-15'
    );
    const lockRecord: BackendIslandLockRecord = {
      local_order_no: 'HD202607150001',
      orderNo: 'SUP202607150001',
      ticketNo: 'TICKET20260715001',
      lock_expire_at: '2026-07-15T19:15:00+08:00'
    };
    const saleRecord: BackendIslandSaleRecord = {
      local_order_no: 'HD202607150001',
      ticketNo: 'TICKET20260715001',
      codeContent: 'HD0105049143806',
      paid_at: '2026-07-15T18:05:00+08:00'
    };

    const lockedOrder = mapIslandLock(lockRecord, 88);
    const ticket = mapIslandSale(saleRecord, voyage, '陈小珠 成人票 x1');

    expect(lockedOrder.status).toBe('pending_payment');
    expect(lockedOrder.amount).toBe(88);
    expect(ticket.status).toBe('ticketed');
    expect(ticket.maskedCode).toBe('HD01********3806');
    expect(ticket.voyageLabel).toContain('07-15 19:30');
  });

  it('exposes required phase three mock scenarios', () => {
    const success = getIslandCruiseMock('phase3-success');
    const noVoyage = getIslandCruiseMock('phase3-no-voyage');
    const insufficientStock = getIslandCruiseMock('phase3-insufficient-stock');
    const expired = getIslandCruiseMock('phase3-lock-expired');
    const invalid = getIslandCruiseMock('phase3-passenger-invalid');
    const saleFailed = getIslandCruiseMock('phase3-sale-failed');

    expect(success.voyages.length).toBeGreaterThan(0);
    expect(noVoyage.voyages).toEqual([]);
    expect(insufficientStock.voyages[0]?.fares[0]?.stock).toBe(0);
    expect(expired.lockedOrder?.status).toBe('expired');
    expect(invalid.draft.passengers).toEqual([]);
    expect(saleFailed.ticket?.status).toBe('sale_failed');
  });

  it('falls back to success for unknown phase three ids', () => {
    expect(normalizeIslandCruiseScenario('phase2-success')).toBe('phase3-success');
  });
});
