import type {
  BackendDriverCommissionRecord,
  BackendDriverProfileRecord,
  BackendDriverWalletRecord,
  BackendDriverWithdrawalRecord,
  DriverCommissionSummary,
  DriverProfileSummary,
  DriverStatus,
  DriverWalletSummary,
  DriverWithdrawalSummary
} from '@/types/driver';
import { maskSensitive } from '@/utils/sensitive';

function money(value = 0): string {
  return `¥${value.toFixed(2)}`;
}

function statusLabel(status?: string): string {
  const labels: Record<string, string> = {
    pending: '待处理',
    settled: '已结算',
    approved: '已审核',
    transferred: '已打款',
    rejected: '已驳回'
  };
  return labels[status ?? ''] ?? status ?? '待处理';
}

function driverStatus(value?: string): DriverStatus {
  if (value === 'active') return 'active';
  if (value === 'rejected') return 'rejected';
  if (value === 'unauthorized') return 'unauthorized';
  return 'pending_review';
}

function dateLabel(value?: string): string {
  if (!value) return '--';
  return value.replace('T', ' ').slice(0, 16);
}

export function mapDriverProfile(record: BackendDriverProfileRecord): DriverProfileSummary {
  const rate = Number(record.commission_rate ?? 0);
  return {
    id: record.driver_id ?? record.id ?? 0,
    driverNo: record.driver_no ?? 'DR-ZH-0000',
    name: record.name ?? '司机',
    maskedPhone: maskSensitive(record.phone ?? '') as string,
    status: driverStatus(record.status),
    vehicleLabel: `${record.vehicle_type ?? '商务车'} · ${record.car_plate_no ?? '待绑定车牌'}`,
    commissionRateLabel: `${(rate * 100).toFixed(1)}%`,
    qrCodeText: record.qr_code ?? record.driver_no ?? 'DR-ZH-0000'
  };
}

export function mapDriverWallet(record: BackendDriverWalletRecord): DriverWalletSummary {
  const available = Number(record.available ?? 0);
  return {
    availableAmount: available,
    availableLabel: money(available),
    pendingLabel: money(Number(record.pending_total ?? 0)),
    settledLabel: money(Number(record.settled_total ?? 0)),
    withdrawnLabel: money(Number(record.withdrawn_total ?? 0))
  };
}

export function mapDriverCommission(record: BackendDriverCommissionRecord): DriverCommissionSummary {
  return {
    id: String(record.id ?? record.order_no ?? ''),
    orderNo: record.order_no ?? '待同步订单',
    amountLabel: money(Number(record.commission_amount ?? 0)),
    statusLabel: statusLabel(record.status),
    createdAtLabel: dateLabel(record.created_at)
  };
}

export function mapDriverWithdrawal(record: BackendDriverWithdrawalRecord): DriverWithdrawalSummary {
  return {
    id: String(record.id ?? record.withdrawal_no ?? ''),
    withdrawalNo: record.withdrawal_no ?? 'WD-待生成',
    amountLabel: money(Number(record.amount ?? 0)),
    statusLabel: statusLabel(record.status),
    accountLabel: maskSensitive(record.account ?? '') as string,
    createdAtLabel: dateLabel(record.created_at)
  };
}
