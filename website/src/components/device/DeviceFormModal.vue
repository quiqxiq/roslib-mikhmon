<script setup lang="ts">
import { ref, watch } from 'vue'
import Modal from '@/components/ui/Modal.vue'
import Field from '@/components/ui/Field.vue'
import Input from '@/components/ui/Input.vue'
import NumberInput from '@/components/ui/NumberInput.vue'
import Icon from '@/components/ui/Icon.vue'
import type { FixtureDevice } from '@/fixtures/devices'

const props = defineProps<{
  open: boolean
  initial?: Partial<FixtureDevice> | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', d: Partial<FixtureDevice> & { port?: number; password?: string }): void
}>()

interface FormState {
  id?: string
  name: string
  address: string
  port: number
  username: string
  password: string
}

const form = ref<FormState>({
  name: '',
  address: '',
  port: 8728,
  username: 'admin',
  password: '',
})
const showPwd = ref(false)

watch(
  () => props.initial,
  (v) => {
    if (v) {
      form.value = {
        id: v.id,
        name: v.name ?? '',
        address: v.address ?? '',
        port: 8728,
        username: 'admin',
        password: '',
      }
    } else {
      form.value = {
        name: '',
        address: '',
        port: 8728,
        username: 'admin',
        password: '',
      }
    }
  },
  { immediate: true },
)

const isEdit = () => Boolean(props.initial?.id)

function submit() {
  if (!form.value.name.trim() || !form.value.address.trim()) return
  emit('save', {
    id: form.value.id,
    name: form.value.name.trim(),
    address: form.value.address.trim(),
    port: form.value.port,
    password: form.value.password,
  })
}
</script>

<template>
  <Modal :open="open" :title="isEdit() ? 'Edit Device' : 'Tambah Device'" @close="emit('close')">
    <div class="space-y-3">
      <Field label="Nama" required>
        <Input v-model="form.name" placeholder="Mikrotik HAP ac²" />
      </Field>
      <div class="grid grid-cols-3 gap-3">
        <Field label="Address" required class="col-span-2">
          <Input v-model="form.address" placeholder="192.168.88.1" />
        </Field>
        <Field label="Port">
          <NumberInput v-model="form.port" :min="1" :max="65535" />
        </Field>
      </div>
      <Field label="Username" required>
        <Input v-model="form.username" placeholder="admin" />
      </Field>
      <Field label="Password" :required="!isEdit()" :hint="isEdit() ? 'kosongkan jika tidak ganti' : undefined">
        <div class="relative">
          <Input
            v-model="form.password"
            :type="showPwd ? 'text' : 'password'"
            placeholder="••••••••"
          />
          <button
            type="button"
            class="btn btn-ghost btn-icon btn-sm absolute"
            style="right: 4px; top: 50%; transform: translateY(-50%)"
            @click="showPwd = !showPwd"
          >
            <Icon :name="showPwd ? 'EyeOff' : 'Eye'" :size="13" />
          </button>
        </div>
      </Field>
    </div>
    <template #footer>
      <button class="btn btn-sm" type="button" @click="emit('close')">Batal</button>
      <button class="btn btn-primary btn-sm" type="button" @click="submit">
        {{ isEdit() ? 'Simpan' : 'Tambah' }}
      </button>
    </template>
  </Modal>
</template>
