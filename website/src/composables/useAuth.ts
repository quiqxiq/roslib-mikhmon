import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/auth'
import type { Credentials } from '@/types/auth'

export function useAuth() {
  const store = useAuthStore()
  const { access_token, refresh_token, user, isAuthenticated, role } = storeToRefs(store)

  async function login(credentials: Credentials) {
    const tokens = await authService.login(credentials)
    store.setTokens(tokens)
    const me = await authService.me()
    store.setUser(me)
  }

  async function logout() {
    if (store.refresh_token) {
      try {
        await authService.logout(store.refresh_token)
      } catch {
        // ignore — local reset still proceeds
      }
    }
    store.reset()
  }

  return {
    access_token,
    refresh_token,
    user,
    isAuthenticated,
    role,
    login,
    logout,
  }
}
