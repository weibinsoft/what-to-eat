import axios, { type AxiosResponse } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：添加 Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器：处理错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 类型定义
export interface User {
  id: number
  username: string
}

export interface Restaurant {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface Menu {
  id: number
  restaurant_id: number
  dish_name: string
  created_at: string
  updated_at: string
  restaurant: Restaurant
}

export interface DecisionRecord {
  id: number
  user_id: number
  menu_id: number
  decided_at: string
  menu: Menu
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface LoginResponse {
  token: string
  user_id: number
  username: string
}

export interface DecideResponse {
  menu: Menu
  message: string
}

export interface HistoryResponse {
  records: DecisionRecord[]
  total: number
}

export interface CreateMenuRequest {
  restaurant_name: string
  dish_name: string
}

export interface CreateMenuResponse {
  menu: Menu
  is_new_restaurant: boolean
}

// API 方法
export const authApi = {
  register: (username: string, password: string): Promise<AxiosResponse<ApiResponse<User>>> => {
    return api.post('/api/auth/register', { username, password })
  },
  login: (username: string, password: string): Promise<AxiosResponse<ApiResponse<LoginResponse>>> => {
    return api.post('/api/auth/login', { username, password })
  },
}

export const restaurantApi = {
  list: (): Promise<AxiosResponse<ApiResponse<Restaurant[]>>> => {
    return api.get('/api/restaurants')
  },
}

export const menuApi = {
  list: (): Promise<AxiosResponse<ApiResponse<Menu[]>>> => {
    return api.get('/api/menus')
  },
  create: (restaurantName: string, dishName: string): Promise<AxiosResponse<ApiResponse<CreateMenuResponse>>> => {
    return api.post('/api/menus', { restaurant_name: restaurantName, dish_name: dishName })
  },
  delete: (id: number): Promise<AxiosResponse<ApiResponse<null>>> => {
    return api.delete(`/api/menus/${id}`)
  },
}

export const decisionApi = {
  decide: (menuIds?: number[]): Promise<AxiosResponse<ApiResponse<DecideResponse>>> => {
    return api.post('/api/decide', { menu_ids: menuIds || [] })
  },
  history: (): Promise<AxiosResponse<ApiResponse<HistoryResponse>>> => {
    return api.get('/api/history')
  },
}

export default api
