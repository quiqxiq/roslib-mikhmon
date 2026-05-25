<script setup lang="ts">
import { useTweaks } from '@/composables/useTweaks'
import { PALETTE_PRESETS } from '@/stores/tweaks'
import Segmented from '@/components/ui/Segmented.vue'
import Icon from '@/components/ui/Icon.vue'

const {
  theme,
  density,
  cardStyle,
  sidebarMode,
  chartKind,
  palette,
  setTheme,
  setDensity,
  setCardStyle,
  setSidebarMode,
  setChartKind,
  setPalette,
} = useTweaks()

function pickPalette(p: (typeof PALETTE_PRESETS)[number]) {
  setPalette(p)
}

function isPaletteActive(p: (typeof PALETTE_PRESETS)[number]) {
  return p[0] === palette.value[0] && p[1] === palette.value[1] && p[2] === palette.value[2]
}
</script>

<template>
  <div class="space-y-6">
    <section class="card p-4">
      <h3
        class="mb-3 text-[10px] font-semibold uppercase"
        style="color: var(--muted); letter-spacing: 0.08em"
      >
        Tampilan
      </h3>
      <div class="space-y-3">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <label class="text-sm" style="color: var(--text-2)">Tema</label>
          <Segmented
            :model-value="theme"
            :options="[
              { value: 'dark', label: 'Dark' },
              { value: 'light', label: 'Light' },
            ]"
            @update:model-value="setTheme"
          />
        </div>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <label class="text-sm" style="color: var(--text-2)">Density</label>
          <Segmented
            :model-value="density"
            :options="[
              { value: 'compact', label: 'Compact' },
              { value: 'regular', label: 'Regular' },
              { value: 'comfy', label: 'Comfy' },
            ]"
            @update:model-value="setDensity"
          />
        </div>
      </div>
    </section>

    <section class="card p-4">
      <h3
        class="mb-3 text-[10px] font-semibold uppercase"
        style="color: var(--muted); letter-spacing: 0.08em"
      >
        Sidebar
      </h3>
      <div class="flex flex-wrap items-center justify-between gap-3">
        <label class="text-sm" style="color: var(--text-2)">Mode</label>
        <Segmented
          :model-value="sidebarMode"
          :options="[
            { value: 'expanded', label: 'Expanded' },
            { value: 'icon', label: 'Icon' },
            { value: 'hidden', label: 'Hidden' },
          ]"
          @update:model-value="setSidebarMode"
        />
      </div>
    </section>

    <section class="card p-4">
      <h3
        class="mb-3 text-[10px] font-semibold uppercase"
        style="color: var(--muted); letter-spacing: 0.08em"
      >
        Cards & Charts
      </h3>
      <div class="space-y-3">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <label class="text-sm" style="color: var(--text-2)">Card style</label>
          <Segmented
            :model-value="cardStyle"
            :options="[
              { value: 'flat', label: 'Flat' },
              { value: 'elevated', label: 'Elevated' },
              { value: 'bordered', label: 'Bordered' },
            ]"
            @update:model-value="setCardStyle"
          />
        </div>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <label class="text-sm" style="color: var(--text-2)">Chart kind</label>
          <Segmented
            :model-value="chartKind"
            :options="[
              { value: 'area', label: 'Area' },
              { value: 'line', label: 'Line' },
              { value: 'bar', label: 'Bar' },
            ]"
            @update:model-value="setChartKind"
          />
        </div>
      </div>
    </section>

    <section class="card p-4">
      <h3
        class="mb-3 text-[10px] font-semibold uppercase"
        style="color: var(--muted); letter-spacing: 0.08em"
      >
        Palette
      </h3>
      <div class="grid grid-cols-5 gap-3 sm:grid-cols-5">
        <button
          v-for="(p, i) in PALETTE_PRESETS"
          :key="i"
          type="button"
          class="relative flex h-12 overflow-hidden rounded-lg"
          :style="{
            border: isPaletteActive(p)
              ? '2px solid var(--accent-cyan)'
              : '1px solid var(--border-strong)',
          }"
          @click="pickPalette(p)"
        >
          <span :style="{ flex: 2, background: p[0] }" />
          <span class="flex flex-1 flex-col">
            <span :style="{ flex: 1, background: p[1] }" />
            <span :style="{ flex: 1, background: p[2] }" />
          </span>
          <span
            v-if="isPaletteActive(p)"
            class="absolute inset-0 flex items-center justify-center"
            style="color: white; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5)"
          >
            <Icon name="Check" :size="16" />
          </span>
        </button>
      </div>
    </section>
  </div>
</template>
