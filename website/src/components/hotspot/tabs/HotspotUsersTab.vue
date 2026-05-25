<script setup lang="ts">
import { computed, h, ref } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import SearchInput from '@/components/ui/SearchInput.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import Avatar from '@/components/ui/Avatar.vue'
import Select from '@/components/ui/Select.vue'
import DataTable from '@/components/ui/DataTable.vue'
import SummaryChip from '@/components/hotspot/SummaryChip.vue'
import HotspotUserDrawer from '@/components/hotspot/HotspotUserDrawer.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { useToast } from '@/composables/useToast'
import { fmtBytes } from '@/utils/fmt'
import { downloadCsv, todayStamp } from '@/utils/export'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useHotspotUsersQuery, useHotspotActiveQuery, useHotspotProfilesQuery } from '@/queries/hotspot.queries'
import { useProfileConfigsQuery } from '@/queries/profile-config.queries'
import { hotspotUsersService } from '@/services/hotspot-users'
import VoucherPrintModal from '@/components/hotspot/VoucherPrintModal.vue'
import type { HotspotUser } from '@/types/hotspot'

const toast = useToast()
const { activeDeviceId } = useActiveDevice()

// Query hooks
const { data: apiUsers, refetch: refetchUsers, isLoading: loadingUsers } = useHotspotUsersQuery(activeDeviceId)
const { data: apiActive, refetch: refetchActive } = useHotspotActiveQuery(activeDeviceId)
const { data: apiProfiles } = useHotspotProfilesQuery(activeDeviceId)
const { data: profileConfigs } = useProfileConfigsQuery(activeDeviceId)

const confirmBulkOpen = ref(false)
const printModalOpen = ref(false)
const search = ref('')
const filterProfile = ref<string>('all')
const filterServer = ref<string>('all')
const filterStatus = ref<string>('all')
const drawerOpen = ref(false)
const editingUser = ref<Partial<HotspotUser> | null>(null)
const selectedIds = ref<string[]>([])

// Hitung user yang online dari list active session
const activeUserNames = computed(() => {
  return new Set((apiActive.value ?? []).map((s) => s.user))
})

const usersList = computed<HotspotUser[]>(() => apiUsers.value ?? [])

const filtered = computed(() => {
  return usersList.value.filter((u) => {
    const isOnline = activeUserNames.value.has(u.name)
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
    if (filterServer.value !== 'all' && u.server !== filterServer.value) return false
    if (filterStatus.value === 'active' && !isOnline) return false
    if (filterStatus.value === 'inactive' && isOnline) return false
    if (filterStatus.value === 'disabled' && !u.disabled) return false
    return true
  })
})

const summary = computed(() => {
  const list = usersList.value
  return {
    total: list.length,
    active: list.filter((u) => activeUserNames.value.has(u.name)).length,
    disabled: list.filter((u) => u.disabled).length,
  }
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
      const isOnline = activeUserNames.value.has(u.name)
      if (u.disabled) return h(Badge, { tone: 'neutral' }, () => 'Disabled')
      return h(Badge, { tone: isOnline ? 'success' : 'neutral', dot: true }, () =>
        isOnline ? 'Online' : 'Offline',
      )
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
  {
    id: '__actions',
    header: '',
    enableSorting: false,
    cell: ({ row }) =>
      h('div', { class: 'flex items-center justify-end gap-1' }, [
        h(
          'button',
          {
            class: 'btn btn-ghost btn-icon btn-xs',
            title: row.original.disabled ? 'Enable' : 'Disable',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              togglePower(row.original.id, row.original.name, row.original.disabled)
            },
          },
          [h(Icon, { name: 'Power', size: 13 })],
        ),
        h(
          'button',
          {
            class: 'btn btn-ghost btn-icon btn-xs',
            title: 'Edit',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              editingUser.value = row.original
              drawerOpen.value = true
            },
          },
          [h(Icon, { name: 'Edit', size: 13 })],
        ),
        h(
          'button',
          {
            class: 'btn btn-ghost btn-icon btn-xs',
            title: 'Reset counters',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              resetCounters(row.original.id, row.original.name)
            },
          },
          [h(Icon, { name: 'Refresh', size: 13 })],
        ),
      ]),
    meta: { align: 'right' },
  },
])

async function togglePower(id: string, name: string, currentlyDisabled: boolean) {
  if (!activeDeviceId.value) return
  try {
    await hotspotUsersService.setDisabled(activeDeviceId.value, id, !currentlyDisabled)
    toast.info(`${name}: ${!currentlyDisabled ? 'disabled' : 'enabled'}`)
    refetchUsers()
  } catch (err) {
    toast.error(`Gagal mengubah status user: ${err instanceof Error ? err.message : String(err)}`)
  }
}

async function resetCounters(id: string, name: string) {
  if (!activeDeviceId.value) return
  try {
    await hotspotUsersService.resetCounters(activeDeviceId.value, id)
    toast.success(`Counter untuk ${name} berhasil direset`)
    refetchUsers()
  } catch (err) {
    toast.error(`Gagal meriset counter: ${err instanceof Error ? err.message : String(err)}`)
  }
}

function exportCsv() {
  if (!filtered.value.length) {
    toast.warning('Tidak ada data untuk di-export')
    return
  }
  downloadCsv(
    filtered.value.map((u) => ({
      id: u.id,
      name: u.name,
      profile: u.profile,
      server: u.server || 'all',
      mac: u.mac_address ?? '',
      bytes_in: u.bytes_in,
      bytes_out: u.bytes_out,
      uptime: u.uptime ?? '—',
      comment: u.comment ?? '',
      disabled: u.disabled,
    })),
    `hotspot-users-${todayStamp()}.csv`,
  )
  toast.success(`${filtered.value.length} user di-export`)
}

