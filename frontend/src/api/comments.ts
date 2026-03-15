import axios from 'axios'
import type { Comment, CreateCommentRequest } from '@/types'

// Use a separate axios instance for comments to avoid the transformUserData
// interceptor in the main apiClient, which incorrectly transforms comment
// objects (they have a `user_id` field that triggers the user transform).
const commentClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token
commentClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Unwrap { data: ... } envelope without running transformUserData
commentClient.interceptors.response.use(
  (response) => {
    if (response.data && 'data' in response.data) {
      return { ...response, data: response.data.data }
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    if (error.response?.data?.error) {
      const apiError = error.response.data.error
      throw new Error(apiError.message || 'An error occurred')
    }
    throw error
  }
)

export const commentsApi = {
  getCommentsByWeek: async (weekId: string): Promise<Comment[]> => {
    const response = await commentClient.get<Comment[]>(`/comments/weeks/${weekId}`)
    return response.data || []
  },

  createComment: async (data: CreateCommentRequest): Promise<Comment> => {
    const response = await commentClient.post<Comment>('/comments', data)
    return response.data
  },

  upvoteComment: async (commentId: string): Promise<void> => {
    await commentClient.post(`/comments/${commentId}/upvote`)
  },

  downvoteComment: async (commentId: string): Promise<void> => {
    await commentClient.post(`/comments/${commentId}/downvote`)
  },
}
