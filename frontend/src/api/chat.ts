import apiClient from './client'

export interface ChatSource {
  source: string
  page?: number
}

export interface ChatReply {
  reply: string
  sources: ChatSource[]
}

export const chatApi = {
  sendMessage: async (message: string): Promise<ChatReply> => {
    const response = await apiClient.post('/chat', { message })
    return {
      reply: response.data.reply ?? '',
      sources: response.data.sources ?? [],
    }
  },
}
