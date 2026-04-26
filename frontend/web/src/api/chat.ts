import { parseSSEStream, type SSEEventHandler } from '@/lib/sse'
import { API_BASE, HttpError } from './client'

export interface ChatPromptItem {
  role: 'user' | 'assistant' | 'system'
  content: string
  fileContentList?: unknown[]
}

export interface ChatStreamRequest {
  thinking: { type: 'enabled' | 'disabled' }
  max_tokens: number
  temperature?: number
  top_p?: number
  tools?: unknown[]
  prompt: ChatPromptItem[]
}

export async function streamChat(
  documentId: number,
  body: ChatStreamRequest,
  onEvent: SSEEventHandler,
  signal?: AbortSignal,
): Promise<void> {
  const res = await fetch(`${API_BASE}/chat/glm/stream/${documentId}`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
    signal,
  })

  if (!res.ok || !res.body) {
    let info: unknown = null
    try {
      info = await res.json()
    } catch {
      info = await res.text()
    }
    throw new HttpError(`Chat stream failed with status ${res.status}`, res.status, info)
  }

  await parseSSEStream(res.body.getReader(), onEvent, signal)
}
