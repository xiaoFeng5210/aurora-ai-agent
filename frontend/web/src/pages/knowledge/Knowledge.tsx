import { useRef, useState, type ChangeEvent, type ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'
import useSWR from 'swr'
import {
  AlertTriangle,
  CheckCircle2,
  Circle,
  Clock,
  Database,
  FileText,
  HardDrive,
  RefreshCw,
  Sparkles,
  Trash2,
  UploadCloud,
} from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import {
  RAG_VECTORIZE_API_VERSION,
  RAG_VECTORIZE_MODE_LABEL,
  baiduKnowledgeListDir,
  createRagFromFile,
  deleteKnowledgeFile,
  listKnowledgeFiles,
  uploadKnowledgeFileToNetworkdisk,
  type NetworkdiskFile,
} from '@/api/knowledge'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { cn } from '@/lib/cn'
import { HttpError } from '@/lib/fetcher'

export function Knowledge() {
  const { user } = useAuth()
  const { show } = useToast()
  const navigate = useNavigate()

  const [uploading, setUploading] = useState(false)
  const [deletingPath, setDeletingPath] = useState<string | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const {
    data: files = [],
    isLoading: loading,
    mutate: mutateFiles,
  } = useSWR<NetworkdiskFile[]>(
    ['knowledge:files', user?.username],
    () =>
      listKnowledgeFiles(baiduKnowledgeListDir(user?.username))
        .then((res) => res.data?.list ?? [])
        .catch((err) => {
          show(getErrorMessage(err, '获取知识库文件失败'), 'error')
          throw err
        }),
    { revalidateOnFocus: false },
  )

  const onSelectFile = (e: ChangeEvent<HTMLInputElement>) => {
    setSelectedFile(e.target.files?.[0] ?? null)
  }

  const onUpload = async () => {
    if (!selectedFile || uploading) return

    setUploading(true)
    try {
      const uploadResult = await uploadKnowledgeFileToNetworkdisk(selectedFile)
      await createRagFromFile(selectedFile, uploadResult.data?.path)
      show(
        RAG_VECTORIZE_API_VERSION === 'v2'
          ? '文件已上传，向量化任务已提交'
          : '文件已上传并写入知识库',
        'success',
      )
      setSelectedFile(null)
      if (inputRef.current) inputRef.current.value = ''
      await mutateFiles()
    } catch (err) {
      show(getErrorMessage(err, '上传或向量化失败'), 'error')
    } finally {
      setUploading(false)
    }
  }

  const onDelete = async (file: NetworkdiskFile) => {
    if (deletingPath) return

    setDeletingPath(file.path)
    try {
      await deleteKnowledgeFile(file.path)
      show('文件和向量数据已删除', 'success')
      await mutateFiles()
    } catch (err) {
      show(getErrorMessage(err, '删除失败'), 'error')
    } finally {
      setDeletingPath(null)
    }
  }

  const fileCount = files.filter((file) => file.isdir !== 1).length
  const completedCount = files.filter(
    (file) => file.isdir !== 1 && file.vector_status === 'completed',
  ).length
  const totalSize = files.reduce((sum, file) => sum + (file.isdir === 1 ? 0 : file.size), 0)

  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <SiteHeader />
      <main className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:py-10">
        <div className="border-b border-ink-200/70 pb-6">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-accent-vermilion">
                Knowledge
              </p>
              <div className="mt-2 flex items-center gap-2">
                <Database className="h-6 w-6 text-accent-vermilion" />
                <h1 className="font-display text-3xl font-black leading-tight text-ink-950 sm:text-4xl">
                  知识库
                </h1>
              </div>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-ink-500">
                上传文档后会写入网盘，并通过 {RAG_VECTORIZE_MODE_LABEL} 接口提交向量化，供 AI 助手检索。
              </p>
            </div>
            <div className="grid grid-cols-2 gap-3 text-xs sm:flex sm:items-center">
              <KnowledgeStat icon={<FileText className="h-4 w-4" />} label="文件" value={`${fileCount}`} />
              <KnowledgeStat icon={<HardDrive className="h-4 w-4" />} label="容量" value={formatBytes(totalSize)} />
            </div>
          </div>

          <div className="mt-6 flex flex-col gap-4 rounded-lg border border-accent-vermilion/20 bg-accent-vermilion/5 p-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex min-w-0 items-start gap-3">
              <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-accent-vermilion/10 text-accent-vermilion">
                <Sparkles className="h-5 w-5" />
              </span>
              <div className="min-w-0">
                <p className="text-sm font-semibold text-ink-950">
                  {completedCount > 0
                    ? `已有 ${completedCount} 个文档入库，可与 AI 对话`
                    : '文档入库后，可与 AI 对话'}
                </p>
                <p className="mt-1 text-sm leading-6 text-ink-500">
                  {completedCount > 0
                    ? '已入库的文档会被 AI 助手检索引用，前往对话页即可基于知识库提问。'
                    : '上传文档并完成向量化入库后，AI 助手即可检索知识库内容与您对话。'}
                </p>
              </div>
            </div>
            <Button
              type="button"
              onClick={() => navigate('/chat')}
              className="w-full shrink-0 sm:w-auto"
            >
              <Sparkles className="h-4 w-4" />
              去 AI 对话
            </Button>
          </div>

          <div className="mt-6 grid gap-3 rounded-lg border border-dashed border-ink-300 bg-paper-100/40 p-4 sm:grid-cols-[1fr_auto] sm:items-center">
            <div className="min-w-0">
              <input
                ref={inputRef}
                type="file"
                accept=".txt,.md,.csv,text/plain,text/markdown,text/csv"
                onChange={onSelectFile}
                className="block w-full min-w-0 rounded-md border border-ink-200 bg-paper-50 px-3 py-2 text-sm text-ink-700 file:mr-3 file:rounded-md file:border-0 file:bg-paper-200 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-ink-900 hover:file:bg-paper-300"
              />
            </div>
            <Button
              type="button"
              onClick={onUpload}
              loading={uploading}
              disabled={!selectedFile || uploading}
              className="w-full sm:w-auto"
            >
              <UploadCloud className="h-4 w-4" />
              上传并向量化
            </Button>
          </div>
        </div>

        <div className="mt-6">
          <div className="mb-3 flex items-center justify-between gap-3">
            <p className="text-sm font-semibold text-ink-950">文件列表</p>
            <Button type="button" variant="ghost" size="sm" onClick={() => mutateFiles()} disabled={loading}>
              <RefreshCw className={cn('h-4 w-4', loading && 'animate-spin')} />
              刷新
            </Button>
          </div>

          <div className="overflow-hidden rounded-lg border border-ink-200 bg-paper-50">
            <div className="grid grid-cols-[minmax(0,1fr)_104px_96px_128px_72px] gap-3 border-b border-ink-200 bg-paper-100/70 px-4 py-3 text-xs font-semibold text-ink-500 max-md:hidden">
              <span>名称</span>
              <span>状态</span>
              <span>大小</span>
              <span>更新时间</span>
              <span className="text-right">操作</span>
            </div>

            {loading ? (
              <div className="px-4 py-10 text-center text-sm text-ink-500">正在加载文件...</div>
            ) : files.length === 0 ? (
              <div className="px-4 py-10 text-center text-sm text-ink-500">暂无文件</div>
            ) : (
              <ul className="divide-y divide-ink-200/70">
                {files.map((file) => (
                  <KnowledgeFileRow
                    key={`${file.fs_id}-${file.path}`}
                    file={file}
                    deleting={deletingPath === file.path}
                    onDelete={onDelete}
                  />
                ))}
              </ul>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}

function KnowledgeStat({
  icon,
  label,
  value,
}: {
  icon: ReactNode
  label: string
  value: string
}) {
  return (
    <div className="min-w-0 border-l border-ink-200 pl-3">
      <div className="flex items-center gap-2 text-ink-500">
        {icon}
        <span className="text-xs">{label}</span>
      </div>
      <p className="mt-3 text-sm font-semibold text-ink-950">{value}</p>
    </div>
  )
}

function KnowledgeFileRow({
  file,
  deleting,
  onDelete,
}: {
  file: NetworkdiskFile
  deleting: boolean
  onDelete: (file: NetworkdiskFile) => void
}) {
  const isDirectory = file.isdir === 1

  return (
    <li className="grid gap-3 px-4 py-3 text-sm transition hover:bg-paper-100/50 md:grid-cols-[minmax(0,1fr)_104px_96px_128px_72px] md:items-center">
      <div className="flex min-w-0 items-center gap-3">
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-ink-200 bg-paper-100 text-ink-500">
          <FileText className="h-4 w-4" />
        </span>
        <div className="min-w-0">
          <p className="truncate font-medium text-ink-950">{file.server_filename}</p>
          <p className="mt-0.5 truncate text-xs text-ink-500">{file.path}</p>
          <div className="mt-2 md:hidden">
            <KnowledgeVectorStatusBadge file={file} />
          </div>
        </div>
      </div>
      <div className="hidden md:block">
        <KnowledgeVectorStatusBadge file={file} />
      </div>
      <span className="text-xs text-ink-500 md:text-sm">
        {isDirectory ? '文件夹' : formatBytes(file.size)}
      </span>
      <span className="text-xs text-ink-500 md:text-sm">{formatUnixDate(file.server_mtime)}</span>
      <div className="flex justify-end">
        <KnowledgeFileDeleteConfirm
          fileName={file.server_filename}
          disabled={isDirectory}
          loading={deleting}
          onConfirm={() => onDelete(file)}
        />
      </div>
    </li>
  )
}

function KnowledgeVectorStatusBadge({ file }: { file: NetworkdiskFile }) {
  if (file.isdir === 1) {
    return <StatusBadge status="not_vectorized" label="文件夹" />
  }
  return (
    <StatusBadge
      status={file.vector_status ?? 'not_vectorized'}
      label={file.vector_status_label || '待入库'}
    />
  )
}

function StatusBadge({
  status,
  label,
}: {
  status: NonNullable<NetworkdiskFile['vector_status']> | 'not_vectorized'
  label: string
}) {
  const styles = {
    not_vectorized: 'border-ink-200 bg-paper-100 text-ink-500',
    vectorizing: 'border-accent-vermilion/25 bg-accent-vermilion/10 text-accent-vermilion',
    completed: 'border-emerald-200 bg-emerald-50 text-emerald-700',
    failed: 'border-rose-200 bg-rose-50 text-rose-700',
  }[status]

  const Icon = {
    not_vectorized: Circle,
    vectorizing: Clock,
    completed: CheckCircle2,
    failed: AlertTriangle,
  }[status]

  return (
    <span
      className={cn(
        'inline-flex h-7 max-w-full items-center gap-1.5 rounded-md border px-2 text-xs font-medium',
        styles,
      )}
      title={label}
    >
      <Icon className="h-3.5 w-3.5 shrink-0" />
      <span className="truncate">{label}</span>
    </span>
  )
}

function KnowledgeFileDeleteConfirm({
  fileName,
  disabled,
  loading,
  onConfirm,
}: {
  fileName: string
  disabled: boolean
  loading: boolean
  onConfirm: () => void | Promise<void>
}) {
  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <Button
          type="button"
          variant="danger"
          size="sm"
          loading={loading}
          disabled={loading || disabled}
          title={disabled ? '暂不支持删除文件夹' : '删除文件'}
          className="w-full md:w-auto"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </AlertDialog.Trigger>
      <AlertDialog.Content className="!rounded-lg" maxWidth="420px">
        <AlertDialog.Title>删除知识库文件</AlertDialog.Title>
        <AlertDialog.Description size="2">
          确认删除「{fileName}」及其向量数据？此操作不可撤销。
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

function formatUnixDate(value?: number) {
  if (!value) return '-'
  const date = new Date(value * 1000)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatBytes(value: number) {
  if (!Number.isFinite(value) || value <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = value
  let unit = 0
  while (size >= 1024 && unit < units.length - 1) {
    size /= 1024
    unit += 1
  }
  return `${size >= 10 || unit === 0 ? size.toFixed(0) : size.toFixed(1)} ${units[unit]}`
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
