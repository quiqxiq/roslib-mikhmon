<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import Segmented from '@/components/ui/Segmented.vue'
import Card from '@/components/ui/Card.vue'
import TrendChart from '@/components/ui/charts/TrendChart.vue'
import DataTable from '@/components/ui/DataTable.vue'
import OverviewKpiCard from '@/components/overview/OverviewKpiCard.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import {
  useSellingQuery,
  useSellingTodayQuery,
  useSellingSummaryQuery,
} from '@/queries/reports.queries'
import { fmtRpShort, fmtRp } from '@/utils/fmt'
import type { Transaction } from '@/types/report'
import { useToast } from '@/composables/useToast'
import { reportsService } from '@/services/reports'
import { todayStamp } from '@/utils/export'

const toast = useToast()
const { activeDeviceId } = useActiveDevice()

const range = ref<'today' | '7d' | '30d' | '90d'>('30d')
const chartView = ref<'rev' | 'vol'>('rev')

// Backend filter via month ("jan2025"). Untuk rentang lebih panjang dari
// satu bulan kita query bulan ini saja lalu narrowed di client. 90d = ambil
// semua transaksi device (tanpa month filter).
const monthParam = computed(() => {
  if (range.value === '90d') return undefined
  const d = new Date()
  const m = d.toLocaleDateString('en-US', { month: 'short' }).toLowerCase()
  return `${m}${d.getFullYear()}`
})

// Helper: extract yyyy-mm-dd dari created_at (preferred — sortable ISO)
// atau fallback ke sale_date format mikhmon "jan/02/2006".
function txDateKey(t: Transaction): string {
  if (t.created_at) return t.created_at.split('T')[0]
  return t.sale_date ?? ''
}

// Queries
const { data: apiSales, isLoading: loadingSales } = useSellingQuery(activeDeviceId, monthParam)
const { data: today } = useSellingTodayQuery(activeDeviceId)
const { data: summary } = useSellingSummaryQuery(activeDeviceId, monthParam)

// Filter ke rentang yang dipilih operator (today/7d/30d/90d).
const cutoffDate = computed(() => {
  const cutoff = new Date()
  if (range.value === 'today') {
    cutoff.setHours(0, 0, 0, 0)
  } else if (range.value === '7d') {
    cutoff.setDate(cutoff.getDate() - 7)
  } else if (range.value === '30d') {
    cutoff.setDate(cutoff.getDate() - 30)
  } else if (range.value === '90d') {
    cutoff.setDate(cutoff.getDate() - 90)
  }
  return cutoff.toISOString().split('T')[0]
})

const salesList = computed<Transaction[]>(() => {
  const all = apiSales.value ?? []
  if (range.value === '90d') {
    return all.filter((s) => txDateKey(s) >= cutoffDate.value)
  }
  return all.filter((s) => txDateKey(s) >= cutoffDate.value)
})

// Revenue pakai sell_price kalau ada, fallback ke price.
function revOf(t: Transaction): number {
  return t.sell_price || t.price
}

const todayRevenue = computed(() => (today.value?.transactions ?? []).reduce((a, s) => a + revOf(s), 0))
const todayCount = computed(() => today.value?.count ?? 0)

const totalRevenue = computed(() => salesList.value.reduce((a, s) => a + revOf(s), 0))
const totalCount = computed(() => salesList.value.length)

// Top profile dari list yang sudah di-filter range.
const topProfiles = computed(() => {
  const groups: Record<string, { count: number; total: number }> = {}
  salesList.value.forEach((s) => {
    const key = s.profile || '(no profile)'
    if (!groups[key]) groups[key] = { count: 0, total: 0 }
    groups[key].count++
    groups[key].total += revOf(s)
  })

  return Object.entries(groups)
    .map(([profile, g]) => ({
      profile,
      count: g.count,
      avg: Math.round(g.total / g.count),
      total: g.total,
      share: totalRevenue.value > 0 ? Math.round((g.total / totalRevenue.value) * 100) : 0,
    }))
    .sort((a, b) => b.total - a.total)
    .slice(0, 5)
})

const txCols = computed<ColumnDef<Transaction>[]>(() => [
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, String(row.original.id)),
  },
  {
    accessorKey: 'created_at',
    header: 'Tanggal',
    cell: ({ row }) => {
      const t = row.original
      const d = t.created_at ? new Date(t.created_at) : null
      return h('div', null, [
        h(
          'div',
          { class: 'text-[12px]' },
          d
            ? d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
            : t.sale_date,
        ),
        h(
          'div',
          { class: 'mono text-[10.5px]', style: 'color: var(--muted)' },
          d ? d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) : t.sale_time,
        ),
      ])
    },
  },
  {
    accessorKey: 'username',
    header: 'Username',
    cell: ({ row }) =>
      h('span', { class: 'mono text-[12px] font-semibold' }, row.original.username),
  },
  {
    accessorKey: 'profile',
    header: 'Profile',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.profile ?? '—'),
  },
  {
    accessorKey: 'validity',
    header: 'Validity',
    cell: ({ row }) =>
      h('span', { class: 'mono text-[11.5px]' }, row.original.validity ?? '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'price',
    header: 'Harga',
    cell: ({ row }) => {
      const t = row.original
      const isProfit = t.sell_price > t.price
      return h(
        'div',
        { class: 'flex flex-col items-end' },
        [
          h(
            'span',
            {
              class: 'mono tabular text-[12.5px] font-semibold',
              style: 'color: var(--accent-lime)',
            },
            fmtRpShort(t.sell_price || t.price),
          ),
          isProfit
            ? h(
                'span',
                { class: 'mono text-[10px]', style: 'color: var(--muted)' },
                `modal ${fmtRpShort(t.price)}`,
              )
            : null,
        ],
      )
    },
    meta: { align: 'right' },
  },
])

