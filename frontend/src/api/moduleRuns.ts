import apiClient from './client'
import type {
  ModuleRun,
  ModuleRunPage,
  CreateModuleRunRequest,
  CreateResponse,
} from '@/types'

export const moduleRunsApi = {
  // List all runs for a specific module
  listModuleRuns: async (moduleId: string): Promise<ModuleRun[]> => {
    const response = await apiClient.get<ModuleRun[]>(`/modules/${moduleId}/runs`)
    return response.data
  },

  // Create a new module run
  createModuleRun: async (
    moduleId: string,
    data: CreateModuleRunRequest
  ): Promise<CreateResponse> => {
    const response = await apiClient.post<CreateResponse>(
      `/modules/${moduleId}/runs`,
      data
    )
    return response.data
  },

  // Get module run by ID with weeks
  getModuleRun: async (id: string): Promise<ModuleRunPage> => {
    const response = await apiClient.get<ModuleRunPage>(`/module-runs/${id}`)
    return response.data
  },

  // Delete a module run
  deleteModuleRun: async (id: string): Promise<void> => {
    await apiClient.delete(`/module-runs/${id}`)
  },
}
