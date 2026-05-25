<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import SearchInput from '@/components/ui/SearchInput.vue'
import Badge from '@/components/ui/Badge.vue'
import Avatar from '@/components/ui/Avatar.vue'
import Select from '@/components/ui/Select.vue'
import DataTable from '@/components/ui/DataTable.vue'
import { fmtBytes } from '@/utils/fmt'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useHotspotUsersQuery, useHotspotActiveQuery, useHotspotProfilesQuery } from '@/queries/hotspot.queries'
import type { HotspotUser } from '@/types/hotspot'

const { activeDeviceId } = useActiveDevice()

const search = ref('')
const filterProfile = ref<string>('all')

const { data: apiUsers, isLoading: loadingUsers } = useHotspotUsersQuery(activeDeviceId)
const { data: apiActive } = useHotspotActiveQuery(activeDeviceId)
const { data: apiProfiles } = useHotspotProfilesQuery(activeDeviceId)

const activeUserNames = computed(() => {
  return new Set((apiActive.value ?? []).map((s) => s.user))
})

const inactive = computed(() => {
  const users = apiUsers.value ?? []
  return users.filter((u) => !activeUserNames.value.has(u.name))
})

const filtered = computed(() => {
  return inactive.value.filter((u) => {
    if (search.value) {
      const s = search.value.toLowerCase()
      if (
        !(
          u.name.toLowerCase().includes(s) ||
          (u.mac_address || '').toLowerCase().includes(s) ||
          u.profile.toLowerCase().includes(s)
        )
      )
        return false
    }
    if (filterProfile.value !== 'all' && u.profile !== filterProfile.value) return false
    return true
  })
})

const columns = computed<ColumnDef<HotspotUser>[]>(() => [
  {
    accessorKey: 'name',
    header: 'User',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2.5' }, [
        h(Avatar, { name: row.original.name, size: 28 }),
        h('div', null, [
          h('div', { class: 'text-[13px] font-medium' }, row.original.name),
          h(
            'div',
            { class: 'mono text-[11px]', style: 'color: var(--muted)' },
            row.original.mac_address || '—',
          ),
        ]),
      ]),
  },
  {
    accessorKey: 'profile',
    header: 'Profile',
    cell: ({ row }) => h(Badge, { tone: 'cyan' }, () => row.original.profile),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'server',
    header: 'Server',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.server || 'all'),
    meta: { mobileHidden: true },
  },
  {
    id: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const u = row.original
      if (u.disabled) return h(Badge, { tone: 'neutral' }, () => 'Disabled')
      return h(Badge, { tone: 'neutral', dot: true }, () => 'Offline')
    },
  },
  {
    accessorKey: 'uptime',
    header: 'Uptime',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.uptime || '—'),
    meta: { mobileHidden: true },
  },
  {
    id: 'usage',
    header: 'Pemakaian',
    cell: ({ row }) => {
      const u = row.original
      const total = u.bytes_in + u.bytes_out
      const pct = u.limit_bytes_total > 0 ? Math.min(100, (total / u.limit_bytes_total) * 100) : 0
      return h('div', { class: 'flex min-w-[140px] flex-col gap-1' }, [
        h('div', { class: 'mono text-[11px]', style: 'color: var(--muted)' }, [
          `↓ ${fmtBytes(u.bytes_in)} · ↑ ${fmtBytes(u.bytes_out)}`,
        ]),
        u.limit_bytes_total > 0 ? h('div', { class: 'bar' }, [h('i', { style: `width: ${pct}%` })]) : null,
      ])
    },
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'comment',
    header: 'Comment',
    cell: ({ row }) => h('span', { class: 'text-xs', style: 'color: var(--muted)' }, row.original.comment || '—'),
    meta: { mobileHidden: true },
  },
])
</script>

<template>
  <div>
    <div v-if="loadingUsers" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading inactive users...</div>
    </div>

    <DataTable
      v-else
      :columns="columns"
      :data="filtered"
      :get-row-id="(u) => u.id"
      :global-filter="search"
      :page-size="10"
      empty-message="Tidak ada user yang cocok dengan filter"
      @update:global-filter="(v) => (search = v)"
    >
      <template #toolbar>
        <SearchInput v-model="search" placeholder="Cari nama, MAC, profile..." />
        <Select
          v-model="filterProfile"
          sm
          :options="[
            { value: 'all', label: 'Semua profile' },
            ...(apiProfiles ?? []).map((p) => ({ value: p.name, label: p.name })),
          ]"
        />
      </template>
    </DataTable>
  </div>
</template>
