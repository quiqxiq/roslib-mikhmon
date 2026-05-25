<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import SearchInput from '@/components/ui/SearchInput.vue'
import Badge from '@/components/ui/Badge.vue'
import DataTable from '@/components/ui/DataTable.vue'
import Icon from '@/components/ui/Icon.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useHotspotCookiesQuery } from '@/queries/hotspot.queries'
import { hotspotSessionsService } from '@/services/hotspot-sessions'
import { useToast } from '@/composables/useToast'
import type { HotspotCookie } from '@/types/hotspot'

const toast = useToast()
const { activeDeviceId } = useActiveDevice()

const search = ref('')

const { data: apiCookies, isLoading: loadingCookies, refetch: refetchCookies } = useHotspotCookiesQuery(activeDeviceId)

const filtered = computed(() => {
  const cookies = apiCookies.value ?? []
  if (!search.value) return cookies
  const s = search.value.toLowerCase()
  return cookies.filter(
    (c) =>
      (c.user || '').toLowerCase().includes(s) ||
      (c.mac_address || '').toLowerCase().includes(s) ||
      (c.domain || '').toLowerCase().includes(s),
  )
})

async function remove(id: string, user: string) {
  if (!activeDeviceId.value) return
  try {
    await hotspotSessionsService.removeCookie(activeDeviceId.value, id)
    toast.success(`Cookie untuk ${user || id} berhasil dihapus`)
    refetchCookies()
  } catch (err) {
    toast.error(`Gagal menghapus cookie: ${err instanceof Error ? err.message : String(err)}`)
  }
}

const columns = computed<ColumnDef<HotspotCookie>[]>(() => [
  {
    accessorKey: 'user',
    header: 'User',
    cell: ({ row }) => h('span', { class: 'text-[13px] font-medium' }, row.original.user || '—'),
  },
  {
    accessorKey: 'domain',
    header: 'Domain',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.domain || '—'),
    meta: { mobileHidden: true },
  },
  {
    accessorKey: 'expires_in',
    header: 'Expire In',
    cell: ({ row }) => h(Badge, { tone: 'warn' }, () => row.original.expires_in || '—'),
  },
  {
    accessorKey: 'mac_address',
    header: 'MAC Address',
    cell: ({ row }) => h('span', { class: 'mono text-[12px]' }, row.original.mac_address || '—'),
    meta: { mobileHidden: true },
  },
  {
    id: '__actions',
    header: '',
    enableSorting: false,
    cell: ({ row }) =>
      h('div', { class: 'flex items-center justify-end' }, [
        h(
          'button',
          {
            class: 'btn btn-ghost btn-icon btn-xs',
            title: 'Hapus Cookie',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              remove(row.original.id, row.original.user || '')
            },
          },
          [h(Icon, { name: 'Trash', size: 13, style: 'color: var(--danger)' })],
        ),
      ]),
    meta: { align: 'right' },
  },
])
</script>

<template>
  <div>
    <div v-if="loadingCookies" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading cookies...</div>
    </div>

    <DataTable
      v-else
      :columns="columns"
      :data="filtered"
      :get-row-id="(c) => c.id"
      :global-filter="search"
      :page-size="10"
      empty-message="Tidak ada cookies"
      @update:global-filter="(v) => (search = v)"
    >
      <template #toolbar>
        <SearchInput v-model="search" placeholder="Cari user, MAC, domain..." />
      </template>
    </DataTable>
  </div>
</template>
