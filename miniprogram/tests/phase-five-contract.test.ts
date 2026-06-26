import { beforeEach, describe, expect, it, vi } from 'vitest';

type UniRequestOptions = {
  url: string;
  method?: string;
  header?: Record<string, string>;
  success: (response: { statusCode: number; data: unknown }) => void;
  fail: (error: unknown) => void;
};

const storage = new Map<string, unknown>();

function stubEnv(allowFallback = false): void {
  vi.stubEnv('VITE_API_MODE', 'local');
  vi.stubEnv('VITE_API_BASE_URL', 'http://127.0.0.1:8080/api/v1');
  vi.stubEnv('VITE_ALLOW_FALLBACK', allowFallback ? 'true' : 'false');
}

function stubUni(handler: (options: UniRequestOptions) => void): void {
  vi.stubGlobal('uni', {
    request: handler,
    getStorageSync: (key: string) => storage.get(key) ?? '',
    setStorageSync: (key: string, value: unknown) => storage.set(key, value),
    removeStorageSync: (key: string) => storage.delete(key)
  });
}

beforeEach(() => {
  storage.clear();
  vi.resetModules();
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
});

describe('phase five local contract adapters', () => {
  it('loads public product and category endpoints without page code changes', async () => {
    stubEnv();
    stubUni((options) => {
      if (options.url.endsWith('/products?size=6')) {
        options.success({
          statusCode: 200,
          data: {
            code: 200,
            message: 'ok',
            total: 1,
            data: [
              {
                id: 12,
                title: '澳门环岛游夜景船票',
                subtitle: '湾仔旅游码头出发',
                product_type: 'ship_ticket',
                category: { id: 2, name: '船票', slug: 'ship' },
                skus: [{ sale_price: 88, status: 'active' }]
              }
            ]
          }
        });
        return;
      }
      if (options.url.endsWith('/categories')) {
        options.success({
          statusCode: 200,
          data: {
            code: 200,
            message: 'ok',
            data: [{ id: 2, name: '船票', slug: 'ship' }]
          }
        });
      }
    });

    const { loadMainShellViewModel } = await import('@/api/main-shell');
    const result = await loadMainShellViewModel();

    expect(result.dataSource).toBe('local');
    expect(result.data.products[0]).toMatchObject({
      title: '澳门环岛游夜景船票',
      category: 'ship',
      priceLabel: '¥88 起',
      route: '/pages/island-cruise/index'
    });
    expect(result.data.orders).toEqual([]);
  });

  it('maps island cruise public smart-search into the booking flow view model', async () => {
    stubEnv();
    stubUni((options) => {
      expect(options.url).toContain('/island-cruise/smart-search');
      options.success({
        statusCode: 200,
        data: {
          code: 200,
          message: 'ok',
          data: {
            recommended: {
              up_port_id: 2312,
              down_port_id: 2332,
              up_port_name: '湾仔码头',
              down_port_name: '澳门环岛游A',
              label: '澳门环岛游',
              date: '2026-07-15',
              count: 1,
              min_price: 88,
              first_time: '19:30',
              voyages: [
                {
                  voyageId: 1001,
                  voyageName: '澳门环岛游夜游',
                  departureDate: '2026-07-15',
                  departureTime: '19:30',
                  cabinPriceList: [
                    {
                      cabinClassId: 1,
                      cabinClassName: '普通舱',
                      cabinTypeCode: 'A',
                      count: 20,
                      fareTypeList: [{ fareTypeId: 10, fareTypeName: '成人票', price: 88 }]
                    }
                  ]
                }
              ]
            },
            routes: []
          }
        }
      });
    });

    const { loadIslandCruiseFlow } = await import('@/api/island-cruise');
    const result = await loadIslandCruiseFlow();
    const voyage = result.data.voyages[0];
    const fare = voyage?.fares[0];

    expect(result.dataSource).toBe('local');
    expect(result.data.recommended?.minPriceLabel).toBe('¥88 起');
    expect(fare?.stock).toBe(20);
    expect(result.data.lockedOrder).toBeNull();
  });

  it('falls back for driver auth boundary when no active driver token is available', async () => {
    stubEnv(true);
    stubUni((options) => {
      expect(options.url).toContain('/driver/me');
      expect(options.header?.Authorization).toBeUndefined();
      options.success({
        statusCode: 401,
        data: { code: 401, message: 'missing token' }
      });
    });

    const { loadDriverViewModel } = await import('@/api/driver');
    const result = await loadDriverViewModel();

    expect(result.dataSource).toBe('fallback');
    expect(result.data.errorMessage).toContain('HTTP 401');
  });
});
