import { redactLogPayload } from '@/utils/sensitive';

type LogLevel = 'debug' | 'info' | 'warn' | 'error';

export interface LogRecord {
  level: LogLevel;
  event: string;
  payload?: unknown;
  createdAt: string;
}

function write(level: LogLevel, event: string, payload?: unknown): LogRecord {
  const record: LogRecord = {
    level,
    event,
    payload: redactLogPayload(payload),
    createdAt: new Date().toISOString()
  };

  if (level === 'error') {
    console.error('[zt]', record);
  } else if (level === 'warn') {
    console.warn('[zt]', record);
  } else {
    console.info('[zt]', record);
  }

  return record;
}

export const logger = {
  debug: (event: string, payload?: unknown) => write('debug', event, payload),
  info: (event: string, payload?: unknown) => write('info', event, payload),
  warn: (event: string, payload?: unknown) => write('warn', event, payload),
  error: (event: string, payload?: unknown) => write('error', event, payload)
};
