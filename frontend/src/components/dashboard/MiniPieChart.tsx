import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts'
import type { NameCount } from '@/types'

const COLORS = ['#6366f1', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#ef4444', '#14b8a6']

export function MiniPieChart({ title, data, maxItems = 6 }: {
  title: string
  data: NameCount[]
  maxItems?: number
}) {
  if (!data || data.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-bg-card p-4">
        <h3 className="mb-3 text-sm font-semibold text-text">{title}</h3>
        <p className="text-sm text-muted">-</p>
      </div>
    )
  }

  const top = data.slice(0, maxItems)
  const otherCount = data.slice(maxItems).reduce((sum, d) => sum + d.count, 0)
  const chartData = otherCount > 0 ? [...top, { name: 'Other', count: otherCount }] : top
  const total = chartData.reduce((sum, d) => sum + d.count, 0)

  return (
    <div className="rounded-xl border border-border bg-bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-text">{title}</h3>
      <div className="flex items-start gap-4">
        <div className="w-24 h-24 shrink-0">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={chartData}
                dataKey="count"
                nameKey="name"
                cx="50%"
                cy="50%"
                innerRadius={20}
                outerRadius={40}
                strokeWidth={1}
              >
                {chartData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
        <div className="space-y-1 min-w-0 flex-1">
          {chartData.map((item, i) => (
            <div key={item.name} className="flex items-center gap-2 text-xs">
              <span
                className="inline-block h-2.5 w-2.5 shrink-0 rounded-sm"
                style={{ backgroundColor: COLORS[i % COLORS.length] }}
              />
              <span className="truncate text-text">{item.name}</span>
              <span className="ml-auto shrink-0 font-mono text-muted">
                {total > 0 ? Math.round((item.count / total) * 100) : 0}%
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
