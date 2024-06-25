package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/web"
)

var cp ICheckPrivilege

type ICheckPrivilege interface {
	Check(path, needRoles, needPermissions string, sessionUser *context.SessionUser) bool
}

func SetCheckPrivilege(check ICheckPrivilege) {
	cp = check
}

func PrivilegeMiddleware(needRoles, needPermissions string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if needRoles == "" && needPermissions == "" {
			c.Next()
		}

		path := c.Request.URL.Path
		if path == "" || path == "/" {
			return
		}

		user := context.NewUserContext(c).GetUser()
		if cp == nil {
			return
		}
		if cp.Check(path, needRoles, needPermissions, user) {
			c.Next()
		} else {
			json := web.NewJsonResult(c)
			json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
			c.Abort()
			return
		}
	}

}
