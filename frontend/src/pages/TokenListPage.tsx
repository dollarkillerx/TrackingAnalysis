import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useTokens } from '@/hooks/useTokens'
import { useToast } from '@/context/ToastContext'
import { useLanguage } from '@/context/LanguageContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { CopyButton } from '@/components/ui/CopyButton'
import { ConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Plus, Trash2 } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Token } from '@/types'

export function TokenListPage() {
  const { t } = useTranslation()
  const { locale } = useLanguage()
  const navigate = useNavigate()
  const { tokens, loading, fetch, remove, deleting } = useTokens()
  const { addToast } = useToast()

  const [deleteToken, setDeleteToken] = useState<Token | null>(null)

  useEffect(() => { fetch() }, [fetch])

  const handleDelete = async () => {
    if (!deleteToken) return
    try {
      await remove({ id: deleteToken.id })
      addToast(t('tokens.tokenDeleted'), 'success')
      setDeleteToken(null)
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToDelete'), 'error')
    }
  }

  const columns: Column<Token>[] = [
    {
      key: 'short_code',
      header: t('tokens.shortCode'),
      render: (r) => (
        <div className="flex items-center gap-2">
          <code className="font-mono text-xs text-primary">{r.short_code}</code>
          <CopyButton text={r.short_code} />
        </div>
      ),
    },
    {
      key: 'tracking_url',
      header: t('tokens.trackingUrl'),
      render: (r) => (
        <div className="flex items-center gap-2">
          <code className="font-mono text-xs text-secondary truncate max-w-[200px]">{r.tracking_url}</code>
          <CopyButton text={r.tracking_url} />
        </div>
      ),
    },
    {
      key: 'mode',
      header: t('tokens.mode'),
      render: (r) => (
        <span className="inline-flex rounded-full bg-accent/15 px-2.5 py-0.5 text-xs font-medium text-accent">{r.mode}</span>
      ),
    },
    {
      key: 'created_at',
      header: t('common.created'),
      render: (r) => <span className="text-muted">{formatDate(r.created_at, locale)}</span>,
    },
    {
      key: 'actions',
      header: '',
      render: (r) => (
        <button
          onClick={() => setDeleteToken(r)}
          className="rounded p-1 text-muted hover:text-error hover:bg-error/10 transition-colors cursor-pointer"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      ),
    },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">{t('tokens.listTitle')}</h2>
        <Button onClick={() => navigate('/tokens/new')}>
          <Plus className="h-4 w-4" /> {t('tokens.generateToken')}
        </Button>
      </div>

      <DataTable columns={columns} data={tokens} loading={loading} rowKey={(r) => r.id} emptyMessage={t('tokens.noTokens')} />

      <ConfirmDialog
        open={!!deleteToken}
        onClose={() => setDeleteToken(null)}
        onConfirm={handleDelete}
        title={t('tokens.deleteToken')}
        message={t('tokens.deleteConfirm', { code: deleteToken?.short_code })}
        loading={deleting}
      />
    </div>
  )
}
