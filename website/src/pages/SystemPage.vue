<script setup lang="ts">
import { computed, h, ref, watch } from 'vue'
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
import { useActiveDevice } from '@/composables/useActiveDevice'
import {
  useIdentityQuery,
  useResourceQuery,
  useLicenseQuery,
  useScriptsQuery,
  useSchedulersQuery,
} from '@/queries/system.queries'
import { systemService } from '@/services/system'
import { buildStreamUrl } from '@/services/stream'
import { useSSE } from '@/composables/useSSE'
import { useToast } from '@/composables/useToast'
import type { SystemScript, SystemScheduler, SystemResource, Routerboard } from '@/types/system'
import type { LogStreamEvent } from '@/types/stream'

const toast = useToast()
const { activeDeviceId } = useActiveDevice()

type TabId = 'overview' | 'scripts' | 'schedulers' | 'logs'
const tab = ref<TabId>('overview')

// Queries
const { data: apiIdentity } = useIdentityQuery(activeDeviceId)
const { data: apiResource, isLoading: loadingResource, refetch: refetchResource } = useResourceQuery(activeDeviceId)
const { data: apiLicense } = useLicenseQuery(activeDeviceId)
const { data: apiScripts, refetch: refetchScripts } = useScriptsQuery(activeDeviceId)
const { data: apiSchedulers, refetch: refetchSchedulers } = useSchedulersQuery(activeDeviceId)

const tabs = computed(() => [
  { id: 'overview' as const, label: 'Overview', icon: 'Server' as const },
  { id: 'scripts' as const, label: 'Scripts', icon: 'Zap' as const, count: (apiScripts.value ?? []).length },
  {
    id: 'schedulers' as const,
    label: 'Schedulers',
    icon: 'Clock' as const,
    count: (apiSchedulers.value ?? []).length,
  },
  { id: 'logs' as const, label: 'Logs', icon: 'Activity' as const, live: true },
])

// Real-time Resources from SSE (event: resource)
const sseResourceUrl = computed(() => activeDeviceId.value ? buildStreamUrl(activeDeviceId.value, 'system/resource') : null)
const { parsed: sseResource } = useSSE<SystemResource>(sseResourceUrl, ['resource'])

const cpu = computed(() => sseResource.value?.cpu_load ?? apiResource.value?.cpu_load ?? 0)
const freeRamBytes = computed(() => sseResource.value?.free_memory ?? apiResource.value?.free_memory ?? 0)
const totalRamBytes = computed(() => sseResource.value?.total_memory ?? apiResource.value?.total_memory ?? 1)
const ramPct = computed(() => Math.round(((totalRamBytes.value - freeRamBytes.value) / totalRamBytes.value) * 100))

const freeDiskBytes = computed(() => sseResource.value?.free_hdd_space ?? apiResource.value?.free_hdd_space ?? 0)
const totalDiskBytes = computed(() => sseResource.value?.total_hdd_space ?? apiResource.value?.total_hdd_space ?? 1)
const diskPct = computed(() => Math.round(((totalDiskBytes.value - freeDiskBytes.value) / totalDiskBytes.value) * 100))

// Real-time Routerboard from SSE (event: routerboard)
const sseRouterboardUrl = computed(() => activeDeviceId.value ? buildStreamUrl(activeDeviceId.value, 'system/routerboard') : null)
const { parsed: sseRouterboard } = useSSE<Routerboard>(sseRouterboardUrl, ['routerboard'])

const rbBoardName = computed(() => sseRouterboard.value?.board_name ?? '')
const rbModel = computed(() => sseRouterboard.value?.model ?? '')
const rbSerial = computed(() => sseRouterboard.value?.serial_number ?? '')
const rbFirmwareType = computed(() => sseRouterboard.value?.firmware_type ?? '')
const rbCurrentFirmware = computed(() => sseRouterboard.value?.current_firmware ?? '')

const scriptCols = computed<ColumnDef<SystemScript>[]>(() => [
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
    cell: ({ row }) => {
      const policyVal = row.original.policy
      const policies = Array.isArray(policyVal)
        ? policyVal
        : typeof policyVal === 'string' && policyVal
          ? (policyVal as string).split(',')
          : []
      return h(
        'div',
        { class: 'flex flex-wrap gap-1' },
        policies.map((p: string) => h(Badge, { tone: 'cyan' }, () => p)),
      )
    },
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'last_started',
    header: 'Last Run',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]', style: 'color: var(--muted)' }, row.original.last_started || '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'run_count',
    header: 'Runs',
    cell: ({ row }) =>
      h(
        'span',
        { class: 'mono tabular text-[12px]' },
        (row.original.run_count ?? 0).toLocaleString('id-ID'),
      ),
  },
])

