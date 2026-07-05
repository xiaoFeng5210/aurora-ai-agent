import { useState } from 'react'
import { Dialog } from '@radix-ui/themes'
import { Eye, Maximize2, Minimize2, Pencil } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Textarea } from '@/components/ui/Textarea'
import { Markdown } from '@/components/chat/Markdown'
import { TagMultiSelect } from './TagMultiSelect'
import {
  createCard,
  updateCard,
  type Card,
  type CreateCardRequest,
  type UpdateCardRequest,
} from '@/api/card'
import type { Tag } from '@/api/tag'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'
import { cn } from '@/lib/cn'

const TITLE_MAX_LENGTH = 100
const CONTENT_MAX_LENGTH = 4000

export interface CardDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  card: Card | null
  tags: Tag[]
  onSaved: (card: Card, mode: 'create' | 'edit') => void
  onTagCreated: (tag: Tag) => void
}

export function CardDialog({ open, onOpenChange, card, tags, onSaved, onTagCreated }: CardDialogProps) {
  const isEdit = !!card
  const [expanded, setExpanded] = useState(false)

  const handleOpenChange = (next: boolean) => {
    onOpenChange(next)
    if (!next) setExpanded(false)
  }

  return (
    <Dialog.Root open={open} onOpenChange={handleOpenChange}>
      <Dialog.Content
        maxWidth={expanded ? undefined : '560px'}
        maxHeight={expanded ? undefined : '85vh'}
        className={cn(
          '!flex !flex-col',
          expanded
            ? '!fixed !inset-0 !m-0 !h-dvh !max-h-dvh !w-screen !max-w-none !rounded-none !p-6 overflow-hidden sm:!p-10'
            : '!rounded-xl',
        )}
      >
        <Dialog.Title className={cn('font-display', expanded && 'sr-only')}>
          {isEdit ? '编辑卡片' : '写一张新卡片'}
        </Dialog.Title>
        <Dialog.Description size="2" className={cn('text-ink-500', expanded && 'sr-only')}>
          {isEdit ? '修改内容与标签，写好后保存即可。' : '记录一段灵感、摘录或想法，贴上标签方便以后检索。'}
        </Dialog.Description>

        {/* key 随卡片切换重新挂载表单，天然重置初始值，无需额外的同步 effect */}
        {open ? (
          <CardDialogForm
            key={card ? `edit-${card.id}` : 'create'}
            card={card}
            tags={tags}
            expanded={expanded}
            onToggleExpanded={() => setExpanded((v) => !v)}
            onCancel={() => handleOpenChange(false)}
            onSaved={(savedCard, mode) => {
              onSaved(savedCard, mode)
              handleOpenChange(false)
            }}
            onTagCreated={onTagCreated}
          />
        ) : null}
      </Dialog.Content>
    </Dialog.Root>
  )
}

