import { useEffect, useState, type FormEvent } from 'react'
import { useChannels } from '@/hooks/useChannels'
import { useTrackers } from '@/hooks/useTrackers'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useToast } from '@/context/ToastContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { SelectField } from '@/components/ui/SelectField'
import { Plus, Upload, X } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Channel } from '@/types'

export function ChannelsPage() {
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
      addToast('Channel created', 'success')
      setShowCreate(false)
      resetForm()
      fetch(filterTracker ? { tracker_id: filterTracker } : {})
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to create', 'error')
    }
  }

  const handleImport = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const parsed = JSON.parse(importJson)
      const result = await batchImport({ channels: parsed })
      addToast(`Imported ${result.imported} channels`, 'success')
      setShowImport(false)
      setImportJson('')
      fetch({})
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Invalid JSON', 'error')
    }
  }

  const trackerName = (id: string) => trackers.find((t) => t.id === id)?.name ?? id.slice(0, 8)
  const campaignName = (id: string) => campaigns.find((c) => c.id === id)?.name ?? id.slice(0, 8)

  const columns: Column<Channel>[] = [
    { key: 'name', header: 'Name', render: (r) => <span className="font-medium">{r.name}</span> },
    { key: 'tracker', header: 'Tracker', render: (r) => <span className="text-muted">{trackerName(r.tracker_id)}</span> },
    { key: 'campaign', header: 'Campaign', render: (r) => <span className="text-muted">{campaignName(r.campaign_id)}</span> },
    { key: 'source', header: 'Source', render: (r) => r.source },
    { key: 'medium', header: 'Medium', render: (r) => r.medium },
    {
      key: 'tags', header: 'Tags', render: (r) => (
        <div className="flex flex-wrap gap-1">
          {Object.entries(r.tags ?? {}).map(([k, v]) => (
            <span key={k} className="rounded bg-bg-elevated px-1.5 py-0.5 text-xs text-muted">
              {k}={v}
            </span>
          ))}
        </div>
      ),
    },
    { key: 'created_at', header: 'Created', render: (r) => <span className="text-muted">{formatDate(r.created_at)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">Channels</h2>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => setShowImport(true)}>
            <Upload className="h-4 w-4" /> Batch Import
          </Button>
          <Button onClick={() => { resetForm(); setShowCreate(true) }}>
            <Plus className="h-4 w-4" /> New Channel
          </Button>
        </div>
      </div>

      <div className="mb-4 flex gap-4">
        <SelectField
          label="Filter by Tracker"
          value={filterTracker}
          onChange={(e) => { setFilterTracker(e.target.value); setFilterCampaign('') }}
          options={trackers.map((t) => ({ value: t.id, label: t.name }))}
          placeholder="All trackers"
        />
        <SelectField
          label="Filter by Campaign"
          value={filterCampaign}
          onChange={(e) => setFilterCampaign(e.target.value)}
          options={filteredCampaigns.map((c) => ({ value: c.id, label: c.name }))}
          placeholder="All campaigns"
        />
      </div>

      <DataTable columns={columns} data={channels} loading={loading} rowKey={(r) => r.id} emptyMessage="No channels yet" />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title="Create Channel">
        <form onSubmit={handleCreate} className="space-y-4">
          <SelectField
            label="Tracker"
            value={trackerId}
            onChange={(e) => { setTrackerId(e.target.value); setCampaignId(''); fetchCampaigns({ tracker_id: e.target.value }) }}
            options={trackers.map((t) => ({ value: t.id, label: t.name }))}
            placeholder="Select tracker"
            required
          />
          <SelectField
            label="Campaign"
            value={campaignId}
            onChange={(e) => setCampaignId(e.target.value)}
            options={createCampaigns.map((c) => ({ value: c.id, label: c.name }))}
            placeholder="Select campaign"
            required
          />
          <FormField label="Name" value={name} onChange={(e) => setName(e.target.value)} required />
          <FormField label="Source" value={source} onChange={(e) => setSource(e.target.value)} />
          <FormField label="Medium" value={medium} onChange={(e) => setMedium(e.target.value)} />

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium text-muted">Tags</label>
              <button type="button" onClick={addTagRow} className="text-xs text-primary hover:text-primary/80 cursor-pointer">+ Add Tag</button>
            </div>
            {tagRows.map((row, i) => (
              <div key={i} className="flex gap-2 items-center">
                <input
                  className="flex-1 rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text outline-none focus:border-primary"
                  placeholder="Key"
                  value={row.key}
                  onChange={(e) => updateTagRow(i, 'key', e.target.value)}
                />
                <input
                  className="flex-1 rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text outline-none focus:border-primary"
                  placeholder="Value"
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
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>Cancel</Button>
            <Button type="submit" loading={creating}>Create</Button>
          </div>
        </form>
      </Modal>

      <Modal open={showImport} onClose={() => setShowImport(false)} title="Batch Import Channels">
        <form onSubmit={handleImport} className="space-y-4">
          <div className="space-y-1">
            <label className="block text-sm font-medium text-muted">JSON Array</label>
            <textarea
              className="w-full rounded-lg border border-border bg-bg-card px-3 py-2 text-sm text-text font-mono outline-none focus:border-primary h-40 resize-y"
              placeholder='[{"tracker_id":"...","campaign_id":"...","name":"...","source":"...","medium":"...","tags":{}}]'
              value={importJson}
              onChange={(e) => setImportJson(e.target.value)}
              required
            />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowImport(false)}>Cancel</Button>
            <Button type="submit" loading={importing}>Import</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
