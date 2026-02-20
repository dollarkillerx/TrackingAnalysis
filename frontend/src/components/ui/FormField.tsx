import { classNames } from '@/lib/utils'
import type { InputHTMLAttributes } from 'react'

interface FormFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export function FormField({ label, error, className, id, ...props }: FormFieldProps) {
  const fieldId = id || label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div className="space-y-1">
      <label htmlFor={fieldId} className="block text-sm font-medium text-muted">
        {label}
      </label>
      <input
        id={fieldId}
        className={classNames(
          'w-full rounded-lg border bg-bg-card px-3 py-2 text-sm text-text placeholder:text-muted/60 outline-none transition-colors duration-150 focus:border-primary focus:ring-1 focus:ring-primary',
          error ? 'border-error' : 'border-border',
          className,
        )}
        {...props}
      />
      {error && <p className="text-xs text-error">{error}</p>}
    </div>
  )
}
