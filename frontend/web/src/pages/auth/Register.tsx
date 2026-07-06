import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { login, register } from '@/api/auth'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

interface FormState {
  username: string
  email: string
  password: string
  phone: string
}

export function Register() {
  const navigate = useNavigate()
  const { refresh } = useAuth()
  const { show } = useToast()

  const [form, setForm] = useState<FormState>({
    username: '',
    email: '',
    password: '',
    phone: '',
  })
  const [errors, setErrors] = useState<Partial<Record<keyof FormState, string>>>({})
  const [submitting, setSubmitting] = useState(false)

  const set = (key: keyof FormState) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((prev) => ({ ...prev, [key]: e.target.value }))

  const validate = () => {
    const next: typeof errors = {}
    if (form.username.trim().length < 2) next.username = '用户名至少 2 位'
    if (!EMAIL_RE.test(form.email)) next.email = '邮箱格式不正确'
    if (form.password.length < 6) next.password = '密码至少 6 位'
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      await register({
        username: form.username.trim(),
        email: form.email.trim(),
        password: form.password,
        phone: form.phone.trim() || undefined,
      })
      await login({ email: form.email.trim(), password: form.password })
      await refresh()
      show('注册成功，已自动登录', 'success')
      navigate('/chat', { replace: true })
    } catch (err) {
      const msg =
        err instanceof HttpError
          ? extractMessage(err.info) || `注册失败 (${err.status})`
          : err instanceof Error
            ? err.message
            : '注册失败'
      show(msg, 'error')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthLayout
      title="注册账号"
      subtitle="加入 Ariadne，开始构建你的 AI Agent"
      footer={
        <span>
          已有账号？
          <span
            role="link"
            tabIndex={0}
            onClick={() => navigate('/login')}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') navigate('/login')
            }}
            className="ml-1 cursor-pointer text-accent-vermilion hover:underline"
          >
            前往登录
          </span>
        </span>
      }
    >
      <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
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
        <Field label="密码" error={errors.password}>
          <Input
            type="password"
            value={form.password}
            onChange={set('password')}
            placeholder="至少 6 位"
            autoComplete="new-password"
            invalid={!!errors.password}
          />
        </Field>
        <Field label="手机号（可选）">
          <Input
            value={form.phone}
            onChange={set('phone')}
            placeholder="可不填"
            autoComplete="tel"
          />
        </Field>

        <Button type="submit" size="lg" loading={submitting}>
          注册并登录
        </Button>
      </form>
    </AuthLayout>
  )
}

function Field({
  label,
  error,
  children,
}: {
  label: string
  error?: string
  children: React.ReactNode
}) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="text-sm text-ink-700">{label}</span>
      {children}
      {error ? <span className="text-xs text-danger">{error}</span> : null}
    </label>
  )
}

function extractMessage(info: unknown): string | null {
  if (info && typeof info === 'object' && 'message' in info) {
    const m = (info as { message?: unknown }).message
    if (typeof m === 'string') return m
  }
  return null
}
