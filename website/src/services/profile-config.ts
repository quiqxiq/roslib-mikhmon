import { http } from '@/plugins/axios'
import type { ApiEnvelope } from '@/types/api'
import type {
  ProfileConfig,
  ProfileConfigUpsertRequest,
  ProfileConfigSyncResponse,
} from '@/types/profile-config'

// Base URL Phase 3 — pindah ke /hotspot/profile-configs (sebelumnya /profile-configs).
const base = (deviceId: string) => `/devices/${deviceId}/hotspot/profile-configs`

export const profileConfigService = {
  async list(deviceId: string): Promise<ProfileConfig[]> {
    const { data } = await http.get<ApiEnvelope<ProfileConfig[]>>(base(deviceId))
    return data.data
  },

  async get(deviceId: string, profileName: string): Promise<ProfileConfig> {
    const { data } = await http.get<ApiEnvelope<ProfileConfig>>(
      `${base(deviceId)}/${encodeURIComponent(profileName)}`,
    )
    return data.data
  },

  // Upsert by profile_name (PUT). Backend mengembalikan 200 (updated) atau 201
  // (created) — keduanya membawa ProfileConfig final di body.
  async upsert(
    deviceId: string,
    profileName: string,
    payload: ProfileConfigUpsertRequest,
  ): Promise<ProfileConfig> {
    const { data } = await http.put<ApiEnvelope<ProfileConfig>>(
      `${base(deviceId)}/${encodeURIComponent(profileName)}`,
      payload,
    )
    return data.data
  },

  async remove(deviceId: string, profileName: string): Promise<void> {
    await http.delete(`${base(deviceId)}/${encodeURIComponent(profileName)}`)
  },

  // Tarik profile router → DB + auto-inject on-login script. Butuh router online.
  async sync(deviceId: string): Promise<ProfileConfigSyncResponse> {
    const { data } = await http.post<ApiEnvelope<ProfileConfigSyncResponse>>(
      `${base(deviceId)}/sync`,
    )
    return data.data
  },
}
