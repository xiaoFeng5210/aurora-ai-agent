import { useMemo, useState, type ChangeEvent, type FormEvent, type ReactNode } from 'react'
import { Save, Settings2, Sparkles } from 'lucide-react'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Textarea } from '@/components/ui/Textarea'
import { updateMe, type User } from '@/api/auth'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'

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
        </section>
      </main>
    </div>
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
