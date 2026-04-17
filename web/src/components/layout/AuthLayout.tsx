import { Link } from 'react-router-dom'
import type { ReactNode } from 'react'
import { Music } from 'lucide-react'
import { Card } from '@/components/ui/Card'

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
  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <div className="mx-auto flex min-h-screen max-w-md flex-col px-6">
        <div className="pt-10">
          <Link to="/" className="inline-flex items-center gap-2">
            <Music className="h-6 w-6 text-accent-olive" strokeWidth={2.2} />
            <span className="font-display text-xl font-extrabold text-accent-vermilion">Aurora</span>
            <span className="text-ink-300">·</span>
          </Link>
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
