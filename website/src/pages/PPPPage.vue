<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Badge from '@/components/ui/Badge.vue'
import Avatar from '@/components/ui/Avatar.vue'
import Card from '@/components/ui/Card.vue'
import DataTable from '@/components/ui/DataTable.vue'
import OverviewKpiCard from '@/components/overview/OverviewKpiCard.vue'
import {
  PPP_SECRETS,
  PPP_ACTIVE,
  PPP_PROFILES,
  PPP_INACTIVE,
  type FixturePPPSecret,
  type FixturePPPActive,
} from '@/fixtures/ppp'
import { fmtBytes, fmtDuration, fmtRate, fmtAgoFromMs } from '@/utils/fmt'

type TabId = 'secret' | 'active' | 'profile' | 'inactive'
const tab = ref<TabId>('secret')

const tabs = computed(() => [
  { id: 'secret' as const, label: 'Secret', icon: 'Lock' as const, count: PPP_SECRETS.length },
  {
    id: 'active' as const,
    label: 'Active',
    icon: 'Activity' as const,
    count: PPP_ACTIVE.length,
    live: true,
  },
  { id: 'profile' as const, label: 'Profile', icon: 'Wifi' as const, count: PPP_PROFILES.length },
  {
    id: 'inactive' as const,
    label: 'Inactive',
    icon: 'Power' as const,
    count: PPP_INACTIVE.length,
  },
])

const secretCols = computed<ColumnDef<FixturePPPSecret>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) =>
      h('div', null, [
        h('div', { class: 'mono text-[13px] font-medium' }, row.original.name),
        row.original.comment
          ? h('div', { class: 'text-[11px]', style: 'color: var(--muted)' }, row.original.comment)
          : null,
      ]),
  },
  {
    accessorKey: 'profile',
    header: 'Profile',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.profile),
  },
  {
    accessorKey: 'service',
    header: 'Service',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.service),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'address',
    header: 'Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.address),
    meta: { mobileHidden: true },
  },
  {
    id: 'status',
    header: 'Status',
    cell: ({ row }) => {
      if (row.original.disabled) return h(Badge, { tone: 'neutral' }, () => 'Disabled')
      return h(Badge, { tone: row.original.isActive ? 'success' : 'neutral', dot: true }, () =>
        row.original.isActive ? 'Online' : 'Offline',
      )
    },
  },
  {
    accessorKey: 'lastLoggedOut',
    header: 'Last logout',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[12px]', style: 'color: var(--muted)' },
        fmtAgoFromMs(row.original.lastLoggedOut),
      ),
    meta: { mobileHidden: true },
  },
  {
    id: '__actions',
    header: '',
    enableSorting: false,
    cell: () =>
      h('div', { class: 'flex justify-end gap-1' }, [
        h('button', { class: 'btn btn-ghost btn-icon btn-xs' }, [
          h(Icon, { name: 'Edit', size: 13 }),
        ]),
        h('button', { class: 'btn btn-ghost btn-icon btn-xs' }, [
          h(Icon, { name: 'More', size: 13 }),
        ]),
      ]),
    meta: { align: 'right' },
  },
])

const activeCols = computed<ColumnDef<FixturePPPActive>[]>(() => [
  {
    accessorKey: 'name',
    header: 'User',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2.5' }, [
        h(Avatar, { name: row.original.name, size: 28 }),
        h('span', { class: 'mono text-[13px] font-medium' }, row.original.name),
      ]),
  },
  {
    accessorKey: 'profile',
    header: 'Profile',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.profile),
  },
  {
    accessorKey: 'address',
    header: 'Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.address),
    meta: { mobileHidden: true },
  },
  {
    id: 'uptime',
    header: 'Uptime',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[12px]' },
        fmtDuration(Math.floor((Date.now() - row.original.uptimeStart) / 1000)),
      ),
    meta: { mobileHidden: true },
  },
  {
    id: 'throughput',
    header: 'Throughput',
    cell: ({ row }) =>
      h('div', { class: 'mono text-[11px]' }, [
        h('div', { style: 'color: var(--accent-cyan)' }, `↓ ${fmtRate(row.original.rxRate)}`),
        h('div', { style: 'color: var(--accent-violet)' }, `↑ ${fmtRate(row.original.txRate)}`),
      ]),
  },
  {
    id: 'bytes',
    header: 'Total',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[12px]' },
        fmtBytes(row.original.bytesIn + row.original.bytesOut),
      ),
    meta: { mobileHidden: true },
  },
])

