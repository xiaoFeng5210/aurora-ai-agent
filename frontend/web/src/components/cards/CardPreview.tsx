import { useState } from 'react'
import type { CSSProperties } from 'react'
import { Dialog } from '@radix-ui/themes'
import { Check, Copy, X } from 'lucide-react'
import type { Card } from '@/api/card'
import { Markdown } from '@/components/chat/Markdown'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'

const GRAIN_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="220" height="220"><filter id="n"><feTurbulence type="fractalNoise" baseFrequency="0.85" numOctaves="2" stitchTiles="stitch"/><feColorMatrix type="saturate" values="0"/></filter><rect width="100%" height="100%" filter="url(#n)"/></svg>`
const GRAIN_TEXTURE = `url("data:image/svg+xml,${encodeURIComponent(GRAIN_SVG)}")`

// 信纸行线：极淡的横线随手写内容一起滚动，像旧笺上的格线
const LETTER_LINES: CSSProperties = {
  backgroundImage:
    'repeating-linear-gradient(to bottom, transparent 0, transparent 35px, rgba(138, 118, 89, 0.13) 35px, rgba(138, 118, 89, 0.13) 36px)',
  backgroundAttachment: 'local',
}

// 邮票齿孔：边缘圆孔 + 中心实心，两层 mask 叠加出打孔效果
const STAMP_HOLES = 'radial-gradient(circle 3px at 4.5px 4.5px, transparent 96%, #000 100%)'
const STAMP_EDGE: CSSProperties = {
  WebkitMaskImage: `${STAMP_HOLES}, linear-gradient(#000 0 0)`,
  WebkitMaskSize: '9px 9px, calc(100% - 8px) calc(100% - 8px)',
  WebkitMaskPosition: '-4.5px -4.5px, 4px 4px',
  WebkitMaskRepeat: 'repeat, no-repeat',
  maskImage: `${STAMP_HOLES}, linear-gradient(#000 0 0)`,
  maskSize: '9px 9px, calc(100% - 8px) calc(100% - 8px)',
  maskPosition: '-4.5px -4.5px, 4px 4px',
  maskRepeat: 'repeat, no-repeat',
}

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

  // 落款印章取标题首字，无标题时落一枚「忆」
  const sealChar = card?.title?.trim().charAt(0) || '忆'

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
          {/* 信纸天头双线 */}
          <div aria-hidden className="pointer-events-none absolute left-5 right-24 top-5 z-0 sm:left-6 sm:right-28 sm:top-6">
            <div className="h-px bg-ink-300/45" />
            <div className="mt-[3px] h-px bg-ink-300/25" />
          </div>
          <span
            aria-hidden
            className="pointer-events-none absolute -bottom-8 -right-5 z-0 select-none font-display text-[150px] font-black leading-none text-ink-900/[0.045] sm:-right-7 sm:text-[170px]"
          >
            忆
          </span>
          {/* 左侧竖排问候 */}
          <span
            aria-hidden
            className="pointer-events-none absolute left-[15px] top-1/2 z-10 -translate-y-1/2 select-none text-[10px] tracking-[0.5em] text-ink-300 sm:left-[19px]"
            style={{ writingMode: 'vertical-rl' }}
          >
            见字如晤
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

          <div
            className="relative z-10 min-h-0 flex-1 overflow-y-auto overscroll-contain pb-9 pl-10 pr-7 pt-12 sm:pb-11 sm:pl-12 sm:pr-10 sm:pt-14"
            style={LETTER_LINES}
          >
            <div className="relative">
              <div className="pr-[62px] sm:pr-[74px]">
                <p className="text-[11px] font-medium tracking-[0.32em] text-ink-500">
                  {formatLetterDate(card.created_at)}
                </p>

                <h2 className="font-display mt-3 text-[26px] font-semibold leading-[1.25] text-ink-950 sm:text-[30px]">
                  {card.title?.trim() || '无标题'}
                </h2>

                <div className="mt-5 flex items-center gap-2.5">
                  <span className="h-px w-10 bg-accent-vermilion/45" />
                  <span className="h-[5px] w-[5px] rotate-45 border border-accent-vermilion/40" />
                  <span className="h-px w-4 bg-accent-vermilion/25" />
                </div>
              </div>

              {/* 邮票，随信头一起滚走 */}
              <div
                aria-hidden
                className="absolute -top-2 right-0 rotate-[5deg] sm:-top-3"
                style={{ filter: 'drop-shadow(0 2px 3px rgba(43, 32, 22, 0.22))' }}
              >
                <div className="relative h-[63px] w-[45px] bg-paper-100 sm:h-[72px] sm:w-[54px]" style={STAMP_EDGE}>
                  <div className="absolute inset-[5px] flex flex-col items-center justify-center gap-[3px] border border-accent-vermilion/40 bg-paper-50/60">
                    <span className="font-display text-base font-semibold leading-none text-accent-vermilion/85 sm:text-lg">
                      忆
                    </span>
                    <span className="text-[6px] tracking-[0.14em] text-ink-500/90">回忆邮局</span>
                  </div>
                </div>
              </div>
            </div>

            <Markdown
              content={card.content}
              className="mt-7 font-serif text-[15.5px] leading-[1.9] text-ink-800 sm:text-[16px] [&_*]:my-4 [&_*:first-child]:mt-0 [&_*:last-child]:mb-0"
            />

            {/* 落款印章 */}
            <div className="mt-9 flex justify-end pr-1">
              <span
                aria-hidden
                className="inline-flex h-9 w-9 -rotate-3 select-none items-center justify-center rounded-[4px] bg-accent-vermilion/90 font-display text-lg font-semibold text-paper-50 shadow-[inset_0_0_0_1.5px_rgba(253,250,241,0.35),0_1px_2px_rgba(43,32,22,0.22)]"
              >
                {sealChar}
              </span>
            </div>

            {resolvedTags.length > 0 ? (
              <div className="mt-8 flex flex-wrap items-center gap-x-4 gap-y-2 border-t border-ink-200/60 pt-6">
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

const CN_DIGITS = ['〇', '一', '二', '三', '四', '五', '六', '七', '八', '九']

function toChineseYear(year: number) {
  return String(year)
    .split('')
    .map((d) => CN_DIGITS[Number(d)])
    .join('')
}

function toChineseNumber(n: number) {
  if (n <= 10) return n === 10 ? '十' : CN_DIGITS[n]
  if (n < 20) return `十${n % 10 === 0 ? '' : CN_DIGITS[n % 10]}`
  return `${CN_DIGITS[Math.floor(n / 10)]}十${n % 10 === 0 ? '' : CN_DIGITS[n % 10]}`
}

// 书信式日期：二〇二六年七月十九日
function formatLetterDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${toChineseYear(date.getFullYear())}年${toChineseNumber(date.getMonth() + 1)}月${toChineseNumber(date.getDate())}日`
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
