/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import useSWR, { useSWRConfig } from 'swr'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import { Trash2 } from 'lucide-react'
import {
  listHistoryMessages,
  updateMessageFeedback,
  type HistoryMessage,
  type MessageFeedback,
} from '@/api/messages'
import { deleteDocument, getDocument, type Document } from '@/api/documents'
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
  messageId?: string
  feedback?: MessageFeedback
}

const toMessageFeedback = (value: number | undefined): MessageFeedback => {
  if (value === 1 || value === -1) return value
  return 0
}

export function Chat() {
  const navigate = useNavigate()
  const { documentId } = useParams<{ documentId?: string }>()
  const docId = documentId ? Number(documentId) : null
  const { show } = useToast()
  const { mutate } = useSWRConfig()

  const docKey = docId ? `/documents/${docId}` : null
  const { data: docRes } = useSWR<ApiEnvelope<Document> | null>(
    docKey,
    () => (docId ? getDocument(docId) : null),
    { revalidateOnFocus: false },
  )
  const doc = docRes?.data

  const historyKey = docId ? `/documents/${docId}/messages/proxy/history?order=asc&pageSize=200` : null
  const { data: historyRes, isLoading: historyLoading, mutate: mutateHistory } = useSWR<ApiEnvelope<HistoryMessage[]> | null>(
    historyKey,
    () => (docId ? listHistoryMessages(docId, { order: 'asc', pageSize: 200 }) : null),
    { revalidateOnFocus: false, revalidateIfStale: false },
  )
  const history = useMemo(() => historyRes?.data ?? [], [historyRes])

  const [pending, setPending] = useState<DisplayMessage[]>([])
  const [feedbackByMessageId, setFeedbackByMessageId] = useState<Record<string, MessageFeedback>>({})
  const { send, abort, streaming } = useChatStream()

  // 切会话时清空本地消息
  useEffect(() => {
    setPending([])
    setFeedbackByMessageId({})
  }, [docId])

  const messages = useMemo<DisplayMessage[]>(() => {
    const base: DisplayMessage[] = history.map((m) => ({
      key: `db-${m.id}`,
      role: (m.role === 'assistant' ? 'assistant' : 'user') as 'user' | 'assistant',
      content: m.content,
      messageId: m.message_id,
      feedback: feedbackByMessageId[m.message_id] ?? toMessageFeedback(m.is_liked),
    }))
    return [...base, ...pending]
  }, [feedbackByMessageId, history, pending])

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

    // 历史由后端按 document/session 维护，前端不再把整个列表塞进请求体
    // const prompt: ChatPromptItem[] = [
    //   ...messages.map<ChatPromptItem>((m) => ({ role: m.role, content: m.content })),
    //   { role: 'user', content: text },
    // ]
    const prompt: ChatPromptItem[] = [{ role: 'user', content: text }]

    let accumulated = ''
    const toolCalls: ToolCallTrace[] = []
    let completed = false
    let failed = false

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
        // 先保持 pending，等流式请求完全结束后再拿带 message_id 的历史记录替换。
        completed = true
        updateAssistant({ content: final || accumulated, streaming: false })
      },
      onError: (msg) => {
        failed = true
        show(msg, 'error')
        setPending((prev) => prev.filter((m) => m.key !== assistantKey && m.key !== userMsg.key))
      },
      onAbort: () => {
        updateAssistant({ streaming: false })
      },
    })

    if (completed && !failed) {
      try {
        const fresh = await mutateHistory()
        if ((fresh?.data?.length ?? 0) > history.length) {
          setPending((prev) => prev.filter((m) => m.key !== assistantKey && m.key !== userMsg.key))
        }
      } catch {
        // 保留 pending，避免保存后刷新失败导致刚生成的回答从 UI 消失。
      }
    }
  }

  const onFeedback = async (messageId: string, feedback: MessageFeedback) => {
    if (!docId) return
    const previous = feedbackByMessageId[messageId]
    setFeedbackByMessageId((prev) => ({ ...prev, [messageId]: feedback }))
    try {
      await updateMessageFeedback(docId, messageId, feedback)
      await mutateHistory()
    } catch (err) {
      setFeedbackByMessageId((prev) => {
        const next = { ...prev }
        if (previous === undefined) delete next[messageId]
        else next[messageId] = previous
        return next
      })
      show(err instanceof Error ? err.message : '反馈失败', 'error')
    }
  }

  const onDeleteDocument = async () => {
    if (!docId) return
    try {
      abort()
      await deleteDocument(docId)
      await mutate('/documents')
      show('文档已删除', 'success')
      navigate('/chat', { replace: true })
    } catch (err) {
      show(err instanceof Error ? err.message : '删除失败', 'error')
    }
  }

  const showNoDoc = !docId
  const showEmpty = !!docId && !historyLoading && messages.length === 0

  return (
    <div className="grid h-screen grid-cols-[280px_1fr] bg-paper-50">
      <Sidebar onDeleteCurrentDocument={onDeleteDocument} />
      <section className="flex h-screen flex-col">
        <header className="z-10 flex shrink-0 items-center justify-between gap-4 border-b border-ink-200/50 bg-paper-50/80 px-6 py-3 backdrop-blur">
          <div className="min-w-0">
            <h1 className="font-display truncate text-base font-semibold text-ink-950">
              {doc?.display_name ?? (docId ? '加载中…' : 'Aurora')}
            </h1>
            {docId ? <p className="text-[11px] text-ink-500">会话 #{docId}</p> : null}
          </div>
          {docId ? (
            <DocumentDeleteConfirm
              docName={doc?.display_name ?? `会话 #${docId}`}
              onConfirm={onDeleteDocument}
            />
          ) : null}
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
                  messageId={m.messageId}
                  feedback={m.feedback}
                  onFeedback={onFeedback}
                />
              ))}
            </div>
          )}
        </div>

        {
          showNoDoc ? null : (
            <Composer disabled={!docId} streaming={streaming} onSend={onSend} onStop={abort} />
          )
        }
      </section>
    </div>
  )
}

function DocumentDeleteConfirm({
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
          className="inline-flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-md border border-danger/30 px-2.5 text-xs font-medium text-danger transition hover:bg-danger/10"
        >
          <Trash2 className="h-3.5 w-3.5" />
          删除
        </button>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="420px">
        <AlertDialog.Title>删除文档</AlertDialog.Title>
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
              确认删除
            </RTButton>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
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
