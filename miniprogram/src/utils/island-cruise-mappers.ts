import type {
  BackendIslandLockRecord,
  BackendIslandSaleRecord,
  BackendIslandVoyageRecord,
  IslandFare,
  IslandLockedOrder,
  IslandTicketResult,
  IslandVoyage
} from '@/types/island-cruise';
import { maskSensitive } from '@/utils/sensitive';

function money(value: number): string {
  return `¥${Number.isInteger(value) ? value : value.toFixed(2)}`;
}

function numberValue(value: number | string | undefined, fallback = 0): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function shortDate(value: string): string {
  const [, month, day] = value.split('-');
  return month && day ? `${month}-${day}` : value;
}

function maskTicketCode(value: string): string {
  const trimmed = value.trim();
  if (trimmed.length <= 8) return maskSensitive(trimmed) as string;
  return `${trimmed.slice(0, 4)}********${trimmed.slice(-4)}`;
}

export function mapIslandVoyage(record: BackendIslandVoyageRecord, fallbackDate: string): IslandVoyage {
  const fares: IslandFare[] = [];
  (record.cabinPriceList ?? []).forEach((cabin) => {
    (cabin.fareTypeList ?? []).forEach((fare) => {
      const price = numberValue(fare.price ?? fare.settlePrice ?? fare.originalPrice);
      if (price <= 0) return;
      fares.push({
        id: String(fare.fareTypeId ?? `${cabin.cabinClassId}-${fare.fareTypeName}`),
        name: fare.fareTypeName ?? '票种',
        cabinClassId: String(cabin.cabinClassId ?? ''),
        cabinClassName: cabin.cabinClassName ?? '普通舱',
        cabinTypeCode: cabin.cabinTypeCode ?? cabin.cabinType ?? 'S',
        price,
        originalPrice: numberValue(fare.originalPrice, price),
        stock: numberValue(cabin.count)
      });
    });
  });

  return {
    id: record.voyageId,
    name: record.voyageName ?? record.lineName ?? '澳门环岛游',
    voyageNo: record.voyageNo ?? '',
    shipName: record.shipName ?? '环岛游船',
    departureDate: record.departureDate ?? fallbackDate,
    departureTime: record.departureTime ?? record.upTime ?? '19:30',
    upPort: {
      id: record.upPortId ?? 2312,
      name: record.upPortName ?? '湾仔码头'
    },
    downPort: {
      id: record.downPortId ?? 2332,
      name: record.downPortName ?? '澳门环岛游A'
    },
    stateLabel: record.voyageState ?? '售票',
    fares: fares.sort((left, right) => left.price - right.price)
  };
}

export function minFareLabel(voyage: IslandVoyage): string {
  const min = voyage.fares[0]?.price ?? 0;
  return min > 0 ? `${money(min)} 起` : '查价';
}

export function totalStock(voyage: IslandVoyage): number {
  return voyage.fares.reduce((max, fare) => Math.max(max, fare.stock), 0);
}

export function mapIslandLock(record: BackendIslandLockRecord, amount: number, expired = false): IslandLockedOrder {
  const expireAt = record.lock_expire_at ? new Date(record.lock_expire_at) : null;
  const expireTimeLabel =
    expireAt && !Number.isNaN(expireAt.getTime())
      ? `${String(expireAt.getHours()).padStart(2, '0')}:${String(expireAt.getMinutes()).padStart(2, '0')} 前支付`
      : '10 分钟内';

  return {
    localOrderNo: record.local_order_no ?? 'HD202607150001',
    supplierOrderNo: record.orderNo ?? 'SUP202607150001',
    ticketNo: record.ticketNo ?? 'TICKET20260715001',
    expireTimeLabel,
    status: expired ? 'expired' : 'pending_payment',
    amount
  };
}

export function mapIslandSale(
  record: BackendIslandSaleRecord,
  voyage: IslandVoyage,
  passengerLabel: string,
  saleFailed = false
): IslandTicketResult {
  const codeContent = record.codeContent ?? record.ticketNo ?? 'HD202607150001';
  return {
    localOrderNo: record.local_order_no ?? 'HD202607150001',
    ticketNo: record.ticketNo ?? 'TICKET20260715001',
    codeContent,
    maskedCode: maskTicketCode(codeContent),
    paidAtLabel: record.paid_at ? record.paid_at.replace('T', ' ').slice(0, 16) : '2026-07-15 18:05',
    voyageLabel: `${shortDate(voyage.departureDate)} ${voyage.departureTime} · ${voyage.name}`,
    passengerLabel,
    status: saleFailed ? 'sale_failed' : 'ticketed',
    failureMessage: saleFailed ? '供应商出票暂时失败，订单保留为待处理状态。' : undefined
  };
}
