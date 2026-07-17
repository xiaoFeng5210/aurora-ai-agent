import { useState } from 'react'
import { Dialog } from '@radix-ui/themes'
import { Check, Copy, X } from 'lucide-react'
import type { Card } from '@/api/card'
import { Markdown } from '@/components/chat/Markdown'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'

export interface CardPreviewProps {
  open: boolean
  card: Card | null
  tagNameById: Map<number, string>
  onOpenChange: (open: boolean) => void
}

export function CardPreview({ open, card, tagNameById, onOpenChange }: CardPreviewProps) {
  const { show } = useToast()
  const [copied, setCopied] = useState(false)

  const resolvedTags =
    card && card.tag_ids && card.tag_ids.length > 0
      ? card.tag_ids.map((id) => tagNameById.get(id)).filter((v): v is string => !!v)
      : (card?.tags ?? [])

  const onCopy = async () => {
    if (!card) return
    try {
      await navigator.clipboard.writeText(formatCardForCopy(card, resolvedTags))
      setCopied(true)
      show('卡片内容已复制', 'success')
      window.setTimeout(() => setCopied(false), 1400)
    } catch {
      show('复制失败', 'error')
    }
  }

  return (
    <Dialog.Root
      open={open && !!card}
      onOpenChange={(next) => {
        onOpenChange(next)
        if (!next) setCopied(false)
      }}
    >
      {card ? (
        <Dialog.Content
          maxWidth="40rem"
          maxHeight="min(88dvh, 44rem)"
          className={cn(
            '!flex !flex-col !overflow-hidden !rounded-2xl !border !border-ink-200/80 !bg-paper-50 !p-0',
            'shadow-[0_24px_80px_-28px_rgba(43,32,22,0.45)]',
          )}
          style={{
            backgroundImage:
              'radial-gradient(circle at 1px 1px, rgba(58,46,31,0.045) 1px, transparent 0)',
            backgroundSize: '18px 18px',
          }}
        >
          <Dialog.Title className="sr-only">{card.title?.trim() || '卡片预览'}</Dialog.Title>
          <Dialog.Description className="sr-only">预览卡片全文，可滚动阅读并复制内容</Dialog.Description>

          <header className="flex shrink-0 items-start justify-between gap-3 border-b border-ink-200/70 px-5 py-4 sm:px-7 sm:py-5">
            <div className="min-w-0 pt-0.5">
              <p className="text-[10px] font-semibold uppercase tracking-[0.22em] text-accent-vermilion/80">
                Preview
              </p>
              <h2 className="font-display mt-1.5 truncate text-xl font-semibold leading-snug text-ink-950 sm:text-2xl">
                {card.title?.trim() || '无标题'}
              </h2>
              <p className="mt-1.5 text-xs text-ink-500">{formatQuietDate(card.created_at)}</p>
            </div>

            <div className="flex shrink-0 items-center gap-1">
              <button
                type="button"
                onClick={onCopy}
                aria-label="复制卡片内容"
                className="inline-flex h-9 items-center gap-1.5 rounded-lg px-2.5 text-sm text-ink-500 transition hover:bg-paper-200/70 hover:text-ink-900"
              >
                {copied ? (
                  <Check className="h-4 w-4 text-accent-vermilion" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
                <span className="hidden sm:inline">{copied ? '已复制' : '复制'}</span>
              </button>
              <Dialog.Close>
                <button
                  type="button"
                  aria-label="关闭预览"
                  className="inline-flex h-9 w-9 items-center justify-center rounded-lg text-ink-400 transition hover:bg-paper-200/70 hover:text-ink-800"
                >
                  <X className="h-4 w-4" />
                </button>
              </Dialog.Close>
            </div>
          </header>

          <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-5 py-5 sm:px-7 sm:py-6">
            <Markdown
              content={card.content}
              className="font-serif text-[15px] leading-[1.85] text-ink-800 sm:text-base sm:leading-[1.9] [&_*]:my-3 [&_*:first-child]:mt-0 [&_*:last-child]:mb-0"
            />

            {resolvedTags.length > 0 ? (
              <div className="mt-8 flex flex-wrap gap-2 border-t border-dashed border-ink-200 pt-5">
                {resolvedTags.map((tag) => (
                  <span
                    key={tag}
                    className="rounded-full border border-ink-200/80 bg-paper-100/60 px-2.5 py-0.5 text-xs text-ink-500"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            ) : null}
          </div>
        </Dialog.Content>
      ) : null}
    </Dialog.Root>
  )
}

function formatQuietDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function formatCardForCopy(card: Card, tags: string[]) {
  const parts: string[] = []

  if (card.title.trim()) {
    parts.push(card.title.trim())
    parts.push('')
  }

  parts.push(card.content)

  if (tags.length > 0) {
    parts.push('')
    parts.push(`标签: ${tags.join(', ')}`)
  }

  const links = [...card.external_links, ...card.internal_links].filter(Boolean)
  if (links.length > 0) {
    parts.push('')
    parts.push('链接:')
    parts.push(...links.map((link) => `- ${link}`))
  }

  return parts.join('\n')
}
