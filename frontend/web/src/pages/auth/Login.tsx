import { useState, type FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { login } from '@/api/auth'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { HttpError } from '@/lib/fetcher'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function Login() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const redirect = params.get('redirect') || '/chat'
  const { refresh } = useAuth()
  const { show } = useToast()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({})

  const validate = () => {
    const next: typeof errors = {}
    if (!EMAIL_RE.test(email)) next.email = '邮箱格式不正确'
    if (password.length < 6) next.password = '密码至少 6 位'
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      const res = await login({ email, password })
      if (res.code === 0) {
        await refresh()
        show('登录成功', 'success')
      } else {
        show(res.message, 'error')
      }
      navigate(redirect, { replace: true })
    } catch (err) {
      const msg =
        err instanceof HttpError
          ? extractMessage(err.info) || `登录失败 (${err.status})`
          : err instanceof Error
            ? err.message
            : '登录失败'
      show(msg, 'error')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthLayout
      title="登录"
      subtitle="使用邮箱与密码登录到 Aurora"
      footer={
        <span>
          还没有账号？
          <span
            role="link"
            tabIndex={0}
            onClick={() => navigate('/register')}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') navigate('/register')
            }}
            className="ml-1 cursor-pointer text-accent-vermilion hover:underline"
          >
            前往注册
          </span>
        </span>
      }
    >
      <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
        <Field label="邮箱" error={errors.email}>
          <Input
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
            invalid={!!errors.email}
          />
        </Field>

        <Field label="密码" error={errors.password}>
          <Input
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="至少 6 位"
            invalid={!!errors.password}
          />
        </Field>

        <Button type="submit" size="lg" loading={submitting}>
          登录
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
