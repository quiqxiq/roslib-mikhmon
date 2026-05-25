<script setup lang="ts">
import { computed } from 'vue'
import Icon from './Icon.vue'
import type { IconName } from './icons'

type Kind = 'info' | 'success' | 'warning' | 'error'

const props = defineProps<{
  kind: Kind
  message: string
}>()

defineEmits<{
  (e: 'dismiss'): void
}>()

const iconName = computed<IconName>(() =>
  props.kind === 'success'
    ? 'Check'
    : props.kind === 'warning'
      ? 'Sparkles'
      : props.kind === 'error'
        ? 'X'
        : 'Bell',
)

const tokens = computed(() => {
  const m = {
    info: { color: 'var(--accent-cyan)', bg: 'var(--accent-cyan-soft)', border: 'rgba(34,211,238,0.3)' },
    success: { color: 'var(--success)', bg: 'rgba(16,185,129,0.12)', border: 'rgba(16,185,129,0.3)' },
    warning: { color: 'var(--warning)', bg: 'rgba(245,158,11,0.12)', border: 'rgba(245,158,11,0.3)' },
    error: { color: 'var(--danger)', bg: 'rgba(244,63,94,0.12)', border: 'rgba(244,63,94,0.3)' },
  }
  return m[props.kind]
})
</script>

<template>
  <div
    class="slide-in-r flex items-start gap-3 rounded-lg px-3.5 py-3 text-sm shadow-lg"
    :style="{
      background: 'var(--bg-1)',
      borderLeft: `3px solid ${tokens.color}`,
      border: `1px solid var(--border-strong)`,
      borderLeftWidth: '3px',
      borderLeftColor: tokens.color,
      minWidth: '280px',
      maxWidth: '420px',
      boxShadow: 'var(--shadow-2)',
    }"
  >
    <div
      class="flex h-6 w-6 shrink-0 items-center justify-center rounded-md"
      :style="{ background: tokens.bg, color: tokens.color }"
    >
      <Icon :name="iconName" :size="14" />
    </div>
    <div class="flex-1 text-[13px]">{{ message }}</div>
    <button
      type="button"
      class="btn btn-ghost btn-icon btn-xs shrink-0"
      title="Tutup"
      @click="$emit('dismiss')"
    >
      <Icon name="X" :size="12" />
    </button>
  </div>
</template>
