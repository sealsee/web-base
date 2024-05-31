package utils

import (
	"github.com/sealsee/web-base/public/jwt"
	cacheUtil "github.com/sealsee/web-base/public/utils/cache"
	"github.com/sealsee/web-base/public/utils/mq/kafka"
	"github.com/sealsee/web-base/public/utils/mq/rabbitmq"
	"github.com/sealsee/web-base/public/utils/redis"
	"github.com/sealsee/web-base/public/utils/snowflake"
	_ "github.com/sealsee/web-base/public/utils/sys"
	"github.com/sealsee/web-base/public/utils/token"
)

func Init() {
	snowflake.Init()
	token.Init()
	jwt.Init()
	redis.Init()
	cacheUtil.Init()
	rabbitmq.Init()
	kafka.Init()
}
