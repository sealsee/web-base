package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/web"
)

var cp ICheckPrivilege

type ICheckPrivilege interface {
	Check(path string, loginUser *context.SessionUser) bool
}

func SetCheckPrivilege(check ICheckPrivilege) {
	cp = check
}

func PrivilegeMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "" || path == "/" {
			return
		}

		user := context.NewUserContext(c).GetUser()
		if cp == nil {
			return
		}
		if cp.Check(path, user) {
			c.Next()
		} else {
			json := web.NewJsonResult(c)
			json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
			c.Abort()
			return
		}
	}

}
