package token

import (
	"fmt"
	"time"

	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/redis"
)

var timeLive time.Duration

func Init() {
	timeLive = time.Duration(setting.Conf.TokenConfig.ExpireTime) * time.Minute
}

func RefreshToken(sessionUser *context.SessionUser) {
	sessionUser.ExpireTime = time.Now().Add(timeLive).Unix()
	tokenKey := fmt.Sprintf(cst.LoginTokenKey+"%s", sessionUser.Token)
	redis.SetStruct(tokenKey, sessionUser, timeLive)
}
