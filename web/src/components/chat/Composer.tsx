import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { Send, Square } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Textarea } from '@/components/ui/Textarea'

export function Composer({
  disabled,
  streaming,
  onSend,
  onStop,
  placeholder = '给 Aurora 发消息（Enter 发送 / Shift+Enter 换行）',
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
    if (!text || disabled) return
    setValue('')
    onSend(text)
  }

  const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey && !e.nativeEvent.isComposing) {
      e.preventDefault()
      submit()
    }
  }

  return (
    <div className="border-t border-ink-200/60 bg-paper-50">
      <div className="mx-auto flex max-w-3xl items-center gap-2 px-4 py-4">
        <Textarea
          ref={ref}
          rows={1}
          value={value}
          placeholder={placeholder}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={onKeyDown}
          disabled={disabled && !streaming}
          className="max-h-60 min-h-[44px]"
        />
        {streaming && onStop ? (
          <Button variant="danger" size="md" onClick={onStop} type="button">
            <Square className="h-4 w-4" />
            停止
          </Button>
        ) : (
          <Button
            variant="primary"
            size="md"
            onClick={submit}
            disabled={disabled || !value.trim()}
            type="button"
          >
            <Send className="h-4 w-4" />
            {/* 发送 */}
          </Button>
        )}
      </div>
    </div>
  )
}
