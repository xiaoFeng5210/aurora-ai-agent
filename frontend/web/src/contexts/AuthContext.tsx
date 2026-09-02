/* eslint-disable react-refresh/only-export-components */
import { createContext, useCallback, useEffect, useState, type ReactNode } from 'react'
import { getMe, logout as requestLogout, type User } from '@/api/auth'
import { HttpError } from '@/lib/fetcher'

export interface AuthContextValue {
  user: User | null
  loading: boolean
  refresh: () => Promise<User | null>
  setUser: (user: User | null) => void
  logout: () => Promise<void>
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

  const logout = useCallback(async () => {
    try {
      await requestLogout()
    } catch {
      // 接口失败也清掉本地登录态，避免卡在已失效会话
    }
    setUser(null)
  }, [])

  useEffect(() => {
    init()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, refresh, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
