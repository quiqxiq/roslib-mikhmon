<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'

defineProps<{
  align?: 'left' | 'right'
}>()

const open = ref(false)
// anchor dipakai oleh Vue template binding (ref="anchor")
const anchor = ref<HTMLElement | null>(null)
  void anchor.value // suppress TS6133
const panel = ref<HTMLElement | null>(null)

onClickOutside(panel, () => {
  open.value = false
})

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

defineExpose({ open, close })
</script>

<template>
  <div ref="anchor" class="relative inline-block">
    <span @click="toggle">
      <slot name="trigger" :open="open" :toggle="toggle" />
    </span>
    <Transition>
      <div
        v-if="open"
        ref="panel"
        class="card scale-in absolute z-50 min-w-[180px] py-1"
        :style="{
          top: 'calc(100% + 4px)',
          [align === 'left' ? 'left' : 'right']: '0',
          padding: '4px',
        }"
      >
        <slot :close="close" />
      </div>
    </Transition>
  </div>
</template>