// Kelompokkan data harian untuk grafik (15 hari terakhir).
const chartData = computed(() => {
  const days: Record<string, { revenue: number; volume: number }> = {}
  for (let i = 14; i >= 0; i--) {
    const d = new Date()
    d.setDate(d.getDate() - i)
    const key = d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
    days[key] = { revenue: 0, volume: 0 }
  }

  salesList.value.forEach((s) => {
    const dateStr = txDateKey(s)
    if (!dateStr) return
    const d = new Date(dateStr)
    const key = d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
    if (days[key]) {
      days[key].revenue += revOf(s)
      days[key].volume++
    }
  })

  return {
    labels: Object.keys(days),
    revenue: Object.values(days).map((d) => d.revenue),
    volume: Object.values(days).map((d) => d.volume),
  }
})

const chartSeries = computed(() => [
  chartView.value === 'rev'
    ? { name: 'Revenue (Rp)', data: chartData.value.revenue, color: 'var(--accent-cyan)' }
    : {
        name: 'Volume (Voucher)',
        data: chartData.value.volume,
        color: 'var(--accent-violet)',
      },
])

async function downloadReport() {
  if (!activeDeviceId.value) return
  try {
    const blob = await reportsService.exportCsv(activeDeviceId.value, {
      month: monthParam.value,
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `reports-sales-${range.value}-${todayStamp()}.csv`
    document.body.appendChild(a)
    a.click()
    a.remove()
    toast.success('Laporan CSV berhasil diunduh')
  } catch (err) {
    toast.error(`Gagal mengunduh CSV: ${(err as Error).message || err}`)
  }
}
</script>

<template>
  <div class="fade-in">
    <PageHeader title="Laporan" subtitle="Penjualan voucher, performa operator, top profiles">
      <template #right>
        <Segmented
          v-model="range"
          :options="[
            { value: 'today', label: 'Hari ini' },
            { value: '7d', label: '7d' },
            { value: '30d', label: '30d' },
            { value: '90d', label: '90d' },
          ]"
        />
        <button class="btn btn-sm" type="button" @click="downloadReport">
          <Icon name="Download" :size="13" />
          CSV
        </button>
      </template>
    </PageHeader>

    <div v-if="loadingSales" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading reports...</div>
    </div>

    <div v-else>
      <div class="mb-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <OverviewKpiCard
          label="Hari ini"
          :value="fmtRpShort(todayRevenue)"
          :delta="`${todayCount} transaksi`"
          trend="up"
          icon="Calendar"
          accent="lime"
        />
        <OverviewKpiCard
          label="Total Revenue"
          :value="fmtRpShort(totalRevenue)"
          :delta="`${totalCount} transaksi`"
          trend="up"
          icon="Report"
          accent="violet"
        />
        <OverviewKpiCard
          label="Penjualan Bulan Ini"
          :value="fmtRpShort(summary?.total_sell_price || 0)"
          :delta="`${summary?.count || 0} voucher`"
          trend="flat"
          icon="Calendar"
          accent="cyan"
        />
        <OverviewKpiCard
          label="Profit Bulan Ini"
          :value="fmtRpShort(summary?.profit || 0)"
          delta="sell − modal"
          trend="up"
          icon="Sparkles"
          accent="cyan"
        />
      </div>

      <Card class="mb-4">
        <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
          <div>
            <div
              class="text-xs font-medium uppercase"
              style="color: var(--muted); letter-spacing: 0.05em"
            >
              Tren Penjualan
            </div>
            <div class="mt-1 flex items-baseline gap-2">
              <span class="text-2xl font-semibold" style="letter-spacing: -0.02em">{{
                fmtRp(totalRevenue)
              }}</span>
              <Badge tone="lime">{{ range }} filter</Badge>
            </div>
          </div>
          <Segmented
            v-model="chartView"
            :options="[
              { value: 'rev', label: 'Revenue' },
              { value: 'vol', label: 'Volume' },
            ]"
          />
        </div>
        <TrendChart :series="chartSeries" :x-labels="chartData.labels" :height="240" />
      </Card>

      <Card class="mb-4">
        <div class="mb-3 text-[13px] font-semibold">Top Profile</div>
        <table class="tbl">
          <thead>
            <tr>
              <th>Profile</th>
              <th>Terjual</th>
              <th>Avg</th>
              <th>Total</th>
              <th>Share</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in topProfiles" :key="p.profile">
              <td>
                <Badge tone="cyan">{{ p.profile }}</Badge>
              </td>
              <td class="mono tabular">{{ p.count }}×</td>
              <td class="mono">{{ fmtRpShort(p.avg) }}</td>
              <td class="mono font-semibold" style="color: var(--accent-lime)">
                {{ fmtRpShort(p.total) }}
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <div class="bar w-20"><i :style="{ width: `${p.share}%` }" /></div>
                  <span class="mono text-[11px]" style="color: var(--muted)">{{ p.share }}%</span>
                </div>
              </td>
            </tr>
            <tr v-if="!topProfiles.length">
              <td colspan="5" class="p-8 text-center text-xs" style="color: var(--muted)">
                Belum ada data penjualan pada periode ini.
              </td>
            </tr>
          </tbody>
        </table>
      </Card>

      <DataTable
        :columns="txCols"
        :data="salesList"
        :get-row-id="(t) => String(t.id)"
        :page-size="10"
      >
        <template #toolbar>
          <span class="text-xs" style="color: var(--muted)"
            >{{ totalCount }} transaksi · {{ range }}</span
          >
        </template>
      </DataTable>
    </div>
  </div>
</template>
