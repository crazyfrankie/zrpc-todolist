package middleware

import (
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Trace(service string, opts ...otelgin.Option) gin.HandlerFunc {
	return otelgin.Middleware(service, opts...)
}
