import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authApi } from '@/api/auth'
import type { User, LoginRequest, RegisterRequest } from '@/types'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (userData: RegisterRequest) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Initialize auth state on mount
  useEffect(() => {
    const initializeAuth = async () => {
      const token = localStorage.getItem('auth_token')

      if (token) {
        try {
          // Always fetch fresh user data from /users/me endpoint
          const userData = await authApi.getMe()
          setUser(userData)
          localStorage.setItem('auth_user', JSON.stringify(userData))
        } catch (error) {
          console.error('Failed to fetch user data:', error)
          // If token is invalid, clear storage
          localStorage.removeItem('auth_token')
          localStorage.removeItem('auth_user')
        }
      }
      setIsLoading(false)
    }

    initializeAuth()
  }, [])

  const login = async (credentials: LoginRequest) => {
    try {
      const { token } = await authApi.login(credentials)
      localStorage.setItem('auth_token', token)

      // Fetch user details using /users/me endpoint
      const userData = await authApi.getMe()
      setUser(userData)
      localStorage.setItem('auth_user', JSON.stringify(userData))
    } catch (error) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      throw error
    }
  }

  const register = async (userData: RegisterRequest) => {
    try {
      await authApi.register(userData)
      // Auto-login after registration
      await login({ email: userData.email, password: userData.password })
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
    setUser(null)
    window.location.href = '/login'
  }

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
