import { http } from '@/plugins/axios'
import type { ApiEnvelope } from '@/types/api'
import type { AuthUser, Credentials, TokenPair } from '@/types/auth'

// Service tipis: pass-through 1:1 ke DTO snake_case backend.
// Tidak ada manual mapping camel↔snake — types sudah snake_case.
export const authService = {
  async login(credentials: Credentials): Promise<TokenPair> {
    const { data } = await http.post<ApiEnvelope<TokenPair>>('/auth/login', credentials)
    return data.data
  },
  async refresh(refreshToken: string): Promise<TokenPair> {
    const { data } = await http.post<ApiEnvelope<TokenPair>>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return data.data
  },
  async logout(refreshToken: string): Promise<void> {
    await http.post('/auth/logout', { refresh_token: refreshToken })
  },
  async me(): Promise<AuthUser> {
    const { data } = await http.get<ApiEnvelope<AuthUser>>('/auth/me')
    return data.data
  },
}
