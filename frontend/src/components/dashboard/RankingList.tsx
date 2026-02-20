import { truncate } from '@/lib/utils'

export function RankingList({ title, data, maxItems = 10 }: {
  title: string
  data: { name: string; count: number }[]
  maxItems?: number
}) {
  const items = data?.slice(0, maxItems) ?? []
  const max = items.length > 0 ? items[0].count : 0

  return (
    <div className="rounded-xl border border-border bg-bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-text">{title}</h3>
      {items.length === 0 ? (
        <p className="text-sm text-muted">-</p>
      ) : (
        <div className="space-y-1">
          {items.map((item) => {
            const pct = max > 0 ? (item.count / max) * 100 : 0
            return (
              <div key={item.name} className="relative flex items-center justify-between rounded px-2 py-1 text-sm">
                <div
                  className="absolute inset-0 rounded bg-primary/10"
                  style={{ width: `${pct}%` }}
                />
                <span className="relative z-10 truncate text-text" title={item.name}>
                  {truncate(item.name, 30)}
                </span>
                <span className="relative z-10 ml-2 shrink-0 font-mono text-xs text-muted">
                  {item.count.toLocaleString()}
                </span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
