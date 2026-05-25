<script setup lang="ts">
import { computed } from 'vue'
import Modal from '@/components/ui/Modal.vue'
import Icon from '@/components/ui/Icon.vue'

const props = defineProps<{
  open: boolean
  routerProfileNames: string[]
  dbConfigNames: string[]
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'confirm'): void
}>()

const toKeep = computed(() =>
  props.routerProfileNames.filter((n) => props.dbConfigNames.includes(n)),
)
const toCreate = computed(() =>
  props.routerProfileNames.filter((n) => !props.dbConfigNames.includes(n)),
)
const orphans = computed(() =>
  props.dbConfigNames.filter((n) => !props.routerProfileNames.includes(n)),
)
</script>

<template>
  <Modal :open="open" title="Sinkronkan Profile Hotspot?" @close="$emit('close')">
    <div class="space-y-4">
      <!-- Stats row -->
      <div
        class="grid grid-cols-2 gap-3 rounded-lg p-3"
        style="background: var(--bg-2); border: 1px solid var(--border)"
      >
        <div class="text-center">
          <div class="text-[11px] uppercase" style="color: var(--muted); letter-spacing: 0.05em">
            Di Router
          </div>
          <div class="mono mt-0.5 text-2xl font-bold" style="color: var(--accent-cyan)">
            {{ routerProfileNames.length }}
          </div>
          <div class="text-[11px]" style="color: var(--muted)">profile</div>
        </div>
        <div class="text-center">
          <div class="text-[11px] uppercase" style="color: var(--muted); letter-spacing: 0.05em">
            Di Database
          </div>
          <div class="mono mt-0.5 text-2xl font-bold" style="color: var(--text)">
            {{ dbConfigNames.length }}
          </div>
          <div class="text-[11px]" style="color: var(--muted)">config</div>
        </div>
      </div>

      <!-- Akan dipertahankan -->
      <div class="rounded-lg p-3" style="background: var(--bg-2); border: 1px solid var(--border)">
        <div class="mb-2 flex items-center gap-2">
          <span
            class="flex h-5 w-5 items-center justify-center rounded-md"
            style="background: rgba(16, 185, 129, 0.15)"
          >
            <Icon name="Check" :size="11" style="color: var(--success)" />
          </span>
          <span class="text-sm font-semibold">{{ toKeep.length }} Dipertahankan</span>
          <span class="text-xs" style="color: var(--muted)">— price & mode tidak ditimpa</span>
        </div>
        <div v-if="toKeep.length" class="flex flex-wrap gap-1">
          <span
            v-for="n in toKeep"
            :key="n"
            class="mono rounded px-1.5 py-0.5 text-[11px]"
            style="background: var(--bg-1); border: 1px solid var(--border); color: var(--text-2)"
          >{{ n }}</span>
        </div>
        <div v-else class="text-[11px] italic" style="color: var(--muted-2)">— tidak ada —</div>
      </div>

      <!-- Akan dibuat baru -->
      <div class="rounded-lg p-3" style="background: var(--bg-2); border: 1px solid var(--border)">
        <div class="mb-2 flex items-center gap-2">
          <span
            class="flex h-5 w-5 items-center justify-center rounded-md"
            style="background: rgba(6, 182, 212, 0.15)"
          >
            <Icon name="Plus" :size="11" style="color: var(--accent-cyan)" />
          </span>
          <span class="text-sm font-semibold">{{ toCreate.length }} Dibuat Baru</span>
          <span class="text-xs" style="color: var(--muted)">— default Free, harga 0</span>
        </div>
        <div v-if="toCreate.length" class="flex flex-wrap gap-1">
          <span
            v-for="n in toCreate"
            :key="n"
            class="mono rounded px-1.5 py-0.5 text-[11px]"
            style="background: rgba(6, 182, 212, 0.08); border: 1px solid rgba(6, 182, 212, 0.3); color: var(--accent-cyan)"
          >{{ n }}</span>
        </div>
        <div v-else class="text-[11px] italic" style="color: var(--muted-2)">— tidak ada yang baru —</div>
      </div>

      <!-- Orphan -->
      <div
        v-if="orphans.length"
        class="rounded-lg p-3"
        style="background: var(--bg-2); border: 1px solid var(--border)"
      >
        <div class="mb-2 flex items-center gap-2">
          <span
            class="flex h-5 w-5 items-center justify-center rounded-md"
            style="background: rgba(245, 158, 11, 0.15)"
          >
            <Icon name="AlertCircle" :size="11" style="color: var(--warning, #f59e0b)" />
          </span>
          <span class="text-sm font-semibold">{{ orphans.length }} Orphan</span>
          <span class="text-xs" style="color: var(--muted)">— ada di DB, tidak ada di router</span>
        </div>
        <div class="flex flex-wrap gap-1">
          <span
            v-for="n in orphans"
            :key="n"
            class="mono rounded px-1.5 py-0.5 text-[11px]"
            style="background: rgba(245, 158, 11, 0.08); border: 1px solid rgba(245, 158, 11, 0.3); color: #f59e0b"
          >{{ n }}</span>
        </div>
      </div>

      <!-- Warning script inject -->
      <div
        class="flex items-start gap-2 rounded-lg p-3"
        style="background: rgba(139, 92, 246, 0.08); border: 1px solid rgba(139, 92, 246, 0.25)"
      >
        <Icon name="Zap" :size="14" class="mt-0.5 shrink-0" style="color: var(--accent-violet)" />
        <p class="text-xs leading-relaxed" style="color: var(--text-2)">
          On-login script akan di-push ke
          <strong>semua {{ routerProfileNames.length }} profile</strong> di router.
          Perubahan manual pada script di router akan <strong>ditimpa</strong>.
        </p>
      </div>
    </div>

    <template #footer>
      <button class="btn btn-sm" type="button" @click="$emit('close')">Batal</button>
      <button class="btn btn-primary btn-sm" type="button" @click="$emit('confirm')">
        <Icon name="Refresh" :size="13" />
        Sync Sekarang
      </button>
    </template>
  </Modal>
</template>
