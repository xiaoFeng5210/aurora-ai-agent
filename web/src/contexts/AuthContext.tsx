/* eslint-disable react-refresh/only-export-components */
import { createContext, useCallback, useEffect, useState, type ReactNode } from 'react'
import { getMe, type User } from '@/api/auth'
import { HttpError } from '@/lib/fetcher'

export interface AuthContextValue {
  user: User | null
  loading: boolean
  refresh: () => Promise<User | null>
  setUser: (user: User | null) => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const init = async () => {
    try {
      const res = await getMe()
      const me = res.data ?? null
      setUser(me)
      return me
    } catch (err) {
      if (err instanceof HttpError && (err.status === 401 || err.status === 403)) {
        setUser(null)
        return null
      }
      setUser(null)
      return null
    } finally {
      setLoading(false)
    }
  }

  const refresh = useCallback(async () => {
    return init()
  }, [init])

  useEffect(() => {
    init()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, refresh, setUser }}>
      {children}
    </AuthContext.Provider>
  )
}
