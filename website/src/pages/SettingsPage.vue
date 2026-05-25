<script setup lang="ts">
import { ref } from 'vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Tabs from '@/components/ui/Tabs.vue'
import TemplateSettingsTab from '@/components/settings/TemplateSettingsTab.vue'
import PreferenceSettingsTab from '@/components/settings/PreferenceSettingsTab.vue'

type TabId = 'template' | 'preferences'

const tab = ref<TabId>('template')

const tabs = [
  { id: 'template' as TabId, label: 'Template Voucher', icon: 'Printer' as const },
  { id: 'preferences' as TabId, label: 'UI Preferences', icon: 'Cog' as const },
]
</script>

<template>
  <div class="fade-in">
    <PageHeader title="Settings" subtitle="Konfigurasi template voucher dan preferensi tampilan" />

    <Tabs v-model="tab" :tabs="tabs" class="mb-4" />

    <TemplateSettingsTab v-if="tab === 'template'" />
    <PreferenceSettingsTab v-else-if="tab === 'preferences'" />
  </div>
</template>
