<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Badge from '@/components/ui/Badge.vue'
import Card from '@/components/ui/Card.vue'
import DataTable from '@/components/ui/DataTable.vue'
import OverviewKpiCard from '@/components/overview/OverviewKpiCard.vue'
import {
  INTERFACES,
  IP_POOLS,
  ARP,
  DHCP_LEASES,
  QUEUES,
  type FixtureInterface,
  type FixtureARP,
  type FixtureDHCPLease,
  type FixtureQueue,
} from '@/fixtures/network'
import { fmtBytes, fmtDuration, fmtRate } from '@/utils/fmt'

type TabId = 'interfaces' | 'pools' | 'arp' | 'dhcp' | 'queues'
const tab = ref<TabId>('interfaces')

const tabs = computed(() => [
  {
    id: 'interfaces' as const,
    label: 'Interfaces',
    icon: 'Network' as const,
    count: INTERFACES.length,
  },
  { id: 'pools' as const, label: 'IP Pools', icon: 'Globe' as const, count: IP_POOLS.length },
  { id: 'arp' as const, label: 'ARP', icon: 'Activity' as const, count: ARP.length },
  { id: 'dhcp' as const, label: 'DHCP', icon: 'Server' as const, count: DHCP_LEASES.length },
  { id: 'queues' as const, label: 'Queues', icon: 'Boot' as const, count: QUEUES.length },
])

const ifaceCols = computed<ColumnDef<FixtureInterface>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Interface',
    cell: ({ row }) =>
      h('div', null, [
        h('div', { class: 'mono text-[13px] font-medium' }, row.original.name),
        row.original.comment
          ? h('div', { class: 'text-[11px]', style: 'color: var(--muted)' }, row.original.comment)
          : null,
      ]),
  },
  {
    accessorKey: 'type',
    header: 'Type',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.type),
  },
  {
    accessorKey: 'link',
    header: 'Link',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.link),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'macAddress',
    header: 'MAC',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[11px]', style: 'color: var(--muted)' },
        row.original.macAddress,
      ),
    meta: { mobileHidden: true },
  },
  {
    id: 'throughput',
    header: 'Throughput',
    cell: ({ row }) =>
      row.original.running
        ? h('div', { class: 'mono text-[11px]' }, [
            h('div', { style: 'color: var(--accent-cyan)' }, `↓ ${fmtRate(row.original.rxRate)}`),
            h('div', { style: 'color: var(--accent-violet)' }, `↑ ${fmtRate(row.original.txRate)}`),
          ])
        : h('span', { style: 'color: var(--muted-2)' }, '—'),
  },
  {
    id: 'total',
    header: 'Total',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[12px]' },
        `${fmtBytes(row.original.rxBytes)} / ${fmtBytes(row.original.txBytes)}`,
      ),
    meta: { mobileHidden: true },
  },
  {
    id: 'status',
    header: 'Status',
    cell: ({ row }) => {
      if (row.original.disabled) return h(Badge, { tone: 'neutral' }, () => 'Disabled')
      if (!row.original.running) return h(Badge, { tone: 'warn' }, () => 'No link')
      return h(Badge, { tone: 'success', dot: true }, () => 'Running')
    },
  },
])

const arpCols = computed<ColumnDef<FixtureARP>[]>(() => [
  {
    accessorKey: 'address',
    header: 'IP Address',
    cell: ({ row }) => h('span', { class: 'mono text-[13px]' }, row.original.address),
  },
  {
    accessorKey: 'macAddress',
    header: 'MAC Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.macAddress),
  },
  {
    accessorKey: 'interface',
    header: 'Interface',
    cell: ({ row }) =>
      h(
        Badge,
        { tone: row.original.interface.startsWith('wlan') ? 'violet' : 'cyan' },
        () => row.original.interface,
      ),
  },
  {
    accessorKey: 'dynamic',
    header: 'Type',
    cell: ({ row }) =>
      h(Badge, { tone: row.original.dynamic ? 'cyan' : 'neutral' }, () =>
        row.original.dynamic ? 'Dynamic' : 'Static',
      ),
    meta: { mobileHidden: true },
  },
])

const dhcpCols = computed<ColumnDef<FixtureDHCPLease>[]>(() => [
  {
    accessorKey: 'address',
    header: 'IP / Host',
    cell: ({ row }) =>
      h('div', null, [
        h('div', { class: 'mono text-[13px]' }, row.original.address),
        h('div', { class: 'text-[11px]', style: 'color: var(--muted)' }, row.original.hostName),
      ]),
  },
  {
    accessorKey: 'macAddress',
    header: 'MAC',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.macAddress),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'interface',
    header: 'Interface',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.interface),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) =>
      h(
        Badge,
        { tone: row.original.status === 'bound' ? 'success' : 'warn', dot: true },
        () => row.original.status,
      ),
  },
  {
    accessorKey: 'expiresAfter',
    header: 'Expires',
    cell: ({ row }) =>
      h('span', { class: 'mono text-[12px]' }, fmtDuration(row.original.expiresAfter)),
    meta: { mobileHidden: true },
  },
])

