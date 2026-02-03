import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type LoginResponse } from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<{ id: number; username: string } | null>(
    JSON.parse(localStorage.getItem('user') || 'null')
  )

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const response = await authApi.login(username, password)
    const data: LoginResponse = response.data.data

    token.value = data.token
    user.value = { id: data.user_id, username: data.username }

    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  async function register(username: string, password: string) {
    await authApi.register(username, password)
    // 注册成功后自动登录
    await login(username, password)
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout,
  }
})
