import type { ApiMode, PlatformName, RuntimeConfig } from '@/types/platform';

const DEFAULT_MOCK_SCENARIO = 'phase1-success';

function readEnv(name: string): string | undefined {
  const env = import.meta.env as Record<string, string | undefined>;
  return env[name];
}

function normalizeApiMode(value: string | undefined): ApiMode {
  if (value === 'local' || value === 'test' || value === 'prod') {
    return value;
  }
  return 'mock';
}

export function detectPlatform(): PlatformName {
  // #ifdef H5
  return 'h5';
  // #endif
  // #ifdef MP-WEIXIN
  return 'mp-weixin';
  // #endif
  return 'unknown';
}

export function getRuntimeConfig(): RuntimeConfig {
  const apiMode = normalizeApiMode(readEnv('VITE_API_MODE'));
  return {
    apiMode,
    baseUrl: readEnv('VITE_API_BASE_URL') ?? (apiMode === 'local' ? 'http://127.0.0.1:8080/api/v1' : '/api/v1'),
    mockScenarioId: readEnv('VITE_MOCK_SCENARIO') ?? DEFAULT_MOCK_SCENARIO,
    allowFallback: readEnv('VITE_ALLOW_FALLBACK') === 'true',
    platform: detectPlatform()
  };
}
