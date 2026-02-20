import type { Tracker } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useTrackers() {
  const { data, loading, fetch } = useRpcList<Tracker>(RPC.TRACKER_LIST)
  const { execute: create, loading: creating } = useRpcCall<Tracker>(RPC.TRACKER_CREATE)
  const { execute: update, loading: updating } = useRpcCall<Tracker>(RPC.TRACKER_UPDATE)
  const { execute: remove, loading: deleting } = useRpcCall<{ ok: boolean }>(RPC.TRACKER_DELETE)

  return { trackers: data, loading, fetch, create, creating, update, updating, remove, deleting }
}
