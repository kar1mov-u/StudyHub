import apiClient from './client'
import type {
  AcademicTerm,
  CreateAcademicTermRequest,
  CreateResponse,
} from '@/types'

export const academicTermsApi = {
  // Get the current academic term
  getCurrentAcademicTerm: async (): Promise<AcademicTerm> => {
    const response = await apiClient.get<AcademicTerm>('/academic-terms/current')
    return response.data
  },

  // Create a new academic term (creates term + all module runs)
  createNewTerm: async (
    data: CreateAcademicTermRequest
  ): Promise<CreateResponse> => {
    const response = await apiClient.post<CreateResponse>('/academic-terms/new-term', data)
    return response.data
  },
}
