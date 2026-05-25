import type { ExpiryMode } from '@/types/profile-config'

// View-model gabungan untuk satu kartu profile di HotspotProfilesTab.
// Menggabungkan data RouterOS (HotspotProfile) + mikhmon-local
// (ProfileConfig) untuk display & edit.
export interface ProfileViewModel {
  id: string // ID RouterOS (mis. "*1")
  name: string
  speed: string
  validity: string
  price: number
  sell_price: number
  expiry_mode: ExpiryMode
  lock_mac: boolean
  shared_users: number
  address_pool: string
  add_mac_cookie: boolean
  transparent_proxy: boolean
  // Decorative
  color: 'cyan' | 'violet' | 'lime'
  sold: number
  // True kalau sudah ada record di hotspot_profile_configs (mikhmon DB).
  has_config: boolean
}

// Label IDR-friendly untuk ExpiryMode dropdown.
export const EXPIRY_MODE_OPTIONS: Array<{ value: ExpiryMode; label: string; hint: string }> = [
  { value: '0', label: 'Free / None', hint: 'Tidak ada selling record (profile gratis)' },
  { value: 'rem', label: 'Remove on expiry', hint: 'Hapus user saat expired' },
  { value: 'ntf', label: 'Notice + kick', hint: 'Set limit-uptime ke 0 + kick' },
  { value: 'remc', label: 'Remove + record', hint: 'Hapus + catat transaksi (recommended)' },
  { value: 'ntfc', label: 'Notice + record', hint: 'Notice + catat transaksi' },
]

export function expiryModeLabel(m: ExpiryMode): string {
  return EXPIRY_MODE_OPTIONS.find((o) => o.value === m)?.label ?? m
}

export function isPaidMode(m: ExpiryMode): boolean {
  return m === 'remc' || m === 'ntfc'
}
