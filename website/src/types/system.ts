export interface SystemIdentity {
  name: string
  identity?: string
}

export interface SystemResource {
  uptime: string
  free_memory: number
  total_memory: number
  cpu_load: number
  cpu_frequency?: number
  version: string
  board_name: string
  architecture_name?: string
  free_hdd_space?: number
  total_hdd_space?: number
  cpu?: string
  cpu_count?: number
}

export interface Routerboard {
  routerboard: boolean
  board_name?: string
  model?: string
  serial_number?: string
  firmware_type?: string
  current_firmware?: string
  upgrade_firmware?: string
}

export interface SystemClock {
  time: string
  date: string
  timezone: string
}

export interface SystemLicense {
  software_id?: string
  n_level?: string
  features?: string
}

export interface SystemScript {
  id: string
  name: string
  source: string
  policy?: string[]
  last_started?: string
  run_count?: number
}

export interface SystemScheduler {
  id: string
  name: string
  start_date: string
  start_time: string
  interval: string
  on_event: string
  disabled?: boolean
  run_count?: number
}
