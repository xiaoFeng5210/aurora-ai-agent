import { Sparkles } from 'lucide-react'
import { Markdown } from './Markdown'

export interface ToolCallDisplay {
  id: string
  name: string
  result?: string
}

export interface MessageBubbleProps {
  role: 'user' | 'assistant'
  content: string
  streaming?: boolean
  toolCalls?: ToolCallDisplay[]
}

export function MessageBubble({ role, content, streaming, toolCalls }: MessageBubbleProps) {
  if (role === 'user') {
    return (
      <div className="flex justify-end">
        <div className="max-w-[85%] whitespace-pre-wrap break-words rounded-2xl bg-paper-200/70 px-4 py-2.5 text-[15px] leading-relaxed text-ink-900 shadow-[0_1px_2px_rgba(90,60,30,0.06)]">
          {content}
        </div>
      </div>
    )
  }

  const empty = !content
  return (
    <div className="group flex gap-3">
      <div className="mt-0.5 flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full bg-accent-vermilion/10 text-accent-vermilion ring-1 ring-accent-vermilion/15">
        <Sparkles className="h-3.5 w-3.5" strokeWidth={2.2} />
      </div>
      <div className="min-w-0 flex-1 pt-0.5">
        {empty && streaming ? (
          <div className="text-[15px] leading-[1.75] text-ink-900">
            <TypingDots />
          </div>
        ) : (
          <div className="relative">
            <Markdown content={content} />
            {streaming ? (
              <span className="ml-0.5 inline-block h-[1.05em] w-[2px] translate-y-[2px] animate-pulse bg-ink-500 align-text-bottom" />
            ) : null}
          </div>
        )}

        {toolCalls && toolCalls.length > 0 ? (
          <div className="mt-3 space-y-1.5">
            {toolCalls.map((tc) => (
              <details
                key={tc.id}
                className="rounded-lg border border-ink-200/70 bg-paper-100 px-3 py-2 text-xs text-ink-700"
              >
                <summary className="cursor-pointer select-none text-ink-800">
                  工具调用 · <span className="font-mono text-ink-950">{tc.name}</span>
                </summary>
                {tc.result ? (
                  <pre className="mt-2 max-h-48 overflow-auto rounded bg-paper-200/70 p-2 font-mono text-[11px] leading-snug text-ink-800">
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

function TypingDots() {
  return (
    <span className="inline-flex items-center gap-1 py-1.5">
      {[0, 1, 2].map((i) => (
        <span
          key={i}
          className="h-1.5 w-1.5 animate-bounce rounded-full bg-ink-500"
          style={{ animationDelay: `${i * 150}ms`, animationDuration: '900ms' }}
        />
      ))}
    </span>
  )
}
