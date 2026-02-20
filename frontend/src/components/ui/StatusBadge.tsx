import { classNames } from '@/lib/utils'

interface StatusBadgeProps {
  status: string
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const normalized = status.toLowerCase()
  return (
    <span
      className={classNames(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        normalized === 'active' && 'bg-success/15 text-success',
        normalized === 'inactive' && 'bg-warning/15 text-warning',
        normalized !== 'active' && normalized !== 'inactive' && 'bg-muted/15 text-muted',
      )}
    >
      {status}
    </span>
  )
}
