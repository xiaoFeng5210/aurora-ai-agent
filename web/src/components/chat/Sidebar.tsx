import { Link, useNavigate, useParams } from 'react-router-dom'
import useSWR from 'swr'
import { Music, Plus, Trash2 } from 'lucide-react'
import { createDocument, deleteDocument, listDocuments, type Document } from '@/api/documents'
import type { ApiEnvelope } from '@/api/client'
import { Button } from '@/components/ui/Button'
import { Spinner } from '@/components/ui/Spinner'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'
import { useMemo } from 'react'

const KEY = '/documents'

export function Sidebar() {
  const navigate = useNavigate()
  const { documentId } = useParams<{ documentId?: string }>()
  const currentId = documentId ? Number(documentId) : null
  const { show } = useToast()

  const { data, error, isLoading, mutate } = useSWR<ApiEnvelope<Document[]>>(
    KEY,
    () => listDocuments(),
  )

  const onCreate = async () => {
    try {
      const name = `新会话 · ${new Date().toLocaleString('zh-CN', { hour12: false })}`
      const created = await createDocument({ display_name: name })
      await mutate()
      if (created.data) navigate(`/chat/${created.data.id}`)
    } catch (e) {
      show(e instanceof Error ? e.message : '创建失败', 'error')
    }
  }

  const onDelete = async (e: React.MouseEvent, id: number) => {
    e.stopPropagation()
    e.preventDefault()
    if (!confirm('确认删除这个会话及全部消息？')) return
    try {
      await deleteDocument(id)
      await mutate()
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
        <Link to="/" className="flex items-center gap-2">
          <Music className="h-5 w-5 text-accent-olive" strokeWidth={2.2} />
          <span className="font-display text-lg font-extrabold text-accent-vermilion">Aurora</span>
        </Link>
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
                <li key={doc.id}>
                  <Link
                    to={`/chat/${doc.id}`}
                    className={cn(
                      'group relative flex items-center gap-2 rounded-md px-3 py-2 text-sm transition',
                      active
                        ? 'bg-paper-200/80 text-ink-950'
                        : 'text-ink-700 hover:bg-paper-200/50 hover:text-ink-950',
                    )}
                  >
                    <span className="flex-1 truncate">{doc.display_name}</span>
                    <button
                      type="button"
                      onClick={(e) => onDelete(e, doc.id)}
                      className="hidden rounded p-1 text-ink-500 hover:bg-paper-300/60 hover:text-danger group-hover:inline-flex"
                      aria-label="删除"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </Link>
                </li>
              )
            })}
          </ul>
        )}
      </div>
    </aside>
  )
}
