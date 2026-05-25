import type { HotspotSession, HotspotUser } from '@/types/hotspot'
import type { PPPActive, PPPInactiveEvent, PPPSecret } from '@/types/ppp'
import type { Routerboard, SystemResource } from '@/types/system'

export type StreamChangeKind = 'added' | 'updated' | 'removed'

export interface ChangeEvent<T> {
  kind: StreamChangeKind
  id: string
  item: T
}

export type HotspotActiveStreamEvent = ChangeEvent<HotspotSession>
export type HotspotUsersStreamEvent = ChangeEvent<HotspotUser>
export type HotspotInactiveStreamEvent = ChangeEvent<HotspotUser>
export type PPPActiveStreamEvent = ChangeEvent<PPPActive>
export type PPPSecretsStreamEvent = ChangeEvent<PPPSecret>
export type PPPInactiveStreamEvent = ChangeEvent<PPPInactiveEvent>

// Log SSE event — backend sends flat fields: time, topics, message.
export interface LogStreamEvent {
  time: string
  topics: string
  message: string
}

// Resource SSE — backend sends flat SystemResource fields directly (event: resource).
export type ResourceStreamEvent = SystemResource

// Traffic monitor-traffic SSE — backend sends flat per-interface fields (event: traffic).
export interface TrafficStreamEvent {
  name: string
  'rx-bits-per-second': number
  'tx-bits-per-second': number
  'rx-packets-per-second': number
  'tx-packets-per-second': number
}

// Interface stats SSE — backend sends flat per-row fields (event: stats).
// Not wrapped in { interfaces: [...] }.
export interface InterfaceStatsStreamEvent {
  id: string
  name: string
  type: string
  rx_byte: number
  tx_byte: number
  rx_packet: number
  tx_packet: number
  running: boolean
  disabled: boolean
}

// Queue stats SSE — backend sends flat QueueStatsEvent fields per row (event: stats).
export interface QueueStatsStreamEvent {
  id: string
  name: string
  target: string
  parent?: string
  disabled: boolean
  dynamic: boolean
  comment?: string
  bytes: string
  packets: string
  rate: string
  total_rate?: string
  packet_rate?: string
  total_packet_rate?: string
  queued_bytes?: string
  total_queued_bytes?: string
  queued_packets?: string
  total_queued_packets?: string
  total_bytes?: string
  total_packets?: string
  dropped?: string
  total_dropped?: string
  max_limit?: string
}

// Routerboard SSE — backend sends flat RouterboardEvent fields (event: routerboard).
export type RouterboardStreamEvent = Routerboard

export interface PingStreamEvent {
  seq: number
  host: string
  size: number
  ttl: number
  time_ms: number
  status: string
}
