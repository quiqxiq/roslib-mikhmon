import { type MaybeRefOrGetter, toValue } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { systemService } from '@/services/system'
import { systemScriptsService } from '@/services/system-scripts'
import { systemSchedulersService } from '@/services/system-schedulers'
import { queryKeys } from '@/queries/query-keys'

export function useIdentityQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.identity(String(toValue(deviceId))),
    queryFn: () => systemService.identity(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useResourceQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.resource(String(toValue(deviceId))),
    queryFn: () => systemService.resource(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useRouterboardQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.routerboard(String(toValue(deviceId))),
    queryFn: () => systemService.routerboard(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useClockQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.clock(String(toValue(deviceId))),
    queryFn: () => systemService.clock(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useLicenseQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.license(String(toValue(deviceId))),
    queryFn: () => systemService.license(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useScriptsQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.scripts(String(toValue(deviceId))),
    queryFn: () => systemScriptsService.list(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}

export function useSchedulersQuery(deviceId: MaybeRefOrGetter<string | null>) {
  return useQuery({
    queryKey: queryKeys.system.schedulers(String(toValue(deviceId))),
    queryFn: () => systemSchedulersService.list(String(toValue(deviceId))),
    enabled: () => Boolean(toValue(deviceId)),
  })
}
