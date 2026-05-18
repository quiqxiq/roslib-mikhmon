package dto

// HistoryRow adalah representasi satu baris hasil query historis dari
// InfluxDB3. Skema kolom bervariasi per endpoint (resource, interfaces,
// hotspot, ppp, queues) sehingga di-model sebagai dinamis tapi tetap
// punya nama tipe agar OpenAPI/dokumentasi bisa merujuk konsisten.
//
// Kontrak: keys & values selalu serializable ke JSON. Values bisa
// scalar (string, number, bool) atau time.Time (untuk kolom DATE_BIN).
type HistoryRow map[string]any
