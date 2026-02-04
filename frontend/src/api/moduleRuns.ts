import apiClient from './client'
import type {
  ModuleRun,
  ModuleRunPage,
} from '@/types'

export const moduleRunsApi = {
  // List all runs for a specific module
  listModuleRuns: async (moduleId: string): Promise<ModuleRun[]> => {
    const response = await apiClient.get<ModuleRun[]>(`/modules/${moduleId}/runs`)
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
