import { API_BASE_URL } from '@/utils/env'
import { useAuthStore } from '@/stores/auth'

export type StreamPath =
  | 'hotspot/active'
  | 'hotspot/users'
  | 'hotspot/inactive'
  | 'ppp/active'
  | 'ppp/secrets'
  | 'ppp/inactive'
  | 'log'
  | 'system/resource'
  | 'system/routerboard'
  | 'network/interfaces/stats'
  | 'network/queues/stats'
  | `network/queues/stats/${string}`
  | 'network/queues/parents/stats'
  | 'ping'

export function buildStreamUrl(deviceId: string, path: StreamPath | string): string {
  const auth = useAuthStore()
  const url = new URL(`${API_BASE_URL}/devices/${deviceId}/stream/${path}`, window.location.origin)
  if (auth.access_token) {
    url.searchParams.set('access_token', auth.access_token)
  }
  return url.toString()
}

export function buildInterfaceTrafficUrl(deviceId: string, ifaceName: string): string {
  return buildStreamUrl(deviceId, `network/interfaces/${encodeURIComponent(ifaceName)}/traffic`)
}
