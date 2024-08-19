package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sealsee/web-base/public/ds"
	cacheUtil "github.com/sealsee/web-base/public/utils/cache"
	"go.uber.org/zap"

	"time"
)

var client *redis.Client
var useMemory bool

func Init() {
	client = ds.GetRedisClient()
	// 兼容内存组件
	if client == nil {
		useMemory = true
	}
}

func handelErr(err error) {
	if err != nil {
		zap.L().Error("Redis存储失败", zap.Error(err))
		panic(fmt.Sprintf("Redis存储失败: %v", zap.Error(err)))
	}
}

func Expire(key string, expiration time.Duration) {
	if useMemory {
		cacheUtil.Expire(key, expiration)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	err := client.Expire(key, expiration).Err()
	handelErr(err)
}

func Exists(key string) bool {
	if useMemory {
		return cacheUtil.Exists(key)
	}
	return client.Exists(key).Val() > 0
}

func SetString(key, str string, expiration time.Duration) {
	if useMemory {
		cacheUtil.SetString(key, str, expiration)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	err := client.Set(key, str, expiration).Err()
	handelErr(err)
}

func GetString(key string) string {
	if useMemory {
		return cacheUtil.GetString(key)
	}
	return client.Get(key).Val()
}

func SetStruct(key string, value interface{}, expiration time.Duration) {
	if useMemory {
		cacheUtil.SetStruct(key, value, expiration)
		return
	}
	marshal, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("Redis存储失败", zap.Error(err))
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	err = client.Set(key, marshal, expiration).Err()
	handelErr(err)
}

func GetStruct[T any](key string, t *T) (*T, error) {
	if useMemory {
		return cacheUtil.GetStruct(key, t)
	}
	newT := new(T)
	structJsonString := GetString(key)
	err := json.Unmarshal([]byte(structJsonString), newT)
	return newT, err
}

func Hset(key, field string, value interface{}) {
	if useMemory {
		cacheUtil.Hset(key, field, value)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	res, _ := json.Marshal(value)
	err := client.HSet(key, field, res).Err()
	handelErr(err)
}

// func Hmset(key string, fields map[string]interface{}) {
// 	if useMemory {
// 		cacheUtil.Hmset(key, fields)
// 		return
// 	}
// 	defer func() {
// 		if err := recover(); err != nil {
// 			zap.L().Error("redisUtils", zap.Any("error", err))
// 			panic(fmt.Sprintf("redisUtils: %v", zap.Any("error", err)))
// 		}
// 	}()
// 	err := client.HMSet(key, fields).Err()
// 	handelErr(err)
// }

func Hget[T any](key, field string) T {
	if useMemory {
		return cacheUtil.Hget[T](key, field)
	}
	newT := new(T)
	val := client.HGet(key, field).Val()
	_ = json.Unmarshal([]byte(val), newT)
	return *newT
}

func HGetAll(key string) map[string]string {
	if useMemory {
		return cacheUtil.HGetAll(key)
	}
	return client.HGetAll(key).Val()
}

func Hdel(key string, field ...string) {
	if useMemory {
		cacheUtil.Hdel(key, field...)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	client.HDel(key, field...)
}

func Del(key ...string) {
	if useMemory {
		cacheUtil.Del(key...)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("redisUtils", zap.Any("error", err))
		}
	}()
	client.Del(key...)
}
