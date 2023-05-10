package redis

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/sealsee/web-base/public/ds"
	"go.uber.org/zap"

	"time"
)

var client *redis.Client

func Init() {
	client = ds.GetRedisClient()
}

func handelErr(err error) {
	if err != nil {
		zap.L().Error("Redis存储失败", zap.Error(err))
	}
}

func Expire(key string, expiration time.Duration) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		err := client.Expire(key, expiration).Err()
		handelErr(err)
	}()
}

func Exists(key string) bool {
	return client.Exists(key).Val() != -1
}

func SetString(key, str string, expiration time.Duration) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		err := client.Set(key, str, expiration).Err()
		handelErr(err)
	}()
}

func GetString(key string) string {
	return client.Get(key).Val()
}

func SetStruct(key string, value interface{}, expiration time.Duration) {
	marshal, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("Redis存储失败", zap.Error(err))
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		err = client.Set(key, marshal, expiration).Err()
		handelErr(err)
	}()
}

func GetStruct[T any](key string, t *T) (*T, error) {
	newT := new(T)
	LoginUserJson := GetString(key)
	err := json.Unmarshal([]byte(LoginUserJson), newT)
	return newT, err
}

func Hset(key, field string, value interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		err := client.HSet(key, field, value).Err()
		handelErr(err)
	}()
}

func Hmset(key string, fields map[string]interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		err := client.HMSet(key, fields).Err()
		handelErr(err)
	}()
}

func Hget[T any](key, field string, _ *T) *T {
	newT := new(T)
	val := client.HGet(key, field).Val()
	_ = json.Unmarshal([]byte(val), newT)
	return newT
}

func HGetAll(key string) map[string]string {
	return client.HGetAll(key).Val()
}

func Hdel(key string, field ...string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		client.HDel(key, field...)
	}()
}

func Del(key ...string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("redisUtils", zap.Any("error", err))
			}
		}()
		client.Del(key...)
	}()
}
