<script setup lang="ts">
import { computed, h, onBeforeUnmount, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Badge from '@/components/ui/Badge.vue'
import Card from '@/components/ui/Card.vue'
import LiveTag from '@/components/ui/LiveTag.vue'
import DataTable from '@/components/ui/DataTable.vue'
import Segmented from '@/components/ui/Segmented.vue'
import SystemRingItem from '@/components/overview/SystemRingItem.vue'
import {
  SCRIPTS,
  SCHEDULERS,
  LOGS,
  SYSTEM_RESOURCE,
  SYSTEM_CLOCK,
  type FixtureScript,
  type FixtureScheduler,
} from '@/fixtures/system'
import { fmtAgoFromMs } from '@/utils/fmt'

type TabId = 'overview' | 'scripts' | 'schedulers' | 'logs'
const tab = ref<TabId>('overview')

const tabs = computed(() => [
  { id: 'overview' as const, label: 'Overview', icon: 'Server' as const },
  { id: 'scripts' as const, label: 'Scripts', icon: 'Zap' as const, count: SCRIPTS.length },
  {
    id: 'schedulers' as const,
    label: 'Schedulers',
    icon: 'Clock' as const,
    count: SCHEDULERS.length,
  },
  { id: 'logs' as const, label: 'Logs', icon: 'Activity' as const, count: LOGS.length, live: true },
])

const cpu = ref(SYSTEM_RESOURCE.cpu)
const ram = ref(SYSTEM_RESOURCE.ram)
const disk = ref(SYSTEM_RESOURCE.disk)
const intervalId = window.setInterval(() => {
  cpu.value = Math.max(2, Math.min(98, cpu.value + (Math.random() - 0.5) * 6))
  ram.value = Math.max(2, Math.min(98, ram.value + (Math.random() - 0.5) * 4))
  disk.value = Math.max(2, Math.min(98, disk.value + (Math.random() - 0.5) * 1))
}, 1500)
onBeforeUnmount(() => window.clearInterval(intervalId))

const scriptCols = computed<ColumnDef<FixtureScript>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2' }, [
        h(Icon, { name: 'Zap', size: 14, style: 'color: var(--accent-cyan)' }),
        h('span', { class: 'mono text-[13px] font-medium' }, row.original.name),
      ]),
  },
  {
    accessorKey: 'policy',
    header: 'Policy',
    cell: ({ row }) =>
      h(
        'div',
        { class: 'flex flex-wrap gap-1' },
        row.original.policy.split(',').map((p) => h(Badge, { tone: 'cyan' }, () => p)),
      ),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'lastStarted',
    header: 'Last Run',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono text-[12px]', style: 'color: var(--muted)' },
        fmtAgoFromMs(row.original.lastStarted),
      ),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'runCount',
    header: 'Runs',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono tabular text-[12px]' },
        row.original.runCount.toLocaleString('id-ID'),
      ),
  },
  {
    id: '__actions',
    header: '',
    enableSorting: false,
    cell: () =>
      h('div', { class: 'flex justify-end gap-1' }, [
        h('button', { class: 'btn btn-ghost btn-icon btn-xs', title: 'Run' }, [
          h(Icon, { name: 'Zap', size: 13 }),
        ]),
        h('button', { class: 'btn btn-ghost btn-icon btn-xs' }, [
          h(Icon, { name: 'Edit', size: 13 }),
        ]),
      ]),
    meta: { align: 'right' },
  },
])

const schedCols = computed<ColumnDef<FixtureScheduler>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) =>
      h('div', null, [
        h('div', { class: 'mono text-[13px] font-medium' }, row.original.name),
        h(
          'div',
          { class: 'mono text-[11px]', style: 'color: var(--muted)' },
          `${row.original.startDate} ${row.original.startTime}`,
        ),
      ]),
  },
  {
    accessorKey: 'interval',
    header: 'Interval',
    cell: ({ row }) => h(Badge, { tone: 'violet' }, () => `every ${row.original.interval}`),
  },
  {
    accessorKey: 'onEvent',
    header: 'On Event',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.onEvent),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'runCount',
    header: 'Runs',
    cell: ({ row }) =>
      h('span', { class: 'mono tabular' }, row.original.runCount.toLocaleString('id-ID')),
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

const logFilter = ref<'all' | 'hotspot' | 'system' | 'firewall'>('all')
const filteredLogs = computed(() => {
  if (logFilter.value === 'all') return LOGS
  return LOGS.filter((l) => l.topics.toLowerCase().includes(logFilter.value))
})
</script>

