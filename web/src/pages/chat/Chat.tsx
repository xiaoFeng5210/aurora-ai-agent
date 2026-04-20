/* eslint-disable react-hooks/set-state-in-effect */
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
  const { data: docRes } = useSWR<ApiEnvelope<Document> | null>(
    docKey,
    () => (docId ? getDocument(docId) : null),
    { revalidateOnFocus: false },
  )
  const doc = docRes?.data

  const historyKey = docId ? `/documents/${docId}/messages/proxy/history?order=asc&pageSize=200` : null
  const { data: historyRes, isLoading: historyLoading } = useSWR<ApiEnvelope<HistoryMessage[]> | null>(
    historyKey,
    () => (docId ? listHistoryMessages(docId, { order: 'asc', pageSize: 200 }) : null),
    { revalidateOnFocus: false, revalidateIfStale: false },
  )
  const history = useMemo(() => historyRes?.data ?? [], [historyRes])

  const [pending, setPending] = useState<DisplayMessage[]>([])
  const { send, abort, streaming } = useChatStream()

  // 切会话时清空本地消息
  useEffect(() => {
    setPending([])
  }, [docId])

  const messages = useMemo<DisplayMessage[]>(() => {
    const base: DisplayMessage[] = history.map((m) => ({
      key: `db-${m.id}`,
      role: (m.role === 'assistant' ? 'assistant' : 'user') as 'user' | 'assistant',
      content: m.content,
    }))
    return [...base, ...pending]
  }, [history, pending])

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
      onDone: (final) => {
        // 不再 revalidate history + 清 pending —— 避免后端异步写入造成的空白闪烁
        updateAssistant({ content: final || accumulated, streaming: false })
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

  const showNoDoc = !docId
  const showEmpty = !!docId && !historyLoading && messages.length === 0

  return (
    <div className="grid h-screen grid-cols-[280px_1fr] bg-paper-50">
      <Sidebar />
      <section className="flex h-screen flex-col">
        <header className="z-10 flex shrink-0 items-center justify-between border-b border-ink-200/50 bg-paper-50/80 px-6 py-3 backdrop-blur">
          <div className="min-w-0">
            <h1 className="font-display truncate text-base font-semibold text-ink-950">
              {doc?.display_name ?? (docId ? '加载中…' : 'Aurora')}
            </h1>
            {docId ? <p className="text-[11px] text-ink-500">会话 #{docId}</p> : null}
          </div>
        </header>

        <div ref={scrollRef} className="flex-1 overflow-y-auto">
          {showNoDoc ? (
            <EmptyFull
              title="你好，我是 Aurora"
              subtitle="在左侧新建或选择一个会话开始对话"
            />
          ) : historyLoading ? (
            <div className="flex h-full items-center justify-center">
              <Spinner />
            </div>
          ) : showEmpty ? (
            <EmptyFull title="今天想聊点什么？" subtitle="向 Aurora 提问或让它帮你生成内容" />
          ) : (
            <div className="mx-auto flex max-w-3xl flex-col gap-7 px-6 py-8">
              {messages.map((m) => (
                <MessageBubble
                  key={m.key}
                  role={m.role}
                  content={m.content}
                  streaming={m.streaming}
                  toolCalls={m.toolCalls}
                />
              ))}
            </div>
          )}
        </div>

        <Composer disabled={!docId} streaming={streaming} onSend={onSend} onStop={abort} />
      </section>
    </div>
  )
}

function EmptyFull({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <div className="flex h-full flex-col items-center justify-center px-6 text-center">
      <h2 className="font-display text-4xl font-bold text-ink-950 md:text-5xl">{title}</h2>
      <p className="mt-3 text-sm text-ink-500">{subtitle}</p>
    </div>
  )
}
