import { useNavigate, useParams } from 'react-router-dom'
import useSWR from 'swr'
import { Database, LayoutGrid, Plus, Trash2 } from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import { createDocument, deleteDocument, listDocuments, type Document } from '@/api/documents'
import type { ApiEnvelope } from '@/api/client'
import { Button } from '@/components/ui/Button'
import { Spinner } from '@/components/ui/Spinner'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'
import { useMemo } from 'react'
import logoUrl from '@/assets/logo.png'

const KEY = '/documents'

const NAV_ITEMS = [
  { label: '卡片', href: '/cards', icon: LayoutGrid },
  { label: '知识库', href: '/knowledge', icon: Database },
]

export function Sidebar({
  onDeleteCurrentDocument,
  onNavigate,
}: {
  onDeleteCurrentDocument?: (id: number) => void | Promise<void>
  onNavigate?: () => void
}) {
  const navigate = useNavigate()
  const { documentId } = useParams<{ documentId?: string }>()
  const currentId = documentId ? Number(documentId) : null
  const { show } = useToast()

  const go = (to: string) => () => {
    navigate(to)
    onNavigate?.()
  }

  const { data, error, isLoading, mutate } = useSWR<ApiEnvelope<Document[]>>(
    KEY,
    () => listDocuments(),
  )

  const onCreate = async () => {
    try {
      const name = `新对话 · ${new Date().toLocaleString('zh-CN', { hour12: false })}`
      const created = await createDocument({ display_name: name })
      await mutate()
      if (created.data) go(`/chat/${created.data.id}`)()
    } catch (e) {
      show(e instanceof Error ? e.message : '创建失败', 'error')
    }
  }

  const onDelete = async (id: number) => {
    try {
      if (currentId === id && onDeleteCurrentDocument) {
        await onDeleteCurrentDocument(id)
        await mutate()
        return
      }
      await deleteDocument(id)
      await mutate()
      show('文档已删除', 'success')
      if (currentId === id) navigate('/chat')
    } catch (err) {
      show(err instanceof Error ? err.message : '删除失败', 'error')
    }
  }

  const sorted = useMemo(() => {
    const list = data?.data ?? []
    return list
      .slice()
      .sort((a, b) => +new Date(b.updated_at) - +new Date(a.updated_at))
  }, [data])

  return (
    <aside className="flex h-screen w-[280px] flex-col border-r border-ink-200/80 bg-paper-100">
      <div className="flex items-center justify-between border-b border-ink-200/60 px-4 py-4">
        <div
          role="link"
          tabIndex={0}
          onClick={go('/')}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') go('/')()
          }}
          className="flex w-fit cursor-pointer items-center gap-2 select-none"
        >
          <img src={logoUrl} alt="Ariadne" className="h-7 w-7" />
          <span className="font-display text-lg font-extrabold text-accent-vermilion">Ariadne</span>
        </div>
      </div>

      <div className="px-3 pt-3">
        <Button size="sm" variant="primary" className="w-full" onClick={onCreate}>
          <Plus className="h-4 w-4" />
          新建会话
        </Button>
      </div>

      <div className="mt-3 flex-1 overflow-y-auto px-2 pb-3">
        {isLoading ? (
          <div className="flex items-center justify-center py-6"><Spinner /></div>
        ) : error ? (
          <p className="px-2 py-3 text-xs text-danger">加载失败</p>
        ) : sorted.length === 0 ? (
          <p className="px-2 py-3 text-xs text-ink-500">还没有会话，点上方按钮开始</p>
        ) : (
          <ul className="flex flex-col gap-0.5">
            {sorted.map((doc) => {
              const active = currentId === doc.id
              return (
                <li key={doc.id} className="group relative">
                  <div
                    role="link"
                    tabIndex={0}
                    onClick={go(`/chat/${doc.id}`)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') go(`/chat/${doc.id}`)()
                    }}
                    className={cn(
                      'flex cursor-pointer items-center gap-2 rounded-md py-2 pl-3 pr-9 text-sm transition',
                      active
                        ? 'bg-paper-200/80 text-ink-950'
                        : 'text-ink-700 hover:bg-paper-200/50 hover:text-ink-950',
                    )}
                  >
                    <span className="flex-1 truncate">{doc.display_name}</span>
                  </div>
                  <DeleteConfirm
                    docName={doc.display_name}
                    onConfirm={() => onDelete(doc.id)}
                  />
                </li>
              )
            })}
          </ul>
        )}
      </div>

      <nav className="shrink-0 border-t border-ink-200/60 px-2 py-2" aria-label="页面导航">
        <p className="px-3 pb-1 text-[11px] font-medium tracking-wide text-ink-500">导航</p>
        <ul className="flex flex-col gap-0.5">
          {NAV_ITEMS.map((item) => (
            <li key={item.href}>
              <button
                type="button"
                onClick={go(item.href)}
                className="flex w-full cursor-pointer items-center gap-2.5 rounded-md px-3 py-2 text-sm text-ink-700 transition hover:bg-paper-200/50 hover:text-ink-950"
              >
                <item.icon className="h-4 w-4 text-ink-500" />
                {item.label}
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  )
}

function DeleteConfirm({
  docName,
  onConfirm,
}: {
  docName: string
  onConfirm: () => void | Promise<void>
}) {
  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <button
          type="button"
          className="absolute right-2 top-1/2 hidden -translate-y-1/2 rounded p-1 text-ink-500 hover:bg-paper-300/60 hover:text-danger group-hover:inline-flex"
          aria-label="删除"
        >
          <Trash2 className="h-3.5 w-3.5" />
        </button>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="420px">
        <AlertDialog.Title>删除会话</AlertDialog.Title>
        <AlertDialog.Description size="2">
          确认删除「{docName}」及其全部消息？此操作不可撤销。
        </AlertDialog.Description>
        <Flex gap="3" mt="4" justify="end">
          <AlertDialog.Cancel>
            <RTButton variant="soft" color="gray">
              取消
            </RTButton>
          </AlertDialog.Cancel>
          <AlertDialog.Action>
            <RTButton color="tomato" onClick={() => onConfirm()}>
              删除
            </RTButton>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
  )
}
