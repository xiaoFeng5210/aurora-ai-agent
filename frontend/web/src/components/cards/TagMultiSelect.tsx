import { useState } from 'react'
import { Popover } from '@radix-ui/themes'
import { Check, ChevronDown, Plus, X } from 'lucide-react'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'
import { createTag, type Tag } from '@/api/tag'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'
import { cn } from '@/lib/cn'

export interface TagMultiSelectProps {
  tags: Tag[]
  selectedTagIds: number[]
  onToggle: (id: number) => void
  onTagCreated: (tag: Tag) => void
}

export function TagMultiSelect({ tags, selectedTagIds, onToggle, onTagCreated }: TagMultiSelectProps) {
  const { show } = useToast()
  const [open, setOpen] = useState(false)
  const [newTagName, setNewTagName] = useState('')
  const [creatingTag, setCreatingTag] = useState(false)

  const selectedTags = tags.filter((tag) => selectedTagIds.includes(tag.id))

  const onCreateTag = async () => {
    const name = newTagName.trim()
    if (!name || creatingTag) return

    setCreatingTag(true)
    try {
      const res = await createTag({ name })
      if (res.data) {
        onTagCreated(res.data)
        if (!selectedTagIds.includes(res.data.id)) onToggle(res.data.id)
      }
      setNewTagName('')
    } catch (err) {
      show(getErrorMessage(err, '创建标签失败'), 'error')
    } finally {
      setCreatingTag(false)
    }
  }

  return (
    <Popover.Root open={open} onOpenChange={setOpen}>
      <Popover.Trigger>
        <button
          type="button"
          className={cn(
            'flex min-h-[2.375rem] w-full flex-wrap items-center gap-1.5 rounded-md border bg-paper-50 px-3 py-1.5 text-left text-sm transition focus:outline-none focus:ring-2 focus:ring-accent-vermilion/30',
            open ? 'border-accent-vermilion/60' : 'border-ink-200 hover:border-ink-300',
          )}
        >
          {selectedTags.length > 0 ? (
            selectedTags.map((tag) => (
              <span
                key={tag.id}
                className="inline-flex items-center gap-1 rounded-full bg-accent-vermilion/10 py-0.5 pl-2 pr-1 text-xs font-medium text-accent-vermilion"
              >
                {tag.name}
                <span
                  role="button"
                  aria-label={`移除标签 ${tag.name}`}
                  tabIndex={-1}
                  onClick={(e) => {
                    e.stopPropagation()
                    onToggle(tag.id)
                  }}
                  className="rounded-full p-0.5 transition hover:bg-accent-vermilion/20"
                >
                  <X className="h-2.5 w-2.5" />
                </span>
              </span>
            ))
          ) : (
            <span className="text-ink-400">选择标签…</span>
          )}
          <ChevronDown
            className={cn(
              'ml-auto h-4 w-4 shrink-0 text-ink-400 transition-transform',
              open && 'rotate-180',
            )}
          />
        </button>
      </Popover.Trigger>

      <Popover.Content width="280px" className="!rounded-lg !p-2">
        {tags.length > 0 ? (
          <div className="max-h-52 overflow-y-auto">
            {tags.map((tag) => {
              const active = selectedTagIds.includes(tag.id)
              return (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => onToggle(tag.id)}
                  className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm text-ink-800 transition hover:bg-paper-200/70"
                >
                  <span
                    className={cn(
                      'inline-flex h-3.5 w-3.5 shrink-0 items-center justify-center rounded-[3px] border transition',
                      active
                        ? 'border-accent-vermilion bg-accent-vermilion text-paper-50'
                        : 'border-ink-300 bg-paper-50',
                    )}
                    aria-hidden
                  >
                    {active ? <Check className="h-2.5 w-2.5" strokeWidth={3} /> : null}
                  </span>
                  <span className="truncate">{tag.name}</span>
                </button>
              )
            })}
          </div>
        ) : (
          <p className="px-2 py-1.5 text-xs text-ink-400">还没有标签，创建一个吧</p>
        )}

        <div className="mt-1 flex items-center gap-1.5 border-t border-ink-100 pt-2">
          <Input
            value={newTagName}
            onChange={(e) => setNewTagName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key !== 'Enter') return
              e.preventDefault()
              void onCreateTag()
            }}
            placeholder="新建标签"
            className="h-8 flex-1 text-xs"
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onCreateTag}
            loading={creatingTag}
            disabled={!newTagName.trim()}
            className="h-8 px-2.5"
          >
            <Plus className="h-3.5 w-3.5" />
          </Button>
        </div>
      </Popover.Content>
    </Popover.Root>
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
