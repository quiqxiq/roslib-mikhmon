import { http } from '@/plugins/axios'
import type { ApiEnvelope } from '@/types/api'
import type { Device, DeviceInput, DeviceWriteResponse } from '@/types/device'

// Service tipis: types FE sudah snake_case 1:1 dengan backend DTO.
// Tidak ada toPayload adapter — kirim DeviceInput apa adanya.
export const devicesService = {
  async list(): Promise<Device[]> {
    const { data } = await http.get<ApiEnvelope<Device[]>>('/devices')
    return data.data
  },

  async get(id: string): Promise<Device> {
    const { data } = await http.get<ApiEnvelope<Device>>(`/devices/${id}`)
    return data.data
  },

  async create(input: DeviceInput): Promise<DeviceWriteResponse> {
    const { data } = await http.post<ApiEnvelope<DeviceWriteResponse | Device>>('/devices', input)
    const raw = data.data
    if (raw && (raw as DeviceWriteResponse).device) {
      return raw as DeviceWriteResponse
    }
    return { device: raw as Device }
  },

  async update(id: string, input: Partial<DeviceInput>): Promise<DeviceWriteResponse> {
    // Gunakan PUT sesuai registrasi handler backend Go (bukan PATCH).
    const { data } = await http.put<ApiEnvelope<DeviceWriteResponse | Device>>(
      `/devices/${id}`,
      input,
    )
    const raw = data.data
    if (raw && (raw as DeviceWriteResponse).device) {
      return raw as DeviceWriteResponse
    }
    return { device: raw as Device }
  },

  async remove(id: string): Promise<void> {
    await http.delete(`/devices/${id}`)
  },
}
