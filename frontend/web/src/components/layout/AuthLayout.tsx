import { useNavigate } from 'react-router-dom'
import type { ReactNode } from 'react'
import { Card } from '@/components/ui/Card'
import logoUrl from '@/assets/logo.png'

export function AuthLayout({
  title,
  subtitle,
  footer,
  children,
}: {
  title: string
  subtitle?: string
  footer?: ReactNode
  children: ReactNode
}) {
  const navigate = useNavigate()
  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <div className="mx-auto flex min-h-screen max-w-md flex-col px-6">
        <div className="pt-10">
          <div
            role="link"
            tabIndex={0}
            onClick={() => navigate('/')}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') navigate('/')
            }}
            className="inline-flex w-fit cursor-pointer items-center gap-2 select-none"
          >
            <img src={logoUrl} alt="Ariadne" className="h-8 w-8" />
            <span className="font-display text-xl font-extrabold text-accent-vermilion">Ariadne</span>
            <span className="text-ink-300">·</span>
          </div>
        </div>

        <div className="flex flex-1 items-center justify-center py-10">
          <Card className="w-full p-8">
            <h1 className="font-display text-3xl font-bold text-ink-950">{title}</h1>
            {subtitle ? <p className="mt-2 text-sm text-ink-500">{subtitle}</p> : null}

            <div className="mt-8">{children}</div>

            {footer ? <div className="mt-6 text-sm text-ink-500">{footer}</div> : null}
          </Card>
        </div>
      </div>
    </div>
  )
}
