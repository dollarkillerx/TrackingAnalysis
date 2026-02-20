import { useState, useCallback } from 'react'
import { rpcCall } from '@/lib/rpc-client'
import { RPC } from '@/lib/constants'
import type { ClickStatsResponse, EventStatsResponse } from '@/types'

export type Period = 'today' | '7d' | '30d' | 'custom'

export function getDateRange(period: Exclude<Period, 'custom'>): { startDate: string; endDate: string } {
  const today = new Date()
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  const endDate = fmt(today)

  switch (period) {
    case 'today':
      return { startDate: endDate, endDate }
    case '7d': {
      const start = new Date(today)
      start.setDate(start.getDate() - 6)
      return { startDate: fmt(start), endDate }
    }
    case '30d': {
      const start = new Date(today)
      start.setDate(start.getDate() - 29)
      return { startDate: fmt(start), endDate }
    }
  }
}

export function useClickStats() {
  const [data, setData] = useState<ClickStatsResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetch = useCallback(async (params: { start_date: string; end_date: string; tracker_id?: string; campaign_id?: string; channel_id?: string }) => {
    setLoading(true)
    setError(null)
    try {
      const result = await rpcCall<ClickStatsResponse>(RPC.STATS_CLICKS, params)
      setData(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      setData(null)
    } finally {
      setLoading(false)
    }
  }, [])

  return { data, loading, error, fetch }
}

export function useEventStats() {
  const [data, setData] = useState<EventStatsResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetch = useCallback(async (params: { start_date: string; end_date: string; site_id?: string }) => {
    setLoading(true)
    setError(null)
    try {
      const result = await rpcCall<EventStatsResponse>(RPC.STATS_EVENTS, params)
      setData(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      setData(null)
    } finally {
      setLoading(false)
    }
  }, [])

  return { data, loading, error, fetch }
}
