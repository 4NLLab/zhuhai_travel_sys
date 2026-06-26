import type { ApiResult } from '@/types/platform';
import { logger } from '@/adapters/logger';
import { getRuntimeConfig } from '@/utils/config';

export interface RequestOptions {
  path: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: string | Record<string, unknown> | ArrayBuffer;
  headers?: Record<string, string>;
  allowFallback?: boolean;
}

function buildUrl(baseUrl: string, path: string): string {
  const normalizedBase = baseUrl.replace(/\/$/, '');
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${normalizedBase}${normalizedPath}`;
}

export async function requestJson<T>(options: RequestOptions): Promise<ApiResult<T>> {
  const config = getRuntimeConfig();
  if (config.apiMode === 'mock') {
    throw new Error('requestJson 不直接处理 mock，请通过 domain adapter 选择 fixture。');
  }

  return new Promise<ApiResult<T>>((resolve, reject) => {
    uni.request({
      url: buildUrl(config.baseUrl, options.path),
      method: options.method ?? 'GET',
      data: options.data,
      header: options.headers,
      success(response) {
        if (response.statusCode >= 200 && response.statusCode < 300) {
          resolve({ data: response.data as T, dataSource: 'local' });
          return;
        }
        logger.warn('request_non_2xx', {
          path: options.path,
          statusCode: response.statusCode,
          body: response.data
        });
        reject(new Error(`HTTP ${response.statusCode}: ${options.path}`));
      },
      fail(error) {
        logger.error('request_failed', {
          path: options.path,
          error
        });
        reject(error);
      }
    });
  });
}
