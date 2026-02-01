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

  // Decode JWT token to get user ID
  const decodeToken = (token: string): string | null => {
    try {
      const base64Url = token.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )
      const payload = JSON.parse(jsonPayload)
      return payload.sub // User ID from 'sub' claim
    } catch (error) {
      console.error('Failed to decode token:', error)
      return null
    }
  }

  // Initialize auth state on mount
  useEffect(() => {
    const initializeAuth = async () => {
      const token = localStorage.getItem('auth_token')
      const cachedUser = localStorage.getItem('auth_user')

      if (token && cachedUser) {
        try {
          const parsedUser = JSON.parse(cachedUser)
          setUser(parsedUser)
        } catch (error) {
          console.error('Failed to parse cached user:', error)
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

      // Decode token to get user ID
      const userId = decodeToken(token)
      if (!userId) {
        throw new Error('Invalid token')
      }

      // Fetch user details
      const userData = await authApi.getCurrentUser(userId)
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
