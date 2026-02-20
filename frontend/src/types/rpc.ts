export interface RPCRequest {
  jsonrpc: '2.0'
  method: string
  params: Record<string, unknown>
  id: number
}

export interface RPCResponse<T> {
  jsonrpc: '2.0'
  result?: T
  error?: RPCError
  id: number
}

export interface RPCError {
  code: number
  message: string
  data?: unknown
}