const queueCols = computed<ColumnDef<FixtureQueue>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Queue',
    cell: ({ row }) => h('span', { class: 'mono text-[13px] font-medium' }, row.original.name),
  },
  {
    accessorKey: 'target',
    header: 'Target',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.target),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'maxLimit',
    header: 'Max',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.maxLimit),
  },
  {
    id: 'current',
    header: 'Current',
    cell: ({ row }) =>
      h('div', { class: 'mono text-[11px]' }, [
        h('div', { style: 'color: var(--accent-cyan)' }, `↓ ${fmtRate(row.original.curRx)}`),
        h('div', { style: 'color: var(--accent-violet)' }, `↑ ${fmtRate(row.original.curTx)}`),
      ]),
  },
  {
    id: 'util',
    header: 'Util',
    cell: ({ row }) => {
      const pct = Math.round(
        ((row.original.curRx + row.original.curTx) / (row.original.maxLimitBps * 2)) * 100,
      )
      return h('div', { class: 'flex min-w-[120px] flex-col gap-1' }, [
        h('div', { class: 'mono text-[11px]', style: 'color: var(--muted)' }, `${pct}%`),
        h('div', { class: 'bar' }, [h('i', { style: `width: ${pct}%` })]),
      ])
    },
    meta: { mobileHidden: true },
  },
  {
    id: 'status',
    header: 'Status',
    cell: ({ row }) =>
      row.original.disabled
        ? h(Badge, { tone: 'neutral' }, () => 'Disabled')
        : h(Badge, { tone: 'success', dot: true }, () => 'Active'),
  },
])

const runningCount = computed(() => INTERFACES.filter((i) => i.running).length)
const dhcpBound = computed(() => DHCP_LEASES.filter((l) => l.status === 'bound').length)
const arpDynamic = computed(() => ARP.filter((a) => a.dynamic).length)
</script>

<template>
  <div class="fade-in">
    <PageHeader title="Network" subtitle="Interfaces, IP pools, ARP, DHCP, queues">
      <template #right>
        <button class="btn btn-sm" type="button">
          <Icon name="Refresh" :size="13" />
          Reload
        </button>
      </template>
    </PageHeader>

    <!-- Health KPIs -->
    <div class="mb-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <OverviewKpiCard
        label="Interfaces running"
        :value="`${runningCount} / ${INTERFACES.length}`"
        :delta="`${INTERFACES.length - runningCount} idle`"
        trend="flat"
        icon="Network"
        accent="cyan"
      />
      <OverviewKpiCard
        label="DHCP Leases"
        :value="DHCP_LEASES.length"
        :delta="`${dhcpBound} bound`"
        trend="up"
        icon="Server"
        accent="violet"
      />
      <OverviewKpiCard
        label="ARP Entries"
        :value="ARP.length"
        :delta="`${arpDynamic} dynamic`"
        trend="flat"
        icon="Activity"
        accent="lime"
      />
      <OverviewKpiCard
        label="Simple Queues"
        :value="QUEUES.length"
        delta="aktif"
        trend="flat"
        icon="Boot"
        accent="cyan"
      />
    </div>

    <Tabs v-model="tab" :tabs="tabs" class="mb-4" />

    <DataTable
      v-if="tab === 'interfaces'"
      :columns="ifaceCols"
      :data="INTERFACES"
      :get-row-id="(i) => i.id"
      :page-size="10"
    />

    <div v-else-if="tab === 'pools'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="p in IP_POOLS" :key="p.id" accent="var(--accent-cyan)">
        <div class="mb-2 flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div
              class="flex h-9 w-9 items-center justify-center rounded-lg"
              style="background: var(--accent-cyan-soft); color: var(--accent-cyan)"
            >
              <Icon name="Globe" :size="16" />
            </div>
            <div>
              <div class="text-sm font-semibold">{{ p.name }}</div>
              <div class="mono text-[11px]" style="color: var(--muted)">{{ p.ranges }}</div>
            </div>
          </div>
        </div>
        <div class="mt-3 flex items-baseline justify-between">
          <span class="text-xs" style="color: var(--muted)">Used</span>
          <span class="mono text-sm font-semibold">{{ p.used }} / {{ p.total }}</span>
        </div>
        <div class="bar mt-1.5">
          <i
            :style="{
              width: `${(p.used / p.total) * 100}%`,
              background:
                p.used / p.total > 0.95
                  ? 'var(--danger)'
                  : p.used / p.total > 0.8
                    ? 'var(--warning)'
                    : 'linear-gradient(90deg, var(--accent-cyan), var(--accent-violet))',
            }"
          />
        </div>
        <div class="mt-2 text-[11px]" style="color: var(--muted)">
          {{ p.total - p.used }} IP tersedia
        </div>
      </Card>
    </div>

    <DataTable
      v-else-if="tab === 'arp'"
      :columns="arpCols"
      :data="ARP"
      :get-row-id="(a) => a.id"
      :page-size="10"
    />
    <DataTable
      v-else-if="tab === 'dhcp'"
      :columns="dhcpCols"
      :data="DHCP_LEASES"
      :get-row-id="(l) => l.id"
      :page-size="10"
    />
    <DataTable
      v-else-if="tab === 'queues'"
      :columns="queueCols"
      :data="QUEUES"
      :get-row-id="(q) => q.id"
      :page-size="10"
    />
  </div>
</template>
