<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import { fmtRpShort } from '@/utils/fmt'
import type { ProfileViewModel } from './profile-vm'
import { expiryModeLabel } from './profile-vm'

const props = defineProps<{
  profile: ProfileViewModel
  maxSold: number
}>()

const tone = computed(
  () =>
    ({
      cyan: { color: 'var(--accent-cyan)', bg: 'var(--accent-cyan-soft)' },
      violet: { color: 'var(--accent-violet)', bg: 'var(--accent-violet-soft)' },
      lime: { color: 'var(--accent-lime)', bg: 'var(--accent-lime-soft)' },
    })[props.profile.color],
)

const pct = computed(() => (props.profile.sold / Math.max(1, props.maxSold)) * 100)
const isFree = computed(() => props.profile.expiry_mode === '0')
</script>

<template>
  <button
    type="button"
    class="card flex w-full flex-col items-stretch text-left transition-transform hover:-translate-y-0.5"
    :style="{ borderTop: `2px solid ${tone.color}` }"
  >
    <div class="mb-2 flex items-center justify-between">
      <div
        class="flex h-9 w-9 items-center justify-center rounded-lg"
        :style="{ background: tone.bg, color: tone.color }"
      >
        <Icon name="Ticket" :size="18" />
      </div>
      <div class="flex items-center gap-1">
        <Icon
          v-if="profile.lock_mac"
          name="Lock"
          :size="12"
          :title="`Lock to MAC`"
          :style="{ color: 'var(--muted)' }"
        />
        <Badge :tone="profile.color">{{ profile.validity || '—' }}</Badge>
      </div>
    </div>
    <div class="text-lg font-bold">{{ profile.name }}</div>
    <div
      class="mt-1 text-xl font-bold tabular"
      :style="{ color: tone.color, letterSpacing: '-0.02em' }"
    >
      <span v-if="isFree" class="text-base">Free</span>
      <span v-else>{{ fmtRpShort(profile.price) }}</span>
    </div>
    <div class="mt-1 text-[10.5px]" style="color: var(--muted)">
      Mode: <span class="font-medium" style="color: var(--text-2)">{{ expiryModeLabel(profile.expiry_mode) }}</span>
    </div>
    <div class="mt-2 flex flex-wrap items-center gap-2 text-xs" style="color: var(--muted)">
      <span class="flex items-center gap-1">
        <Icon name="Activity" :size="12" />
        <span class="mono">{{ profile.speed }}</span>
      </span>
      <span>·</span>
      <span class="flex items-center gap-1">
        <Icon name="Clock" :size="12" />
        <span class="mono">{{ profile.validity || 'no expiry' }}</span>
      </span>
    </div>
    <div class="mt-3 flex items-center justify-between text-xs" style="color: var(--muted)">
      <span>Terjual bulan ini</span>
      <span class="mono font-semibold" style="color: var(--text-2)">{{ profile.sold }}×</span>
    </div>
    <div class="bar mt-1.5">
      <i
        :style="{
          width: `${pct}%`,
          background: `linear-gradient(90deg, ${tone.color}, ${tone.color}aa)`,
        }"
      />
    </div>
  </button>
</template>
