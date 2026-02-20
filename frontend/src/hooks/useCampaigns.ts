import type { Campaign } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useCampaigns() {
  const { data, loading, fetch } = useRpcList<Campaign>(RPC.CAMPAIGN_LIST)
  const { execute: create, loading: creating } = useRpcCall<Campaign>(RPC.CAMPAIGN_CREATE)

  return { campaigns: data, loading, fetch, create, creating }
}
