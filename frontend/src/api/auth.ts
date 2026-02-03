import apiClient from './client'
import type { LoginRequest, LoginResponse, RegisterRequest, User } from '@/types'

export const authApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  register: async (userData: RegisterRequest): Promise<{ id: string }> => {
    const response = await apiClient.post('/users', userData)
    return response.data
  },

  getCurrentUser: async (userId: string): Promise<User> => {
    const response = await apiClient.get(`/users/${userId}`)
    return response.data
  },

  getMe: async (): Promise<User> => {
    const response = await apiClient.get('/users/me')
    return response.data
  },
}
