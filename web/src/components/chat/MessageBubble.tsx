import { cn } from '@/lib/cn'
import { Spinner } from '@/components/ui/Spinner'

export interface MessageBubbleProps {
  role: 'user' | 'assistant'
  content: string
  streaming?: boolean
  toolCalls?: Array<{ id: string; name: string; result?: string }>
}

export function MessageBubble({ role, content, streaming, toolCalls }: MessageBubbleProps) {
  const isUser = role === 'user'
  return (
    <div className={cn('flex w-full', isUser ? 'justify-end' : 'justify-start')}>
      <div
        className={cn(
          'max-w-[78%] whitespace-pre-wrap break-words rounded-lg px-4 py-3 text-sm leading-relaxed',
          isUser
            ? 'border border-accent-vermilion/30 bg-accent-vermilion/5 text-ink-950'
            : 'border border-ink-200 bg-paper-50 text-ink-900',
        )}
      >
        {content || (streaming ? <Spinner /> : <span className="text-ink-300">（空消息）</span>)}
        {streaming && content ? (
          <span className="ml-1 inline-block h-3 w-1.5 animate-pulse bg-accent-vermilion align-middle" />
        ) : null}

        {toolCalls && toolCalls.length > 0 ? (
          <div className="mt-3 space-y-1.5 border-t border-ink-200/60 pt-2">
            {toolCalls.map((tc) => (
              <details key={tc.id} className="text-xs text-ink-500">
                <summary className="cursor-pointer text-ink-700">
                  🛠 工具调用: <span className="font-mono">{tc.name}</span>
                </summary>
                {tc.result ? (
                  <pre className="mt-1 max-h-48 overflow-auto rounded bg-paper-200/60 p-2 font-mono text-[11px] text-ink-800">
                    {tc.result}
                  </pre>
                ) : null}
              </details>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  )
}
