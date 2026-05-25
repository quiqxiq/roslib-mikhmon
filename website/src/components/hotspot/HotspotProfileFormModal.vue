<script setup lang="ts">
import { ref, watch } from 'vue'
import Modal from '@/components/ui/Modal.vue'
import Field from '@/components/ui/Field.vue'
import Input from '@/components/ui/Input.vue'
import NumberInput from '@/components/ui/NumberInput.vue'
import Select from '@/components/ui/Select.vue'
import Toggle from '@/components/ui/Toggle.vue'
import type { ExpiryMode } from '@/types/profile-config'
import { EXPIRY_MODE_OPTIONS } from './profile-vm'

// Bentuk emit saat operator submit form. Parent component bertanggung
// jawab call hotspotProfilesService.create + profileConfigService.upsert
// secara berurutan.
export interface ProfileFormPayload {
  name: string
  rate_limit: string
  expiry_mode: ExpiryMode
  validity: string
  price: number
  sell_price: number
  lock_mac: boolean
}

const props = defineProps<{
  open: boolean
  // Kalau initial null → mode create. Kalau ada → mode edit (name di-disable).
  initial?: ProfileFormPayload | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', p: ProfileFormPayload): void
}>()

const DEFAULT: ProfileFormPayload = {
  name: '',
  rate_limit: '5M/2M',
  expiry_mode: '0',
  validity: '1d',
  price: 0,
  sell_price: 0,
  lock_mac: false,
}

const form = ref<ProfileFormPayload>({ ...DEFAULT })

watch(
  () => props.initial,
  (v) => {
    form.value = v ? { ...v } : { ...DEFAULT }
  },
  { immediate: true },
)

const isEdit = () => Boolean(props.initial)

function submit() {
  if (!form.value.name.trim()) return
  // Validasi: mode dengan record butuh validity & price.
  if (
    (form.value.expiry_mode === 'remc' || form.value.expiry_mode === 'ntfc') &&
    !form.value.validity.trim()
  ) {
    return
  }
  emit('save', { ...form.value })
}
</script>

<template>
  <Modal
    :open="open"
    :title="isEdit() ? 'Edit Profile' : 'Tambah Hotspot Profile'"
    @close="emit('close')"
  >
    <div class="space-y-4">
      <section>
        <h3
          class="mb-2 text-[10px] font-semibold uppercase"
          style="color: var(--muted); letter-spacing: 0.08em"
        >
          Router Profile
        </h3>
        <div class="space-y-3">
          <Field label="Nama profile" required>
            <Input v-model="form.name" placeholder="1HR-3K" :disabled="isEdit()" />
          </Field>
          <Field label="Rate limit (RX/TX)" hint="Format MikroTik, mis. 10M/3M">
            <Input v-model="form.rate_limit" placeholder="10M/3M" />
          </Field>
        </div>
      </section>

      <section>
        <h3
          class="mb-2 text-[10px] font-semibold uppercase"
          style="color: var(--muted); letter-spacing: 0.08em"
        >
          Mikhmon Config
        </h3>
        <div class="space-y-3">
          <Field
            label="Expiry mode"
            hint="`Free` = tidak record selling. `*record` = otomatis catat transaksi via webhook on-login."
            required
          >
            <Select
              v-model="form.expiry_mode"
              :options="EXPIRY_MODE_OPTIONS.map((o) => ({ value: o.value, label: o.label }))"
            />
          </Field>
          <Field label="Validity" hint="Durasi RouterOS (mis. `7d`, `1w`, `1mo`).">
            <Select
              v-model="form.validity"
              :options="[
                { value: '', label: '— tidak ada —' },
                { value: '1h', label: '1 jam' },
                { value: '6h', label: '6 jam' },
                { value: '1d', label: '1 hari' },
                { value: '7d', label: '7 hari' },
                { value: '30d', label: '30 hari' },
              ]"
            />
          </Field>
          <div class="grid grid-cols-2 gap-3">
            <Field label="Harga modal (Rp)">
              <NumberInput v-model="form.price" :min="0" :step="500" />
            </Field>
            <Field label="Harga jual (Rp)">
              <NumberInput v-model="form.sell_price" :min="0" :step="500" />
            </Field>
          </div>
          <div
            class="flex items-center justify-between rounded-lg p-3"
            style="background: var(--bg-2)"
          >
            <div>
              <div class="text-sm font-medium">Lock to MAC</div>
              <div class="text-xs" style="color: var(--muted)">
                Set MAC address user otomatis saat login pertama
              </div>
            </div>
            <Toggle v-model="form.lock_mac" />
          </div>
        </div>
      </section>
    </div>
    <template #footer>
      <button class="btn btn-sm" type="button" @click="emit('close')">Batal</button>
      <button class="btn btn-primary btn-sm" type="button" @click="submit">
        {{ isEdit() ? 'Simpan' : 'Tambah' }}
      </button>
    </template>
  </Modal>
</template>
