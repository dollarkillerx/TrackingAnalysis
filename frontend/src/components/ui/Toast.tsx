import { useToast } from '@/context/ToastContext'
import { CheckCircle, XCircle, AlertTriangle, X } from 'lucide-react'
import { classNames } from '@/lib/utils'

const icons = {
  success: CheckCircle,
  error: XCircle,
  warning: AlertTriangle,
}

const styles = {
  success: 'border-success/30 bg-success/10 text-success',
  error: 'border-error/30 bg-error/10 text-error',
  warning: 'border-warning/30 bg-warning/10 text-warning',
}

export function ToastContainer() {
  const { toasts, removeToast } = useToast()

  if (toasts.length === 0) return null

  return (
    <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2">
      {toasts.map((toast) => {
        const Icon = icons[toast.type]
        return (
          <div
            key={toast.id}
            className={classNames(
              'flex items-center gap-3 rounded-lg border px-4 py-3 shadow-lg transition-all duration-300',
              styles[toast.type],
            )}
          >
            <Icon className="h-5 w-5 shrink-0" />
            <span className="text-sm text-text">{toast.message}</span>
            <button
              onClick={() => removeToast(toast.id)}
              className="ml-2 rounded p-0.5 hover:bg-bg-elevated/50 transition-colors cursor-pointer"
            >
              <X className="h-4 w-4 text-muted" />
            </button>
          </div>
        )
      })}
    </div>
  )
}
