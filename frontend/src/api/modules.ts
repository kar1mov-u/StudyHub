import apiClient from './client'
import type {
  Module,
  ModulePage,
  CreateModuleRequest,
  CreateResponse,
} from '@/types'

export const modulesApi = {
  // Get all modules
  listModules: async (): Promise<Module[]> => {
    const response = await apiClient.get<Module[]>('/modules')
    return response.data
  },

  // Get module by ID with full details
  getModuleFull: async (id: string): Promise<ModulePage> => {
    const response = await apiClient.get<ModulePage>(`/modules/${id}`)
    return response.data
  },

  // Create a new module
  createModule: async (data: CreateModuleRequest): Promise<CreateResponse> => {
    const response = await apiClient.post<CreateResponse>('/modules', data)
    return response.data
  },

  // Delete a module
  deleteModule: async (id: string): Promise<void> => {
    await apiClient.delete(`/modules/${id}`)
  },
}
