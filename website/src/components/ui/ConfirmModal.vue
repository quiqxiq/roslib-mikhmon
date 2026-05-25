<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Modal from './Modal.vue'
import Input from './Input.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    message?: string
    confirmText?: string
    cancelText?: string
    variant?: 'default' | 'danger'
    typeToConfirm?: string
    loading?: boolean
  }>(),
  {
    confirmText: 'OK',
    cancelText: 'Batal',
    variant: 'default',
    loading: false,
  },
)

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'close'): void
}>()

const typed = ref('')

watch(
  () => props.open,
  (v) => {
    if (v) typed.value = ''
  },
)

const canConfirm = computed(() => {
  if (!props.typeToConfirm) return true
  return typed.value === props.typeToConfirm
})
</script>

<template>
  <Modal :open="open" :title="title" @close="emit('close')">
    <p v-if="message" class="text-sm" style="color: var(--text-2)">{{ message }}</p>
    <div v-if="typeToConfirm" class="mt-4">
      <p class="mb-1.5 text-xs" style="color: var(--muted)">
        Untuk konfirmasi, ketik
        <b class="mono" style="color: var(--text)">{{ typeToConfirm }}</b>
        di bawah.
      </p>
      <Input v-model="typed" :placeholder="typeToConfirm" />
    </div>
    <template #footer>
      <button class="btn btn-sm" type="button" :disabled="loading" @click="emit('close')">
        {{ cancelText }}
      </button>
      <button
        type="button"
        class="btn btn-sm"
        :class="variant === 'danger' ? 'btn-danger' : 'btn-primary'"
        :disabled="!canConfirm || loading"
        @click="emit('confirm')"
      >
        {{ confirmText }}
      </button>
    </template>
  </Modal>
</template>
