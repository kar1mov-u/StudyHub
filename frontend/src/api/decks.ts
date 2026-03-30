import apiClient from './client'
import type { 
  UserDeckCard, 
  DeckStats, 
  AddCardToDeckRequest,
  CreateCustomCardRequest,
  UpdateCardRequest,
  RecordReviewRequest
} from '@/types'

export const decksApi = {
  // Get user's deck for a week
  getUserDeck: async (weekId: string): Promise<UserDeckCard[]> => {
    const response = await apiClient.get(`/decks/weeks/${weekId}/cards`)
    return response.data.data || []
  },

  // Add auto-generated card to deck
  addCardToDeck: async (weekId: string, request: AddCardToDeckRequest): Promise<void> => {
    await apiClient.post(`/decks/weeks/${weekId}/cards`, request)
  },

  // Create custom card
  createCustomCard: async (weekId: string, request: CreateCustomCardRequest): Promise<UserDeckCard> => {
    const response = await apiClient.post(`/decks/weeks/${weekId}/cards/custom`, request)
    return response.data.data
  },

  // Update card
  updateCard: async (cardId: string, request: UpdateCardRequest): Promise<void> => {
    await apiClient.patch(`/decks/cards/${cardId}`, request)
  },

  // Remove card from deck
  removeCard: async (cardId: string): Promise<void> => {
    await apiClient.delete(`/decks/cards/${cardId}`)
  },

  // Record card review
  recordReview: async (cardId: string, request: RecordReviewRequest): Promise<void> => {
    await apiClient.post(`/decks/cards/${cardId}/review`, request)
  },

  // Get deck statistics
  getDeckStats: async (weekId: string): Promise<DeckStats> => {
    const response = await apiClient.get(`/decks/weeks/${weekId}/stats`)
    return response.data.data
  },
}
