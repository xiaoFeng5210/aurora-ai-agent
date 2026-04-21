import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { ArrowUp, Square } from 'lucide-react'
import { cn } from '@/lib/cn'

export function Composer({
  disabled,
  streaming,
  onSend,
  onStop,
  placeholder = '给 Aurora 发消息…',
}: {
  disabled?: boolean
  streaming?: boolean
  onSend: (text: string) => void
  onStop?: () => void
  placeholder?: string
}) {
  const [value, setValue] = useState('')
  const ref = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 240) + 'px'
  }, [value])

  const submit = () => {
    const text = value.trim()
    if (!text || disabled || streaming) return
    setValue('')
    onSend(text)
  }

  const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey && !e.nativeEvent.isComposing) {
      e.preventDefault()
      submit()
    }
  }

  const canSend = !!value.trim() && !disabled && !streaming

  return (
    <div className="bg-linear-to-b from-transparent via-paper-50/70 to-paper-50 px-4 pb-6">
      <div className="mx-auto w-full max-w-3xl">
        <div className="relative rounded-3xl border border-ink-200 bg-paper-50 ">
          <textarea
            ref={ref}
            rows={1}
            value={value}
            placeholder={placeholder}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={onKeyDown}
            disabled={disabled && !streaming}
            className="block w-full resize-none rounded-3xl bg-transparent px-5 pb-14 pt-4 text-[15px] leading-relaxed text-ink-900 placeholder:text-ink-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-60"
            style={{ minHeight: 56 }}
          />
          <div className="absolute bottom-2.5 right-2.5 flex items-center gap-2">
            {streaming && onStop ? (
              <button
                type="button"
                onClick={onStop}
                className="flex h-9 w-9 items-center justify-center rounded-full bg-ink-950 text-paper-50 transition hover:bg-ink-900"
                aria-label="停止生成"
              >
                <Square className="h-3 w-3 fill-paper-50" />
              </button>
            ) : (
              <button
                type="button"
                onClick={submit}
                disabled={!canSend}
                className={cn(
                  'flex h-9 w-9 items-center justify-center rounded-full transition',
                  canSend
                    ? 'bg-accent-vermilion text-paper-50 hover:bg-accent-vermilion/90 shadow-sm'
                    : 'bg-ink-200 text-ink-500',
                )}
                aria-label="发送"
              >
                <ArrowUp className="h-4 w-4" strokeWidth={2.4} />
              </button>
            )}
          </div>
        </div>
        <p className="mt-2 text-center text-[11px] text-ink-500">
          Aurora 可能会出错 · Enter 发送 · Shift + Enter 换行
        </p>
      </div>
    </div>
  )
}
