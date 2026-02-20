import type { RPCResponse } from '@/types/rpc'
import { TOKEN_KEY, RPC } from './constants'

let requestId = 0

export class RPCCallError extends Error {
  code: number
  data?: unknown

  constructor(code: number, message: string, data?: unknown) {
    super(message)
    this.name = 'RPCCallError'
    this.code = code
    this.data = data
  }
}

export async function rpcCall<T>(
  method: string,
  params: Record<string, unknown> = {},
): Promise<T> {
  const token = localStorage.getItem(TOKEN_KEY)

  if (token && method !== RPC.LOGIN) {
    params = { ...params, admin_token: token }
  }

  const res = await fetch('https://asn.siliconnexus.cc/rpc', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      jsonrpc: '2.0',
      method,
      params,
      id: ++requestId,
    }),
  })

  const body: RPCResponse<T> = await res.json()

  if (body.error) {
    if (body.error.code === 4001 || body.error.code === 4002) {
      localStorage.removeItem(TOKEN_KEY)
      window.location.href = '/login'
    }
    throw new RPCCallError(body.error.code, body.error.message, body.error.data)
  }

  return body.result as T
}
