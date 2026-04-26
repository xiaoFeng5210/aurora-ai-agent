import { useCallback, useRef, useState } from 'react'
import { streamChat, type ChatPromptItem } from '@/api/chat'

export type StreamPhase = 'idle' | 'streaming'

export interface ToolCallTrace {
  id: string
  name: string
  arguments?: string
  result?: string
}

export interface UseChatStreamHandlers {
  onStart?: () => void
  onDelta?: (chunk: string, accumulated: string) => void
  onToolCall?: (calls: ToolCallTrace[]) => void
  onToolResult?: (result: ToolCallTrace) => void
  onDone?: (finalContent: string) => void
  onError?: (message: string) => void
  onAbort?: () => void
}

export function useChatStream() {
  const [phase, setPhase] = useState<StreamPhase>('idle')
  const abortRef = useRef<AbortController | null>(null)
  const handlersRef = useRef<UseChatStreamHandlers>({})

  const send = useCallback(
    async (
      documentId: number,
      prompt: ChatPromptItem[],
      handlers: UseChatStreamHandlers,
      opts?: { maxTokens?: number; temperature?: number; topP?: number; thinking?: 'enabled' | 'disabled' },
    ) => {
      handlersRef.current = handlers
      const controller = new AbortController()
      abortRef.current = controller

      setPhase('streaming')
      let accumulated = ''
      let finished = false

      try {
        await streamChat(
          documentId,
          {
            thinking: { type: opts?.thinking ?? 'disabled' },
            max_tokens: opts?.maxTokens ?? 4096,
            temperature: opts?.temperature ?? 0.7,
            top_p: opts?.topP ?? 0.9,
            tools: [],
            prompt,
          },
          (event, data) => {
            const obj = (data ?? {}) as Record<string, unknown>
            switch (event) {
              case 'start':
                handlers.onStart?.()
                break
              case 'delta': {
                const chunk = typeof obj.content === 'string' ? obj.content : ''
                if (chunk) {
                  accumulated += chunk
                  handlers.onDelta?.(chunk, accumulated)
                }
                break
              }
              case 'tool_call': {
                const list = (obj.tool_calls as Array<Record<string, unknown>> | undefined) ?? []
                const traces: ToolCallTrace[] = list.map((c) => {
                  const fn = (c.function as Record<string, unknown> | undefined) ?? {}
                  return {
                    id: String(c.id ?? ''),
                    name: String(fn.name ?? ''),
                    arguments: typeof fn.arguments === 'string' ? fn.arguments : undefined,
                  }
                })
                handlers.onToolCall?.(traces)
                break
              }
              case 'tool_result':
                handlers.onToolResult?.({
                  id: String(obj.tool_call_id ?? ''),
                  name: String(obj.name ?? ''),
                  result: typeof obj.content === 'string' ? obj.content : JSON.stringify(obj.content),
                })
                break
              case 'error':
                handlers.onError?.(typeof obj.message === 'string' ? obj.message : '未知错误')
                break
              case 'done': {
                finished = true
                const final = typeof obj.content === 'string' ? obj.content : accumulated
                handlers.onDone?.(final)
                break
              }
              default:
                break
            }
          },
          controller.signal,
        )
        if (!finished) handlers.onDone?.(accumulated)
      } catch (err) {
        if (controller.signal.aborted) {
          handlers.onAbort?.()
        } else {
          const msg = err instanceof Error ? err.message : '请求失败'
          handlers.onError?.(msg)
        }
      } finally {
        abortRef.current = null
        setPhase('idle')
      }
    },
    [],
  )

  const abort = useCallback(() => {
    abortRef.current?.abort()
  }, [])

  return { phase, streaming: phase === 'streaming', send, abort }
}
