<script setup lang="ts">
import Icon from './Icon.vue'
import type { IconName } from './icons'

withDefaults(
  defineProps<{
    icon?: IconName
    label: string
    danger?: boolean
    disabled?: boolean
  }>(),
  { danger: false, disabled: false },
)

defineEmits<{
  (e: 'click'): void
}>()
</script>

<template>
  <button
    type="button"
    class="row-hover flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-left text-[13px]"
    :style="{
      color: danger ? 'var(--danger)' : 'var(--text-2)',
      opacity: disabled ? 0.5 : 1,
      cursor: disabled ? 'not-allowed' : 'pointer',
      background: 'transparent',
      border: 0,
    }"
    :disabled="disabled"
    @click="$emit('click')"
  >
    <Icon v-if="icon" :name="icon" :size="14" />
    <span class="flex-1">{{ label }}</span>
  </button>
</template>
