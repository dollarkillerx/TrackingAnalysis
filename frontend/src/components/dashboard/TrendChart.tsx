import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from 'recharts'

export interface TrendSeries {
  key: string
  color: string
  name: string
}

export function TrendChart({ title, data, series, height = 300, stacked = false }: {
  title: string
  data: Record<string, unknown>[]
  series: TrendSeries[]
  height?: number
  stacked?: boolean
}) {
  return (
    <div className="rounded-xl border border-border bg-bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-text">{title}</h3>
      {data.length > 0 ? (
        <ResponsiveContainer width="100%" height={height}>
          <AreaChart data={data} margin={{ top: 10, right: 20, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
            <XAxis dataKey="date" tick={{ fontSize: 11 }} />
            <YAxis tick={{ fontSize: 11 }} />
            <Tooltip />
            {series.map((s) => (
              <Area
                key={s.key}
                type="monotone"
                dataKey={s.key}
                stackId={stacked ? '1' : undefined}
                stroke={s.color}
                fill={s.color}
                fillOpacity={0.2}
                name={s.name}
              />
            ))}
          </AreaChart>
        </ResponsiveContainer>
      ) : (
        <p className="text-sm text-muted">-</p>
      )}
    </div>
  )
}
