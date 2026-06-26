import type { DataSource } from './platform';

export type DriverScenarioId =
  | 'phase4-active'
  | 'phase4-unauthorized'
  | 'phase4-login-failed'
  | 'phase4-pending-review'
  | 'phase4-empty-wallet'
  | 'phase4-insufficient-balance'
  | 'phase4-withdraw-success'
  | 'phase4-failure';

export type DriverStatus = 'active' | 'pending_review' | 'rejected' | 'unauthorized';

export interface DriverProfileSummary {
  id: number;
  driverNo: string;
  name: string;
  maskedPhone: string;
  status: DriverStatus;
  vehicleLabel: string;
  commissionRateLabel: string;
  qrCodeText: string;
}

export interface DriverWalletSummary {
  availableLabel: string;
  pendingLabel: string;
  settledLabel: string;
  withdrawnLabel: string;
  availableAmount: number;
}

export interface DriverCommissionSummary {
  id: string;
  orderNo: string;
  amountLabel: string;
  statusLabel: string;
  createdAtLabel: string;
}

export interface DriverWithdrawalSummary {
  id: string;
  withdrawalNo: string;
  amountLabel: string;
  statusLabel: string;
  accountLabel: string;
  createdAtLabel: string;
}

export interface DriverWithdrawDraft {
  amount: number;
  account: string;
  realName: string;
  channel: 'alipay' | 'wechat';
}

export interface DriverViewModel {
  dataSource: DataSource;
  scenarioId: DriverScenarioId;
  profile: DriverProfileSummary | null;
  wallet: DriverWalletSummary | null;
  commissions: DriverCommissionSummary[];
  withdrawals: DriverWithdrawalSummary[];
  withdrawDraft: DriverWithdrawDraft;
  statusMessage?: string;
  errorMessage?: string;
}

export interface BackendDriverProfileRecord {
  driver_id?: number;
  id?: number;
  driver_no?: string;
  name?: string;
  phone?: string;
  status?: string;
  car_plate_no?: string;
  vehicle_type?: string;
  commission_rate?: number | string;
  qr_code?: string;
}

export interface BackendDriverWalletRecord {
  available?: number;
  pending_total?: number;
  settled_total?: number;
  withdrawn_total?: number;
}

export interface BackendDriverCommissionRecord {
  id?: number | string;
  order_no?: string;
  commission_amount?: number;
  status?: string;
  created_at?: string;
}

export interface BackendDriverWithdrawalRecord {
  id?: number | string;
  withdrawal_no?: string;
  amount?: number;
  status?: string;
  account?: string;
  created_at?: string;
}
