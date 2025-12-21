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
}

export interface LoginRequest {
  email: string
  password: string
}

export interface UpdateProfileRequest {
  username?: string
  email?: string
}

export interface ChangePasswordRequest {
  new_password: string
}

// Response types
export interface RegisterResponse {
  message: string
  user: User
}

export interface LoginResponse {
  message: string
  user: User
}

export interface GetUserResponse {
  user: User
}

export interface UpdateProfileResponse {
  message: string
  user: User
}

export interface ChangePasswordResponse {
  message: string
}

// API methods
export const usersApi = {
  /**
   * Register a new user
   * POST /api/v1/users/register
   */
  register(data: RegisterRequest): Promise<RegisterResponse> {
    return fetcher<RegisterResponse>('/v1/users/register', {
      method: 'POST',
      body: data,
    })
  },

  /**
   * Login with email and password
   * POST /api/v1/users/login
   */
  login(data: LoginRequest): Promise<LoginResponse> {
    return fetcher<LoginResponse>('/v1/users/login', {
      method: 'POST',
      body: data,
    })
  },

  /**
   * Get user by ID
   * GET /api/v1/users/:id
   */
  getUser(id: number): Promise<GetUserResponse> {
    return fetcher<GetUserResponse>(`/v1/users/${id}`, {
      method: 'GET',
    })
  },

  /**
   * Update user profile
   * PUT /api/v1/users/:id
   */
  updateProfile(id: number, data: UpdateProfileRequest): Promise<UpdateProfileResponse> {
    return fetcher<UpdateProfileResponse>(`/v1/users/${id}`, {
      method: 'PUT',
      body: data,
    })
  },

  /**
   * Change user password
   * PUT /api/v1/users/:id/password
   */
  changePassword(id: number, data: ChangePasswordRequest): Promise<ChangePasswordResponse> {
    return fetcher<ChangePasswordResponse>(`/v1/users/${id}/password`, {
      method: 'PUT',
      body: data,
    })
  },
}
