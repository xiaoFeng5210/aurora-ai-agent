import { useState, type ReactNode } from 'react'
import { Check, Copy, Sparkles, ThumbsDown, ThumbsUp } from 'lucide-react'
import { Markdown } from './Markdown'
import { cn } from '@/lib/cn'

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
  messageId?: string
  feedback?: -1 | 0 | 1
  onFeedback?: (messageId: string, feedback: -1 | 0 | 1) => void
}

export function MessageBubble({
  role,
  content,
  streaming,
  toolCalls,
  messageId,
  feedback = 0,
  onFeedback,
}: MessageBubbleProps) {
  const [copied, setCopied] = useState(false)

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
  const canFeedback = !!messageId && !streaming && !!onFeedback
  const onCopy = async () => {
    if (!content) return
    try {
      await navigator.clipboard.writeText(content)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1400)
    } catch {
      // ignore clipboard errors
    }
  }
  const onReact = (value: -1 | 1) => {
    if (!messageId || !onFeedback || streaming) return
    onFeedback(messageId, feedback === value ? 0 : value)
  }

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

        {!streaming || content ? (
          <div className="mt-2 flex h-8 items-center gap-1 opacity-0 transition-opacity duration-150 group-focus-within:opacity-100 group-hover:opacity-100">
            <MessageActionButton
              label={copied ? '已复制' : '复制'}
              onClick={onCopy}
              disabled={!content}
              active={copied}
            >
              {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
            </MessageActionButton>
            <MessageActionButton
              label="点赞"
              onClick={() => onReact(1)}
              disabled={!canFeedback}
              active={feedback === 1}
            >
              <ThumbsUp className="h-3.5 w-3.5" />
            </MessageActionButton>
            <MessageActionButton
              label="点踩"
              onClick={() => onReact(-1)}
              disabled={!canFeedback}
              active={feedback === -1}
            >
              <ThumbsDown className="h-3.5 w-3.5" />
            </MessageActionButton>
          </div>
        ) : null}
      </div>
    </div>
  )
}

function MessageActionButton({
  label,
  active,
  disabled,
  onClick,
  children,
}: {
  label: string
  active?: boolean
  disabled?: boolean
  onClick: () => void
  children: ReactNode
}) {
  return (
    <button
      type="button"
      title={label}
      aria-label={label}
      disabled={disabled}
      onClick={onClick}
      className={cn(
        'inline-flex h-7 w-7 items-center justify-center rounded-md text-ink-500 transition',
        'hover:bg-paper-200/70 hover:text-ink-900 focus-visible:bg-paper-200/70 focus-visible:text-ink-900 focus-visible:outline-none',
        active && 'bg-paper-200/80 text-accent-vermilion',
        disabled && 'cursor-not-allowed opacity-35 hover:bg-transparent hover:text-ink-500',
      )}
    >
      {children}
    </button>
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
