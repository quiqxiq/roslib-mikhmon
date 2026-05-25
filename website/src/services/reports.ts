import { http } from '@/plugins/axios'
import type { ApiEnvelope } from '@/types/api'
import type { Transaction, ReportSummary, ReportTodayResponse } from '@/types/report'

// Endpoint backend (sesuai api/handler/report.go):
//   GET  /reports/selling?month=jan2025
//   GET  /reports/selling/today
//   GET  /reports/selling/summary?month=jan2025&include_transactions=true
//   GET  /reports/selling.csv?month=jan2025  atau  ?date=jan/05/2025
const base = (deviceId: string) => `/devices/${deviceId}/reports`

export const reportsService = {
  // List transaksi. Filter `month` format mikhmon ("jan2025"). Kalau kosong,
  // backend kembalikan semua transaksi device.
  async listSelling(deviceId: string, month?: string): Promise<Transaction[]> {
    const { data } = await http.get<ApiEnvelope<Transaction[]>>(`${base(deviceId)}/selling`, {
      params: month ? { month } : undefined,
    })
    return data.data
  },

  async today(deviceId: string): Promise<ReportTodayResponse> {
    const { data } = await http.get<ApiEnvelope<ReportTodayResponse>>(
      `${base(deviceId)}/selling/today`,
    )
    return data.data
  },

  async summary(
    deviceId: string,
    month?: string,
    includeTransactions = false,
  ): Promise<ReportSummary> {
    const params: Record<string, string> = {}
    if (month) params.month = month
    if (includeTransactions) params.include_transactions = 'true'
    const { data } = await http.get<ApiEnvelope<ReportSummary>>(
      `${base(deviceId)}/selling/summary`,
      { params: Object.keys(params).length ? params : undefined },
    )
    return data.data
  },

  // CSV export. `month` (mis. "jan2025") atau `date` (mis. "jan/05/2025").
  async exportCsv(deviceId: string, opts: { month?: string; date?: string } = {}): Promise<Blob> {
    const { data } = await http.get(`${base(deviceId)}/selling.csv`, {
      params: opts,
      responseType: 'blob',
    })
    return data
  },
}
