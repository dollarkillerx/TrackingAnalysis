package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRPC(t *testing.T, d *Dispatcher, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/rpc", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	d.GinHandler()(c)
	return w
}

func TestDispatcher_Register(t *testing.T) {
	d := NewDispatcher()
	called := false
	d.Register("test.method", func(ctx context.Context, params json.RawMessage) (interface{}, *RPCError) {
		called = true
		return "ok", nil
	})

	body, _ := json.Marshal(Request{
		JSONRPC: "2.0",
		Method:  "test.method",
		ID:      1,
	})
	performRPC(t, d, body)
	if !called {
		t.Error("registered handler was not called")
	}
}

func TestDispatcher_ValidRequest(t *testing.T) {
	d := NewDispatcher()
	d.Register("echo", func(ctx context.Context, params json.RawMessage) (interface{}, *RPCError) {
		return map[string]string{"echo": "hello"}, nil
	})

	body, _ := json.Marshal(Request{
		JSONRPC: "2.0",
		Method:  "echo",
		Params:  json.RawMessage(`{}`),
		ID:      42,
	})
	w := performRPC(t, d, body)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.JSONRPC != "2.0" {
		t.Errorf("jsonrpc = %q, want 2.0", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Errorf("unexpected error: %+v", resp.Error)
	}
	if resp.Result == nil {
		t.Fatal("result is nil")
	}
}

func TestDispatcher_InvalidJSON(t *testing.T) {
	d := NewDispatcher()
	w := performRPC(t, d, []byte(`{invalid json`))

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error == nil {
		t.Fatal("expected parse error")
	}
	if resp.Error.Code != ErrCodeParseError {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeParseError)
	}
}

func TestDispatcher_WrongVersion(t *testing.T) {
	d := NewDispatcher()
	body, _ := json.Marshal(Request{
		JSONRPC: "1.0",
		Method:  "test",
		ID:      1,
	})
	w := performRPC(t, d, body)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == nil {
		t.Fatal("expected invalid request error")
	}
	if resp.Error.Code != ErrCodeInvalidRequest {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeInvalidRequest)
	}
}

func TestDispatcher_MethodNotFound(t *testing.T) {
	d := NewDispatcher()
	body, _ := json.Marshal(Request{
		JSONRPC: "2.0",
		Method:  "nonexistent.method",
		ID:      1,
	})
	w := performRPC(t, d, body)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == nil {
		t.Fatal("expected method not found error")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestDispatcher_HandlerError(t *testing.T) {
	d := NewDispatcher()
	d.Register("fail", func(ctx context.Context, params json.RawMessage) (interface{}, *RPCError) {
		return nil, NewRPCError(ErrCodeRateLimited, nil)
	})

	body, _ := json.Marshal(Request{
		JSONRPC: "2.0",
		Method:  "fail",
		ID:      1,
	})
	w := performRPC(t, d, body)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == nil {
		t.Fatal("expected error from handler")
	}
	if resp.Error.Code != ErrCodeRateLimited {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeRateLimited)
	}
	if resp.Error.Message != "rate_limited" {
		t.Errorf("error message = %q, want %q", resp.Error.Message, "rate_limited")
	}
}
