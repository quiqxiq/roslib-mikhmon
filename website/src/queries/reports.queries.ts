import { type MaybeRefOrGetter, toValue, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { reportsService } from '@/services/reports'
import { queryKeys } from '@/queries/query-keys'

// List transaksi penjualan. `month` format mikhmon ("jan2025").
export function useSellingQuery(
  deviceId: MaybeRefOrGetter<string | null>,
  month?: MaybeRefOrGetter<string | undefined>,
) {
  return useQuery({
    queryKey: computed(() =>
      queryKeys.reports.selling(String(toValue(deviceId)), toValue(month)),
    ),
    queryFn: () => reportsService.listSelling(String(toValue(deviceId)), toValue(month)),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useSellingTodayQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: computed(() => queryKeys.reports.today(String(toValue(deviceId)))),
    queryFn: () => reportsService.today(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useSellingSummaryQuery(
  deviceId: MaybeRefOrGetter<string | null>,
  month?: MaybeRefOrGetter<string | undefined>,
) {
  return useQuery({
    queryKey: computed(() =>
      queryKeys.reports.summary(String(toValue(deviceId)), toValue(month)),
    ),
    queryFn: () => reportsService.summary(String(toValue(deviceId)), toValue(month)),
    enabled: () => Boolean(toValue(deviceId)),
  })
}
