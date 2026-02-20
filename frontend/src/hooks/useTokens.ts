import type { Token } from '@/types'
import { RPC } from '@/lib/constants'
import { useRpcList, useRpcCall } from './useRpc'

export function useTokens() {
  const { data, loading, fetch } = useRpcList<Token>(RPC.TOKEN_LIST)
  const { execute: generate, loading: generating } = useRpcCall<Token>(RPC.TOKEN_GENERATE)
  const { execute: remove, loading: deleting } = useRpcCall<{ ok: boolean }>(RPC.TOKEN_DELETE)

  return { tokens: data, loading, fetch, generate, generating, remove, deleting }
}
