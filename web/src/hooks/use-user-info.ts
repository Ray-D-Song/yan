import { useLocalStorage } from '@vueuse/core'
import { usersApi, type User, type LoginRequest, type RegisterRequest } from '@/api/users'

const userInfo = useLocalStorage('user-info', {})

/**
 * Login with email and password
 * Session cookie is automatically managed by the browser
 */
async function login(credentials: LoginRequest) {
  try {
    const user = await usersApi.login(credentials)

    userInfo.value = user

    return user
  }
  catch (error) {
    // Clear any existing data on login failure
    userInfo.value = null
    throw error
  }
}

/**
 * Register a new user
 */
async function register(data: RegisterRequest) {
  try {
    const user = await usersApi.register(data)

    // Note: After registration, user still needs to login
    // Backend doesn't auto-login after registration

    return user
  }
  catch (error) {
    throw error
  }
}

/**
 * Logout current user
 * Clears session on backend and removes cookie
 */
async function logout() {
  try {
    await usersApi.logout()
  }
  catch (error) {
    console.error('Logout API call failed:', error)
    // Continue with local cleanup even if API fails
  }
  finally {
    // Clear local user info
    userInfo.value = null

    // Redirect to login page
    window.location.href = '/login'
  }
}

/**
 * Check if user is logged in
 */
function isLoggedIn() {
  return !!userInfo.value
}

/**
 * Get current user info
 */
function getCurrentUser() {
  return userInfo.value
}

export function useUserInfo() {

  return {
    userInfo,
    login,
    register,
    logout,
    isLoggedIn,
    getCurrentUser,
  }
}
