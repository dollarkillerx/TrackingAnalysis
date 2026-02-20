package rpc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type HandlerFunc func(ctx context.Context, params json.RawMessage) (interface{}, *RPCError)

type Dispatcher struct {
	methods map[string]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		methods: make(map[string]HandlerFunc),
	}
}

func (d *Dispatcher) Register(method string, handler HandlerFunc) {
	d.methods[method] = handler
}

func (d *Dispatcher) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, Response{
				JSONRPC: "2.0",
				Error:   NewRPCError(ErrCodeParseError, err.Error()),
				ID:      nil,
			})
			return
		}

		if req.JSONRPC != "2.0" {
			c.JSON(http.StatusOK, Response{
				JSONRPC: "2.0",
				Error:   NewRPCError(ErrCodeInvalidRequest, "jsonrpc must be 2.0"),
				ID:      req.ID,
			})
			return
		}

		handler, ok := d.methods[req.Method]
		if !ok {
			c.JSON(http.StatusOK, Response{
				JSONRPC: "2.0",
				Error:   NewRPCError(ErrCodeMethodNotFound, req.Method),
				ID:      req.ID,
			})
			return
		}

		// Pass gin context values into context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ginContextKey, c)

		result, rpcErr := handler(ctx, req.Params)
		resp := Response{
			JSONRPC: "2.0",
			ID:      req.ID,
		}
		if rpcErr != nil {
			resp.Error = rpcErr
		} else {
			resp.Result = result
		}
		c.JSON(http.StatusOK, resp)
	}
}

type contextKey string

const ginContextKey contextKey = "gin"

func GinContext(ctx context.Context) *gin.Context {
	val := ctx.Value(ginContextKey)
	if c, ok := val.(*gin.Context); ok {
		return c
	}
	return nil
}
