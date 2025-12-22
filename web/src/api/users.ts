import { fetcher } from '../lib/fetcher'

// User model
export interface User {
  id: number
  username: string
  email: string
  status: number
  is_admin: number
  created_at: string
}

// Request types
export interface RegisterRequest {
  username: string
  password: string
  email: string
  [key: string]: string
}

export interface LoginRequest {
  email: string
  password: string
  [key: string]: string
}

export interface UpdateProfileRequest {
  username?: string
  email?: string
  [key: string]: string | undefined
}

export interface ChangePasswordRequest {
  new_password: string
  [key: string]: string
}

// API methods
export const usersApi = {
  /**
   * Register a new user
   * POST /api/v1/users/register
   */
  register(data: RegisterRequest): Promise<User> {
    return fetcher<User>('/v1/users/register', {
      method: 'POST',
      body: data,
    })
  },

  /**
   * Login with email and password
   * POST /api/v1/users/login
   */
  login(data: LoginRequest): Promise<User> {
    return fetcher<User>('/v1/users/login', {
      method: 'POST',
      body: data,
    })
  },

  /**
   * Logout current user
   * POST /api/v1/users/logout
   */
  logout(): Promise<null> {
    return fetcher<null>('/v1/users/logout', {
      method: 'POST',
    })
  },

  /**
   * Get user by ID
   * GET /api/v1/users/:id
   */
  getUser(id: number): Promise<User> {
    return fetcher<User>(`/v1/users/${id}`, {
      method: 'GET',
    })
  },

  /**
   * Update user profile
   * PUT /api/v1/users/:id
   */
  updateProfile(id: number, data: UpdateProfileRequest): Promise<User> {
    return fetcher<User>(`/v1/users/${id}`, {
      method: 'PUT',
      body: data,
    })
  },

  /**
   * Change user password
   * PUT /api/v1/users/:id/password
   */
  changePassword(id: number, data: ChangePasswordRequest): Promise<null> {
    return fetcher<null>(`/v1/users/${id}/password`, {
      method: 'PUT',
      body: data,
    })
  },
}
