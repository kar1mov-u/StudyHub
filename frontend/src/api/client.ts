import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Helper function to transform user object from backend format to frontend format
const transformUserData = (data: any): any => {
  if (data && typeof data === 'object') {
    // If this is a user object (has user_id field), transform it
    if ('user_id' in data) {
      return {
        ID: data.user_id,
        FirstName: data.first_name,
        LastName: data.last_name,
        Email: data.email,
        IsAdmin: data.is_admin,
        CreatedAt: data.created_at || data.CreatedAt,
        UpdatedAt: data.updated_at || data.UpdatedAt,
      }
    }
    
    // If data is an array, transform each item
    if (Array.isArray(data)) {
      return data.map(transformUserData)
    }
    
    // If data is an object, recursively transform nested objects
    const transformed: any = {}
    for (const [key, value] of Object.entries(data)) {
      transformed[key] = transformUserData(value)
    }
    return transformed
  }
  
  return data
}

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
      // Transform user_id to ID for user objects
      const unwrappedData = transformUserData(response.data.data)
      return { ...response, data: unwrappedData }
    }
    // Transform even if data is not wrapped
    return { ...response, data: transformUserData(response.data) }
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
