import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import { FormField } from '@/components/ui/FormField'
import { Button } from '@/components/ui/Button'
import { Crosshair } from 'lucide-react'

export function LoginPage() {
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  if (isAuthenticated) {
    navigate('/', { replace: true })
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(username, password)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg-deep p-4">
      <div className="w-full max-w-sm rounded-xl border border-border bg-bg-card p-8 shadow-2xl">
        <div className="mb-8 text-center">
          <Crosshair className="mx-auto mb-3 h-10 w-10 text-primary" />
          <h1 className="text-xl font-bold font-mono text-text">TrackingAnalysis</h1>
          <p className="mt-1 text-sm text-muted">Admin Dashboard</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <FormField
            label="Username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="admin"
            required
            autoFocus
          />
          <FormField
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="********"
            required
          />

          {error && (
            <div className="rounded-lg border border-error/30 bg-error/10 px-4 py-2 text-sm text-error">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" loading={loading}>
            Sign In
          </Button>
        </form>
      </div>
    </div>
  )
}
