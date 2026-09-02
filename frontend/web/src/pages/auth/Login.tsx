import { useState, type FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { Button } from '@/components/ui/Button'
import { Checkbox } from '@/components/ui/Checkbox'
import { Input } from '@/components/ui/Input'
import { login } from '@/api/auth'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { clearCredentials, loadCredentials, saveCredentials } from '@/lib/credentials'
import { HttpError } from '@/lib/fetcher'
import { safeRedirectPath } from '@/router/redirect'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function Login() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const redirect = safeRedirectPath(params.get('redirect'))
  const { refresh } = useAuth()
  const { show } = useToast()

  const saved = loadCredentials()
  const [email, setEmail] = useState(saved?.email ?? '')
  const [password, setPassword] = useState(saved?.password ?? '')
  const [remember, setRemember] = useState(true)
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
      await login({ email, password })
      if (remember) {
        saveCredentials({ email, password })
      } else {
        clearCredentials()
      }
      await refresh()
      show('登录成功', 'success')
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
      subtitle="使用邮箱与密码登录到 Ariadne"
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

        <label className="flex cursor-pointer items-center gap-2 text-sm text-ink-700 select-none">
          <Checkbox checked={remember} onChange={(e) => setRemember(e.target.checked)} />
          记住账号和密码
        </label>

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
