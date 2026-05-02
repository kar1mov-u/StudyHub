import apiClient from './client'

export const chatApi = {
  sendMessage: async (message: string): Promise<string> => {
    const response = await apiClient.post('/chat', { message })
    return response.data.reply
  },
}
