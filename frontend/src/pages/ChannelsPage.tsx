import { useEffect, useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { useChannels } from '@/hooks/useChannels'
import { useTrackers } from '@/hooks/useTrackers'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useToast } from '@/context/ToastContext'
import { useLanguage } from '@/context/LanguageContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { SelectField } from '@/components/ui/SelectField'
import { Plus, Upload, X } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Channel } from '@/types'

export function ChannelsPage() {
  const { t } = useTranslation()
  const { locale } = useLanguage()
  const { channels, loading, fetch, create, creating, batchImport, importing } = useChannels()
  const { trackers, fetch: fetchTrackers } = useTrackers()
  const { campaigns, fetch: fetchCampaigns } = useCampaigns()
  const { addToast } = useToast()

  const [filterTracker, setFilterTracker] = useState('')
  const [filterCampaign, setFilterCampaign] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [showImport, setShowImport] = useState(false)

  const [name, setName] = useState('')
  const [trackerId, setTrackerId] = useState('')
  const [campaignId, setCampaignId] = useState('')
  const [source, setSource] = useState('')
  const [medium, setMedium] = useState('')
  const [tagRows, setTagRows] = useState<{ key: string; value: string }[]>([])
  const [importJson, setImportJson] = useState('')

  useEffect(() => { fetchTrackers() }, [fetchTrackers])

  useEffect(() => {
    if (filterTracker) {
      fetchCampaigns({ tracker_id: filterTracker })
    }
  }, [filterTracker, fetchCampaigns])

  useEffect(() => {
    const params: Record<string, unknown> = {}
    if (filterTracker) params.tracker_id = filterTracker
    if (filterCampaign) params.campaign_id = filterCampaign
    fetch(params)
  }, [fetch, filterTracker, filterCampaign])

  const filteredCampaigns = campaigns.filter((c) => !filterTracker || c.tracker_id === filterTracker)
  const createCampaigns = campaigns.filter((c) => !trackerId || c.tracker_id === trackerId)

  const addTagRow = () => setTagRows([...tagRows, { key: '', value: '' }])
  const removeTagRow = (i: number) => setTagRows(tagRows.filter((_, idx) => idx !== i))
  const updateTagRow = (i: number, field: 'key' | 'value', val: string) => {
    const next = [...tagRows]
    next[i][field] = val
    setTagRows(next)
  }

  const resetForm = () => {
    setName('')
    setTrackerId('')
    setCampaignId('')
    setSource('')
    setMedium('')
    setTagRows([])
  }

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    const tags: Record<string, string> = {}
    tagRows.forEach((r) => { if (r.key) tags[r.key] = r.value })
    try {
      await create({ tracker_id: trackerId, campaign_id: campaignId, name, source, medium, tags })
      addToast(t('channels.channelCreated'), 'success')
      setShowCreate(false)
      resetForm()
      fetch(filterTracker ? { tracker_id: filterTracker } : {})
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToCreate'), 'error')
    }
  }

  const handleImport = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const parsed = JSON.parse(importJson)
      const result = await batchImport({ channels: parsed })
      addToast(t('channels.importedCount', { count: result.imported }), 'success')
      setShowImport(false)
      setImportJson('')
      fetch({})
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.invalidJson'), 'error')
    }
  }

  const trackerName = (id: string) => trackers.find((t) => t.id === id)?.name ?? id.slice(0, 8)
  const campaignName = (id: string) => campaigns.find((c) => c.id === id)?.name ?? id.slice(0, 8)

  const columns: Column<Channel>[] = [
    { key: 'name', header: t('common.name'), render: (r) => <span className="font-medium">{r.name}</span> },
    { key: 'tracker', header: t('common.tracker'), render: (r) => <span className="text-muted">{trackerName(r.tracker_id)}</span> },
    { key: 'campaign', header: t('common.campaign'), render: (r) => <span className="text-muted">{campaignName(r.campaign_id)}</span> },
    { key: 'source', header: t('channels.source'), render: (r) => r.source },
    { key: 'medium', header: t('channels.medium'), render: (r) => r.medium },
    {
      key: 'tags', header: t('channels.tags'), render: (r) => (
        <div className="flex flex-wrap gap-1">
          {Object.entries(r.tags ?? {}).map(([k, v]) => (
            <span key={k} className="rounded bg-bg-elevated px-1.5 py-0.5 text-xs text-muted">
              {k}={v}
            </span>
          ))}
        </div>
      ),
    },
    { key: 'created_at', header: t('common.created'), render: (r) => <span className="text-muted">{formatDate(r.created_at, locale)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">{t('channels.title')}</h2>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => setShowImport(true)}>
            <Upload className="h-4 w-4" /> {t('channels.batchImport')}
          </Button>
          <Button onClick={() => { resetForm(); setShowCreate(true) }}>
            <Plus className="h-4 w-4" /> {t('channels.newChannel')}
          </Button>
        </div>
      </div>

      <div className="mb-4 flex gap-4">
        <SelectField
          label={t('channels.filterByTracker')}
          value={filterTracker}
          onChange={(e) => { setFilterTracker(e.target.value); setFilterCampaign('') }}
          options={trackers.map((t) => ({ value: t.id, label: t.name }))}
          placeholder={t('common.allTrackers')}
        />
        <SelectField
          label={t('channels.filterByCampaign')}
          value={filterCampaign}
          onChange={(e) => setFilterCampaign(e.target.value)}
          options={filteredCampaigns.map((c) => ({ value: c.id, label: c.name }))}
          placeholder={t('common.allCampaigns')}
        />
      </div>

      <DataTable columns={columns} data={channels} loading={loading} rowKey={(r) => r.id} emptyMessage={t('channels.noChannels')} />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title={t('channels.createChannel')}>
        <form onSubmit={handleCreate} className="space-y-4">
          <SelectField
            label={t('common.tracker')}
            value={trackerId}
            onChange={(e) => { setTrackerId(e.target.value); setCampaignId(''); fetchCampaigns({ tracker_id: e.target.value }) }}
            options={trackers.map((t) => ({ value: t.id, label: t.name }))}
            placeholder={t('common.selectTracker')}
            required
          />
          <SelectField
            label={t('common.campaign')}
            value={campaignId}
            onChange={(e) => setCampaignId(e.target.value)}
            options={createCampaigns.map((c) => ({ value: c.id, label: c.name }))}
            placeholder={t('common.selectCampaign')}
            required
          />
          <FormField label={t('common.name')} value={name} onChange={(e) => setName(e.target.value)} required />
          <FormField label={t('channels.source')} value={source} onChange={(e) => setSource(e.target.value)} />
          <FormField label={t('channels.medium')} value={medium} onChange={(e) => setMedium(e.target.value)} />

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium text-muted">{t('channels.tags')}</label>
              <button type="button" onClick={addTagRow} className="text-xs text-primary hover:text-primary/80 cursor-pointer">{t('channels.addTag')}</button>
            </div>
            {tagRows.map((row, i) => (
              <div key={i} className="flex gap-2 items-center">
                <input
                  className="flex-1 rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text outline-none focus:border-primary"
                  placeholder={t('channels.key')}
                  value={row.key}
                  onChange={(e) => updateTagRow(i, 'key', e.target.value)}
                />
                <input
                  className="flex-1 rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text outline-none focus:border-primary"
                  placeholder={t('channels.value')}
                  value={row.value}
                  onChange={(e) => updateTagRow(i, 'value', e.target.value)}
                />
                <button type="button" onClick={() => removeTagRow(i)} className="text-muted hover:text-error cursor-pointer">
                  <X className="h-4 w-4" />
                </button>
              </div>
            ))}
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>{t('common.cancel')}</Button>
            <Button type="submit" loading={creating}>{t('common.create')}</Button>
          </div>
        </form>
      </Modal>

      <Modal open={showImport} onClose={() => setShowImport(false)} title={t('channels.batchImportChannels')}>
        <form onSubmit={handleImport} className="space-y-4">
          <div className="space-y-1">
            <label className="block text-sm font-medium text-muted">{t('channels.jsonArray')}</label>
            <textarea
              className="w-full rounded-lg border border-border bg-bg-card px-3 py-2 text-sm text-text font-mono outline-none focus:border-primary h-40 resize-y"
              placeholder={t('channels.jsonPlaceholder')}
              value={importJson}
              onChange={(e) => setImportJson(e.target.value)}
              required
            />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowImport(false)}>{t('common.cancel')}</Button>
            <Button type="submit" loading={importing}>{t('common.import')}</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
