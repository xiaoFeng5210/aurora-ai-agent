import { apiGet } from './client'

export interface ToolCallFunction {
  name: string
  arguments: string
}

export interface ToolCall {
  id: string
  type: string
  function: ToolCallFunction
}

export interface HistoryMessage {
  id: number
  message_id: string
  document_id: number
  role: 'user' | 'assistant' | 'tool' | 'system'
  content: string
  tool_calls: ToolCall[] | null
  created_at: string
  updated_at: string
}

export interface HistoryQuery {
  startTime?: string
  endTime?: string
  keyword?: string
  page?: number
  pageSize?: number
  order?: 'asc' | 'desc'
}

export const listHistoryMessages = (documentId: number, query: HistoryQuery = {}) => {
  const params = new URLSearchParams()
  Object.entries(query).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') params.set(k, String(v))
  })
  const qs = params.toString()
  const path = `/documents/${documentId}/messages/proxy/history${qs ? `?${qs}` : ''}`
  return apiGet<HistoryMessage[]>(path)
}