const totalRX = computed(() => PPP_ACTIVE.reduce((a, s) => a + s.rxRate, 0))
const totalTX = computed(() => PPP_ACTIVE.reduce((a, s) => a + s.txRate, 0))
</script>

<template>
  <div class="fade-in">
    <PageHeader title="PPP" subtitle="PPPoE secrets, active sessions, profiles & inactive">
      <template #right>
        <button class="btn btn-sm" type="button">
          <Icon name="Refresh" :size="13" />
          Reload
        </button>
        <button class="btn btn-primary btn-sm" type="button">
          <Icon name="Plus" :size="14" />
          Tambah {{ tab === 'profile' ? 'Profile' : 'Secret' }}
        </button>
      </template>
    </PageHeader>

    <Tabs v-model="tab" :tabs="tabs" class="mb-4" />

    <!-- Secret -->
    <div v-if="tab === 'secret'">
      <DataTable
        :columns="secretCols"
        :data="PPP_SECRETS"
        :get-row-id="(s) => s.id"
        :page-size="10"
        enable-row-selection
      >
        <template #toolbar>
          <span class="text-xs" style="color: var(--muted)">{{ PPP_SECRETS.length }} secrets</span>
        </template>
      </DataTable>
    </div>

    <!-- Active -->
    <div v-else-if="tab === 'active'">
      <div class="mb-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <OverviewKpiCard
          label="Active Sessions"
          :value="PPP_ACTIVE.length"
          delta="live"
          trend="flat"
          icon="Activity"
          accent="cyan"
          live
        />
        <OverviewKpiCard
          label="Throughput RX"
          :value="fmtRate(totalRX)"
          delta="total"
          trend="up"
          icon="ArrowUpRight"
          accent="cyan"
        />
        <OverviewKpiCard
          label="Throughput TX"
          :value="fmtRate(totalTX)"
          delta="total"
          trend="up"
          icon="ArrowUpRight"
          accent="violet"
        />
        <OverviewKpiCard
          label="Encoding"
          value="mppe128"
          delta="all"
          trend="flat"
          icon="Lock"
          accent="lime"
        />
      </div>
      <DataTable
        :columns="activeCols"
        :data="PPP_ACTIVE"
        :get-row-id="(s) => s.id"
        :page-size="10"
      />
    </div>

    <!-- Profile -->
    <div v-else-if="tab === 'profile'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="p in PPP_PROFILES" :key="p.id" accent="var(--accent-violet)">
        <div class="mb-3 flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div
              class="flex h-9 w-9 items-center justify-center rounded-lg"
              style="background: var(--accent-violet-soft); color: var(--accent-violet)"
            >
              <Icon name="Wifi" :size="16" />
            </div>
            <div>
              <div class="text-sm font-semibold">{{ p.name }}</div>
              <div class="mono text-[11px]" style="color: var(--muted)">{{ p.id }}</div>
            </div>
          </div>
          <Badge tone="violet">{{ p.sessions }} sesi</Badge>
        </div>
        <div class="grid grid-cols-2 gap-2 text-xs">
          <div>
            <div style="color: var(--muted)">Rate Limit</div>
            <div class="mono font-medium">{{ p.rateLimit }}</div>
          </div>
          <div>
            <div style="color: var(--muted)">Local</div>
            <div class="mono font-medium">{{ p.localAddress }}</div>
          </div>
          <div>
            <div style="color: var(--muted)">Pool</div>
            <div class="mono font-medium">{{ p.remoteAddress }}</div>
          </div>
          <div>
            <div style="color: var(--muted)">Parent</div>
            <div class="mono font-medium">{{ p.parentQueue }}</div>
          </div>
          <div class="col-span-2">
            <div style="color: var(--muted)">DNS</div>
            <div class="mono font-medium">{{ p.dnsServer }}</div>
          </div>
        </div>
        <div class="mt-3 flex gap-2">
          <button class="btn btn-sm flex-1" type="button">Edit</button>
          <button class="btn btn-sm flex-1" type="button">View users</button>
        </div>
      </Card>
    </div>

    <!-- Inactive -->
    <div v-else-if="tab === 'inactive'">
      <DataTable
        :columns="secretCols"
        :data="PPP_INACTIVE"
        :get-row-id="(s) => s.id"
        :page-size="10"
        empty-message="Tidak ada secret inactive"
      />
    </div>
  </div>
</template>
