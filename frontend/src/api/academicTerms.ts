import apiClient from './client'
import type {
  AcademicTerm,
  CreateAcademicTermRequest,
  CreateResponse,
} from '@/types'

export const academicTermsApi = {
  // List all academic terms
  listAcademicTerms: async (): Promise<AcademicTerm[]> => {
    const response = await apiClient.get<AcademicTerm[]>('/academic-terms')
    return response.data
  },

  // Get the active academic term
  getActiveAcademicTerm: async (): Promise<AcademicTerm> => {
    const response = await apiClient.get<AcademicTerm>('/academic-terms/active')
    return response.data
  },

  // Create a new academic term
  createAcademicTerm: async (
    data: CreateAcademicTermRequest
  ): Promise<CreateResponse> => {
    const response = await apiClient.post<CreateResponse>('/academic-terms', data)
    return response.data
  },

  // Activate an academic term
  activateAcademicTerm: async (id: string): Promise<void> => {
    await apiClient.patch(`/academic-terms/${id}/activate`)
  },

  // Deactivate an academic term
  deactivateAcademicTerm: async (id: string): Promise<void> => {
    await apiClient.patch(`/academic-terms/${id}/deactivate`)
  },
}
