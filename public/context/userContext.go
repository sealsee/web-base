package context

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gconv"
	"github.com/sealsee/web-base/public/cst"
)

type UserContext struct {
	*gin.Context
}

func NewUserContext(c *gin.Context) *UserContext {
	return &UserContext{c}
}

func (uc *UserContext) GetUser() (loginUser *SessionUser) {
	loginUserKey, _ := uc.Get(cst.LoginUserKey)
	if loginUserKey != nil {
		loginUser = loginUserKey.(*SessionUser)
	}
	return
}

func (uc *UserContext) GetUserName() string {
	user := uc.GetUser()
	if user == nil {
		return ""
	}
	return user.UserName
}
func (uc *UserContext) GetUserId() int64 {
	user := uc.GetUser()
	if user == nil {
		return 0
	}
	return user.UserId
}

func (uc *UserContext) QueryInt64(key string) int64 {
	return gconv.Int64(uc.Query(key))
}

// func (uc *UserContext) SetUserAgent(login *monitorModels.Logininfor) {
// 	login.InfoId = snowflake.GenID()
// 	ua := user_agent.New(uc.Request.Header.Get("User-Agent"))
// 	ip := uc.ClientIP()
// 	login.IpAddr = ip
// 	login.Os = ua.OS()
// 	login.LoginLocation = ipUtils.GetRealAddressByIP(ip)
// 	login.Browser, _ = ua.Browser()
// }

// func (bzc *UserContext) ParamInt64(key string) int64 {
// 	return gconv.Int64(bzc.Param(key))
// }
// func (bzc *UserContext) ParamInt64Array(key string) []int64 {
// 	return gconv.Int64s(strings.Split(bzc.Param(key), ","))
// }
// func (bzc *UserContext) ParamStringArray(key string) []string {
// 	return strings.Split(bzc.Param(key), ",")
// }

// func (bzc *UserContext) QueryInt64(key string) int64 {
// 	return gconv.Int64(bzc.Query(key))
// }
// func (bzc *UserContext) QueryInt64Array(key string) []int64 {
// 	return gconv.Int64s(strings.Split(bzc.Query(key), ","))

// }
