import { describe, expect, it } from 'vitest';
import type {
  BackendDriverCommissionRecord,
  BackendDriverProfileRecord,
  BackendDriverWalletRecord,
  BackendDriverWithdrawalRecord
} from '@/types/driver';
import { getDriverMock, normalizeDriverScenario } from '@/mock/driver';
import { mapDriverCommission, mapDriverProfile, mapDriverWallet, mapDriverWithdrawal } from '@/utils/driver-mappers';

describe('phase four driver adapters', () => {
  it('maps driver profile and masks phone', () => {
    const record: BackendDriverProfileRecord = {
      driver_id: 42,
      driver_no: 'DR-ZH-0420',
      name: '陈师傅',
      phone: '13800138000',
      status: 'active',
      car_plate_no: '粤C·D4208',
      vehicle_type: '七座商务车',
      commission_rate: 0.12
    };

    const profile = mapDriverProfile(record);

    expect(profile.driverNo).toBe('DR-ZH-0420');
    expect(profile.maskedPhone).toBe('138****8000');
    expect(profile.status).toBe('active');
    expect(profile.commissionRateLabel).toBe('12.0%');
  });

  it('maps wallet, commission and withdrawal records', () => {
    const wallet: BackendDriverWalletRecord = {
      available: 1268,
      pending_total: 436,
      settled_total: 1888,
      withdrawn_total: 620
    };
    const commission: BackendDriverCommissionRecord = {
      id: 1,
      order_no: 'ZH20260715001028',
      commission_amount: 88,
      status: 'settled',
      created_at: '2026-07-15T21:20:00+08:00'
    };
    const withdrawal: BackendDriverWithdrawalRecord = {
      id: 1,
      withdrawal_no: 'WD202607160001',
      amount: 500,
      status: 'approved',
      account: 'driver-alipay@example.com',
      created_at: '2026-07-16T10:30:00+08:00'
    };

    expect(mapDriverWallet(wallet).availableLabel).toBe('¥1268.00');
    expect(mapDriverCommission(commission)).toMatchObject({
      orderNo: 'ZH20260715001028',
      amountLabel: '¥88.00',
      statusLabel: '已结算'
    });
    expect(mapDriverWithdrawal(withdrawal)).toMatchObject({
      withdrawalNo: 'WD202607160001',
      amountLabel: '¥500.00',
      statusLabel: '已审核'
    });
  });

  it('exposes required phase four mock scenarios', () => {
    const active = getDriverMock('phase4-active');
    const unauthorized = getDriverMock('phase4-unauthorized');
    const pending = getDriverMock('phase4-pending-review');
    const empty = getDriverMock('phase4-empty-wallet');
    const insufficient = getDriverMock('phase4-insufficient-balance');
    const success = getDriverMock('phase4-withdraw-success');

    expect(active.profile?.status).toBe('active');
    expect(unauthorized.profile).toBeNull();
    expect(pending.profile?.status).toBe('pending_review');
    expect(empty.commissions).toEqual([]);
    expect(insufficient.withdrawDraft.amount).toBeGreaterThan(insufficient.wallet?.availableAmount ?? 0);
    expect(success.statusMessage).toContain('提现申请');
  });

  it('falls back to active for unknown phase four ids', () => {
    expect(normalizeDriverScenario('phase3-success')).toBe('phase4-active');
  });
});
