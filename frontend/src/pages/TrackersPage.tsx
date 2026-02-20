import { useEffect, useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { useTrackers } from '@/hooks/useTrackers'
import { useToast } from '@/context/ToastContext'
import { useLanguage } from '@/context/LanguageContext'
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
  const { t } = useTranslation()
  const { locale } = useLanguage()
  const { trackers, loading, fetch, create, creating, update, updating, remove, deleting } = useTrackers()
  const { addToast } = useToast()

  const [showCreate, setShowCreate] = useState(false)
  const [editTracker, setEditTracker] = useState<Tracker | null>(null)
  const [deleteTracker, setDeleteTracker] = useState<Tracker | null>(null)

  const [name, setName] = useState('')
  const [type, setType] = useState<string>('ad')
  const [status, setStatus] = useState<string>('active')

  useEffect(() => { fetch() }, [fetch])

  const resetForm = () => {
    setName('')
    setType('ad')
    setStatus('active')
  }

  const openEdit = (t: Tracker) => {
    setEditTracker(t)
    setName(t.name)
    setStatus(t.status)
  }

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await create({ name, type })
      addToast(t('trackers.trackerCreated'), 'success')
      setShowCreate(false)
      resetForm()
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToCreate'), 'error')
    }
  }

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault()
    if (!editTracker) return
    try {
      await update({ id: editTracker.id, name, status })
      addToast(t('trackers.trackerUpdated'), 'success')
      setEditTracker(null)
      resetForm()
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToUpdate'), 'error')
    }
  }

  const handleDelete = async () => {
    if (!deleteTracker) return
    try {
      await remove({ id: deleteTracker.id })
      addToast(t('trackers.trackerDeleted'), 'success')
      setDeleteTracker(null)
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('common.failedToDelete'), 'error')
    }
  }

  const columns: Column<Tracker>[] = [
    { key: 'name', header: t('common.name'), render: (r) => <span className="font-medium">{r.name}</span> },
    {
      key: 'type', header: t('trackers.type'), render: (r) => (
        <span className="inline-flex rounded-full bg-secondary/15 px-2.5 py-0.5 text-xs font-medium text-secondary">{r.type}</span>
      ),
    },
    { key: 'status', header: t('common.status'), render: (r) => <StatusBadge status={r.status} /> },
    { key: 'created_at', header: t('common.created'), render: (r) => <span className="text-muted">{formatDate(r.created_at, locale)}</span> },
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
        <h2 className="text-2xl font-bold font-mono text-text">{t('trackers.title')}</h2>
        <Button onClick={() => { resetForm(); setShowCreate(true) }}>
          <Plus className="h-4 w-4" /> {t('trackers.newTracker')}
        </Button>
      </div>

      <DataTable columns={columns} data={trackers} loading={loading} rowKey={(r) => r.id} emptyMessage={t('trackers.noTrackers')} />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title={t('trackers.createTracker')}>
        <form onSubmit={handleCreate} className="space-y-4">
          <FormField label={t('common.name')} value={name} onChange={(e) => setName(e.target.value)} required />
          <SelectField label={t('trackers.type')} value={type} onChange={(e) => setType(e.target.value)} options={[{ value: 'ad', label: t('trackers.ad') }, { value: 'web', label: t('trackers.web') }]} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>{t('common.cancel')}</Button>
            <Button type="submit" loading={creating}>{t('common.create')}</Button>
          </div>
        </form>
      </Modal>

      <Modal open={!!editTracker} onClose={() => setEditTracker(null)} title={t('trackers.editTracker')}>
        <form onSubmit={handleUpdate} className="space-y-4">
          <FormField label={t('common.name')} value={name} onChange={(e) => setName(e.target.value)} required />
          <SelectField label={t('common.status')} value={status} onChange={(e) => setStatus(e.target.value)} options={[{ value: 'active', label: t('common.active') }, { value: 'inactive', label: t('common.inactive') }]} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" type="button" onClick={() => setEditTracker(null)}>{t('common.cancel')}</Button>
            <Button type="submit" loading={updating}>{t('common.update')}</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog
        open={!!deleteTracker}
        onClose={() => setDeleteTracker(null)}
        onConfirm={handleDelete}
        title={t('trackers.deleteTracker')}
        message={t('trackers.deleteConfirm', { name: deleteTracker?.name })}
        loading={deleting}
      />
    </div>
  )
}
