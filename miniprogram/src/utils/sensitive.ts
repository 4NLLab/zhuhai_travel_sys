const TOKEN_PATTERN = /(Bearer\s+)?[A-Za-z0-9_-]{24,}/g;
const PHONE_PATTERN = /(?<!\d)1[3-9]\d{9}(?!\d)/g;
const ID_CARD_PATTERN = /(?<![A-Za-z0-9])\d{6}(?:19|20)\d{2}(?:0[1-9]|1[0-2])(?:0[1-9]|[12]\d|3[01])\d{3}[\dXx](?![A-Za-z0-9])/g;
const BANK_ACCOUNT_PATTERN = /(?<!\d)\d{12,19}(?!\d)/g;

export function maskSensitive(value: unknown): unknown {
  if (typeof value === 'string') {
    return value
      .replace(TOKEN_PATTERN, '[MASKED_TOKEN]')
      .replace(PHONE_PATTERN, (phone) => `${phone.slice(0, 3)}****${phone.slice(-4)}`)
      .replace(ID_CARD_PATTERN, (id) => `${id.slice(0, 4)}**********${id.slice(-4)}`)
      .replace(BANK_ACCOUNT_PATTERN, (account) => `${account.slice(0, 4)}********${account.slice(-4)}`);
  }

  if (Array.isArray(value)) {
    return value.map((item) => maskSensitive(item));
  }

  if (value && typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>).map(([key, item]) => {
        if (/token|password|secret|credential/i.test(key)) {
          return [key, '[MASKED_SECRET]'];
        }
        return [key, maskSensitive(item)];
      })
    );
  }

  return value;
}

export function redactLogPayload(payload: unknown): unknown {
  return maskSensitive(payload);
}
