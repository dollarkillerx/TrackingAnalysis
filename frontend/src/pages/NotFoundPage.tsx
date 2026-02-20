import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/Button'
import { AlertTriangle } from 'lucide-react'

export function NotFoundPage() {
  const navigate = useNavigate()
  const { t } = useTranslation()

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-bg-deep text-center">
      <AlertTriangle className="mb-4 h-16 w-16 text-warning" />
      <h1 className="mb-2 text-4xl font-bold font-mono text-text">{t('notFound.title')}</h1>
      <p className="mb-6 text-muted">{t('notFound.message')}</p>
      <Button onClick={() => navigate('/')}>{t('notFound.goHome')}</Button>
    </div>
  )
}
