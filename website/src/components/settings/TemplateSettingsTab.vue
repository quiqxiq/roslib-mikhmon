<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/ui/Icon.vue'
import { useTweaksStore, type VoucherTemplate } from '@/stores/tweaks'
import VoucherTemplateDefault from '@/components/hotspot/voucher-templates/VoucherTemplateDefault.vue'
import VoucherTemplateSmall from '@/components/hotspot/voucher-templates/VoucherTemplateSmall.vue'
import VoucherTemplateThermal from '@/components/hotspot/voucher-templates/VoucherTemplateThermal.vue'

const store = useTweaksStore()
const selected = ref<VoucherTemplate>(store.defaultVoucherTemplate)

const templates: { id: VoucherTemplate; name: string; desc: string; component: unknown }[] = [
  {
    id: 'default',
    name: 'Default',
    desc: '220px — Full layout dengan logo, QR, dan login URL.',
    component: VoucherTemplateDefault,
  },
  {
    id: 'small',
    name: 'Small',
    desc: '160px — Compact, no logo, minimal label size.',
    component: VoucherTemplateSmall,
  },
  {
    id: 'thermal',
    name: 'Thermal',
    desc: '180px — Optimized untuk thermal receipt paper.',
    component: VoucherTemplateThermal,
  },
]

const sample = {
  hotspotName: 'Mikhmon Hotspot',
  num: 1,
  userMode: 'vc' as const,
  username: 'ABCD1234',
  password: '1234',
  validity: '1d',
  timeLimit: '',
  dataLimit: '',
  price: 'Rp 10.000',
  dnsName: 'hotspot.local',
}

function save() {
  store.setDefaultVoucherTemplate(selected.value)
}
</script>

<template>
  <div class="space-y-6">
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="t in templates"
        :key="t.id"
        class="cursor-pointer rounded-xl border-2 p-4 transition-all"
        :class="selected === t.id ? 'border-[var(--accent-cyan)]' : 'border-transparent'"
        style="background: var(--bg-2)"
        @click="selected = t.id"
      >
        <div class="mb-2 flex items-center gap-2">
          <div
            class="flex h-5 w-5 items-center justify-center rounded-full border"
            :class="selected === t.id ? 'border-[var(--accent-cyan)]' : 'border-[var(--border)]'"
          >
            <div
              v-if="selected === t.id"
              class="h-2.5 w-2.5 rounded-full"
              style="background: var(--accent-cyan)"
            />
          </div>
          <span class="text-sm font-semibold">{{ t.name }}</span>
        </div>
        <div class="text-xs" style="color: var(--muted)">{{ t.desc }}</div>
      </div>
    </div>

    <div class="rounded-xl border p-4" style="border-color: var(--border)">
      <div class="mb-3 text-sm font-medium">Preview</div>
      <div class="flex flex-wrap gap-2">
        <component
          :is="templates.find((t) => t.id === selected)?.component ?? VoucherTemplateDefault"
          v-bind="sample"
        />
      </div>
    </div>

    <div class="flex justify-end">
      <button class="btn btn-primary btn-sm" type="button" @click="save">
        <Icon name="Check" :size="13" />
        Simpan Default
      </button>
    </div>
  </div>
</template>
