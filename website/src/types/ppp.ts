export interface PPPSecret {
  id: string
  name: string
  password?: string
  service: string
  profile: string
  local_address?: string
  remote_address?: string
  caller_id?: string
  disabled?: boolean
  comment?: string
}

export interface PPPProfile {
  id: string
  name: string
  local_address?: string
  remote_address?: string
  rate_limit?: string
  session_timeout?: string
  only_one?: 'yes' | 'no' | 'default'
}

export interface PPPActive {
  id: string
  name: string
  service: string
  caller_id?: string
  address?: string
  uptime: string
  encoding?: string
}

export interface PPPInactiveEvent {
  name: string
  profile: string
  last_seen_address?: string
  last_seen_at?: string
}
