import { useEffect, useMemo, useState } from 'react'
import useSWR from 'swr'
import useSWRInfinite from 'swr/infinite'
import { Check, PenLine, Plus, Search, Trash2 } from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { Postcard } from '@/components/cards/Postcard'
import { CardDialog } from '@/components/cards/CardDialog'
import { TagChip } from '@/components/cards/TagChip'
import { deleteCard, queryCards, type Card } from '@/api/card'
import { deleteTag, queryTags, type Tag } from '@/api/tag'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'
import { cn } from '@/lib/cn'

const PAGE_SIZE = 12
const SEARCH_DEBOUNCE_MS = 350

type CardsPageKey = readonly ['cards', string, string, number]

export function Cards() {
  const { show } = useToast()

  const [keyword, setKeyword] = useState('')
  const [debouncedKeyword, setDebouncedKeyword] = useState('')
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>([])
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [deletingTagId, setDeletingTagId] = useState<number | null>(null)

  const [dialogOpen, setDialogOpen] = useState(false)
  const [activeCard, setActiveCard] = useState<Card | null>(null)

  useEffect(() => {
    const timer = window.setTimeout(() => setDebouncedKeyword(keyword.trim()), SEARCH_DEBOUNCE_MS)
    return () => window.clearTimeout(timer)
  }, [keyword])

  const { data: tags = [], mutate: mutateTags } = useSWR<Tag[]>(
    'aurora:tags:all',
    () => queryTags({ page: 1, page_size: 200 }).then((res) => res.data ?? []),
    { revalidateOnFocus: false },
  )

  const tagIdsKey = selectedTagIds.slice().sort((a, b) => a - b).join(',')

  const {
    data: cardPages,
    size,
    setSize,
    isLoading: loading,
    isValidating,
    mutate: mutateCards,
  } = useSWRInfinite<Card[], Error, (index: number, previous: Card[] | null) => CardsPageKey | null>(
    (pageIndex, previousPageData) => {
      if (previousPageData && previousPageData.length < PAGE_SIZE) return null
      return ['cards', debouncedKeyword, tagIdsKey, pageIndex + 1]
    },
    ([, contentFilter, tagIdsFilter, page]) =>
      queryCards({
        content: contentFilter || undefined,
        tag_ids: tagIdsFilter ? tagIdsFilter.split(',').map(Number) : undefined,
        page,
        page_size: PAGE_SIZE,
      })
        .then((res) => res.data ?? [])
        .catch((err) => {
          show(getErrorMessage(err, '卡片加载失败'), 'error')
          throw err
        }),
    { revalidateOnFocus: false, revalidateFirstPage: false },
  )

  const cards = useMemo(() => (cardPages ?? []).flat(), [cardPages])
  const hasMore = !!cardPages && (cardPages[cardPages.length - 1]?.length ?? 0) === PAGE_SIZE
  const loadingMore = isValidating && size > 1

  const tagNameById = useMemo(() => {
    const map = new Map<number, string>()
    tags.forEach((t) => map.set(t.id, t.name))
    return map
  }, [tags])

  const toggleTagFilter = (id: number) => {
    setSelectedTagIds((prev) => (prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id]))
  }

  const openCreate = () => {
    setActiveCard(null)
    setDialogOpen(true)
  }

  const openEdit = (card: Card) => {
    setActiveCard(card)
    setDialogOpen(true)
  }

  const onSaved = () => {
    mutateCards()
  }

  const onTagCreated = (tag: Tag) => {
    mutateTags((prev) => (prev?.some((t) => t.id === tag.id) ? prev : [...(prev ?? []), tag]), {
      revalidate: false,
    })
  }

  const onTagDeleted = async (tag: Tag) => {
    setDeletingTagId(tag.id)
    try {
      await deleteTag(tag.id)
      setSelectedTagIds((prev) => prev.filter((id) => id !== tag.id))
      setActiveCard((prev) => (prev ? removeTagFromCard(prev, tag) : prev))
      mutateTags((prev) => prev?.filter((t) => t.id !== tag.id) ?? [], { revalidate: false })
      await mutateCards((prev) => removeTagFromCardPages(prev, tag), { revalidate: true })
      show('标签已删除', 'success')
    } catch (err) {
      show(getErrorMessage(err, '删除标签失败'), 'error')
    } finally {
      setDeletingTagId(null)
    }
  }

  const onDelete = async (card: Card) => {
    setDeletingId(card.id)
    try {
      await deleteCard(card.id)
      show('卡片已删除', 'success')
      await mutateCards()
    } catch (err) {
      show(getErrorMessage(err, '删除失败'), 'error')
    } finally {
      setDeletingId(null)
    }
  }

  const onLoadMore = () => {
    if (loadingMore) return
    setSize(size + 1)
  }

  const filtered = !!debouncedKeyword || selectedTagIds.length > 0

  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <SiteHeader />

      <div className="sticky top-16 z-30 border-b border-ink-200/60 bg-paper-50/90 backdrop-blur">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center gap-3 px-4 py-3 sm:px-6">
          <h1 className="font-display text-lg font-bold text-ink-950 sm:text-xl">我的卡片</h1>

          <div className="ml-auto flex items-center gap-2">
            <Button type="button" onClick={openCreate} size="sm">
              <Plus className="h-4 w-4" />
              <span className="hidden sm:inline">新建卡片</span>
            </Button>
          </div>
        </div>

        <div className="mx-auto flex max-w-6xl items-center gap-3 px-4 pb-3 sm:px-6">
          <div className="relative w-full max-w-xs">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ink-300" />
            <input
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              placeholder="搜索卡片内容…"
              className="w-full rounded-full border border-ink-200 bg-paper-50 py-1.5 pl-9 pr-3 text-sm text-ink-900 placeholder:text-ink-300 transition focus:border-accent-vermilion/60 focus:outline-none focus:ring-2 focus:ring-accent-vermilion/30"
            />
          </div>
        </div>

        {tags.length > 0 ? (
          <div
            className="mx-auto flex max-w-6xl items-center gap-2 overflow-x-auto px-4 pb-3 sm:px-6 [&::-webkit-scrollbar]:hidden"
            style={{ scrollbarWidth: 'none' }}
          >
            <TagChip active={selectedTagIds.length === 0} showCheckbox onClick={() => setSelectedTagIds([])}>
              全部
            </TagChip>
            {tags.map((tag) => (
              <TagFilterPill
                key={tag.id}
                tag={tag}
                active={selectedTagIds.includes(tag.id)}
                onClick={() => toggleTagFilter(tag.id)}
                onDelete={() => onTagDeleted(tag)}
                deleting={deletingTagId === tag.id}
              />
            ))}
          </div>
        ) : null}
      </div>

      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 sm:py-10">
        {loading ? (
          <SkeletonGrid />
        ) : cards.length === 0 ? (
          <EmptyState onCreate={openCreate} filtered={filtered} />
        ) : (
          <>
            <div className="grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-3">
              {cards.map((card, i) => (
                <Postcard
                  key={card.id}
                  card={card}
                  index={i}
                  tagNameById={tagNameById}
                  onOpen={openEdit}
                  onDelete={onDelete}
                  deleting={deletingId === card.id}
                />
              ))}
            </div>

            {hasMore ? (
              <div className="mt-10 flex justify-center">
                <Button type="button" variant="ghost" onClick={onLoadMore} loading={loadingMore}>
                  加载更多
                </Button>
              </div>
            ) : null}
          </>
        )}
      </main>

      <CardDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        card={activeCard}
        tags={tags}
        onSaved={onSaved}
        onTagCreated={onTagCreated}
      />
    </div>
  )
}