function reload() {
  if (activeDeviceId.value) {
    refetchUsers()
    refetchActive()
    toast.info('Reload data dari router…')
  }
}

function askDeleteSelected() {
  if (!selectedIds.value.length) return
  confirmBulkOpen.value = true
}

async function deleteSelected() {
  if (!activeDeviceId.value) return
  const n = selectedIds.value.length
  try {
    await hotspotUsersService.bulkRemove(activeDeviceId.value, selectedIds.value)
    selectedIds.value = []
    confirmBulkOpen.value = false
    toast.warning(`${n} user dihapus`)
    refetchUsers()
  } catch (err) {
    toast.error(`Gagal menghapus user: ${err instanceof Error ? err.message : String(err)}`)
  }
}

function openCreate() {
  editingUser.value = {}
  drawerOpen.value = true
}

async function onSave(u: Partial<HotspotUser>) {
  if (!activeDeviceId.value) return
  try {
    if (u.id) {
      // Edit
      await hotspotUsersService.update(activeDeviceId.value, u.id, {
        name: u.name,
        profile: u.profile,
        server: u.server,
        disabled: u.disabled,
        comment: u.comment,
        mac_address: u.mac_address,
        limit_uptime: u.limit_uptime,
        limit_bytes_total: u.limit_bytes_total,
      })
      toast.success(`User ${u.name} disimpan`)
    } else {
      // Create new
      await hotspotUsersService.create(activeDeviceId.value, {
        name: u.name ?? '',
        password: u.comment ?? '', // standard mikhmon password/comment mapping
        profile: u.profile,
        server: u.server,
        disabled: u.disabled,
        limit_uptime: u.limit_uptime,
        limit_bytes_total: u.limit_bytes_total,
        comment: u.comment,
      })
      toast.success(`User ${u.name} ditambahkan`)
    }
    refetchUsers()
    drawerOpen.value = false
  } catch (err) {
    toast.error(`Gagal menyimpan user: ${err instanceof Error ? err.message : String(err)}`)
  }
}

async function onDelete(id: string) {
  if (!activeDeviceId.value) return
  try {
    await hotspotUsersService.remove(activeDeviceId.value, id)
    drawerOpen.value = false
    toast.warning(`User berhasil dihapus`)
    refetchUsers()
  } catch (err) {
    toast.error(`Gagal menghapus user: ${err instanceof Error ? err.message : String(err)}`)
  }
}
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <button class="btn btn-sm" type="button" @click="exportCsv">
        <Icon name="Download" :size="13" />
        Export CSV
      </button>
      <button class="btn btn-sm" type="button" @click="reload">
        <Icon name="Refresh" :size="13" />
        Reload
      </button>
      <button class="btn btn-sm" type="button" @click="printModalOpen = true">
        <Icon name="Printer" :size="13" />
        Print
      </button>
      <div class="flex-1" />
      <button class="btn btn-primary btn-sm" type="button" @click="openCreate">
        <Icon name="Plus" :size="14" />
        Tambah User
      </button>
    </div>

    <div v-if="loadingUsers" class="mb-4 flex items-center justify-center p-8">
      <div class="text-sm" style="color: var(--muted)">Loading users...</div>
    </div>

    <div v-else class="mb-4 flex flex-wrap gap-2.5">
      <SummaryChip label="Total" :value="summary.total" accent="cyan" />
      <SummaryChip label="Aktif sekarang" :value="summary.active" accent="lime" />
      <SummaryChip label="Disabled" :value="summary.disabled" accent="violet" />
    </div>

    <DataTable
      v-if="!loadingUsers"
      :columns="columns"
      :data="filtered"
      :get-row-id="(u) => u.id"
      :global-filter="search"
      :page-size="9"
      enable-row-selection
      empty-message="Tidak ada user yang cocok dengan filter"
      @update:global-filter="(v) => (search = v)"
      @selection-change="(ids) => (selectedIds = ids)"
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
        <Select
          v-model="filterServer"
          sm
          :options="[
            { value: 'all', label: 'Semua server' },
            { value: 'all', label: 'all' },
            { value: 'hotspot1', label: 'hotspot1' },
            { value: 'hotspot2', label: 'hotspot2' },
          ]"
        />
        <Select
          v-model="filterStatus"
          sm
          :options="[
            { value: 'all', label: 'Semua status' },
            { value: 'active', label: 'Aktif' },
            { value: 'inactive', label: 'Tidak aktif' },
            { value: 'disabled', label: 'Disabled' },
          ]"
        />
      </template>

      <template #bulkBar="{ selectedCount, clear }">
        <span class="text-xs font-medium" style="color: var(--accent-cyan)">
          {{ selectedCount }} dipilih
        </span>
        <div class="flex items-center gap-2">
          <button class="btn btn-xs btn-ghost" type="button" @click="clear">Batal</button>
          <button class="btn btn-xs btn-danger" type="button" @click="askDeleteSelected">
            <Icon name="Trash" :size="11" />
            Hapus
          </button>
        </div>
      </template>
    </DataTable>

    <HotspotUserDrawer
      :open="drawerOpen"
      :user="editingUser"
      @close="drawerOpen = false"
      @save="onSave"
      @delete="onDelete"
    />

    <ConfirmModal
      :open="confirmBulkOpen"
      title="Hapus user terpilih?"
      :message="`${selectedIds.length} user akan dihapus dari sistem.`"
      confirm-text="Hapus"
      variant="danger"
      @close="confirmBulkOpen = false"
      @confirm="deleteSelected"
    />

    <VoucherPrintModal
      :open="printModalOpen"
      :profile-configs="profileConfigs ?? []"
      @close="printModalOpen = false"
    />
  </div>
</template>
