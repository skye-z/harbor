import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo } from '../../types'

interface LoginResponse {
  token: string
  user_id: number
  username: string
  is_admin: boolean
}

export const useUserStore = defineStore('user', () => {
  const user = ref<UserInfo | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  const setUser = (userData: UserInfo) => {
    user.value = userData
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserFromLogin = (response: LoginResponse) => {
    if (response.token) {
      token.value = response.token
      localStorage.setItem('token', response.token)
    }
    user.value = {
      id: response.user_id?.toString() || '',
      username: response.username || '',
      is_admin: response.is_admin ?? false
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    user,
    token,
    isAuthenticated,
    isAdmin,
    setUser,
    setUserFromLogin,
    setToken,
    logout
  }
})
