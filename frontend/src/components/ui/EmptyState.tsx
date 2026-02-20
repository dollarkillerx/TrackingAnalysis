import type { ReactNode } from 'react'
import type { LucideIcon } from 'lucide-react'
import { Inbox } from 'lucide-react'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  action?: ReactNode
}

export function EmptyState({ icon: Icon = Inbox, title, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <Icon className="h-12 w-12 text-muted/40 mb-4" />
      <p className="text-muted mb-4">{title}</p>
      {action}
    </div>
  )
}
