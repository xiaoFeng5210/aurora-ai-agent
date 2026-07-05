import { useState, type MouseEvent } from 'react'
import { Check, Copy, Flower2, Pencil, Trash2 } from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import type { Card } from '@/api/card'
import { Markdown } from '@/components/chat/Markdown'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'

const ROTATIONS = ['sm:-rotate-1', 'sm:rotate-1', 'sm:-rotate-2', 'sm:rotate-2', 'sm:rotate-0', 'sm:rotate-1']

const SEAL_GLYPHS = ['笺', '忆', '简', '笔']

export interface PostcardProps {
  card: Card
  index: number
  tagNameById: Map<number, string>
  onOpen: (card: Card) => void
  onDelete: (card: Card) => Promise<void> | void
  deleting?: boolean
}

export function Postcard({ card, index, tagNameById, onOpen, onDelete, deleting }: PostcardProps) {
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
      onClick={() => onOpen(card)}
      className={cn(
        'group relative flex cursor-pointer flex-col overflow-hidden rounded-xl border border-ink-200 bg-paper-50 shadow-[0_1px_0_0_rgba(58,46,31,0.04),0_8px_24px_-12px_rgba(58,46,31,0.15)] transition-all duration-200 hover:-translate-y-1 hover:rotate-0 hover:shadow-[0_1px_0_0_rgba(58,46,31,0.06),0_18px_36px_-14px_rgba(58,46,31,0.28)]',
        rotate,
      )}
      style={{
        backgroundImage:
          'radial-gradient(circle at 1px 1px, rgba(58,46,31,0.06) 1px, transparent 0)',
        backgroundSize: '16px 16px',
      }}
    >
      <div className="flex items-center justify-between gap-2 border-b border-dashed border-ink-200 px-4 pt-3 pb-2 sm:px-5">
        <span className="inline-flex min-w-0 items-center gap-1.5 text-[10px] font-semibold uppercase tracking-[0.2em] text-ink-300">
          <span className="shrink-0 text-accent-vermilion">●</span>
          <span className="truncate text-[16px]" title={card.title || undefined}>
            {card.title ? card.title : '明信片 · POSTCARD'}
          </span>
        </span>
        <div className="flex items-center gap-1 opacity-70 transition group-hover:opacity-100">
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
              onOpen(card)
            }}
            className="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded-md text-ink-400 transition hover:bg-paper-200/70 hover:text-ink-800"
          >
            <Pencil className="h-3.5 w-3.5" />
          </button>
          <PostcardDeleteConfirm deleting={!!deleting} onConfirm={() => onDelete(card)} />
        </div>
      </div>

      <div className="flex flex-1 gap-3 px-4 py-4 sm:px-5 sm:py-5">
        <div className="relative min-w-0 flex-1">
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
          <span
            className="pointer-events-none absolute -bottom-1 right-0 flex h-7 w-7 rotate-[-8deg] items-center justify-center rounded-sm border-2 border-accent-vermilion/70 bg-accent-vermilion/5 font-display text-[13px] font-bold text-accent-vermilion/80"
            aria-hidden
          >
            {seal}
          </span>
        </div>

        <div className="flex w-[76px] shrink-0 flex-col items-center gap-2 border-l border-dashed border-ink-300/80 pl-3 sm:w-20">
          <div className="flex h-14 w-11 shrink-0 items-center justify-center rounded-[2px] border-2 border-dashed border-accent-vermilion/45 bg-accent-vermilion/5 p-1">
            <div className="flex h-full w-full items-center justify-center rounded-[1px] border border-accent-vermilion/30">
              <Flower2 className="h-4 w-4 text-accent-vermilion/70" />
            </div>
          </div>

          <div
            className="flex h-14 w-14 shrink-0 rotate-[-10deg] flex-col items-center justify-center rounded-full border-[1.5px] border-ink-400/60 text-ink-500"
            title={new Date(card.created_at).toLocaleString('zh-CN')}
          >
            <span className="text-[7px] font-semibold uppercase tracking-[0.14em]">Aurora</span>
            <span className="my-0.5 h-px w-6 bg-ink-400/50" />
            <span className="text-[8px] font-medium">{formatPostmark(card.created_at)}</span>
          </div>

          <div className="mt-1 flex w-full flex-col gap-1.5">
            {resolvedTags.length > 0 ? (
              resolvedTags.slice(0, 3).map((tag, i) => (
                <span
                  key={`${tag}-${i}`}
                  className="truncate border-b border-ink-200 pb-0.5 text-center text-[10px] text-ink-500"
                  title={tag}
                >
                  {tag}
                </span>
              ))
            ) : (
              <span className="truncate border-b border-ink-200 pb-0.5 text-center text-[10px] italic text-ink-300">
                无标签
              </span>
            )}
            {resolvedTags.length > 3 ? (
              <span className="text-center text-[10px] text-ink-300">+{resolvedTags.length - 3}</span>
            ) : null}
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
