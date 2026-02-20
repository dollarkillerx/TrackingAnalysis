import { LogOut } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/context/AuthContext'
import { Button } from '@/components/ui/Button'
import { LanguageSwitcher } from '@/components/ui/LanguageSwitcher'

export function Header() {
  const { logout } = useAuth()
  const { t } = useTranslation()

  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-bg-card px-6 shrink-0">
      <h1 className="text-lg font-semibold font-mono text-text">{t('common.adminDashboard')}</h1>
      <div className="flex items-center gap-3">
        <LanguageSwitcher />
        <Button variant="ghost" size="sm" onClick={logout}>
          <LogOut className="h-4 w-4" />
          {t('header.logout')}
        </Button>
      </div>
    </header>
  )
}
