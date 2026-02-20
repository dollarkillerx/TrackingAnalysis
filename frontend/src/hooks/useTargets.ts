import type { Target } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useTargets() {
  const { data, loading, fetch } = useRpcList<Target>(RPC.TARGET_LIST)
  const { execute: create, loading: creating } = useRpcCall<Target>(RPC.TARGET_CREATE)

  return { targets: data, loading, fetch, create, creating }
}
