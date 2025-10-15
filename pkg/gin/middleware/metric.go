package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/response"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
)

func Metric() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		path := c.FullPath()
		if c.Writer.Status() == http.StatusNotFound {
			metrics.HttpCall("<404>", c.Request.Method, c.Writer.Status())
		} else {
			metrics.HttpCall(path, c.Request.Method, c.Writer.Status())
		}
		if resp := response.GetApiResp(c); resp != nil {
			metrics.APICall(path, c.Request.Method, int(resp.Code))
		}
	}
}