const schedCols = computed<ColumnDef<SystemScheduler>[]>(() => [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) =>
      h('div', null, [
        h('div', { class: 'mono text-[13px] font-medium' }, row.original.name),
        h(
          'div',
          { class: 'mono text-[11px]', style: 'color: var(--muted)' },
          `${row.original.start_date || ''} ${row.original.start_time || ''}`,
        ),
      ]),
  },
  {
    accessorKey: 'interval',
    header: 'Interval',
    cell: ({ row }) => h(Badge, { tone: 'violet' }, () => `every ${row.original.interval}`),
  },
  {
    accessorKey: 'on_event',
    header: 'On Event',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.on_event || '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'run_count',
    header: 'Runs',
    cell: ({ row }) =>
      h('span', { class: 'mono tabular' }, (row.original.run_count ?? 0).toLocaleString('id-ID')),
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

// Real-time Logs via SSE (event: log)
const paused = ref(false)
const liveLogs = ref<LogStreamEvent[]>([])
const sseLogsUrl = computed(() => activeDeviceId.value ? buildStreamUrl(activeDeviceId.value, 'log') : null)
const { parsed: latestLog } = useSSE<LogStreamEvent>(sseLogsUrl, ['log'])

watch(latestLog, (log) => {
  if (log && !paused.value) {
    liveLogs.value = [log, ...liveLogs.value].slice(0, 150) // keep last 150 log entries
  }
})

const logFilter = ref<'all' | 'hotspot' | 'system' | 'firewall'>('all')
const filteredLogs = computed(() => {
  const logs = liveLogs.value
  if (logFilter.value === 'all') return logs
  return logs.filter((l) => (l.topics || '').toLowerCase().includes(logFilter.value))
})

async function reboot() {
  if (!activeDeviceId.value) return
  if (!confirm('Apakah Anda yakin ingin me-reboot router?')) return
  try {
    await systemService.reboot(activeDeviceId.value)
    toast.success('Router sedang reboot…')
  } catch (err) {
    toast.error(`Gagal me-reboot: ${err instanceof Error ? err.message : String(err)}`)
  }
}

async function shutdown() {
  if (!activeDeviceId.value) return
  if (!confirm('Apakah Anda yakin ingin mematikan router (shutdown)?')) return
  try {
    await systemService.shutdown(activeDeviceId.value)
    toast.warning('Router sedang shutdown…')
  } catch (err) {
    toast.error(`Gagal mematikan router: ${err instanceof Error ? err.message : String(err)}`)
  }
}

function reload() {
  if (activeDeviceId.value) {
    refetchResource()
    refetchScripts()
    refetchSchedulers()
    toast.success('Data System disinkronkan!')
  }
}
</script>

<template>
  <div class="fade-in">
    <PageHeader
      :title="`System · ${apiIdentity?.identity || 'MikroTik'}`"
      subtitle="Identity, resource, scripts, schedulers, logs"
    >
      <template #right>
        <button class="btn btn-sm" type="button" @click="reload">
          <Icon name="Refresh" :size="13" />
          Reload
        </button>
        <button class="btn btn-danger btn-sm" type="button" @click="reboot">
          <Icon name="Power" :size="13" />
          Reboot
        </button>
      </template>
    </PageHeader>

    <div v-if="loadingResource" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading system resource...</div>
    </div>

    <div v-else>
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
                <div class="text-base font-semibold">{{ apiResource?.board_name || rbBoardName || 'Mikrotik Board' }}</div>
                <div class="text-xs" style="color: var(--muted)">
                  {{ apiResource?.version || 'RouterOS' }} · {{ apiResource?.architecture_name || 'unknown' }}
                </div>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3 text-xs sm:grid-cols-3">
              <div>
                <div style="color: var(--muted)">Board</div>
                <div class="mono mt-0.5 font-medium">{{ apiResource?.board_name || rbBoardName || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">Architecture</div>
                <div class="mono mt-0.5 font-medium">{{ apiResource?.architecture_name || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">CPU</div>
                <div class="mono mt-0.5 font-medium">{{ apiResource?.cpu || '—' }} @ {{ apiResource?.cpu_frequency || 0 }} MHz</div>
              </div>
              <div>
                <div style="color: var(--muted)">CPU Count</div>
                <div class="mono mt-0.5 font-medium">{{ apiResource?.cpu_count || 1 }} core(s)</div>
              </div>
              <div>
                <div style="color: var(--muted)">Uptime</div>
                <div class="mono mt-0.5 font-medium">{{ apiResource?.uptime || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">Model</div>
                <div class="mono mt-0.5 font-medium">{{ rbModel || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">Serial</div>
                <div class="mono mt-0.5 font-medium">{{ rbSerial || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">Firmware</div>
                <div class="mono mt-0.5 font-medium">{{ rbCurrentFirmware || '—' }}</div>
              </div>
              <div>
                <div style="color: var(--muted)">FW Type</div>
                <div class="mono mt-0.5 font-medium">{{ rbFirmwareType || '—' }}</div>
              </div>
            </div>
            <!-- License info -->
            <div v-if="apiLicense?.software_id" class="mt-3 rounded-lg p-2.5" style="background: var(--bg-2)">
              <div class="text-[11px] font-medium uppercase" style="color: var(--muted); letter-spacing: 0.05em">License</div>
              <div class="mt-1 grid grid-cols-3 gap-2 text-xs">
                <div>
                  <div style="color: var(--muted)">Software ID</div>
                  <div class="mono mt-0.5 font-medium">{{ apiLicense.software_id }}</div>
                </div>
                <div>
                  <div style="color: var(--muted)">N-Level</div>
                  <div class="mono mt-0.5 font-medium">{{ apiLicense.n_level || '—' }}</div>
                </div>
                <div>
                  <div style="color: var(--muted)">Features</div>
                  <div class="mono mt-0.5 font-medium">{{ apiLicense.features || '—' }}</div>
                </div>
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
                <div class="mt-1 text-lg font-semibold">Active Streams</div>
              </div>
              <LiveTag label="LIVE · SSE" />
            </div>
            <div class="grid grid-cols-3 gap-2.5">
              <SystemRingItem
                label="CPU"
                :value="cpu"
                color="var(--accent-cyan)"
                :detail="`${apiResource?.cpu || 'CPU'}`"
              />
              <SystemRingItem
                label="RAM"
                :value="ramPct"
                color="var(--accent-violet)"
                :detail="`${Math.round((totalRamBytes - freeRamBytes)/1024/1024)}MB / ${Math.round(totalRamBytes/1024/1024)}MB`"
              />
              <SystemRingItem
                label="Disk"
                :value="diskPct"
                color="var(--accent-lime)"
                :detail="`${Math.round((totalDiskBytes - freeDiskBytes)/1024/1024)}MB / ${Math.round(totalDiskBytes/1024/1024)}MB`"
              />
            </div>
          </Card>
        </div>

        <div class="space-y-4">
          <Card>
            <div
              class="text-xs font-medium uppercase"
              style="color: var(--muted); letter-spacing: 0.05em"
            >
              System Controls
            </div>
            <div class="divider" />
            <div class="space-y-2">
              <button class="btn btn-danger w-full" type="button" @click="reboot">
                <Icon name="Power" :size="14" />
                Reboot Router
              </button>
              <button class="btn btn-danger w-full btn-ghost" type="button" @click="shutdown">
                <Icon name="Power" :size="14" />
                Shutdown Router
              </button>
            </div>
          </Card>
        </div>
      </div>

      <DataTable
        v-else-if="tab === 'scripts'"
        :columns="scriptCols"
        :data="apiScripts ?? []"
        :get-row-id="(s) => s.id"
        :page-size="10"
      />
      <DataTable
        v-else-if="tab === 'schedulers'"
        :columns="schedCols"
        :data="apiSchedulers ?? []"
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
          <button class="btn btn-xs" type="button" @click="paused = !paused">
            {{ paused ? 'Resume Stream' : 'Pause Stream' }}
          </button>
          <button class="btn btn-xs btn-ghost" type="button" @click="liveLogs = []">
            Clear Logs
          </button>
        </div>
        <div class="mono max-h-[480px] overflow-y-auto text-xs">
          <div
            v-for="(l, idx) in filteredLogs"
            :key="idx"
            class="row-hover flex items-start gap-3 px-4 py-1.5"
            style="border-bottom: 1px solid var(--border)"
          >
            <span class="shrink-0 tabular" style="color: var(--muted-2); width: 120px">{{
              l.time
            }}</span>
            <span
              class="shrink-0 rounded px-1.5 text-[10px]"
              :style="{
                background: (l.topics || '').includes('warning') ? 'rgba(245,158,11,0.12)' : 'var(--bg-2)',
                color: (l.topics || '').includes('warning') ? 'var(--warning)' : 'var(--muted)',
              }"
            >
              {{ l.topics }}
            </span>
            <span class="flex-1">{{ l.message }}</span>
          </div>
          <div v-if="!filteredLogs.length" class="p-8 text-center" style="color: var(--muted)">
            Menunggu log streaming masuk dari router MikroTik...
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
