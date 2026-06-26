import type {
  BackendDriverCommissionRecord,
  BackendDriverProfileRecord,
  BackendDriverWalletRecord,
  BackendDriverWithdrawalRecord,
  DriverScenarioId,
  DriverViewModel
} from '@/types/driver';
import { mapDriverCommission, mapDriverProfile, mapDriverWallet, mapDriverWithdrawal } from '@/utils/driver-mappers';

const activeProfile: BackendDriverProfileRecord = {
  driver_id: 42,
  driver_no: 'DR-ZH-0420',
  name: '陈师傅',
  phone: '13800138000',
  status: 'active',
  car_plate_no: '粤C·D4208',
  vehicle_type: '七座商务车',
  commission_rate: 0.12,
  qr_code: 'DR-ZH-0420-SEAT-08'
};

const wallet: BackendDriverWalletRecord = {
  available: 1268,
  pending_total: 436,
  settled_total: 1888,
  withdrawn_total: 620
};

const commissions: BackendDriverCommissionRecord[] = [
  { id: 1, order_no: 'ZH20260715001028', commission_amount: 88, status: 'settled', created_at: '2026-07-15T21:20:00+08:00' },
  { id: 2, order_no: 'ZH20260714000982', commission_amount: 64, status: 'pending', created_at: '2026-07-14T18:40:00+08:00' }
];

const withdrawals: BackendDriverWithdrawalRecord[] = [
  { id: 1, withdrawal_no: 'WD202607160001', amount: 500, status: 'approved', account: 'driver-alipay@example.com', created_at: '2026-07-16T10:30:00+08:00' }
];

export function normalizeDriverScenario(scenarioId: string): DriverScenarioId {
  if (
    scenarioId === 'phase4-active' ||
    scenarioId === 'phase4-unauthorized' ||
    scenarioId === 'phase4-login-failed' ||
    scenarioId === 'phase4-pending-review' ||
    scenarioId === 'phase4-empty-wallet' ||
    scenarioId === 'phase4-insufficient-balance' ||
    scenarioId === 'phase4-withdraw-success' ||
    scenarioId === 'phase4-failure'
  ) {
    return scenarioId;
  }
  return 'phase4-active';
}

export function getDriverMock(scenarioId: string): DriverViewModel {
  const normalizedScenarioId = normalizeDriverScenario(scenarioId);
  const empty = normalizedScenarioId === 'phase4-empty-wallet';
  const unauthorized = normalizedScenarioId === 'phase4-unauthorized';
  const pendingReview = normalizedScenarioId === 'phase4-pending-review';
  const loginFailed = normalizedScenarioId === 'phase4-login-failed';
  const failure = normalizedScenarioId === 'phase4-failure';
  const profile = pendingReview
    ? mapDriverProfile({ ...activeProfile, status: 'pending_review', driver_no: 'DR-ZH-待审核', qr_code: '审核后生成推广码' })
    : unauthorized || loginFailed || failure
      ? null
      : mapDriverProfile(activeProfile);

  return {
    dataSource: 'mock',
    scenarioId: normalizedScenarioId,
    profile,
    wallet: profile && !pendingReview ? mapDriverWallet(empty ? {} : wallet) : null,
    commissions: profile && !pendingReview && !empty ? commissions.map(mapDriverCommission) : [],
    withdrawals: profile && !pendingReview && !empty ? withdrawals.map(mapDriverWithdrawal) : [],
    withdrawDraft: {
      amount: normalizedScenarioId === 'phase4-insufficient-balance' ? 2000 : 300,
      account: 'driver-alipay@example.com',
      realName: '陈师傅',
      channel: 'alipay'
    },
    statusMessage: pendingReview
      ? '注册申请已提交，等待后台审核后生成车座二维码。'
      : normalizedScenarioId === 'phase4-withdraw-success'
        ? '提现申请已提交，待管理员审核打款。'
        : undefined,
    errorMessage: failure
      ? '司机服务暂时不可用，请稍后重试。'
      : loginFailed
        ? '手机号或密码错误，请检查后重试。'
        : unauthorized
          ? '请先登录司机端。'
          : undefined
  };
}
