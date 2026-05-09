import {
  useCallback,
  useMemo,
  useRef,
  useState,
  type ChangeEvent,
  type FormEvent,
  type RefObject,
  type ReactNode,
} from 'react'
import {
  Database,
  FileText,
  HardDrive,
  RefreshCw,
  Save,
  Settings2,
  Sparkles,
  Trash2,
  UploadCloud,
  UserRound,
} from 'lucide-react'
import { AlertDialog, Button as RTButton, Flex } from '@radix-ui/themes'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Textarea } from '@/components/ui/Textarea'
import { updateMe, type User } from '@/api/auth'
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

type ProfileTab = 'profile' | 'knowledge'

interface ProfileForm {
  username: string
  email: string
  phone: string
  birthday: string
  user_prompt: string
}

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const emptyForm: ProfileForm = {
  username: '',
  email: '',
  phone: '',
  birthday: '',
  user_prompt: '',
}

const formFromUser = (user: User | null): ProfileForm => {
  if (!user) return emptyForm
  return {
    username: user.username ?? '',
    email: user.email ?? '',
    phone: user.phone ?? '',
    birthday: user.birthday ?? '',
    user_prompt: user.user_prompt ?? '',
  }
}

export function Profile() {
  const { user, setUser } = useAuth()
  const { show } = useToast()
  const [activeTab, setActiveTab] = useState<ProfileTab>('profile')
  const [form, setForm] = useState<ProfileForm>(() => formFromUser(user))
  const [errors, setErrors] = useState<Partial<Record<keyof ProfileForm, string>>>({})
  const [saving, setSaving] = useState(false)
  const [knowledgeFiles, setKnowledgeFiles] = useState<NetworkdiskFile[]>([])
  const [knowledgeLoading, setKnowledgeLoading] = useState(false)
  const [uploadingKnowledge, setUploadingKnowledge] = useState(false)
  const [deletingKnowledgePath, setDeletingKnowledgePath] = useState<string | null>(null)
  const [selectedKnowledgeFile, setSelectedKnowledgeFile] = useState<File | null>(null)
  const knowledgeInputRef = useRef<HTMLInputElement>(null)

  const initials = useMemo(() => {
    const name = form.username.trim() || user?.email || 'A'
    return name.slice(0, 1).toUpperCase()
  }, [form.username, user?.email])

  const dirty = user ? !sameForm(form, formFromUser(user)) : false
  const promptLength = form.user_prompt.length

  const set = (key: keyof ProfileForm) => (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const value = e.target.value
    setForm((prev) => ({ ...prev, [key]: value }))
    setErrors((prev) => ({ ...prev, [key]: undefined }))
  }

  const loadKnowledgeFiles = useCallback(async () => {
    setKnowledgeLoading(true)
    try {
      const res = await listKnowledgeFiles(baiduKnowledgeListDir(user?.username))
      setKnowledgeFiles(res.data?.list ?? [])
    } catch (err) {
      const msg = getErrorMessage(err, '获取知识库文件失败')
      show(msg, 'error')
    } finally {
      setKnowledgeLoading(false)
    }
  }, [show, user?.username])

  const onSelectKnowledgeFile = (e: ChangeEvent<HTMLInputElement>) => {
    setSelectedKnowledgeFile(e.target.files?.[0] ?? null)
  }

  const onUploadKnowledgeFile = async () => {
    if (!selectedKnowledgeFile || uploadingKnowledge) return

    setUploadingKnowledge(true)
    try {
      await uploadKnowledgeFileToNetworkdisk(selectedKnowledgeFile)
      await createRagFromFile(selectedKnowledgeFile)
      show(
        RAG_VECTORIZE_API_VERSION === 'v2'
          ? '文件已上传，向量化任务已提交'
          : '文件已上传并写入知识库',
        'success',
      )
      setSelectedKnowledgeFile(null)
      if (knowledgeInputRef.current) knowledgeInputRef.current.value = ''
      await loadKnowledgeFiles()
    } catch (err) {
      const msg = getErrorMessage(err, '上传或向量化失败')
      show(msg, 'error')
    } finally {
      setUploadingKnowledge(false)
    }
  }

  const onDeleteKnowledgeFile = async (file: NetworkdiskFile) => {
    if (deletingKnowledgePath) return

    setDeletingKnowledgePath(file.path)
    try {
      await deleteKnowledgeFile(file.path)
      show('文件和向量数据已删除', 'success')
      await loadKnowledgeFiles()
    } catch (err) {
      const msg = getErrorMessage(err, '删除失败')
      show(msg, 'error')
    } finally {
      setDeletingKnowledgePath(null)
    }
  }

  const reset = () => {
    setForm(formFromUser(user))
    setErrors({})
  }

  const validate = () => {
    const next: Partial<Record<keyof ProfileForm, string>> = {}
    if (form.username.trim().length < 2) next.username = '用户名至少 2 位'
    if (!EMAIL_RE.test(form.email.trim())) next.email = '邮箱格式不正确'
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!user || !dirty || saving) return
    if (!validate()) return

    setSaving(true)
    try {
      const res = await updateMe({
        username: form.username.trim(),
        email: form.email.trim(),
        phone: form.phone.trim(),
        birthday: form.birthday,
        user_prompt: form.user_prompt,
      })
      if (res.data) {
        setUser(res.data)
        setForm(formFromUser(res.data))
      }
      show('个人信息已保存', 'success')
    } catch (err) {
      const msg =
        err instanceof HttpError
          ? extractMessage(err.info) || `保存失败 (${err.status})`
          : err instanceof Error
            ? err.message
            : '保存失败'
      show(msg, 'error')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <SiteHeader />
      <main className="mx-auto grid max-w-7xl gap-7 px-4 py-6 sm:px-6 lg:grid-cols-[260px_1fr] lg:gap-8 lg:py-10">
        <aside className="lg:sticky lg:top-24 lg:h-fit">
          <div className="border-b border-ink-200/70 pb-5 sm:pb-6">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-ink-950 text-lg font-bold text-paper-50 shadow-sm sm:h-12 sm:w-12 sm:text-xl">
                {initials}
              </div>
              <div className="min-w-0">
                <p className="truncate text-base font-semibold text-ink-950">
                  {form.username || 'Aurora 用户'}
                </p>
                <p className="truncate text-xs text-ink-500">{form.email || '未设置邮箱'}</p>
              </div>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3 text-xs sm:mt-5">
              <AccountFact label="账户编号" value={user ? `#${user.id}` : '-'} />
              <AccountFact label="最后更新" value={formatDate(user?.updated_at)} />
            </div>
          </div>

          <nav className="mt-4 grid grid-cols-2 gap-2 sm:mt-5 lg:flex lg:flex-col">
            <NavItem
              active={activeTab === 'profile'}
              icon={<UserRound className="h-4 w-4" />}
              label="个人信息"
              description="资料与 User Prompt"
              onClick={() => setActiveTab('profile')}
            />
            <NavItem
              active={activeTab === 'knowledge'}
              icon={<Database className="h-4 w-4" />}
              label="知识库"
              description="文档与索引"
              onClick={() => {
                setActiveTab('knowledge')
                void loadKnowledgeFiles()
              }}
            />
          </nav>
        </aside>

        <section className="min-w-0">
          {activeTab === 'profile' ? (
            <div className="border-b border-ink-200/70 pb-5 sm:pb-6">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-accent-vermilion">
                Account
              </p>
              <div className="mt-3 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
                <div>
                  <h1 className="font-display text-3xl font-black leading-tight text-ink-950 sm:text-4xl md:text-5xl">
                    个人中心
                  </h1>
                  <p className="mt-3 max-w-2xl text-sm leading-6 text-ink-500">
                    管理账户资料和 Aurora 的默认交互偏好。
                  </p>
                </div>
                <Button variant="ghost" onClick={() => window.history.back()} className="w-full sm:w-auto">
                  返回
                </Button>
              </div>
            </div>
          ) : null}

          {activeTab === 'profile' ? (
            <form onSubmit={onSubmit} className="mt-6 sm:mt-8" noValidate>
              <section className="border-b border-ink-200/70 pb-7 sm:pb-8">
                <div className="mb-5 flex items-center gap-2">
                  <Settings2 className="h-4 w-4 text-accent-vermilion" />
                  <h2 className="text-base font-semibold text-ink-950">基础资料</h2>
                </div>
                <div className="grid gap-5 md:grid-cols-2">
                  <Field label="用户名" error={errors.username}>
                    <Input
                      value={form.username}
                      onChange={set('username')}
                      placeholder="昵称"
                      autoComplete="username"
                      invalid={!!errors.username}
                    />
                  </Field>
                  <Field label="邮箱" error={errors.email}>
                    <Input
                      type="email"
                      value={form.email}
                      onChange={set('email')}
                      placeholder="you@example.com"
                      autoComplete="email"
                      invalid={!!errors.email}
                    />
                  </Field>
                  <Field label="手机号">
                    <Input
                      value={form.phone}
                      onChange={set('phone')}
                      placeholder="可不填"
                      autoComplete="tel"
                    />
                  </Field>
                  <Field label="生日">
                    <Input
                      type="date"
                      value={form.birthday}
                      onChange={set('birthday')}
                    />
                  </Field>
                </div>
              </section>

              <section className="border-b border-ink-200/70 py-7 sm:py-8">
                <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
                  <div>
                    <div className="flex items-center gap-2">
                      <Sparkles className="h-4 w-4 text-accent-vermilion" />
                      <h2 className="text-base font-semibold text-ink-950">User Prompt</h2>
                    </div>
                    <p className="mt-2 max-w-2xl text-sm leading-6 text-ink-500">
                      这段内容会作为你的默认偏好，帮助 Aurora 贴近你的表达方式和工作习惯。
                    </p>
                  </div>
                  <span className="text-xs text-ink-500">{promptLength}/2000</span>
                </div>
                <Textarea
                  className="mt-5 min-h-44 text-[15px] leading-7 sm:min-h-52"
                  value={form.user_prompt}
                  onChange={set('user_prompt')}
                  maxLength={2000}
                  placeholder="例如：回答请先给结论，再列出依据；代码建议尽量贴近当前项目风格。"
                />
              </section>

              <div className="flex flex-col-reverse gap-3 pt-6 sm:flex-row sm:justify-end">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={reset}
                  disabled={!dirty || saving}
                  className="w-full sm:w-auto"
                >
                  重置
                </Button>
                <Button type="submit" loading={saving} disabled={!dirty} className="w-full sm:w-auto">
                  <Save className="h-4 w-4" />
                  保存更改
                </Button>
              </div>
            </form>
          ) : (
            <KnowledgePanel
              files={knowledgeFiles}
              loading={knowledgeLoading}
              uploading={uploadingKnowledge}
              deletingPath={deletingKnowledgePath}
              selectedFile={selectedKnowledgeFile}
              vectorizeModeLabel={RAG_VECTORIZE_MODE_LABEL}
              inputRef={knowledgeInputRef}
              onSelectFile={onSelectKnowledgeFile}
              onUpload={onUploadKnowledgeFile}
              onRefresh={loadKnowledgeFiles}
              onDelete={onDeleteKnowledgeFile}
            />
          )}
        </section>
      </main>
    </div>
  )
}

