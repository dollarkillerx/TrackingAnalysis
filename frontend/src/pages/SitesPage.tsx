import { useEffect, useState, type FormEvent } from 'react'
import { useSites } from '@/hooks/useSites'
import { useToast } from '@/context/ToastContext'
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
      addToast('Site created', 'success')
      setCreatedSite(site)
      setName('')
      setDomain('')
      fetch()
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to create', 'error')
    }
  }

  const columns: Column<Site>[] = [
    { key: 'name', header: 'Name', render: (r) => <span className="font-medium">{r.name}</span> },
    { key: 'domain', header: 'Domain', render: (r) => <span className="font-mono text-sm">{r.domain}</span> },
    {
      key: 'site_key', header: 'Site Key', render: (r) => (
        <div className="flex items-center gap-1">
          <span className="font-mono text-xs text-muted">{r.site_key.slice(0, 12)}...</span>
          <CopyButton text={r.site_key} />
        </div>
      ),
    },
    { key: 'status', header: 'Status', render: (r) => <StatusBadge status={r.status} /> },
    { key: 'created_at', header: 'Created', render: (r) => <span className="text-muted">{formatDate(r.created_at)}</span> },
  ]

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2 className="text-2xl font-bold font-mono text-text">Sites</h2>
        <Button onClick={() => { setShowCreate(true); setCreatedSite(null) }}>
          <Plus className="h-4 w-4" /> New Site
        </Button>
      </div>

      <DataTable columns={columns} data={sites} loading={loading} rowKey={(r) => r.id} emptyMessage="No sites yet" />

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title={createdSite ? 'Site Created' : 'Create Site'}>
        {createdSite ? (
          <div className="space-y-4">
            <p className="text-sm text-muted">Your site has been created. Copy the site key below:</p>
            <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
              <code className="flex-1 font-mono text-sm text-primary break-all">{createdSite.site_key}</code>
              <CopyButton text={createdSite.site_key} />
            </div>
            <div className="flex justify-end">
              <Button onClick={() => setShowCreate(false)}>Done</Button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleCreate} className="space-y-4">
            <FormField label="Name" value={name} onChange={(e) => setName(e.target.value)} required />
            <FormField label="Domain" value={domain} onChange={(e) => setDomain(e.target.value)} placeholder="example.com" required />
            <div className="flex justify-end gap-3 pt-2">
              <Button variant="ghost" type="button" onClick={() => setShowCreate(false)}>Cancel</Button>
              <Button type="submit" loading={creating}>Create</Button>
            </div>
          </form>
        )}
      </Modal>
    </div>
  )
}
