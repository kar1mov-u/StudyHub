import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

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
    // Handle error responses
    if (error.response?.data?.error) {
      const apiError = error.response.data.error
      throw new Error(apiError.message || 'An error occurred')
    }
    throw error
  }
)

export default apiClient
