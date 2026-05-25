<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageHeader from '@/components/ui/PageHeader.vue'
import Icon from '@/components/ui/Icon.vue'
import Tabs from '@/components/ui/Tabs.vue'
import HotspotUsersTab from '@/components/hotspot/tabs/HotspotUsersTab.vue'
import HotspotProfilesTab from '@/components/hotspot/tabs/HotspotProfilesTab.vue'
import HotspotActiveTab from '@/components/hotspot/tabs/HotspotActiveTab.vue'
import HotspotInactiveTab from '@/components/hotspot/tabs/HotspotInactiveTab.vue'
import HotspotIpBindingTab from '@/components/hotspot/tabs/HotspotIpBindingTab.vue'
import HotspotCookiesTab from '@/components/hotspot/tabs/HotspotCookiesTab.vue'
import HotspotHostsTab from '@/components/hotspot/tabs/HotspotHostsTab.vue'
import { useActiveDevice } from '@/composables/useActiveDevice'
import { useHotspotUsersQuery, useHotspotActiveQuery, useHotspotProfilesQuery } from '@/queries/hotspot.queries'

type TabId = 'users' | 'profiles' | 'active' | 'inactive' | 'ip-binding' | 'cookies' | 'hosts'

const route = useRoute()
const router = useRouter()
const { activeDeviceId } = useActiveDevice()

// Fetch data untuk counts di tab header
const { data: hotspotUsers } = useHotspotUsersQuery(activeDeviceId)
const { data: hotspotActive } = useHotspotActiveQuery(activeDeviceId)
const { data: hotspotProfiles } = useHotspotProfilesQuery(activeDeviceId)

const tab = computed<TabId>({
  get() {
    const t = route.params.tab as string
    const valid: TabId[] = ['users', 'profiles', 'active', 'inactive', 'ip-binding', 'cookies', 'hosts']
    return valid.includes(t as TabId) ? (t as TabId) : 'users'
  },
  set(v) {
    router.replace({ name: 'hotspot', params: { tab: v } })
  },
})

const tabs = computed(() => [
  {
    id: 'users' as TabId,
    label: 'Users',
    icon: 'Users' as const,
    count: hotspotUsers.value?.length ?? 0,
  },
  {
    id: 'profiles' as TabId,
    label: 'Profiles',
    icon: 'Wifi' as const,
    count: hotspotProfiles.value?.length ?? 0,
  },
  {
    id: 'active' as TabId,
    label: 'Active',
    icon: 'Activity' as const,
    count: hotspotActive.value?.length ?? 0,
    live: true,
  },
  {
    id: 'inactive' as TabId,
    label: 'Inactive',
    icon: 'Power' as const,
    count: hotspotUsers.value?.filter((u) => u.disabled).length ?? 0,
  },
  {
    id: 'ip-binding' as TabId,
    label: 'IP-Binding',
    icon: 'Link2' as const,
  },
  {
    id: 'cookies' as TabId,
    label: 'Cookies',
    icon: 'Cookie' as const,
  },
  {
    id: 'hosts' as TabId,
    label: 'Hosts',
    icon: 'Server' as const,
  },
])

watch(
  () => route.params.tab,
  (t) => {
    const valid: TabId[] = ['users', 'profiles', 'active', 'inactive', 'ip-binding', 'cookies', 'hosts']
    if (!valid.includes(t as TabId)) {
      router.replace({ name: 'hotspot', params: { tab: 'users' } })
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="fade-in">
    <PageHeader title="Hotspot" subtitle="Manajemen user, profile, session, dan binding">
      <template #right>
        <button
          v-if="!activeDeviceId"
          class="btn btn-xs"
          style="color: var(--warn)"
          type="button"
          @click="$router.push('/devices')"
        >
          <Icon name="AlertCircle" :size="13" />
          Pilih device
        </button>
        <button class="btn btn-sm" type="button" @click="$router.push('/hotspot/voucher')">
          <Icon name="Ticket" :size="13" />
          Voucher Baru
        </button>
      </template>
    </PageHeader>

    <Tabs v-model="tab" :tabs="tabs" class="mb-4" />

    <HotspotUsersTab v-if="tab === 'users'" />
    <HotspotProfilesTab v-else-if="tab === 'profiles'" />
    <HotspotActiveTab v-else-if="tab === 'active'" />
    <HotspotInactiveTab v-else-if="tab === 'inactive'" />
    <HotspotIpBindingTab v-else-if="tab === 'ip-binding'" />
    <HotspotCookiesTab v-else-if="tab === 'cookies'" />
    <HotspotHostsTab v-else-if="tab === 'hosts'" />
  </div>
</template>
