<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder?: string
    rows?: number
    mono?: boolean
    disabled?: boolean
  }>(),
  { rows: 4, mono: false, disabled: false, modelValue: '' },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <textarea
    class="input"
    :class="{ mono: props.mono }"
    :rows="rows"
    :placeholder="placeholder"
    :value="modelValue"
    :disabled="disabled"
    style="height: auto; padding: 10px 12px; line-height: 1.5; resize: vertical"
    @input="onInput"
  />
</template>
