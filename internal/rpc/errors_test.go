package rpc

import (
	"testing"
)

func TestNewRPCError_KnownCode(t *testing.T) {
	tests := []struct {
		code    int
		wantMsg string
	}{
		{ErrCodeInvalidToken, "invalid_token"},
		{ErrCodeExpiredToken, "expired_token"},
		{ErrCodeRateLimited, "rate_limited"},
		{ErrCodeReplayDetected, "replay_detected"},
		{ErrCodeDecryptFailed, "decrypt_failed"},
		{ErrCodeBotBlocked, "bot_blocked"},
		{ErrCodeDBError, "database_error"},
		{ErrCodeInvalidRequest, "invalid_request"},
		{ErrCodeMethodNotFound, "method_not_found"},
		{ErrCodeInvalidParams, "invalid_params"},
		{ErrCodeInternalError, "internal_error"},
		{ErrCodeParseError, "parse_error"},
	}
	for _, tt := range tests {
		e := NewRPCError(tt.code, nil)
		if e.Code != tt.code {
			t.Errorf("code %d: Code = %d", tt.code, e.Code)
		}
		if e.Message != tt.wantMsg {
			t.Errorf("code %d: Message = %q, want %q", tt.code, e.Message, tt.wantMsg)
		}
	}
}

func TestNewRPCError_UnknownCode(t *testing.T) {
	e := NewRPCError(9999, nil)
	if e.Message != "unknown_error" {
		t.Errorf("Message = %q, want %q", e.Message, "unknown_error")
	}
}

func TestNewRPCErrorWithMessage(t *testing.T) {
	e := NewRPCErrorWithMessage(4001, "custom message")
	if e.Code != 4001 {
		t.Errorf("Code = %d, want 4001", e.Code)
	}
	if e.Message != "custom message" {
		t.Errorf("Message = %q, want %q", e.Message, "custom message")
	}
}

func TestRPCError_ErrorInterface(t *testing.T) {
	e := NewRPCError(ErrCodeRateLimited, nil)
	var err error = e // verify it implements error interface
	if err.Error() != "rate_limited" {
		t.Errorf("Error() = %q, want %q", err.Error(), "rate_limited")
	}
}
