import type { ReactNode } from 'react'
import { Navigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { Spinner } from '@/components/ui/Spinner'
import { safeRedirectPath } from './redirect'

export function GuestGuard({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth()
  const [params] = useSearchParams()
  const redirect = safeRedirectPath(params.get('redirect'))

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-paper-50">
        <Spinner className="h-6 w-6" />
      </div>
    )
  }

  if (user) {
    return <Navigate to={redirect} replace />
  }

  return <>{children}</>
}