function NavItem({
  active,
  icon,
  label,
  description,
  onClick,
}: {
  active: boolean
  icon: ReactNode
  label: string
  description: string
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        'group flex min-w-0 items-center gap-2 rounded-lg px-2.5 py-2.5 text-left transition sm:gap-3 sm:px-3 sm:py-3',
        active
          ? 'bg-paper-200/80 text-ink-950'
          : 'text-ink-700 hover:bg-paper-100 hover:text-ink-950',
      )}
    >
      <span
        className={cn(
          'flex h-9 w-9 shrink-0 items-center justify-center rounded-md border transition',
          active
            ? 'border-accent-vermilion/20 bg-accent-vermilion/10 text-accent-vermilion'
            : 'border-ink-200 bg-paper-50 text-ink-500 group-hover:text-ink-800',
        )}
      >
        {icon}
      </span>
      <span className="min-w-0">
        <span className="block text-sm font-semibold">{label}</span>
        <span className="hidden truncate text-xs text-ink-500 sm:block">{description}</span>
      </span>
    </button>
  )
}

function Field({
  label,
  error,
  children,
}: {
  label: string
  error?: string
  children: ReactNode
}) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="text-sm font-medium text-ink-800">{label}</span>
      {children}
      {error ? <span className="text-xs text-danger">{error}</span> : null}
    </label>
  )
}

