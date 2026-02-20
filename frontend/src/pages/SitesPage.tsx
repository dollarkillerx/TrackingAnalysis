import { useEffect, useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { useSites } from '@/hooks/useSites'
import { useToast } from '@/context/ToastContext'
import { useLanguage } from '@/context/LanguageContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { StatusBadge } from '@/components/ui/StatusBadge'
import { CopyButton } from '@/components/ui/CopyButton'
import { Plus } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Site } from '@/types'

export function SitesPage() {
  const { t } = useTranslation()
  const { locale } = useLanguage()
  const { sites, loading, fetch, create, creating } = useSites()
  const { addToast } = useToast()

  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [domain, setDomain] = useState('')
  const [createdSite, setCreatedSite] = useState<Site | null>(null)

  useEffect(() => { fetch() }, [fetch])

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const site = await create({ name, domain })
      addToast(t('sites.siteCreated'), 'success')
      setCreatedSite(site)
      setName('')
      setDomain('')
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToCreate'), 'error')
    }
  }

  const columns: Column<Site>[] = [
    { key: 'name', header: t('common.name'), render: (r) => <span className="font-medium">{r.name}</span> },
    { key: 'domain', header: t('sites.domain'), render: (r) => <span className="font-mono text-sm">{r.domain}</span> },
    {
      key: 'site_key', header: t('sites.siteKey'), render: (r) => (
        <div className="flex items-center gap-1">
          <span className="font-mono text-xs text-muted">{r.site_key.slice(0, 12)}...</span>
          <CopyButton text={r.site_key} />
        </div>
      ),
    },
    { key: 'status', header: t('common.status'), render: (r) => <StatusBadge status={r.status} /> },
    { key: 'created_at', header: t('common.created'), render: (r) => <span className="text-muted">{formatDate(r.created_at, locale)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">{t('sites.title')}</h2>
        <Button onClick={() => { setShowCreate(true); setCreatedSite(null) }}>
          <Plus className="h-4 w-4" /> {t('sites.newSite')}
        </Button>
      </div>

      <DataTable columns={columns} data={sites} loading={loading} rowKey={(r) => r.id} emptyMessage={t('sites.noSites')} />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title={createdSite ? t('sites.siteCreatedTitle') : t('sites.createSite')}>
        {createdSite ? (
          <div className="space-y-4">
            <p className="text-sm text-muted">{t('sites.siteKeyInstruction')}</p>
            <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
              <code className="flex-1 font-mono text-sm text-primary break-all">{createdSite.site_key}</code>
              <CopyButton text={createdSite.site_key} />
            </div>
            <div className="flex justify-end">
              <Button onClick={() => setShowCreate(false)}>{t('common.done')}</Button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleCreate} className="space-y-4">
            <FormField label={t('common.name')} value={name} onChange={(e) => setName(e.target.value)} required />
            <FormField label={t('sites.domain')} value={domain} onChange={(e) => setDomain(e.target.value)} placeholder={t('sites.domainPlaceholder')} required />
            <div className="flex justify-end gap-3 pt-2">
              <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>{t('common.cancel')}</Button>
              <Button type="submit" loading={creating}>{t('common.create')}</Button>
            </div>
          </form>
        )}
      </Modal>
    </div>
  )
}
