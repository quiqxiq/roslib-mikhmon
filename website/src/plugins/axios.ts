import axios, { AxiosError, type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'
import { API_BASE_URL } from '@/utils/env'
import { useAuthStore } from '@/stores/auth'
import type { ApiEnvelope } from '@/types/api'
import type { TokenPair } from '@/types/auth'

export const http: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30_000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const auth = useAuthStore()
  if (auth.access_token) {
    config.headers.set('Authorization', `Bearer ${auth.access_token}`)
  }
  return config
})

let refreshInFlight: Promise<string | null> | null = null

async function refreshAccessToken(): Promise<string | null> {
  const auth = useAuthStore()
  if (!auth.refresh_token) return null
  try {
    const res = await axios.post<ApiEnvelope<TokenPair>>(`${API_BASE_URL}/auth/refresh`, {
      refresh_token: auth.refresh_token,
    })
    const tokens = res.data?.data
    if (tokens?.access_token) {
      auth.setTokens(tokens)
      return tokens.access_token
    }
  } catch {
    auth.reset()
  }
  return null
}

http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined
    if (error.response?.status === 401 && original && !original._retry) {
      original._retry = true
      refreshInFlight ??= refreshAccessToken().finally(() => {
        refreshInFlight = null
      })
      const newToken = await refreshInFlight
      if (newToken) {
        original.headers.set('Authorization', `Bearer ${newToken}`)
        return http(original)
      }
    }
    return Promise.reject(error)
  },
)
