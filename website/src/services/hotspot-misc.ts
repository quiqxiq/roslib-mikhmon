import { http } from '@/plugins/axios'
import type { ApiEnvelope } from '@/types/api'
import type { VoucherGenerateRequest, VoucherGenerateResponse } from '@/types/hotspot'

export const hotspotMiscService = {
  async generateVouchers(
    deviceId: string,
    req: VoucherGenerateRequest,
  ): Promise<VoucherGenerateResponse> {
    const { data } = await http.post<ApiEnvelope<VoucherGenerateResponse>>(
      `/devices/${deviceId}/hotspot/vouchers/generate`,
      req,
    )
    return data.data
  },
}
