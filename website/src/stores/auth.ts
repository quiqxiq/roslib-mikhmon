import { defineStore } from 'pinia'
import type { AuthUser, TokenPair } from '@/types/auth'

// Pinia store auth — internal state pakai snake_case agar match 1:1 dengan
// TokenPair backend. Persisted lewat pinia-plugin-persistedstate.
interface AuthState {
  access_token: string | null
  refresh_token: string | null
  user: AuthUser | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    access_token: null,
    refresh_token: null,
    user: null,
  }),
  getters: {
    isAuthenticated: (s) => Boolean(s.access_token),
    role: (s) => s.user?.role ?? null,
    // Alias camelCase untuk backward-compat dengan caller existing.
    accessToken: (s) => s.access_token,
    refreshToken: (s) => s.refresh_token,
  },
  actions: {
    setTokens(tokens: TokenPair) {
      this.access_token = tokens.access_token
      this.refresh_token = tokens.refresh_token
    },
    setUser(user: AuthUser | null) {
      this.user = user
    },
    reset() {
      this.access_token = null
      this.refresh_token = null
      this.user = null
    },
  },
  persist: {
    pick: ['access_token', 'refresh_token', 'user'],
  },
})
