package middlewares

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/web"
)

var regx *regexp.Regexp

func init() {
	regx, _ = regexp.Compile(`[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}`)
}

func AuthDebugMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		clientIP := c.Request.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.Request.Header.Get("X-Real-IP")
		}

		if clientIP == "" {
			clientIP = c.ClientIP()
		}

		if clientIP != "" && regx.Match([]byte(clientIP)) {
			if strings.HasPrefix(clientIP, "192.168") ||
				strings.HasPrefix(clientIP, "10.0") ||
				strings.HasPrefix(clientIP, "127.0.0.1") {
				fmt.Println("debug...")

			} else {
				json := web.NewJsonResult(c)
				json.SetErrs(errs.REFUSE_VISIT_ERR).RenderWithCode(http.StatusUnauthorized)
				c.Abort()
			}
		}
	}
}
