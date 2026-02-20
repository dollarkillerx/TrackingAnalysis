import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { rpcCall } from '@/lib/rpc-client'
import { RPC, TOKEN_KEY } from '@/lib/constants'

interface AuthContextValue {
  token: string | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY))

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
  }, [])

  const login = useCallback(async (username: string, password: string) => {
    const result = await rpcCall<{ admin_token: string }>(RPC.LOGIN, { username, password })
    localStorage.setItem(TOKEN_KEY, result.admin_token)
    setToken(result.admin_token)
  }, [])

  useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === TOKEN_KEY) {
        setToken(e.newValue)
      }
    }
    window.addEventListener('storage', onStorage)
    return () => window.removeEventListener('storage', onStorage)
  }, [])

  return (
    <AuthContext value={{ token, isAuthenticated: !!token, login, logout }}>
      {children}
    </AuthContext>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
