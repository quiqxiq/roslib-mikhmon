<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import Icon from '@/components/ui/Icon.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useDevicesQuery } from '@/queries/devices.queries'
import { useToast } from '@/composables/useToast'

const emit = defineEmits<{
  (e: 'select'): void
}>()

const toast = useToast()
const { activeDeviceId, setActiveDevice } = useActiveDevice()
const open = ref(false)
const triggerRef = ref<HTMLElement | null>(null)
const popoverStyle = ref({ left: '0px', top: '0px' })

// Fetch real devices
const { data: devicesList, isLoading } = useDevicesQuery()

const activeDevice = computed(() => {
  const list = devicesList.value ?? []
  if (!list.length) return null
  return list.find((d) => String(d.id) === String(activeDeviceId.value)) ?? list[0]
})

function toggle() {
  if (isLoading.value || !(devicesList.value?.length)) return
  open.value = !open.value
  if (open.value) {
    nextTick(() => {
      const rect = triggerRef.value?.getBoundingClientRect()
      if (rect) {
        popoverStyle.value = {
          left: `${rect.right + 8}px`,
          top: `${rect.top}px`,
        }
      }
    })
  }
}

function selectDevice(id: string) {
  setActiveDevice(id)
  toast.success('Device aktif diubah')
  open.value = false
  emit('select')
}

function statusColor(status: string | undefined): string {
  if (!status) return 'var(--muted-2)'
  const m: Record<string, string> = {
    connected: 'var(--success)',
    online: 'var(--success)',
    warn: 'var(--warning)',
    connecting: 'var(--warning)',
    offline: 'var(--muted-2)',
    disconnected: 'var(--muted-2)',
    danger: 'var(--danger)',
    error: 'var(--danger)',
  }
  return m[status] || 'var(--muted-2)'
}

onClickOutside(triggerRef, () => {
  if (open.value) open.value = false
})
</script>

<template>
  <div ref="triggerRef" class="relative">
    <!-- Trigger -->
    <button
      type="button"
      class="row-hover relative flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2.5 text-left transition-colors"
      :style="{
        background: 'var(--bg-2)',
        color: 'var(--text)',
        cursor: (isLoading || !devicesList?.length) ? 'default' : 'pointer',
        border: '1px solid var(--border)',
      }"
      :disabled="isLoading || !devicesList?.length"
      @click="toggle"
    >
      <span
        class="relative flex shrink-0 items-center justify-center"
        style="width: 28px; height: 28px; border-radius: 7px; background: var(--bg-3); border: 1px solid var(--border)"
      >
        <Icon name="Server" :size="14" style="color: var(--accent-cyan)" />
        <span
          class="absolute rounded-full"
          :style="{
            bottom: '-2px',
            right: '-2px',
            width: '9px',
            height: '9px',
            background: statusColor(activeDevice?.status),
            border: '2px solid var(--bg)',
            boxShadow: activeDevice?.status === 'connected' || activeDevice?.status === 'online' ? `0 0 6px ${statusColor(activeDevice?.status)}` : 'none',
          }"
        />
      </span>
      <span class="min-w-0 flex-1">
        <span class="block truncate text-[13px] font-medium">
          {{ isLoading ? 'Memuat...' : (activeDevice?.display_name || 'Pilih Router...') }}
        </span>
        <span class="mono block text-[10.5px]" style="color: var(--muted)">
          {{ activeDevice?.address || 'Belum ada router' }}
        </span>
      </span>
      <Icon
        v-if="!isLoading && devicesList?.length"
        name="ChevronDown"
        :size="12"
        :style="{
          color: 'var(--muted)',
          transform: open ? 'rotate(180deg)' : 'rotate(0deg)',
          transition: 'transform 200ms',
        }"
      />
    </button>

    <!-- Popover panel via Teleport -->
    <Teleport to="body">
      <div
        v-if="open && devicesList?.length"
        class="fixed z-[100] w-64 rounded-xl border p-2 shadow-lg"
        :style="{
          left: popoverStyle.left,
          top: popoverStyle.top,
          background: 'var(--bg)',
          borderColor: 'var(--border)',
        }"
      >
        <div
          class="mb-1 px-2 py-1 text-[10px] font-semibold uppercase"
          style="color: var(--muted); letter-spacing: 0.08em"
        >
          Daftar Router
        </div>
        <div class="flex flex-col gap-0.5">
          <button
            v-for="d in devicesList"
            :key="d.id"
            type="button"
            class="row-hover flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-left"
            :style="{
              background: String(activeDeviceId) === String(d.id) ? 'var(--bg-active)' : 'transparent',
              color: String(activeDeviceId) === String(d.id) ? 'var(--text)' : 'var(--text-2)',
              cursor: 'pointer',
              border: 0,
            }"
            @click="selectDevice(String(d.id))"
          >
            <span
              class="relative flex shrink-0 items-center justify-center"
              style="width: 28px; height: 28px; border-radius: 7px; background: var(--bg-2); border: 1px solid var(--border)"
              :style="{ background: String(activeDeviceId) === String(d.id) ? 'var(--bg-3)' : 'var(--bg-2)' }"
            >
              <Icon
                name="Server"
                :size="14"
                :style="{ color: String(activeDeviceId) === String(d.id) ? 'var(--accent-cyan)' : 'var(--muted)' }"
              />
              <span
                class="absolute rounded-full"
                :style="{
                  bottom: '-2px',
                  right: '-2px',
                  width: '9px',
                  height: '9px',
                  background: statusColor(d.status),
                  border: '2px solid var(--bg)',
                  boxShadow: d.status === 'connected' || d.status === 'online' ? `0 0 6px ${statusColor(d.status)}` : 'none',
                }"
              />
            </span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-[13px] font-medium">{{ d.display_name }}</span>
              <span class="mono block text-[10.5px]" style="color: var(--muted)">{{ d.address }}</span>
            </span>
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>
