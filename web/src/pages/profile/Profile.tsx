import { useMemo, useState, type ChangeEvent, type FormEvent, type ReactNode } from 'react'
import {
  BookOpen,
  CalendarDays,
  Database,
  Mail,
  Save,
  Settings2,
  Sparkles,
  UserRound,
} from 'lucide-react'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Textarea } from '@/components/ui/Textarea'
import { updateMe, type User } from '@/api/auth'
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
              onClick={() => setActiveTab('knowledge')}
            />
          </nav>
        </aside>

        <section className="min-w-0">
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
            <KnowledgePlaceholder />
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

function KnowledgePlaceholder() {
  return (
    <section className="mt-6 sm:mt-8">
      <div className="rounded-xl border border-dashed border-ink-300 bg-paper-100/40 p-5 sm:p-8 md:p-10">
        <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent-indigo/10 text-accent-indigo sm:h-12 sm:w-12">
          <BookOpen className="h-5 w-5" />
        </div>
        <h2 className="font-display mt-5 text-2xl font-bold text-ink-950 sm:mt-6 sm:text-3xl">知识库</h2>
        <p className="mt-3 max-w-2xl text-sm leading-6 text-ink-500">
          这里会承载文档、索引状态和知识库设置。当前先保留入口，详细管理能力后续接入。
        </p>
        <div className="mt-6 grid gap-4 border-t border-ink-200/70 pt-5 sm:mt-8 sm:gap-5 sm:pt-6 md:grid-cols-3">
          <PlaceholderMetric icon={<Database className="h-4 w-4" />} label="文档" value="待接入" />
          <PlaceholderMetric icon={<CalendarDays className="h-4 w-4" />} label="同步" value="待接入" />
          <PlaceholderMetric icon={<Mail className="h-4 w-4" />} label="权限" value="待接入" />
        </div>
      </div>
    </section>
  )
}

function PlaceholderMetric({
  icon,
  label,
  value,
}: {
  icon: ReactNode
  label: string
  value: string
}) {
  return (
    <div className="min-w-0">
      <div className="flex items-center gap-2 text-ink-500">
        {icon}
        <span className="text-xs">{label}</span>
      </div>
      <p className="mt-3 text-sm font-semibold text-ink-950">{value}</p>
    </div>
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

function extractMessage(info: unknown): string | null {
  if (info && typeof info === 'object' && 'message' in info) {
    const m = (info as { message?: unknown }).message
    if (typeof m === 'string') return m
  }
  return null
}
