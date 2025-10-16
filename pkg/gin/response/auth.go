package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetAuthorization(c *gin.Context, access, refresh string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.Header("x-access-token", access)
	c.SetCookie("zrpc-todolist-refresh", refresh, int(time.Hour*24), "/", "", false, true)
}
