import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add JWT token to all requests
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to unwrap the data from the API response
apiClient.interceptors.response.use(
  (response) => {
    // Unwrap { data: ... } from the response
    if (response.data && 'data' in response.data) {
      return { ...response, data: response.data.data }
    }
    return response
  },
  (error) => {
    // Handle 401 Unauthorized - clear token and redirect to login
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      // Redirect to login page if not already there
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    
    // Handle error responses
    if (error.response?.data?.error) {
      const apiError = error.response.data.error
      throw new Error(apiError.message || 'An error occurred')
    }
    throw error
  }
)

export default apiClient
