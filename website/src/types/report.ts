// Sesuai dto.TransactionResponse backend (snake_case 1:1).
//
// Sumber data: tabel `transactions` (di-tulis oleh service/expiry atau
// webhook /hook/hotspot/login). Field `cashier`/`voucher` lama
// di-deprecate — backend tidak punya konsep itu.
export interface Transaction {
  id: number
  sale_date: string  // "jan/02/2006" (lowercase month, mikhmon format)
  sale_time: string  // "15:04:05"
  sale_month: string // "jan2025"
  username: string
  price: number
  sell_price: number
  ip?: string
  mac?: string
  validity?: string
  profile?: string
  comment?: string
  created_at: string // ISO 8601, sortable
}

// Sesuai dto.ReportSummary backend.
export interface ReportSummary {
  count: number
  total_price: number
  total_sell_price: number
  profit: number
  by_profile?: Record<string, number>
  by_sell_price?: Record<string, number>
  transactions?: Transaction[]
}

// Sesuai dto.ReportTodayResponse backend.
export interface ReportTodayResponse {
  date: string // "jan/02/2006"
  count: number
  transactions: Transaction[]
}