function CardDialogForm({
  card,
  tags,
  expanded,
  onToggleExpanded,
  onCancel,
  onSaved,
  onTagCreated,
}: {
  card: Card | null
  tags: Tag[]
  expanded: boolean
  onToggleExpanded: () => void
  onCancel: () => void
  onSaved: (card: Card, mode: 'create' | 'edit') => void
  onTagCreated: (tag: Tag) => void
}) {
  const { show } = useToast()
  const isEdit = !!card

  const [title, setTitle] = useState(card?.title ?? '')
  const [content, setContent] = useState(card?.content ?? '')
  const [contentMode, setContentMode] = useState<'write' | 'preview'>('write')
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>(card?.tag_ids ?? [])
  const [links, setLinks] = useState((card?.external_links ?? []).join('\n'))
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const toggleTag = (id: number) => {
    setSelectedTagIds((prev) => (prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id]))
  }

  const onSubmit = async () => {
    const trimmedTitle = title.trim()
    const trimmed = content.trim()
    if (!trimmed) {
      setError('请填写卡片内容')
      setContentMode('write')
      return
    }
    if (saving) return

    setSaving(true)
    setError('')

    const externalLinks = links
      .split('\n')
      .map((v) => v.trim())
      .filter(Boolean)

    try {
      if (isEdit && card) {
        const body: UpdateCardRequest = {
          title: trimmedTitle,
          content: trimmed,
          tag_ids: selectedTagIds,
          external_links: externalLinks,
        }
        const res = await updateCard(card.id, body)
        if (res.data) onSaved(res.data, 'edit')
        show('卡片已更新', 'success')
      } else {
        const body: CreateCardRequest = {
          title: trimmedTitle,
          content: trimmed,
          tag_ids: selectedTagIds,
          external_links: externalLinks,
        }
        const res = await createCard(body)
        if (res.data) onSaved(res.data, 'create')
        show('卡片已写下', 'success')
      }
    } catch (err) {
      show(getErrorMessage(err, isEdit ? '更新失败' : '创建失败'), 'error')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className={cn('flex flex-1 flex-col', expanded && 'mx-auto min-h-0 w-full max-w-4xl')}>
      <div className={cn('mt-4 flex flex-col gap-5', expanded && 'min-h-0 flex-1 overflow-y-auto pr-1')}>
        <div>
          <Input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="给卡片起个标题（选填）"
            className={cn('font-display', expanded && 'h-11 text-base')}
            maxLength={TITLE_MAX_LENGTH}
          />
        </div>

        <div className={cn(expanded && 'flex min-h-0 flex-1 flex-col')}>
          <div className="mb-2 flex items-center justify-between">
            <p className="text-sm font-medium text-ink-800">内容</p>
            <div className="flex items-center gap-1">
              <button
                type="button"
                onClick={() => setContentMode((m) => (m === 'write' ? 'preview' : 'write'))}
                className="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-ink-500 transition hover:bg-paper-200/70 hover:text-ink-800"
              >
                {contentMode === 'write' ? (
                  <>
                    <Eye className="h-3.5 w-3.5" />
                    预览
                  </>
                ) : (
                  <>
                    <Pencil className="h-3.5 w-3.5" />
                    编辑
                  </>
                )}
              </button>
              <button
                type="button"
                onClick={onToggleExpanded}
                aria-label={expanded ? '缩小编辑区' : '放大编辑区'}
                className="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-ink-500 transition hover:bg-paper-200/70 hover:text-ink-800"
              >
                {expanded ? (
                  <>
                    <Minimize2 className="h-3.5 w-3.5" />
                    缩小
                  </>
                ) : (
                  <>
                    <Maximize2 className="h-3.5 w-3.5" />
                    放大
                  </>
                )}
              </button>
            </div>
          </div>

          {contentMode === 'write' ? (
            <Textarea
              autoFocus
              value={content}
              onChange={(e) => {
                setContent(e.target.value)
                if (error) setError('')
              }}
              placeholder="写下你的想法、摘录或灵感，支持 Markdown 语法…"
              className={cn(
                'text-[15px] leading-relaxed font-mono',
                expanded ? 'min-h-0 flex-1 resize-none' : 'min-h-36',
              )}
              maxLength={CONTENT_MAX_LENGTH}
            />
          ) : (
            <div
              className={cn(
                'w-full overflow-y-auto rounded-md border border-ink-200 bg-paper-50 px-3 py-2',
                expanded ? 'min-h-0 flex-1' : 'min-h-36',
              )}
            >
              {content.trim() ? (
                <Markdown content={content} />
              ) : (
                <p className="text-sm text-ink-300">还没有内容，切换到编辑查看输入…</p>
              )}
            </div>
          )}

          <div className="mt-1 flex items-center justify-between text-xs">
            {error ? <span className="text-danger">{error}</span> : <span />}
            <span className="text-ink-400">
              {content.length}/{CONTENT_MAX_LENGTH}
            </span>
          </div>
        </div>

        <div>
          <p className="mb-2 text-sm font-medium text-ink-800">标签</p>
          <TagMultiSelect
            tags={tags}
            selectedTagIds={selectedTagIds}
            onToggle={toggleTag}
            onTagCreated={onTagCreated}
          />
        </div>

        <div>
          <p className="mb-2 text-sm font-medium text-ink-800">
            参考链接 <span className="text-xs font-normal text-ink-400">（选填，一行一个）</span>
          </p>
          <Textarea
            value={links}
            onChange={(e) => setLinks(e.target.value)}
            placeholder="https://example.com"
            className="min-h-16 text-sm"
          />
        </div>
      </div>

      <div className={cn('mt-6 flex shrink-0 justify-end gap-3', expanded && 'border-t border-ink-100 pt-4')}>
        <Button type="button" variant="ghost" onClick={onCancel}>
          取消
        </Button>
        <Button type="button" onClick={onSubmit} loading={saving}>
          {isEdit ? '保存修改' : '写下这张卡片'}
        </Button>
      </div>
    </div>
  )
}

function getErrorMessage(err: unknown, fallback: string) {
  if (err instanceof HttpError) {
    return extractMessage(err.info) || `${fallback} (${err.status})`
  }
  return err instanceof Error ? err.message : fallback
}

function extractMessage(info: unknown): string | null {
  if (info && typeof info === 'object' && 'message' in info) {
    const m = (info as { message?: unknown }).message
    if (typeof m === 'string') return m
  }
  return null
}
