export interface NetworkInterface {
  id: string
  name: string
  type: string
  running: boolean
  disabled: boolean
  mac_address?: string
  rx_bytes?: number
  tx_bytes?: number
}

export interface IPPool {
  id: string
  name: string
  ranges: string
  next_pool?: string
}

export interface ARPEntry {
  id: string
  address: string
  mac_address: string
  interface: string
  dynamic: boolean
}

export interface DHCPLease {
  id: string
  address: string
  mac_address: string
  client_id?: string
  host_name?: string
  server?: string
  status: string
}

export interface SimpleQueue {
  id: string
  name: string
  target: string
  max_limit?: string
  burst_limit?: string
  disabled?: boolean
}
