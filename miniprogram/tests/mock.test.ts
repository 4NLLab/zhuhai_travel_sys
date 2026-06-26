import { describe, expect, it } from 'vitest';
import { getPhaseOneStatus } from '@/mock/scenario';

describe('phase one mock scenario', () => {
  it('returns declared mock data source', () => {
    const result = getPhaseOneStatus('phase1-success');

    expect(result.dataSource).toBe('mock');
    expect(result.data.enabledAdapters).toContain('request');
  });

  it('falls back to success scenario for unknown ids', () => {
    const result = getPhaseOneStatus('missing-scenario');

    expect(result.data.title).toContain('Mock');
  });
});