<template>
  <div class="fade-in">
    <PageHeader
      :title="`System · ${SYSTEM_RESOURCE.version}`"
      subtitle="Identity, resource, scripts, schedulers, logs"
    >
      <template #right>
        <button class="btn btn-sm" type="button">
          <Icon name="Download" :size="13" />
          Backup
        </button>
        <button class="btn btn-danger btn-sm" type="button">
          <Icon name="Power" :size="13" />
          Reboot
        </button>
      </template>
    </PageHeader>

    <Tabs v-model="tab" :tabs="tabs" class="mb-4" />

    <!-- Overview tab -->
    <div v-if="tab === 'overview'" class="grid gap-4 lg:grid-cols-[1.6fr_1fr]">
      <div class="space-y-4">
        <Card>
          <div class="mb-4 flex items-center gap-3">
            <div
              class="flex h-11 w-11 items-center justify-center rounded-lg"
              style="background: var(--accent-cyan-soft); color: var(--accent-cyan)"
            >
              <Icon name="Server" :size="22" />
            </div>
            <div>
              <div class="text-base font-semibold">{{ SYSTEM_RESOURCE.board }}</div>
              <div class="text-xs" style="color: var(--muted)">
                {{ SYSTEM_RESOURCE.version }} · {{ SYSTEM_RESOURCE.arch }}
              </div>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 text-xs sm:grid-cols-3">
            <div>
              <div style="color: var(--muted)">Board</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.board }}</div>
            </div>
            <div>
              <div style="color: var(--muted)">Architecture</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.arch }}</div>
            </div>
            <div>
              <div style="color: var(--muted)">Serial</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.serial }}</div>
            </div>
            <div>
              <div style="color: var(--muted)">Firmware</div>
              <div class="mono mt-0.5 font-medium">
                {{ SYSTEM_RESOURCE.firmwareType }} · {{ SYSTEM_RESOURCE.firmware }}
              </div>
            </div>
            <div>
              <div style="color: var(--muted)">IP/Mgmt</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.ipMgmt }}</div>
            </div>
            <div>
              <div style="color: var(--muted)">Uptime</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.uptime }}</div>
            </div>
          </div>
        </Card>

        <Card>
          <div class="mb-3 flex items-center justify-between">
            <div>
              <div
                class="text-xs font-medium uppercase"
                style="color: var(--muted); letter-spacing: 0.05em"
              >
                Resource
              </div>
              <div class="mt-1 text-lg font-semibold">Sehat</div>
            </div>
            <LiveTag label="LIVE · 1s" />
          </div>
          <div class="grid grid-cols-3 gap-2.5">
            <SystemRingItem
              label="CPU"
              :value="Math.round(cpu)"
              color="var(--accent-cyan)"
              :detail="`${SYSTEM_RESOURCE.arch} 880MHz`"
            />
            <SystemRingItem
              label="RAM"
              :value="Math.round(ram)"
              color="var(--accent-violet)"
              detail="159 / 256 MB"
            />
            <SystemRingItem
              label="Disk"
              :value="Math.round(disk)"
              color="var(--accent-lime)"
              detail="14 / 80 MB"
            />
          </div>
          <div class="mt-3 grid grid-cols-2 gap-3 text-xs">
            <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
              <div style="color: var(--muted)">CPU Temp</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.temperature }} °C</div>
            </div>
            <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
              <div style="color: var(--muted)">Voltage</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.voltage }} V</div>
            </div>
            <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
              <div style="color: var(--muted)">Power</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.powerW }} W</div>
            </div>
            <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
              <div style="color: var(--muted)">Firmware</div>
              <div class="mono mt-0.5 font-medium">{{ SYSTEM_RESOURCE.firmware }}</div>
            </div>
          </div>
        </Card>
      </div>

      <div class="space-y-4">
        <Card>
          <div
            class="text-xs font-medium uppercase"
            style="color: var(--muted); letter-spacing: 0.05em"
          >
            Clock
          </div>
          <div class="mono mt-2 text-3xl font-bold tabular">{{ SYSTEM_CLOCK.time }}</div>
          <div class="mt-1 text-sm">{{ SYSTEM_CLOCK.date }}</div>
          <div class="divider" />
          <div class="space-y-2 text-xs">
            <div class="flex items-center justify-between">
              <span style="color: var(--muted)">Timezone</span>
              <span class="mono">{{ SYSTEM_CLOCK.timezone }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span style="color: var(--muted)">NTP</span>
              <Badge tone="success" dot>{{ SYSTEM_CLOCK.ntp }}</Badge>
            </div>
          </div>
        </Card>

        <Card>
          <div class="mb-2 text-[13px] font-semibold">Aksi Berbahaya</div>
          <div class="space-y-2">
            <button class="btn w-full" type="button">
              <Icon name="Power" :size="14" />
              Reboot
            </button>
            <button class="btn btn-danger w-full" type="button">
              <Icon name="Power" :size="14" />
              Shutdown
            </button>
            <button class="btn w-full" type="button">
              <Icon name="Download" :size="14" />
              Backup
            </button>
          </div>
        </Card>
      </div>
    </div>

    <DataTable
      v-else-if="tab === 'scripts'"
      :columns="scriptCols"
      :data="SCRIPTS"
      :get-row-id="(s) => s.id"
      :page-size="10"
    />
    <DataTable
      v-else-if="tab === 'schedulers'"
      :columns="schedCols"
      :data="SCHEDULERS"
      :get-row-id="(s) => s.id"
      :page-size="10"
    />

    <Card v-else-if="tab === 'logs'" style="padding: 0">
      <div class="flex items-center gap-3 p-3" style="border-bottom: 1px solid var(--border)">
        <LiveTag label="STREAMING" />
        <Segmented
          v-model="logFilter"
          :options="[
            { value: 'all', label: 'Semua' },
            { value: 'hotspot', label: 'Hotspot' },
            { value: 'system', label: 'System' },
            { value: 'firewall', label: 'Firewall' },
          ]"
        />
      </div>
      <div class="mono max-h-[480px] overflow-y-auto text-xs">
        <div
          v-for="l in filteredLogs"
          :key="l.id"
          class="row-hover flex items-start gap-3 px-4 py-1.5"
          style="border-bottom: 1px solid var(--border)"
        >
          <span class="shrink-0 tabular" style="color: var(--muted-2); width: 64px">{{
            l.time
          }}</span>
          <span
            class="shrink-0 rounded px-1.5 text-[10px]"
            :style="{
              background: l.topics.includes('warning') ? 'rgba(245,158,11,0.12)' : 'var(--bg-2)',
              color: l.topics.includes('warning') ? 'var(--warning)' : 'var(--muted)',
            }"
          >
            {{ l.topics }}
          </span>
          <span class="flex-1">{{ l.message }}</span>
        </div>
      </div>
    </Card>
  </div>
</template>
