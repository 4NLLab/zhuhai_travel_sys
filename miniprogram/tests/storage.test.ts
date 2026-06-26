import { beforeEach, describe, expect, it, vi } from 'vitest';
import { storageAdapter } from '@/adapters/storage';

const store = new Map<string, unknown>();

beforeEach(() => {
  store.clear();
  vi.stubGlobal('uni', {
    setStorageSync: (key: string, value: unknown) => store.set(key, value),
    getStorageSync: (key: string) => store.get(key) ?? '',
    removeStorageSync: (key: string) => store.delete(key)
  });
});

describe('storageAdapter', () => {
  it('stores minimal session fields and clears logout state', () => {
    storageAdapter.setSession({
      accessToken: 'dev-token',
      profile: {
        id: 'user-1',
        role: 'user',
        displayName: '测试用户'
      }
    });

    expect(storageAdapter.getSession()).toEqual({
      accessToken: 'dev-token',
      profile: {
        id: 'user-1',
        role: 'user',
        displayName: '测试用户'
      }
    });

    storageAdapter.clearSession();
    expect(storageAdapter.getSession()).toBeNull();
  });
});
