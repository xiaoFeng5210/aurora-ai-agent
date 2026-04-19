import { useEffect, useMemo, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import useSWR from 'swr'
import { listHistoryMessages, type HistoryMessage } from '@/api/messages'
import { getDocument, type Document } from '@/api/documents'
import type { ApiEnvelope } from '@/api/client'
import { Sidebar } from '@/components/chat/Sidebar'
import { MessageBubble } from '@/components/chat/MessageBubble'
import { Composer } from '@/components/chat/Composer'
import { Spinner } from '@/components/ui/Spinner'
import { useChatStream, type ToolCallTrace } from '@/hooks/useChatStream'
import { useToast } from '@/hooks/useToast'
import type { ChatPromptItem } from '@/api/chat'

interface DisplayMessage {
  key: string
  role: 'user' | 'assistant'
  content: string
  streaming?: boolean
  toolCalls?: ToolCallTrace[]
}

export function Chat() {
  const { documentId } = useParams<{ documentId?: string }>()
  const docId = documentId ? Number(documentId) : null
  const { show } = useToast()

  const docKey = docId ? `/documents/${docId}` : null
  const { data: docRes } = useSWR<ApiEnvelope<Document> | null>(docKey, () =>
    docId ? getDocument(docId) : null,
  )
  const doc = docRes?.data

  const historyKey = docId ? `/documents/${docId}/messages/proxy/history?order=asc&pageSize=200` : null
  const {
    data: historyRes,
    isLoading: historyLoading,
    mutate: mutateHistory,
  } = useSWR<ApiEnvelope<HistoryMessage[]> | null>(historyKey, () =>
    docId ? listHistoryMessages(docId, { order: 'asc', pageSize: 200 }) : null,
  )
  const history = historyRes?.data

  const [pending, setPending] = useState<DisplayMessage[]>([])
  const { send, abort, streaming } = useChatStream()

  // 切换会话时清空临时态
  useEffect(() => {
    const clear = () => {
      setPending([])
    }
    clear()
  }, [docId])

  const messages = useMemo<DisplayMessage[]>(() => {
    const base: DisplayMessage[] =
      history?.map((m) => ({
        key: `db-${m.id}`,
        role: (m.role === 'assistant' ? 'assistant' : 'user') as 'user' | 'assistant',
        content: m.content,
      })) ?? []
    return [...base, ...pending]
  }, [historyRes, pending])

  // 自动滚到底
  const scrollRef = useRef<HTMLDivElement>(null)
  useEffect(() => {
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [messages])

  const onSend = async (text: string) => {
    if (!docId) {
      show('请先在左侧新建一个会话', 'info')
      return
    }
    if (streaming) return

    const userMsg: DisplayMessage = { key: `tmp-u-${Date.now()}`, role: 'user', content: text }
    const assistantKey = `tmp-a-${Date.now()}`
    const assistantMsg: DisplayMessage = {
      key: assistantKey,
      role: 'assistant',
      content: '',
      streaming: true,
      toolCalls: [],
    }
    setPending((prev) => [...prev, userMsg, assistantMsg])

    const updateAssistant = (patch: Partial<DisplayMessage>) =>
      setPending((prev) => prev.map((m) => (m.key === assistantKey ? { ...m, ...patch } : m)))

    const prompt: ChatPromptItem[] = [
      ...messages.map<ChatPromptItem>((m) => ({ role: m.role, content: m.content })),
      { role: 'user', content: text },
    ]

    let accumulated = ''
    const toolCalls: ToolCallTrace[] = []

    await send(docId, prompt, {
      onDelta: (_chunk, acc) => {
        accumulated = acc
        updateAssistant({ content: acc })
      },
      onToolCall: (calls) => {
        for (const c of calls) {
          if (!toolCalls.find((t) => t.id === c.id)) toolCalls.push(c)
        }
        updateAssistant({ toolCalls: [...toolCalls] })
      },
      onToolResult: (r) => {
        const exist = toolCalls.find((t) => t.id === r.id)
        if (exist) Object.assign(exist, r)
        else toolCalls.push(r)
        updateAssistant({ toolCalls: [...toolCalls] })
      },
      onDone: async (final) => {
        updateAssistant({ content: final || accumulated, streaming: false })
        // 拉一次最新历史，再清掉乐观态
        try {
          await mutateHistory()
        } catch {
          // ignore
        }
        setPending([])
      },
      onError: (msg) => {
        show(msg, 'error')
        setPending((prev) => prev.filter((m) => m.key !== assistantKey && m.key !== userMsg.key))
      },
      onAbort: () => {
        updateAssistant({ streaming: false })
      },
    })
  }

  return (
    <div className="grid h-screen grid-cols-[280px_1fr] bg-paper-100">
      <Sidebar />
      <section className="flex h-screen flex-col bg-paper-50">
        <header className="flex items-center justify-between border-b border-ink-200/60 bg-paper-50 px-6 py-4">
          <div>
            <h1 className="font-display text-lg font-bold text-ink-950">
              {doc?.display_name ?? (docId ? '加载中…' : '未选择会话')}
            </h1>
            <p className="text-xs text-ink-500">
              {docId ? `会话 #${docId}` : '从左侧选择或新建一个会话开始对话'}
            </p>
          </div>
        </header>

        <div ref={scrollRef} className="flex-1 overflow-y-auto px-6 py-6">
          <div className="mx-auto flex max-w-3xl flex-col gap-4">
            {!docId ? (
              <EmptyState text="点击左侧 + 新建会话 开始与 Aurora 对话" />
            ) : historyLoading ? (
              <div className="flex justify-center py-10"><Spinner /></div>
            ) : messages.length === 0 ? (
              <EmptyState text="发送第一条消息开始这个会话" />
            ) : (
              messages.map((m) => (
                <MessageBubble
                  key={m.key}
                  role={m.role}
                  content={m.content}
                  streaming={m.streaming}
                  toolCalls={m.toolCalls}
                />
              ))
            )}
          </div>
        </div>

        <Composer
          disabled={!docId}
          streaming={streaming}
          onSend={onSend}
          onStop={abort}
        />
      </section>
    </div>
  )
}

function EmptyState({ text }: { text: string }) {
  return (
    <div className="mt-20 text-center text-sm text-ink-500">{text}</div>
  )
}
