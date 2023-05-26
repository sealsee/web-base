package middlewares

import (
	"time"

	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/jwt"
	"github.com/sealsee/web-base/public/utils/token"
	"github.com/sealsee/web-base/public/web"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		json := web.NewJsonResult(c)
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			json.SetErrs(errs.UNAUTHORIZED).Render()
			c.Abort()
			return
		}
		sessionUser, err := jwt.ParseToken(authHeader)
		if err != nil {
			json.SetErrs(errs.UNAUTHORIZED).Render()
			c.Abort()
			return
		}
		// 当前时间距离session过期时间不足15分钟时，刷新session，重新设置为30分钟
		if sessionUser.ExpireTime < time.Now().Add(time.Duration(15)*time.Minute).Unix() {
			go token.RefreshToken(sessionUser)
		}
		// 将sessionUser放入上下文
		c.Set(cst.LoginUserKey, sessionUser)
		c.Next()
	}
}
