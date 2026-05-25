<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import SearchInput from '@/components/ui/SearchInput.vue'
import DataTable from '@/components/ui/DataTable.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useHotspotHostsQuery } from '@/queries/hotspot.queries'
import { fmtBytes } from '@/utils/fmt'
import type { HotspotHost } from '@/types/hotspot'

const { activeDeviceId } = useActiveDevice()

const search = ref('')

const { data: apiHosts, isLoading: loadingHosts } = useHotspotHostsQuery(activeDeviceId)

const filtered = computed(() => {
  const hosts = apiHosts.value ?? []
  if (!search.value) return hosts
  const s = search.value.toLowerCase()
  return hosts.filter(
    (h) =>
      (h.mac_address || '').toLowerCase().includes(s) ||
      (h.address || '').toLowerCase().includes(s) ||
      (h.server || '').toLowerCase().includes(s),
  )
})

const columns = computed<ColumnDef<HotspotHost>[]>(() => [
  {
    accessorKey: 'mac_address',
    header: 'MAC Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.mac_address || '—'),
  },
  {
    accessorKey: 'address',
    header: 'Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.address || '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'server',
    header: 'Server',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.server || 'all'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'to_address',
    header: 'To Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.to_address || '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'idle_time',
    header: 'Idle Time',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.idle_time || '—'),
  },
  {
    id: 'usage',
    header: 'Pemakaian (RX/TX)',
    cell: ({ row }) => {
      const u = row.original
      return h('span', { class: 'mono text-[11px]', style: 'color: var(--muted)' }, [
        `↓ ${fmtBytes(u.bytes_in)} · ↑ ${fmtBytes(u.bytes_out)}`,
      ])
    },
    meta: { mobileHidden: true },
  },
])
</script>

<template>
  <div>
    <div v-if="loadingHosts" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading hosts...</div>
    </div>

    <DataTable
      v-else
      :columns="columns"
      :data="filtered"
      :get-row-id="(h) => h.id"
      :global-filter="search"
      :page-size="10"
      empty-message="Tidak ada hosts"
      @update:global-filter="(v) => (search = v)"
    >
      <template #toolbar>
        <SearchInput v-model="search" placeholder="Cari MAC, address, server..." />
      </template>
    </DataTable>
  </div>
</template>