function SkeletonGrid() {
  return (
    <div className="grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <div
          key={i}
          className="h-64 animate-pulse rounded-xl border border-ink-200 bg-paper-100"
        />
      ))}
    </div>
  )
}

function EmptyState({ onCreate, filtered }: { onCreate: () => void; filtered: boolean }) {
  return (
    <div className="mx-auto flex max-w-md flex-col items-center rounded-2xl border-2 border-dashed border-ink-200 bg-paper-100/50 px-8 py-16 text-center">
      <span className="flex h-14 w-14 items-center justify-center rounded-full bg-accent-vermilion/10 text-accent-vermilion">
        <PenLine className="h-6 w-6" />
      </span>
      <h2 className="font-display mt-5 text-xl font-bold text-ink-950">
        {filtered ? '没有找到匹配的卡片' : '还没有卡片'}
      </h2>
      <p className="mt-2 text-sm text-ink-500">
        {filtered ? '换个关键词或标签试试看。' : '写下你的第一段灵感、摘录或想法吧。'}
      </p>
      {!filtered ? (
        <Button type="button" onClick={onCreate} className="mt-6">
          <Plus className="h-4 w-4" />
          写第一张卡片
        </Button>
      ) : null}
    </div>
  )
}

function getErrorMessage(err: unknown, fallback: string) {
  if (err instanceof HttpError) {
    return extractMessage(err.info) || `${fallback} (${err.status})`
  }
  return err instanceof Error ? err.message : fallback
}

