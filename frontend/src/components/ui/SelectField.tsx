import { classNames } from '@/lib/utils'
import type { SelectHTMLAttributes } from 'react'

interface Option {
  value: string
  label: string
}

interface SelectFieldProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string
  options: Option[]
  error?: string
  placeholder?: string
}

export function SelectField({ label, options, error, placeholder, className, id, ...props }: SelectFieldProps) {
  const fieldId = id || label.toLowerCase().replace(/\s+/g, '-')
  return (
    <div className="space-y-1">
      <label htmlFor={fieldId} className="block text-sm font-medium text-muted">
        {label}
      </label>
      <select
        id={fieldId}
        className={classNames(
          'w-full rounded-lg border bg-bg-card px-3 py-2 text-sm text-text outline-none transition-colors duration-150 focus:border-primary focus:ring-1 focus:ring-primary',
          error ? 'border-error' : 'border-border',
          className,
        )}
        {...props}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && <p className="text-xs text-error">{error}</p>}
    </div>
  )
}
