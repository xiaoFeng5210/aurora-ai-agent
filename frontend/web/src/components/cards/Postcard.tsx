import { useState, type CSSProperties, type MouseEvent } from 'react'
import { Check, Copy, Pencil, Trash2 } from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import type { Card } from '@/api/card'
import { Markdown } from '@/components/chat/Markdown'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'

const ROTATIONS = ['sm:-rotate-1', 'sm:rotate-1', 'sm:-rotate-2', 'sm:rotate-2', 'sm:rotate-0', 'sm:rotate-1']

const SEAL_GLYPHS = ['笺', '忆', '简', '笔']

// 信纸行线：极淡的横线，与卡片预览同一语言
const LETTER_LINES: CSSProperties = {
  backgroundImage:
    'repeating-linear-gradient(to bottom, transparent 0, transparent 27px, rgba(138, 118, 89, 0.10) 27px, rgba(138, 118, 89, 0.10) 28px)',
}

// 邮票齿孔（小尺寸版）：边缘圆孔 + 中心实心，两层 mask 叠加
const STAMP_HOLES = 'radial-gradient(circle 2px at 3.5px 3.5px, transparent 96%, #000 100%)'
const STAMP_EDGE: CSSProperties = {
  WebkitMaskImage: `${STAMP_HOLES}, linear-gradient(#000 0 0)`,
  WebkitMaskSize: '7px 7px, calc(100% - 6px) calc(100% - 6px)',
  WebkitMaskPosition: '-3.5px -3.5px, 3px 3px',
  WebkitMaskRepeat: 'repeat, no-repeat',
  maskImage: `${STAMP_HOLES}, linear-gradient(#000 0 0)`,
  maskSize: '7px 7px, calc(100% - 6px) calc(100% - 6px)',
  maskPosition: '-3.5px -3.5px, 3px 3px',
  maskRepeat: 'repeat, no-repeat',
}

export interface PostcardProps {
  card: Card
  index: number
  tagNameById: Map<number, string>
  onPreview: (card: Card) => void
  onEdit: (card: Card) => void
  onDelete: (card: Card) => Promise<void> | void
  deleting?: boolean
}

