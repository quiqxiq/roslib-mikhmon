<script setup lang="ts">
import { computed } from 'vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Badge from '@/components/ui/Badge.vue'
import StatusDot from '@/components/ui/StatusDot.vue'
import Card from '@/components/ui/Card.vue'
import { DEVICES, type FixtureDevice } from '@/fixtures/devices'
import { useActiveDevice } from '@/composables/useActiveDevice'

const { activeDeviceId, setActiveDevice } = useActiveDevice()

function go(d: FixtureDevice) {
  setActiveDevice(d.id)
}

const totals = computed(() => ({
  devices: DEVICES.length,
  online: DEVICES.filter((d) => d.status === 'online').length,
  users: DEVICES.reduce((a, d) => a + d.users, 0),
  active: DEVICES.reduce((a, d) => a + d.active, 0),
}))
</script>

<template>
  <div class="fade-in">
    <PageHeader title="Devices" subtitle="Daftar router yang dikelola">
      <template #right>
        <button class="btn btn-sm" type="button">
          <Icon name="Refresh" :size="13" />
          Reload
        </button>
        <button class="btn btn-primary btn-sm" type="button">
          <Icon name="Plus" :size="14" />
          Tambah Device
        </button>
      </template>
    </PageHeader>

    <div class="mb-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <Card
        ><div class="text-xs" style="color: var(--muted)">Total</div>
        <div class="mt-1 text-2xl font-semibold tabular">{{ totals.devices }}</div></Card
      >
      <Card
        ><div class="text-xs" style="color: var(--muted)">Online</div>
        <div class="mt-1 text-2xl font-semibold tabular" style="color: var(--success)">
          {{ totals.online }}
        </div></Card
      >
      <Card
        ><div class="text-xs" style="color: var(--muted)">Total Users</div>
        <div class="mt-1 text-2xl font-semibold tabular">
          {{ totals.users.toLocaleString('id-ID') }}
        </div></Card
      >
      <Card
        ><div class="text-xs" style="color: var(--muted)">Active Sessions</div>
        <div class="mt-1 text-2xl font-semibold tabular" style="color: var(--accent-cyan)">
          {{ totals.active }}
        </div></Card
      >
    </div>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <button
        v-for="d in DEVICES"
        :key="d.id"
        type="button"
        class="card flex w-full flex-col items-stretch text-left transition-transform hover:-translate-y-0.5"
        :style="{
          borderTop: `2px solid ${activeDeviceId === d.id ? 'var(--accent-cyan)' : 'transparent'}`,
        }"
        @click="go(d)"
      >
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-2.5">
            <div
              class="relative flex h-10 w-10 items-center justify-center rounded-lg"
              :style="{
                background: activeDeviceId === d.id ? 'var(--accent-cyan-soft)' : 'var(--bg-2)',
              }"
            >
              <Icon
                name="Server"
                :size="18"
                :style="{ color: activeDeviceId === d.id ? 'var(--accent-cyan)' : 'var(--muted)' }"
              />
            </div>
            <div>
              <div class="text-sm font-semibold">{{ d.slug }}</div>
              <div class="mono text-[11px]" style="color: var(--muted)">{{ d.address }}</div>
            </div>
          </div>
          <StatusDot :status="d.status" />
        </div>
        <div class="mt-3 text-[11px]" style="color: var(--muted)">{{ d.name }}</div>
        <div class="mt-3 grid grid-cols-3 gap-2 text-xs">
          <div>
            <div style="color: var(--muted)">Users</div>
            <div class="mono mt-0.5 font-semibold">{{ d.users }}</div>
          </div>
          <div>
            <div style="color: var(--muted)">Active</div>
            <div class="mono mt-0.5 font-semibold" style="color: var(--accent-cyan)">
              {{ d.active }}
            </div>
          </div>
          <div>
            <div style="color: var(--muted)">Uptime</div>
            <div class="mono mt-0.5 font-semibold">{{ d.uptime }}</div>
          </div>
        </div>
        <div
          class="mt-3 flex items-center justify-between border-t pt-3"
          style="border-color: var(--border)"
        >
          <Badge tone="neutral">{{ d.version }}</Badge>
          <Icon name="Chevron" :size="14" :style="{ color: 'var(--muted)' }" />
        </div>
      </button>
    </div>
  </div>
</template>
