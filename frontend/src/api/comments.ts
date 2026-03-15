import apiClient from './client'
import type { Comment, CreateCommentRequest } from '@/types'

export const commentsApi = {
  getCommentsByWeek: async (weekId: string): Promise<Comment[]> => {
    const response = await apiClient.get<Comment[]>(`/comments/weeks/${weekId}`)
    return response.data || []
  },

  createComment: async (data: CreateCommentRequest): Promise<Comment> => {
    const response = await apiClient.post<Comment>('/comments', data)
    return response.data
  },

  upvoteComment: async (commentId: string): Promise<void> => {
    await apiClient.post(`/comments/${commentId}/upvote`)
  },

  downvoteComment: async (commentId: string): Promise<void> => {
    await apiClient.post(`/comments/${commentId}/downvote`)
  },
}
