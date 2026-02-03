import apiClient from './client'
import type { Resource, UserResource } from '@/types'

export const resourcesApi = {
  getResourcesByWeek: async (weekId: string): Promise<Resource[]> => {
    const response = await apiClient.get(`/resources/weeks/${weekId}`)
    return response.data
  },

  uploadFile: async (weekId: string, formData: FormData): Promise<void> => {
    await apiClient.post(`/resources/file/${weekId}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },

  uploadLink: async (weekId: string, name: string, url: string): Promise<void> => {
    await apiClient.post(`/resources/link/${weekId}`, {
      name,
      url,
    })
  },

  downloadResource: async (objectId: string): Promise<void> => {
    try {
      // Backend now returns the presigned URL as JSON
      const response = await apiClient.get<{ url: string }>(`/resources/${objectId}`)
      const presignedUrl = response.data.url
      
      if (presignedUrl) {
        // Open the presigned S3 URL directly
        window.open(presignedUrl, '_blank')
      } else {
        console.error('No URL returned from backend')
      }
    } catch (error) {
      console.error('Download failed:', error)
      throw error
    }
  },

  getUserResources: async (userId: string): Promise<UserResource[]> => {
    const response = await apiClient.get(`/resources/users/${userId}`)
    return response.data
  },

  deleteResource: async (resourceId: string): Promise<void> => {
    await apiClient.delete(`/resources/${resourceId}`)
  },
}
