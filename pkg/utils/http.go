package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	requestid "github.com/sumit-tembe/gin-requestid"
)

// Get request id from echo context
func GetRequestID(c *gin.Context) string {
	return requestid.GetRequestIDFromHeaders(c)
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// Get context  with request id
func GetRequestCtx(c *gin.Context) context.Context {
	ctx := context.WithValue(c.Request.Context(), ReqIDCtxKey{}, GetRequestID(c))
	for key, value := range c.Keys {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}

func WithTelemetry(c *gin.Context, key string, mainFunc func(spanCtx context.Context)) {
	span, ctx := opentracing.StartSpanFromContext(GetRequestCtx(c), key)
	span.SetTag("requestId", GetRequestID(c))
	defer span.Finish()
	mainFunc(ctx)
}
