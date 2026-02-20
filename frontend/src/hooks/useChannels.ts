import type { Channel } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useChannels() {
  const { data, loading, fetch } = useRpcList<Channel>(RPC.CHANNEL_LIST)
  const { execute: create, loading: creating } = useRpcCall<Channel>(RPC.CHANNEL_CREATE)
  const { execute: batchImport, loading: importing } = useRpcCall<{ imported: number }>(RPC.CHANNEL_BATCH_IMPORT)

  return { channels: data, loading, fetch, create, creating, batchImport, importing }
}
