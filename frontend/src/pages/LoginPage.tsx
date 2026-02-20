import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/context/AuthContext'
import { FormField } from '@/components/ui/FormField'
import { Button } from '@/components/ui/Button'
import { Crosshair } from 'lucide-react'

export function LoginPage() {
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const { t } = useTranslation()
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
      setError(err instanceof Error ? err.message : t('login.loginFailed'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg-deep p-4">
      <div className="w-full max-w-sm rounded-xl border border-border bg-bg-card p-8 shadow-2xl">
        <div className="mb-8 text-center">
          <Crosshair className="mx-auto mb-3 h-10 w-10 text-primary" />
          <h1 className="text-xl font-bold font-mono text-text">{t('login.title')}</h1>
          <p className="mt-1 text-sm text-muted">{t('login.subtitle')}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <FormField
            label={t('login.username')}
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder={t('login.usernamePlaceholder')}
            required
            autoFocus
          />
          <FormField
            label={t('login.password')}
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder={t('login.passwordPlaceholder')}
            required
          />

          {error && (
            <div className="rounded-lg border border-error/30 bg-error/10 px-4 py-2 text-sm text-error">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" loading={loading}>
            {t('login.signIn')}
          </Button>
        </form>
      </div>
    </div>
  )
}
