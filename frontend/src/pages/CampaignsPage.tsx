import { useEffect, useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useTrackers } from '@/hooks/useTrackers'
import { useToast } from '@/context/ToastContext'
import { useLanguage } from '@/context/LanguageContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { SelectField } from '@/components/ui/SelectField'
import { StatusBadge } from '@/components/ui/StatusBadge'
import { Plus } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Campaign } from '@/types'

export function CampaignsPage() {
  const { t } = useTranslation()
  const { locale } = useLanguage()
  const { campaigns, loading, fetch, create, creating } = useCampaigns()
  const { trackers, fetch: fetchTrackers } = useTrackers()
  const { addToast } = useToast()

  const [filterTracker, setFilterTracker] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [trackerId, setTrackerId] = useState('')

  const adTrackers = trackers.filter((t) => t.type === 'ad')

  useEffect(() => { fetchTrackers() }, [fetchTrackers])

  useEffect(() => {
    fetch(filterTracker ? { tracker_id: filterTracker } : {})
  }, [fetch, filterTracker])

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await create({ tracker_id: trackerId, name })
      addToast(t('campaigns.campaignCreated'), 'success')
      setShowCreate(false)
      setName('')
      setTrackerId('')
      fetch(filterTracker ? { tracker_id: filterTracker } : {})
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToCreate'), 'error')
    }
  }

  const trackerName = (id: string) => trackers.find((t) => t.id === id)?.name ?? id.slice(0, 8)

  const columns: Column<Campaign>[] = [
    { key: 'name', header: t('common.name'), render: (r) => <span className="font-medium">{r.name}</span> },
    { key: 'tracker', header: t('common.tracker'), render: (r) => <span className="text-muted">{trackerName(r.tracker_id)}</span> },
    { key: 'status', header: t('common.status'), render: (r) => <StatusBadge status={r.status} /> },
    { key: 'created_at', header: t('common.created'), render: (r) => <span className="text-muted">{formatDate(r.created_at, locale)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">{t('campaigns.title')}</h2>
        <Button onClick={() => setShowCreate(true)}>
          <Plus className="h-4 w-4" /> {t('campaigns.newCampaign')}
        </Button>
      </div>

      <div className="mb-4">
        <SelectField
          label={t('campaigns.filterByTracker')}
          value={filterTracker}
          onChange={(e) => setFilterTracker(e.target.value)}
          options={adTrackers.map((t) => ({ value: t.id, label: t.name }))}
          placeholder={t('common.allTrackers')}
        />
      </div>

      <DataTable columns={columns} data={campaigns} loading={loading} rowKey={(r) => r.id} emptyMessage={t('campaigns.noCampaigns')} />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title={t('campaigns.createCampaign')}>
        <form onSubmit={handleCreate} className="space-y-4">
          <SelectField
            label={t('campaigns.trackerAdOnly')}
            value={trackerId}
            onChange={(e) => setTrackerId(e.target.value)}
            options={adTrackers.map((t) => ({ value: t.id, label: t.name }))}
            placeholder={t('common.selectTracker')}
            required
          />
          <FormField label={t('common.name')} value={name} onChange={(e) => setName(e.target.value)} required />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>{t('common.cancel')}</Button>
            <Button type="submit" loading={creating}>{t('common.create')}</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
