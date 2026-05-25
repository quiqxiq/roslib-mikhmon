export type Role = 'admin' | 'operator' | 'viewer'

export interface Credentials {
  username: string
  password: string
}

// Sesuai dto.LoginResponse / dto.RefreshResponse backend (snake_case).
export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

// Sesuai dto.UserResponse backend.
export interface AuthUser {
  id: string
  username: string
  role: Role
  created_at: string
}
