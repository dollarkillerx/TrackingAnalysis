import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts'
import type { HourlyCount } from '@/types'

export function HourlyChart({ title, clickData, eventData }: {
  title: string
  clickData: HourlyCount[]
  eventData: HourlyCount[]
}) {
  const clickMap = new Map<number, number>()
  const eventMap = new Map<number, number>()
  clickData?.forEach((d) => clickMap.set(d.hour, d.count))
  eventData?.forEach((d) => eventMap.set(d.hour, d.count))

  const data = Array.from({ length: 24 }, (_, h) => ({
    hour: `${h}:00`,
    clicks: clickMap.get(h) ?? 0,
    events: eventMap.get(h) ?? 0,
  }))

  const hasData = data.some((d) => d.clicks > 0 || d.events > 0)

  return (
    <div className="rounded-xl border border-border bg-bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-text">{title}</h3>
      {hasData ? (
        <ResponsiveContainer width="100%" height={250}>
          <BarChart data={data} margin={{ top: 10, right: 20, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
            <XAxis dataKey="hour" tick={{ fontSize: 10 }} />
            <YAxis tick={{ fontSize: 11 }} />
            <Tooltip />
            <Legend />
            <Bar dataKey="clicks" fill="#6366f1" radius={[2, 2, 0, 0]} />
            <Bar dataKey="events" fill="#8b5cf6" radius={[2, 2, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      ) : (
        <p className="text-sm text-muted">-</p>
      )}
    </div>
  )
}
