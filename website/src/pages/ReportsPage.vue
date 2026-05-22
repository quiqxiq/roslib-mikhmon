<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import Avatar from '@/components/ui/Avatar.vue'
import Segmented from '@/components/ui/Segmented.vue'
import Card from '@/components/ui/Card.vue'
import TrendChart from '@/components/ui/charts/TrendChart.vue'
import DataTable from '@/components/ui/DataTable.vue'
import OverviewKpiCard from '@/components/overview/OverviewKpiCard.vue'
import {
  TRANSACTIONS,
  REVENUE_30D,
  OPERATORS,
  TOP_PROFILES,
  type FixtureTransaction,
} from '@/fixtures/reports'
import { fmtRpShort, fmtRp } from '@/utils/fmt'

const range = ref<'today' | '7d' | '30d' | '90d'>('30d')
const chartView = ref<'rev' | 'vol'>('rev')

const todayRevenue = computed(() => 478_000)
const yesterdayRevenue = computed(() => 412_000)
const monthRevenue = computed(() => TRANSACTIONS.reduce((a, t) => a + t.price, 0))
const monthCount = computed(() => TRANSACTIONS.length)
const topOp = computed(() => OPERATORS[0])

const txCols = computed<ColumnDef<FixtureTransaction>[]>(() => [
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.id),
  },
  {
    accessorKey: 'date',
    header: 'Tanggal',
    cell: ({ row }) => {
      const d = new Date(row.original.date)
      return h('div', null, [
        h(
          'div',
          { class: 'text-[12px]' },
          d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }),
        ),
        h(
          'div',
          { class: 'mono text-[10.5px]', style: 'color: var(--muted)' },
          d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }),
        ),
      ])
    },
  },
  {
    accessorKey: 'user',
    header: 'User',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.user),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'profile',
    header: 'Profile',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.profile),
  },
  {
    accessorKey: 'source',
    header: 'Sumber',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2 text-[12px]' }, [
        h(Icon, {
          name: row.original.source === 'voucher-gen' ? 'Sparkles' : 'Edit',
          size: 13,
          style: 'color: var(--muted)',
        }),
        h('span', null, row.original.source),
      ]),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'operator',
    header: 'Operator',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2' }, [
        h(Avatar, { name: row.original.operator, size: 22 }),
        h('span', { class: 'text-[12px]' }, row.original.operator),
      ]),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'price',
    header: 'Harga',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono tabular text-[12.5px] font-semibold', style: 'color: var(--accent-lime)' },
        fmtRpShort(row.original.price),
      ),
    meta: { align: 'right' },
  },
])

const chartSeries = computed(() => [
  chartView.value === 'rev'
    ? { name: 'Revenue', data: REVENUE_30D.values, color: 'var(--accent-cyan)' }
    : {
        name: 'Volume',
        data: REVENUE_30D.values.map((v) => Math.round(v / 10000)),
        color: 'var(--accent-violet)',
      },
])
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
        <button class="btn btn-sm" type="button">
          <Icon name="Download" :size="13" />
          CSV
        </button>
        <button class="btn btn-sm" type="button">
          <Icon name="Print" :size="13" />
          Cetak
        </button>
      </template>
    </PageHeader>

    <div class="mb-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <OverviewKpiCard
        label="Hari ini"
        :value="fmtRpShort(todayRevenue)"
        delta="32 transaksi"
        trend="up"
        icon="Calendar"
        accent="lime"
      />
      <OverviewKpiCard
        label="Kemarin"
        :value="fmtRpShort(yesterdayRevenue)"
        delta="+16% vs today"
        trend="up"
        icon="Calendar"
        accent="cyan"
      />
      <OverviewKpiCard
        label="Bulan ini"
        :value="fmtRpShort(monthRevenue)"
        :delta="`${monthCount} tx · avg ${fmtRpShort(monthRevenue / Math.max(1, monthCount))}`"
        trend="up"
        icon="Report"
        accent="violet"
      />
      <OverviewKpiCard
        label="Operator terbaik"
        :value="topOp.name"
        :delta="`${topOp.count}× · ${fmtRpShort(topOp.revenue)}`"
        trend="up"
        icon="Users"
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
            Penjualan 30 Hari Terakhir
          </div>
          <div class="mt-1 flex items-baseline gap-2">
            <span class="text-2xl font-semibold" style="letter-spacing: -0.02em">{{
              fmtRp(monthRevenue)
            }}</span>
            <Badge tone="lime">↑ 12.4% MoM</Badge>
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
      <TrendChart :series="chartSeries" :x-labels="REVENUE_30D.labels" :height="240" />
    </Card>

    <div class="mb-4 grid gap-4 lg:grid-cols-2">
      <Card>
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
            <tr v-for="p in TOP_PROFILES" :key="p.profile">
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
          </tbody>
        </table>
      </Card>

      <Card>
        <div class="mb-3 text-[13px] font-semibold">Performa Operator</div>
        <div class="space-y-3">
          <div
            v-for="op in OPERATORS"
            :key="op.name"
            class="row-hover flex items-center gap-3 rounded-lg p-2"
          >
            <Avatar :name="op.name" :size="32" />
            <div class="flex-1">
              <div class="text-[13px] font-medium">{{ op.name }}</div>
              <div class="text-[11px]" style="color: var(--muted)">{{ op.count }} transaksi</div>
            </div>
            <div class="text-right">
              <div class="mono text-sm font-semibold" style="color: var(--accent-lime)">
                {{ fmtRpShort(op.revenue) }}
              </div>
              <div class="bar mt-1 w-20">
                <i :style="{ width: `${(op.revenue / OPERATORS[0].revenue) * 100}%` }" />
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>

    <DataTable :columns="txCols" :data="TRANSACTIONS" :get-row-id="(t) => t.id" :page-size="10">
      <template #toolbar>
        <span class="text-xs" style="color: var(--muted)"
          >{{ TRANSACTIONS.length }} transaksi · {{ range }}</span
        >
      </template>
    </DataTable>
  </div>
</template>
