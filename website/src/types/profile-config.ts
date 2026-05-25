// Sesuai dto.ProfileConfigResponse backend (snake_case 1:1).
//
// expiry_mode enum:
//   - "0"   = none / disable (free profile, tidak ada selling record)
//   - "rem" = hapus user saat expired
//   - "ntf" = set limit-uptime ke 0 + kick (notice)
//   - "remc"= hapus + catat transaksi
//   - "ntfc"= notice + catat transaksi
export type ExpiryMode = '0' | 'rem' | 'ntf' | 'remc' | 'ntfc'

export interface ProfileConfig {
  id: number
  device_id: number
  profile_name: string
  expiry_mode: ExpiryMode
  validity: string
  price: number
  sell_price: number
  lock_mac: boolean
  created_at: string
  updated_at: string
}

// Body untuk PUT /devices/{id}/hotspot/profile-configs/{profile_name}.
export interface ProfileConfigUpsertRequest {
  expiry_mode: ExpiryMode
  validity: string
  price: number
  sell_price: number
  lock_mac: boolean
}

// Sesuai dto.ProfileConfigSyncResponse backend.
export interface ProfileConfigSyncResponse {
  synced: string[]
  created: string[]
  orphan: string[]
  injected: string[]
  inject_failed: string[]
}
