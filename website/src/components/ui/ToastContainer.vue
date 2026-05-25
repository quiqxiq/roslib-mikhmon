<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import Toast from './Toast.vue'

const { toasts, dismiss } = useToast()
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed flex flex-col gap-2"
      style="bottom: 16px; right: 16px; z-index: 9999; max-width: calc(100vw - 32px)"
    >
      <Toast
        v-for="t in toasts"
        :key="t.id"
        :kind="t.kind"
        :message="t.message"
        @dismiss="dismiss(t.id)"
      />
    </div>
  </Teleport>
</template>
