import type { Site } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useSites() {
  const { data, loading, fetch } = useRpcList<Site>(RPC.SITE_LIST)
  const { execute: create, loading: creating } = useRpcCall<Site>(RPC.SITE_CREATE)

  return { sites: data, loading, fetch, create, creating }
}
