/**
 * Convert array of objects to CSV string.
 * Handles strings, numbers, booleans. Escapes quotes/commas/newlines per RFC 4180.
 */
export function toCsv<T extends Record<string, unknown>>(rows: T[], columns?: (keyof T)[]): string {
  if (rows.length === 0) return ''
  const cols = (columns ?? (Object.keys(rows[0]) as (keyof T)[])) as (keyof T)[]
  const escape = (v: unknown): string => {
    if (v == null) return ''
    const s = typeof v === 'string' ? v : String(v)
    if (/[",\n\r]/.test(s)) return `"${s.replace(/"/g, '""')}"`
    return s
  }
  const header = cols.map((c) => escape(String(c))).join(',')
  const body = rows.map((r) => cols.map((c) => escape(r[c])).join(',')).join('\n')
  return `${header}\n${body}`
}

/**
 * Trigger browser download of CSV content.
 */
export function downloadCsv<T extends Record<string, unknown>>(
  rows: T[],
  filename: string,
  columns?: (keyof T)[],
): void {
  const csv = toCsv(rows, columns)
  const blob = new Blob([`﻿${csv}`], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}

export function todayStamp(): string {
  const d = new Date()
  return `${d.getFullYear()}${String(d.getMonth() + 1).padStart(2, '0')}${String(d.getDate()).padStart(2, '0')}`
}