export function Postcard({
  card,
  index,
  tagNameById,
  onPreview,
  onEdit,
  onDelete,
  deleting,
}: PostcardProps) {
  const { show } = useToast()
  const [copied, setCopied] = useState(false)
  const rotate = ROTATIONS[index % ROTATIONS.length]
  const seal = SEAL_GLYPHS[card.id % SEAL_GLYPHS.length]

  const resolvedTags = (card.tag_ids && card.tag_ids.length > 0
    ? card.tag_ids.map((id) => tagNameById.get(id)).filter((v): v is string => !!v)
    : card.tags) ?? []

  const onCopy = async (e: MouseEvent) => {
    e.stopPropagation()
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
    <article
      onClick={() => onPreview(card)}
      className={cn(
        'group relative flex cursor-pointer flex-col overflow-hidden rounded-[14px] border border-ink-200/80 bg-paper-50 shadow-[0_1px_0_0_rgba(58,46,31,0.04),0_10px_26px_-12px_rgba(58,46,31,0.18)] transition-all duration-200 hover:-translate-y-1 hover:rotate-0 hover:shadow-[0_1px_0_0_rgba(58,46,31,0.06),0_20px_38px_-14px_rgba(58,46,31,0.30)]',
        rotate,
      )}
    >
      <div className="flex items-center justify-between gap-2 border-b border-dashed border-ink-200 px-4 pb-2.5 pt-3 sm:px-5">
        <h3
          className="min-w-0 flex-1 truncate font-display text-[15px] font-semibold leading-7 text-ink-950"
          title={card.title || undefined}
        >
          {card.title?.trim() || '无题'}
        </h3>
        <div className="flex shrink-0 items-center gap-0.5 opacity-60 transition group-hover:opacity-100">
          <button
            type="button"
            aria-label="复制卡片内容"
            onClick={onCopy}
            className="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded-md text-ink-400 transition hover:bg-paper-200/70 hover:text-ink-800"
          >
            {copied ? <Check className="h-3.5 w-3.5 text-accent-vermilion" /> : <Copy className="h-3.5 w-3.5" />}
          </button>
          <button
            type="button"
            aria-label="编辑卡片"
            onClick={(e) => {
              e.stopPropagation()
              onEdit(card)
            }}
            className="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded-md text-ink-400 transition hover:bg-paper-200/70 hover:text-ink-800"
          >
            <Pencil className="h-3.5 w-3.5" />
          </button>
          <PostcardDeleteConfirm deleting={!!deleting} onConfirm={() => onDelete(card)} />
        </div>
      </div>

      <div className="flex flex-1 gap-3 px-4 py-4 sm:px-5 sm:py-5">
        <div className="relative min-w-0 flex-1" style={LETTER_LINES}>
          <div
            className="h-32 overflow-hidden sm:h-36"
            style={{
              maskImage: 'linear-gradient(to bottom, black 70%, transparent 100%)',
              WebkitMaskImage: 'linear-gradient(to bottom, black 70%, transparent 100%)',
            }}
          >
            <Markdown
              content={card.content}
              className="font-serif text-[14px] leading-relaxed text-ink-800 sm:text-[15px] [&_*]:my-1.5 [&_*:first-child]:mt-0"
            />
          </div>
        </div>

        <div className="flex w-[68px] shrink-0 flex-col items-center justify-center gap-3 border-l border-dashed border-ink-200 pl-3 sm:w-[72px]">
          {/* 邮票：每张卡片按 id 落一个字 */}
          <div
            aria-hidden
            className="rotate-[4deg]"
            style={{ filter: 'drop-shadow(0 1px 2px rgba(43, 32, 22, 0.18))' }}
          >
            <div className="relative h-[49px] w-[35px] bg-paper-100" style={STAMP_EDGE}>
              <div className="absolute inset-[4px] flex items-center justify-center border border-accent-vermilion/35 bg-paper-50/70">
                <span className="font-display text-sm font-semibold text-accent-vermilion/80">{seal}</span>
              </div>
            </div>
          </div>

          {/* 邮戳 */}
          <div
            aria-hidden
            className="flex h-12 w-12 rotate-[-8deg] items-center justify-center rounded-full border-[1.5px] border-ink-400/40 mix-blend-multiply"
            title={new Date(card.created_at).toLocaleString('zh-CN')}
          >
            <div className="flex h-[38px] w-[38px] flex-col items-center justify-center rounded-full border border-ink-400/25">
              <span className="text-[5px] leading-none tracking-[0.1em] text-ink-500/90">回忆邮局</span>
              <span className="mt-[3px] text-[6px] font-medium leading-none text-ink-500/85">
                {formatPostmark(card.created_at)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </article>
  )
}

function PostcardDeleteConfirm({
  deleting,
  onConfirm,
}: {
  deleting: boolean
  onConfirm: () => void
}) {
  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <button
          type="button"
          aria-label="删除卡片"
          disabled={deleting}
          onClick={(e) => e.stopPropagation()}
          className="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded-md text-ink-400 transition hover:bg-danger/10 hover:text-danger disabled:opacity-50"
        >
          <Trash2 className="h-3.5 w-3.5" />
        </button>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="420px" onClick={(e) => e.stopPropagation()}>
        <AlertDialog.Title>删除卡片</AlertDialog.Title>
        <AlertDialog.Description size="2">
          确认删除这张卡片？删除后会同步移除它的标签关系，且不可撤销。
        </AlertDialog.Description>
        <Flex gap="3" mt="4" justify="end">
          <AlertDialog.Cancel>
            <RTButton variant="soft" color="gray">
              取消
            </RTButton>
          </AlertDialog.Cancel>
          <AlertDialog.Action>
            <RTButton color="tomato" onClick={onConfirm}>
              确认删除
            </RTButton>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
  )
}

function formatPostmark(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '--.--.--'
  const y = date.getFullYear()
  const m = `${date.getMonth() + 1}`.padStart(2, '0')
  const d = `${date.getDate()}`.padStart(2, '0')
  return `${y}.${m}.${d}`
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
