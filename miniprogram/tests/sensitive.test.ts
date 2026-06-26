import { describe, expect, it } from 'vitest';
import { maskSensitive, redactLogPayload } from '@/utils/sensitive';

describe('sensitive masking', () => {
  it('masks phone, id card, bank account and token-like values', () => {
    const masked = maskSensitive({
      phone: '13800138000',
      idCard: '440402199001012345',
      bankAccount: '6222021202001234567',
      accessToken: 'Bearer abcdefghijklmnopqrstuvwxyz123456'
    });

    expect(masked).toEqual({
      phone: '138****8000',
      idCard: '4404**********2345',
      bankAccount: '6222********4567',
      accessToken: '[MASKED_SECRET]'
    });
  });

  it('redacts nested log payloads', () => {
    const redacted = redactLogPayload({
      traveler: {
        mobile: '13900139000',
        credential: '440402199202023456'
      }
    });

    expect(JSON.stringify(redacted)).not.toContain('13900139000');
    expect(JSON.stringify(redacted)).not.toContain('440402199202023456');
  });
});
