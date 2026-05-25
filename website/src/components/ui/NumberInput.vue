<script setup lang="ts">
import Icon from './Icon.vue'

const props = withDefaults(
  defineProps<{
    modelValue?: number
    min?: number
    max?: number
    step?: number
    placeholder?: string
    disabled?: boolean
    sm?: boolean
  }>(),
  { step: 1, sm: false, disabled: false, modelValue: 0 },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()

function clamp(n: number) {
  let v = n
  if (props.min != null) v = Math.max(props.min, v)
  if (props.max != null) v = Math.min(props.max, v)
  return v
}

function set(v: number) {
  emit('update:modelValue', clamp(v))
}

function onInput(e: Event) {
  const v = Number((e.target as HTMLInputElement).value)
  set(Number.isFinite(v) ? v : 0)
}

function inc() {
  set((props.modelValue ?? 0) + props.step)
}
function dec() {
  set((props.modelValue ?? 0) - props.step)
}
</script>

<template>
  <div class="relative flex items-stretch">
    <input
      type="number"
      class="input"
      :class="{ 'input-sm': sm }"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :min="min"
      :max="max"
      :step="step"
      style="padding-right: 60px"
      @input="onInput"
    />
    <div class="absolute top-0 right-1 flex h-full items-center gap-0.5">
      <button
        type="button"
        class="btn btn-ghost btn-icon"
        :class="sm ? 'btn-xs' : 'btn-sm'"
        :disabled="disabled || (min != null && modelValue <= min)"
        @click="dec"
      >
        <Icon name="Down" :size="11" />
      </button>
      <button
        type="button"
        class="btn btn-ghost btn-icon"
        :class="sm ? 'btn-xs' : 'btn-sm'"
        :disabled="disabled || (max != null && modelValue >= max)"
        @click="inc"
      >
        <Icon name="Up" :size="11" />
      </button>
    </div>
  </div>
</template>
