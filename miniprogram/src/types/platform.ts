export type ApiMode = 'mock' | 'local' | 'test' | 'prod';
export type DataSource = 'mock' | 'local' | 'fallback';
export type PlatformName = 'h5' | 'mp-weixin' | 'unknown';

export interface RuntimeConfig {
  apiMode: ApiMode;
  baseUrl: string;
  mockScenarioId: string;
  allowFallback: boolean;
  platform: PlatformName;
}

export interface ApiResult<T> {
  data: T;
  dataSource: DataSource;
}

export interface SessionProfile {
  id: string;
  role: 'user' | 'driver' | 'admin';
  displayName: string;
}

export interface StoredSession {
  accessToken: string;
  profile: SessionProfile;
}
