import apiClient from './client'

export interface ChatSource {
  source: string
  page?: number
}

export interface ChatReply {
  reply: string
  sources: ChatSource[]
}

export interface HistoryMessage {
  role: 'user' | 'assistant'
  content: string
}

export const chatApi = {
  sendMessage: async (message: string, history: HistoryMessage[]): Promise<ChatReply> => {
    const response = await apiClient.post('/chat', { message, history })
    return {
      reply: response.data.reply ?? '',
      sources: response.data.sources ?? [],
    }
  },
}