function TagFilterPill({
  tag,
  active,
  deleting,
  onClick,
  onDelete,
}: {
  tag: Tag
  active: boolean
  deleting: boolean
  onClick: () => void
  onDelete: () => Promise<void> | void
}) {
  return (
    <span
      className={cn(
        'inline-flex shrink-0 items-center overflow-hidden rounded-full border text-xs font-medium transition',
        active
          ? 'border-accent-vermilion bg-accent-vermilion/10 text-accent-vermilion'
          : 'border-ink-200 bg-paper-50 text-ink-600 hover:border-ink-300 hover:text-ink-900',
      )}
    >
      <button
        type="button"
        onClick={onClick}
        className="inline-flex cursor-pointer items-center gap-2 whitespace-nowrap py-1.5 pl-3 pr-2"
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
        {tag.name}
      </button>
      <TagDeleteConfirm tag={tag} deleting={deleting} onConfirm={onDelete} />
    </span>
  )
}

function TagDeleteConfirm({
  tag,
  deleting,
  onConfirm,
}: {
  tag: Tag
  deleting: boolean
  onConfirm: () => Promise<void> | void
}) {
  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <button
          type="button"
          aria-label={`删除标签 ${tag.name}`}
          disabled={deleting}
          className="inline-flex h-7 w-7 cursor-pointer items-center justify-center border-l border-current/10 text-ink-300 transition hover:bg-danger/10 hover:text-danger disabled:cursor-not-allowed disabled:opacity-50"
        >
          <Trash2 className="h-3.5 w-3.5" />
        </button>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="420px">
        <AlertDialog.Title>删除标签</AlertDialog.Title>
        <AlertDialog.Description size="2">
          确认删除「{tag.name}」？删除后会从卡片和筛选条件中移除这个标签关系，且不可撤销。
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

function removeTagFromCard(card: Card, tag: Tag): Card {
  return {
    ...card,
    tag_ids: card.tag_ids?.filter((id) => id !== tag.id),
    tags: card.tags.filter((name) => name !== tag.name),
  }
}

function removeTagFromCardPages(pages: Card[][] | undefined, tag: Tag) {
  return pages?.map((page) => page.map((card) => removeTagFromCard(card, tag)))
}

function extractMessage(info: unknown): string | null {
  if (info && typeof info === 'object' && 'message' in info) {
    const m = (info as { message?: unknown }).message
    if (typeof m === 'string') return m
  }
  return null
}
