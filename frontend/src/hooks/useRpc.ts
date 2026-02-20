import { useState, useCallback } from 'react'
import { rpcCall } from '@/lib/rpc-client'
import i18n from '@/i18n'

export function useRpcCall<TResult>(method: string) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const execute = useCallback(
    async (params: Record<string, unknown> = {}): Promise<TResult> => {
      setLoading(true)
      setError(null)
      try {
        const result = await rpcCall<TResult>(method, params)
        return result
      } catch (err) {
        const msg = err instanceof Error ? err.message : i18n.t('common.unknownError')
        setError(msg)
        throw err
      } finally {
        setLoading(false)
      }
    },
    [method],
  )

  return { execute, loading, error }
}

export function useRpcList<TResult>(method: string) {
  const [data, setData] = useState<TResult[]>([])
  const [loading, setLoading] = useState(false)

  const fetch = useCallback(
    async (params: Record<string, unknown> = {}) => {
      setLoading(true)
      try {
        const result = await rpcCall<TResult[]>(method, params)
        setData(result ?? [])
      } catch {
        setData([])
      } finally {
        setLoading(false)
      }
    },
    [method],
  )

  return { data, loading, fetch }
}
