import { useEffect, useState, type FormEvent } from 'react'
import { useTrackers } from '@/hooks/useTrackers'
import { useToast } from '@/context/ToastContext'
import { DataTable, type Column } from '@/components/ui/DataTable'
import { Button } from '@/components/ui/Button'
import { Modal } from '@/components/ui/Modal'
import { FormField } from '@/components/ui/FormField'
import { SelectField } from '@/components/ui/SelectField'
import { StatusBadge } from '@/components/ui/StatusBadge'
import { ConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import type { Tracker } from '@/types'

export function TrackersPage() {
  const { trackers, loading, fetch, create, creating, update, updating, remove, deleting } = useTrackers()
  const { addToast } = useToast()

  const [showCreate, setShowCreate] = useState(false)
  const [editTracker, setEditTracker] = useState<Tracker | null>(null)
  const [deleteTracker, setDeleteTracker] = useState<Tracker | null>(null)

  const [name, setName] = useState('')
  const [type, setType] = useState<string>('ad')
  const [mode, setMode] = useState<string>('302')
  const [status, setStatus] = useState<string>('active')

  useEffect(() => { fetch() }, [fetch])

  const resetForm = () => {
    setName('')
    setType('ad')
    setMode('302')
    setStatus('active')
  }

  const openEdit = (t: Tracker) => {
    setEditTracker(t)
    setName(t.name)
    setMode(t.mode)
    setStatus(t.status)
  }

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await create({ name, type, mode })
      addToast('Tracker created', 'success')
      setShowCreate(false)
      resetForm()
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to create', 'error')
    }
  }

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault()
    if (!editTracker) return
    try {
      await update({ id: editTracker.id, name, mode, status })
      addToast('Tracker updated', 'success')
      setEditTracker(null)
      resetForm()
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to update', 'error')
    }
  }

  const handleDelete = async () => {
    if (!deleteTracker) return
    try {
      await remove({ id: deleteTracker.id })
      addToast('Tracker deleted', 'success')
      setDeleteTracker(null)
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to delete', 'error')
    }
  }

  const columns: Column<Tracker>[] = [
    { key: 'name', header: 'Name', render: (r) => <span className="font-medium">{r.name}</span> },
    {
      key: 'type', header: 'Type', render: (r) => (
        <span className="inline-flex rounded-full bg-secondary/15 px-2.5 py-0.5 text-xs font-medium text-secondary">{r.type}</span>
      ),
    },
    {
      key: 'mode', header: 'Mode', render: (r) => (
        <span className="inline-flex rounded-full bg-accent/15 px-2.5 py-0.5 text-xs font-medium text-accent">{r.mode}</span>
      ),
    },
    { key: 'status', header: 'Status', render: (r) => <StatusBadge status={r.status} /> },
    { key: 'created_at', header: 'Created', render: (r) => <span className="text-muted">{formatDate(r.created_at)}</span> },
    {
      key: 'actions', header: '', render: (r) => (
        <div className="flex gap-1">
          <button onClick={() => openEdit(r)} className="rounded p-1 text-muted hover:text-text hover:bg-bg-elevated transition-colors cursor-pointer">
            <Pencil className="h-4 w-4" />
          </button>
          <button onClick={() => setDeleteTracker(r)} className="rounded p-1 text-muted hover:text-error hover:bg-error/10 transition-colors cursor-pointer">
            <Trash2 className="h-4 w-4" />
          </button>
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">Trackers</h2>
        <Button onClick={() => { resetForm(); setShowCreate(true) }}>
          <Plus className="h-4 w-4" /> New Tracker
        </Button>
      </div>

      <DataTable columns={columns} data={trackers} loading={loading} rowKey={(r) => r.id} emptyMessage="No trackers yet" />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title="Create Tracker">
        <form onSubmit={handleCreate} className="space-y-4">
          <FormField label="Name" value={name} onChange={(e) => setName(e.target.value)} required />
          <SelectField label="Type" value={type} onChange={(e) => setType(e.target.value)} options={[{ value: 'ad', label: 'Ad' }, { value: 'web', label: 'Web' }]} />
          <SelectField label="Mode" value={mode} onChange={(e) => setMode(e.target.value)} options={[{ value: '302', label: '302 Redirect' }, { value: 'js', label: 'JavaScript' }]} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>Cancel</Button>
            <Button type="submit" loading={creating}>Create</Button>
          </div>
        </form>
      </Modal>

      <Modal open={!!editTracker} onClose={() => setEditTracker(null)} title="Edit Tracker">
        <form onSubmit={handleUpdate} className="space-y-4">
          <FormField label="Name" value={name} onChange={(e) => setName(e.target.value)} required />
          <SelectField label="Mode" value={mode} onChange={(e) => setMode(e.target.value)} options={[{ value: '302', label: '302 Redirect' }, { value: 'js', label: 'JavaScript' }]} />
          <SelectField label="Status" value={status} onChange={(e) => setStatus(e.target.value)} options={[{ value: 'active', label: 'Active' }, { value: 'inactive', label: 'Inactive' }]} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setEditTracker(null)}>Cancel</Button>
            <Button type="submit" loading={updating}>Update</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog
        open={!!deleteTracker}
        onClose={() => setDeleteTracker(null)}
        onConfirm={handleDelete}
        title="Delete Tracker"
        message={`Are you sure you want to delete "${deleteTracker?.name}"? This action cannot be undone.`}
        loading={deleting}
      />
    </div>
  )
}
