import { type MaybeRefOrGetter, toValue } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { profileConfigService } from '@/services/profile-config'
import { queryKeys } from '@/queries/query-keys'
import type { ProfileConfigUpsertRequest } from '@/types/profile-config'

export function useProfileConfigsQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.profileConfig.all(String(toValue(deviceId))),
    queryFn: () => profileConfigService.list(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useProfileConfigQuery(
  deviceId: MaybeRefOrGetter<string | null>,
  profileName: MaybeRefOrGetter<string | null>,
) {
  return useQuery({
    queryKey: queryKeys.profileConfig.byName(
      String(toValue(deviceId)),
      String(toValue(profileName)),
    ),
    queryFn: () =>
      profileConfigService.get(String(toValue(deviceId)), String(toValue(profileName))),
    enabled: () => Boolean(toValue(deviceId)) && Boolean(toValue(profileName)),
  })
}

// Upsert by profile name. onSuccess invalidate juga hotspot.profiles karena
// auto-inject backend mengubah on-login script di router (drawer profile bisa
// menampilkan source script terbaru).
export function useUpsertProfileConfigMutation(deviceId: MaybeRefOrGetter<string | null>) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (args: { profile_name: string; payload: ProfileConfigUpsertRequest }) =>
      profileConfigService.upsert(String(toValue(deviceId)), args.profile_name, args.payload),
    onSuccess: () => {
      const id = String(toValue(deviceId))
      qc.invalidateQueries({ queryKey: queryKeys.profileConfig.all(id) })
      qc.invalidateQueries({ queryKey: queryKeys.hotspot.profiles(id) })
    },
  })
}

export function useDeleteProfileConfigMutation(deviceId: MaybeRefOrGetter<string | null>) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (profileName: string) =>
      profileConfigService.remove(String(toValue(deviceId)), profileName),
    onSuccess: () => {
      const id = String(toValue(deviceId))
      qc.invalidateQueries({ queryKey: queryKeys.profileConfig.all(id) })
    },
  })
}

// Sync profile router → DB + auto-inject. Mengembalikan summary 5-list.
export function useSyncProfilesMutation(deviceId: MaybeRefOrGetter<string | null>) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => profileConfigService.sync(String(toValue(deviceId))),
    onSuccess: () => {
      const id = String(toValue(deviceId))
      qc.invalidateQueries({ queryKey: queryKeys.profileConfig.all(id) })
      qc.invalidateQueries({ queryKey: queryKeys.hotspot.profiles(id) })
    },
  })
}
