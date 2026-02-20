import type { LucideIcon } from 'lucide-react'

export function StatCard({ icon: Icon, label, value, loading }: {
  icon: LucideIcon
  label: string
  value: string | number
  loading: boolean
}) {
  return (
    <div className="rounded-xl border border-border bg-bg-card px-4 py-3">
      <div className="flex items-center gap-3">
        <div className="rounded-lg bg-primary/10 p-2">
          <Icon className="h-5 w-5 text-primary" />
        </div>
        <div className="min-w-0">
          <p className="text-xl font-bold font-mono text-text truncate">
            {loading ? '...' : value}
          </p>
          <p className="text-xs text-muted truncate">{label}</p>
        </div>
      </div>
    </div>
  )
}
