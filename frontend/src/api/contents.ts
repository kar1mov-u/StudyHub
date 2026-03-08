import apiClient from './client'
import type { Flashcard } from '@/types'

export const contentsApi = {
  getFlashcards: async (objectIds: string[]): Promise<Flashcard[]> => {
    // Backend route has a typo ("conents" not "contents") - matching it exactly
    // Using request() because this is a GET with a JSON body
    const response = await apiClient.post('/conents/objects', { ids: objectIds })
    return response.data
  },
}
