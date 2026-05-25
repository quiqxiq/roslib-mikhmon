<script setup lang="ts">
import Drawer from '@/components/ui/Drawer.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import StatusDot from '@/components/ui/StatusDot.vue'
import type { FixtureDevice } from '@/fixtures/devices'

defineProps<{
  open: boolean
  device: FixtureDevice | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'edit'): void
  (e: 'delete'): void
  (e: 'set-active'): void
}>()
</script>

<template>
  <Drawer :open="open" :title="device?.name ?? 'Device'" @close="emit('close')">
    <div v-if="device" class="space-y-4">
      <div class="flex items-center gap-3">
        <div
          class="flex h-12 w-12 items-center justify-center rounded-lg"
          style="background: var(--accent-cyan-soft); color: var(--accent-cyan)"
        >
          <Icon name="Server" :size="24" />
        </div>
        <div>
          <div class="text-base font-semibold">{{ device.name }}</div>
          <div class="mono text-xs" style="color: var(--muted)">{{ device.id }}</div>
        </div>
        <div class="ml-auto">
          <StatusDot :status="device.status" />
        </div>
      </div>

      <div class="grid grid-cols-2 gap-2.5 text-xs">
        <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
          <div style="color: var(--muted)">Address</div>
          <div class="mono mt-0.5 font-medium">{{ device.address }}</div>
        </div>
        <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
          <div style="color: var(--muted)">Uptime</div>
          <div class="mono mt-0.5 font-medium">{{ device.uptime }}</div>
        </div>
        <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
          <div style="color: var(--muted)">Users</div>
          <div class="mono mt-0.5 font-medium">{{ device.users }}</div>
        </div>
        <div class="rounded-lg p-2.5" style="background: var(--bg-2)">
          <div style="color: var(--muted)">Active</div>
          <div class="mono mt-0.5 font-medium" style="color: var(--accent-cyan)">{{ device.active }}</div>
        </div>
        <div class="col-span-2 rounded-lg p-2.5" style="background: var(--bg-2)">
          <div style="color: var(--muted)">Version</div>
          <div class="mt-0.5 font-medium">
            <Badge tone="neutral">{{ device.version }}</Badge>
          </div>
        </div>
      </div>

      <div class="space-y-2">
        <button class="btn btn-primary w-full" type="button" @click="emit('set-active')">
          <Icon name="Check" :size="14" />
          Jadikan device aktif
        </button>
        <button class="btn w-full" type="button" @click="emit('edit')">
          <Icon name="Edit" :size="14" />
          Edit device
        </button>
      </div>
    </div>
    <template #footer>
      <button class="btn btn-danger btn-sm" type="button" @click="emit('delete')">
        <Icon name="Trash" :size="13" />
        Hapus device
      </button>
      <div class="flex-1" />
      <button class="btn btn-sm" type="button" @click="emit('close')">Tutup</button>
    </template>
  </Drawer>
</template>