function AccountFact({ label, value }: { label: string; value: string }) {
  return (
    <div className="border-l border-ink-200 pl-3">
      <p className="text-ink-500">{label}</p>
      <p className="mt-1 truncate font-medium text-ink-900">{value}</p>
    </div>
  )
}

function KnowledgePanel({
  files,
  loading,
  uploading,
  deletingPath,
  selectedFile,
  vectorizeModeLabel,
  inputRef,
  onSelectFile,
  onUpload,
  onRefresh,
  onDelete,
}: {
  files: NetworkdiskFile[]
  loading: boolean
  uploading: boolean
  deletingPath: string | null
  selectedFile: File | null
  vectorizeModeLabel: string
  inputRef: RefObject<HTMLInputElement | null>
  onSelectFile: (e: ChangeEvent<HTMLInputElement>) => void
  onUpload: () => void
  onRefresh: () => void
  onDelete: (file: NetworkdiskFile) => void
}) {
  const fileCount = files.filter((file) => file.isdir !== 1).length
  const totalSize = files.reduce((sum, file) => sum + (file.isdir === 1 ? 0 : file.size), 0)

  return (
    <section>
      <div className="border-b border-ink-200/70 pb-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <div className="flex items-center gap-2">
              <Database className="h-4 w-4 text-accent-vermilion" />
              <h2 className="text-base font-semibold text-ink-950">知识库文件</h2>
            </div>
            <p className="mt-2 max-w-2xl text-sm leading-6 text-ink-500">
              上传后会写入网盘，并通过 {vectorizeModeLabel} 接口提交向量化。
            </p>
          </div>
          <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-end">
            <div className="grid grid-cols-2 gap-3 text-xs sm:flex sm:items-center">
              <KnowledgeStat icon={<FileText className="h-4 w-4" />} label="文件" value={`${fileCount}`} />
              <KnowledgeStat icon={<HardDrive className="h-4 w-4" />} label="容量" value={formatBytes(totalSize)} />
            </div>
            <Button variant="ghost" onClick={() => window.history.back()} className="w-full sm:w-auto">
              返回
            </Button>
          </div>
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
            {/* <p className="mt-2 text-xs text-ink-500">当前向量化：{vectorizeModeLabel}</p> */}
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
          <Button type="button" variant="ghost" size="sm" onClick={onRefresh} disabled={loading}>
            <RefreshCw className={cn('h-4 w-4', loading && 'animate-spin')} />
            刷新
          </Button>
        </div>

        <div className="overflow-hidden rounded-lg border border-ink-200 bg-paper-50">
          <div className="grid grid-cols-[minmax(0,1fr)_96px_128px_72px] gap-3 border-b border-ink-200 bg-paper-100/70 px-4 py-3 text-xs font-semibold text-ink-500 max-md:hidden">
            <span>名称</span>
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
    </section>
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
    <li className="grid gap-3 px-4 py-3 text-sm transition hover:bg-paper-100/50 md:grid-cols-[minmax(0,1fr)_96px_128px_72px] md:items-center">
      <div className="flex min-w-0 items-center gap-3">
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-ink-200 bg-paper-100 text-ink-500">
          <FileText className="h-4 w-4" />
        </span>
        <div className="min-w-0">
          <p className="truncate font-medium text-ink-950">{file.server_filename}</p>
          <p className="mt-0.5 truncate text-xs text-ink-500">{file.path}</p>
        </div>
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

function sameForm(a: ProfileForm, b: ProfileForm) {
  return (
    a.username === b.username &&
    a.email === b.email &&
    a.phone === b.phone &&
    a.birthday === b.birthday &&
    a.user_prompt === b.user_prompt
  )
}

function formatDate(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
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
