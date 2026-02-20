package rpc

// Application error codes
const (
	ErrCodeInvalidToken   = 4001
	ErrCodeExpiredToken   = 4002
	ErrCodeRateLimited    = 4003
	ErrCodeReplayDetected = 4004
	ErrCodeDecryptFailed  = 4005
	ErrCodeBotBlocked     = 4006
	ErrCodeDBError        = 5001
)

// Standard JSON-RPC 2.0 error codes
const (
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603
	ErrCodeParseError     = -32700
)

// Error messages
var errMessages = map[int]string{
	ErrCodeInvalidToken:   "invalid_token",
	ErrCodeExpiredToken:   "expired_token",
	ErrCodeRateLimited:    "rate_limited",
	ErrCodeReplayDetected: "replay_detected",
	ErrCodeDecryptFailed:  "decrypt_failed",
	ErrCodeBotBlocked:     "bot_blocked",
	ErrCodeDBError:        "database_error",
	ErrCodeInvalidRequest: "invalid_request",
	ErrCodeMethodNotFound: "method_not_found",
	ErrCodeInvalidParams:  "invalid_params",
	ErrCodeInternalError:  "internal_error",
	ErrCodeParseError:     "parse_error",
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	return e.Message
}

func NewRPCError(code int, data interface{}) *RPCError {
	msg, ok := errMessages[code]
	if !ok {
		msg = "unknown_error"
	}
	return &RPCError{Code: code, Message: msg, Data: data}
}

func NewRPCErrorWithMessage(code int, message string) *RPCError {
	return &RPCError{Code: code, Message: message}
}
