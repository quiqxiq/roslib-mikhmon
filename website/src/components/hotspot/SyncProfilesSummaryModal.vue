<script setup lang="ts">
import { computed } from 'vue'
import Modal from '@/components/ui/Modal.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import type { ProfileConfigSyncResponse } from '@/types/profile-config'

// Modal yang menampilkan ringkasan hasil POST /hotspot/profile-configs/sync.
//
// Lima kategori (semua array of profile name string):
//   - synced        : profile yang sudah ada di DB (price/mode dipertahankan)
//   - created       : profile router baru yang di-insert ke DB (default
//                     expiry_mode="0", price=0 — operator perlu PUT untuk set)
//   - orphan        : config DB yang profile-nya hilang dari router
//   - injected      : on-login script berhasil di-push ke router
//   - inject_failed : "<profile>: <error>" — non-fatal
const props = defineProps<{
  open: boolean
  result: ProfileConfigSyncResponse | null
}>()

defineEmits<{
  (e: 'close'): void
}>()

interface Group {
  key: keyof ProfileConfigSyncResponse
  label: string
  description: string
  icon: 'Check' | 'Plus' | 'AlertCircle' | 'Zap' | 'X'
  tone: 'success' | 'cyan' | 'warn' | 'violet' | 'danger'
}

const GROUPS: Group[] = [
  {
    key: 'synced',
    label: 'Synced',
    description: 'Profile yang sudah ada di DB — value operator dipertahankan.',
    icon: 'Check',
    tone: 'success',
  },
  {
    key: 'created',
    label: 'Created',
    description:
      'Profile router baru — di-insert ke DB dengan default mode "Free". Operator perlu set price/mode via PUT.',
    icon: 'Plus',
    tone: 'cyan',
  },
  {
    key: 'orphan',
    label: 'Orphan',
    description:
      'Config DB yang profile-nya sudah hilang dari router (mungkin di-delete via WinBox/CLI). Tidak dihapus otomatis.',
    icon: 'AlertCircle',
    tone: 'warn',
  },
  {
    key: 'injected',
    label: 'Script Injected',
    description: 'On-login script berhasil di-push ke router.',
    icon: 'Zap',
    tone: 'violet',
  },
  {
    key: 'inject_failed',
    label: 'Inject Failed',
    description: 'Profile yang gagal di-inject (router error). Non-fatal — bisa retry manual.',
    icon: 'X',
    tone: 'danger',
  },
]

const totals = computed(() => {
  if (!props.result) {
    return { synced: 0, created: 0, orphan: 0, injected: 0, inject_failed: 0 }
  }
  return {
    synced: props.result.synced.length,
    created: props.result.created.length,
    orphan: props.result.orphan.length,
    injected: props.result.injected.length,
    inject_failed: props.result.inject_failed.length,
  }
})

const totalProfiles = computed(() => totals.value.synced + totals.value.created)
</script>

<template>
  <Modal :open="open" title="Hasil Sync Profile" @close="$emit('close')">
    <div v-if="result" class="space-y-4">
      <!-- KPI bar ringkas -->
      <div
        class="grid grid-cols-2 gap-2 rounded-lg p-3 sm:grid-cols-5"
        style="background: var(--bg-2); border: 1px solid var(--border)"
      >
        <div v-for="g in GROUPS" :key="g.key" class="text-center">
          <div class="text-[11px] uppercase" style="color: var(--muted); letter-spacing: 0.05em">
            {{ g.label }}
          </div>
          <div class="mono mt-0.5 text-xl font-semibold tabular">
            {{ totals[g.key] }}
          </div>
        </div>
      </div>

      <p class="text-xs" style="color: var(--muted)">
        Total <span class="font-semibold" style="color: var(--text)">{{ totalProfiles }}</span>
        profile router ter-sync ke DB. Profile baru dibuat dengan mode <Badge tone="neutral">Free</Badge> —
        set harga & mode lewat tombol Edit di kartu profile.
      </p>

      <!-- Per-kategori list -->
      <section v-for="g in GROUPS" :key="g.key">
        <div class="mb-2 flex items-center gap-2">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-md"
            :style="{
              background: `var(--accent-${g.tone === 'success' ? 'lime' : g.tone === 'warn' ? 'amber' : g.tone === 'danger' ? 'red' : g.tone}-soft, var(--bg-2))`,
            }"
          >
            <Icon :name="g.icon" :size="13" />
          </span>
          <h3 class="text-sm font-semibold">
            {{ g.label }}
            <span class="mono ml-1 text-xs" style="color: var(--muted)">
              ({{ totals[g.key] }})
            </span>
          </h3>
        </div>
        <p class="mb-2 text-[11px]" style="color: var(--muted)">{{ g.description }}</p>
        <div v-if="result[g.key].length" class="flex flex-wrap gap-1.5">
          <span
            v-for="item in result[g.key]"
            :key="item"
            class="mono rounded-md px-2 py-1 text-[11px]"
            :style="{
              background: 'var(--bg-2)',
              border: '1px solid var(--border)',
              color: g.tone === 'danger' ? 'var(--danger)' : 'var(--text-2)',
            }"
          >
            {{ item }}
          </span>
        </div>
        <div v-else class="text-[11px] italic" style="color: var(--muted-2)">— kosong —</div>
      </section>
    </div>
    <div v-else class="text-sm" style="color: var(--muted)">Tidak ada data sync.</div>

    <template #footer>
      <button class="btn btn-primary btn-sm" type="button" @click="$emit('close')">Tutup</button>
    </template>
  </Modal>
</template>
