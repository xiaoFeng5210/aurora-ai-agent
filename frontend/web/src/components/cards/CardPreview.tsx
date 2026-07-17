import { useState } from 'react'
import { Dialog } from '@radix-ui/themes'
import { Check, Copy, X } from 'lucide-react'
import type { Card } from '@/api/card'
import { Markdown } from '@/components/chat/Markdown'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'

const GRAIN_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="220" height="220"><filter id="n"><feTurbulence type="fractalNoise" baseFrequency="0.85" numOctaves="2" stitchTiles="stitch"/><feColorMatrix type="saturate" values="0"/></filter><rect width="100%" height="100%" filter="url(#n)"/></svg>`
const GRAIN_TEXTURE = `url("data:image/svg+xml,${encodeURIComponent(GRAIN_SVG)}")`

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
          maxWidth="34rem"
          maxHeight="min(86dvh, 46rem)"
          className={cn(
            'relative !flex !flex-col !overflow-hidden !rounded-[22px] !border !border-ink-200/70 !bg-paper-50 !p-0',
            'shadow-[0_2px_0_0_rgba(58,46,31,0.03),0_32px_64px_-24px_rgba(43,32,22,0.38)]',
          )}
        >
          <Dialog.Title className="sr-only">{card.title?.trim() || '卡片预览'}</Dialog.Title>
          <Dialog.Description className="sr-only">预览卡片全文，可滚动阅读并复制内容</Dialog.Description>

          <div
            className="pointer-events-none absolute inset-0 z-0 mix-blend-multiply opacity-[0.05]"
            style={{ backgroundImage: GRAIN_TEXTURE, backgroundSize: '220px 220px' }}
          />
          <span
            aria-hidden
            className="pointer-events-none absolute -bottom-8 -right-5 z-0 select-none font-display text-[150px] font-black leading-none text-ink-900/[0.045] sm:-right-7 sm:text-[170px]"
          >
            忆
          </span>
          <div className="pointer-events-none absolute inset-[10px] z-20 rounded-[15px] border border-ink-200/55 sm:inset-3" />

          <div className="absolute right-3 top-3 z-30 flex items-center gap-0.5 sm:right-4 sm:top-4">
            <button
              type="button"
              onClick={onCopy}
              aria-label="复制卡片内容"
              title="复制"
              className="inline-flex h-9 w-9 items-center justify-center rounded-full text-ink-400 transition hover:bg-paper-200/80 hover:text-ink-800"
            >
              {copied ? (
                <Check className="h-4 w-4 text-accent-vermilion" />
              ) : (
                <Copy className="h-[15px] w-[15px]" />
              )}
            </button>
            <Dialog.Close>
              <button
                type="button"
                aria-label="关闭预览"
                className="inline-flex h-9 w-9 items-center justify-center rounded-full text-ink-400 transition hover:bg-paper-200/80 hover:text-ink-800"
              >
                <X className="h-4 w-4" />
              </button>
            </Dialog.Close>
          </div>

          <div className="relative z-10 min-h-0 flex-1 overflow-y-auto overscroll-contain px-7 pb-9 pt-12 sm:px-10 sm:pb-11 sm:pt-14">
            <p className="text-[11px] font-medium uppercase tracking-[0.2em] text-ink-400">
              {formatQuietDate(card.created_at)}
            </p>

            <h2 className="font-display mt-2.5 text-[26px] font-semibold leading-[1.2] text-ink-950 sm:text-[30px]">
              {card.title?.trim() || '无标题'}
            </h2>

            <div className="mt-5 h-px w-10 bg-accent-vermilion/45" />

            <Markdown
              content={card.content}
              className="mt-7 font-serif text-[15.5px] leading-[1.9] text-ink-800 sm:text-[16px] [&_*]:my-4 [&_*:first-child]:mt-0 [&_*:last-child]:mb-0"
            />

            {resolvedTags.length > 0 ? (
              <div className="mt-10 flex flex-wrap items-center gap-x-4 gap-y-2 border-t border-ink-200/60 pt-6">
                {resolvedTags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1.5 text-[11px] uppercase tracking-[0.08em] text-ink-400"
                  >
                    <span className="h-1 w-1 rounded-full bg-accent-vermilion/50" />
                    {tag}
                  </span>
                ))}
              </div>
            ) : null}
          </div>

          <div className="pointer-events-none absolute inset-x-0 bottom-0 z-20 h-9 bg-gradient-to-t from-paper-50 to-transparent" />
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
