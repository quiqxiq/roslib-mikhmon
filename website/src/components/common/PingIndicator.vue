<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Spark from '@/components/ui/Spark.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useDeviceQuery } from '@/queries/devices.queries'
import { useSSE } from '@/composables/useSSE'
import { buildStreamUrl } from '@/services/stream'
import type { PingStreamEvent } from '@/types/stream'

const props = withDefaults(
  defineProps<{
    baseMs?: number
  }>(),
  { baseMs: 12 },
)

const { activeDeviceId } = useActiveDevice()
const { data: device } = useDeviceQuery(activeDeviceId.value ?? '')

// Target: IP router tanpa port (default 8.8.8.8 kalau tidak ada)
const targetAddress = computed(() => {
  const addr = device.value?.address
  if (!addr) return '8.8.8.8'
  // Hapus port (mis. "192.168.88.1:8728" → "192.168.88.1")
  return addr.split(':')[0]
})

const sseUrl = computed(() => {
  if (!activeDeviceId.value) return null
  const url = buildStreamUrl(activeDeviceId.value, 'ping')
  const u = new URL(url)
  u.searchParams.set('address', targetAddress.value)
  return u.toString()
})

const { parsed: pingEvent, status: connStatus } = useSSE<PingStreamEvent>(sseUrl, ['ping'])

const ping = ref(props.baseMs)
const history = ref<number[]>(Array.from({ length: 12 }, () => props.baseMs))

watch(pingEvent, (ev) => {
  if (!ev?.time_ms) return
  const ms = Math.round(ev.time_ms)
  history.value = [...history.value.slice(1), ms]
  ping.value = ms
})

// Kalau SSE disconnect → set ping ke 999
watch(
  () => connStatus.value,
  (s) => {
    if (s !== 'OPEN') {
      ping.value = 999
    }
  },
  { immediate: true },
)

const status = computed(() => (ping.value < 30 ? 'good' : ping.value < 80 ? 'ok' : 'bad'))
const color = computed(() =>
  status.value === 'good'
    ? 'var(--success)'
    : status.value === 'ok'
      ? 'var(--warning)'
      : 'var(--danger)',
)
const label = computed(() =>
  status.value === 'good' ? 'Stabil' : status.value === 'ok' ? 'Lambat' : 'Tinggi',
)
</script>

<template>
  <div
    class="inline-flex items-center gap-2 rounded-full"
    :title="`Round-trip ke router · ${label}`"
    :style="{
      padding: '4px 10px 4px 8px',
      background: 'var(--bg-2)',
      border: '1px solid var(--border)',
      height: '30px',
    }"
  >
    <span class="dot dot-live" :style="{ background: color }" />
    <Spark :data="history" :color="color" kind="line" :width="36" :height="14" />
    <span
      class="mono tabular text-xs font-semibold"
      :style="{ color: 'var(--text)', minWidth: '32px', textAlign: 'right' }"
    >
      {{ ping }}<span style="color: var(--muted); margin-left: 2px; font-size: 10px; font-weight: 500">ms</span>
    </span>
  </div>
</template>
