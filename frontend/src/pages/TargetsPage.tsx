import { useEffect, useState, type FormEvent } from 'react'
import { useTargets } from '@/hooks/useTargets'
import { useTrackers } from '@/hooks/useTrackers'
import { useToast } from '@/context/ToastContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { SelectField } from '@/components/ui/SelectField'
import { Plus } from 'lucide-react'
import { formatDate, truncate } from '@/lib/utils'
import type { Target } from '@/types'

export function TargetsPage() {
  const { targets, loading, fetch, create, creating } = useTargets()
  const { trackers, fetch: fetchTrackers } = useTrackers()
  const { addToast } = useToast()

  const [showCreate, setShowCreate] = useState(false)
  const [trackerId, setTrackerId] = useState('')
  const [url, setUrl] = useState('')

  useEffect(() => { fetchTrackers() }, [fetchTrackers])
  useEffect(() => { fetch() }, [fetch])

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await create({ tracker_id: trackerId, url })
      addToast('Target created', 'success')
      setShowCreate(false)
      setTrackerId('')
      setUrl('')
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to create', 'error')
    }
  }

  const trackerName = (id: string) => trackers.find((t) => t.id === id)?.name ?? id.slice(0, 8)

  const columns: Column<Target>[] = [
    { key: 'tracker', header: 'Tracker', render: (r) => <span className="text-muted">{trackerName(r.tracker_id)}</span> },
    {
      key: 'url', header: 'URL', render: (r) => (
        <span className="font-mono text-xs" title={r.url}>{truncate(r.url, 60)}</span>
      ),
    },
    { key: 'created_at', header: 'Created', render: (r) => <span className="text-muted">{formatDate(r.created_at)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">Targets</h2>
        <Button onClick={() => setShowCreate(true)}>
          <Plus className="h-4 w-4" /> New Target
        </Button>
      </div>

      <DataTable columns={columns} data={targets} loading={loading} rowKey={(r) => r.id} emptyMessage="No targets yet" />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title="Create Target">
        <form onSubmit={handleCreate} className="space-y-4">
          <SelectField
            label="Tracker"
            value={trackerId}
            onChange={(e) => setTrackerId(e.target.value)}
            options={trackers.map((t) => ({ value: t.id, label: t.name }))}
            placeholder="Select tracker"
            required
          />
          <FormField label="URL" type="url" value={url} onChange={(e) => setUrl(e.target.value)} placeholder="https://example.com/landing" required />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>Cancel</Button>
            <Button type="submit" loading={creating}>Create</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
