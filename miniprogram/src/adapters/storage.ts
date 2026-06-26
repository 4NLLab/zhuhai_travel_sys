import type { StoredSession } from '@/types/platform';

const SESSION_KEY = 'zt_session_minimal';

function getUniStorage() {
  return uni;
}

export const storageAdapter = {
  setSession(session: StoredSession): void {
    getUniStorage().setStorageSync(SESSION_KEY, {
      accessToken: session.accessToken,
      profile: {
        id: session.profile.id,
        role: session.profile.role,
        displayName: session.profile.displayName
      }
    });
  },

  getSession(): StoredSession | null {
    const value = getUniStorage().getStorageSync(SESSION_KEY) as StoredSession | '';
    return value || null;
  },

  clearSession(): void {
    getUniStorage().removeStorageSync(SESSION_KEY);
  }
};
