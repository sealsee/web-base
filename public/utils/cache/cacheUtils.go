package cacheUtil

import (
	"encoding/json"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/sealsee/web-base/public/ds"
	"go.uber.org/zap"

	"time"
)

var client *cache.Cache

func Init() {
	client = ds.GetCacheMemory()
}

func Expire(key string, expiration time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils Expire", zap.Any("error", err))
		}
	}()
	val, found := client.Get(key)
	if found {
		client.Set(key, val, expiration)
	}
}

func Exists(key string) bool {
	_, found := client.Get(key)
	return found
}

func SetString(key, str string, expiration time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils SetString", zap.Any("error", err))
		}
	}()
	client.Set(key, str, expiration)
}

func GetString(key string) string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils GetString", zap.Any("error", err))
		}
	}()
	val, found := client.Get(key)
	if found {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func SetStruct(key string, value interface{}, expiration time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils SetStruct", zap.Any("error", err))
		}
	}()
	marshal, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("cache存储失败", zap.Error(err))
	}
	client.Set(key, string(marshal), expiration)
}

func GetStruct[T any](key string, t *T) (*T, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils GetStruct", zap.Any("error", err))
		}
	}()
	newT := new(T)
	structJsonString := GetString(key)
	err := json.Unmarshal([]byte(structJsonString), newT)
	return newT, err
}

func Hset(key, field string, value interface{}) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils Hset", zap.Any("error", err))
		}
	}()
	res, _ := json.Marshal(value)
	if val, found := client.Get(key); found {
		if valMap, ok := val.(map[string]string); ok {
			valMap[field] = string(res)
			client.Set(key, valMap, cache.NoExpiration)
		}
	} else {
		hmap := make(map[string]string)
		hmap[field] = string(res)
		client.Set(key, hmap, cache.NoExpiration)
	}
}

// func Hmset(key string, fields map[string]interface{}) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			zap.L().Error("cacheUtils Hmset", zap.Any("error", err))
// 		}
// 	}()
// 	client.Set(key, fields, cache.NoExpiration)
// }

func Hget[T any](key, field string) T {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils HGet", zap.Any("error", err))
		}
	}()
	// var res T
	newT := new(T)
	if val, found := client.Get(key); found {
		if valMap, ok := val.(map[string]string); ok {
			hval := valMap[field]
			// if reflect.ValueOf(hval).Kind() == reflect.Ptr {
			// 	if res, ok := hval.(*T); ok {
			// 		return *res
			// 	}
			// } else {
			// 	if res, ok := hval.(T); ok {
			// 		return res
			// 	}
			// }
			_ = json.Unmarshal([]byte(hval), newT)
		}
	}
	return *newT
}

func HGetAll(key string) map[string]string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils HGetAll", zap.Any("error", err))
		}
	}()
	if val, found := client.Get(key); found {
		if valMap, ok := val.(map[string]string); ok {
			return valMap
		}
	}
	return nil
}

func Hdel(key string, field ...string) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils Hdel", zap.Any("error", err))
		}
	}()
	if val, found := client.Get(key); found {
		if valMap, ok := val.(map[string]string); ok {
			for _, f := range field {
				delete(valMap, f)
			}
			client.Set(key, valMap, cache.NoExpiration)
		}
	}
}

func Del(key ...string) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("cacheUtils Del", zap.Any("error", err))
		}
	}()
	for _, k := range key {
		client.Delete(k)
	}
}
