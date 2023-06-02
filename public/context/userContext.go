package context

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gconv"
	"github.com/mssola/user_agent"
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

func (uc *UserContext) GetUserAgent() *user_agent.UserAgent {
	ua := user_agent.New(uc.Request.Header.Get("User-Agent"))
	//uc.ClientIP()
	//ua.OS()
	//ipUtils.GetRealAddressByIP(ip)
	//ua.Browser()
	return ua
}

func (uc *UserContext) ParamInt64(key string) int64 {
	return gconv.Int64(uc.Param(key))
}

func (uc *UserContext) ParamInt64Array(key string) []int64 {
	return gconv.Int64s(strings.Split(uc.Param(key), ","))
}

func (uc *UserContext) ParamStringArray(key string) []string {
	return strings.Split(uc.Param(key), ",")
}

func (uc *UserContext) QueryInt64(key string) int64 {
	return gconv.Int64(uc.Query(key))
}

func (uc *UserContext) QueryInt64Array(key string) []int64 {
	return gconv.Int64s(strings.Split(uc.Query(key), ","))
}

func (uc *UserContext) QueryStringArray(key string) []string {
	return strings.Split(uc.Query(key), ",")
}
