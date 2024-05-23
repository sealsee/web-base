package jsonUtils

import (
	"fmt"
	"reflect"

	"github.com/sealsee/web-base/public/utils/stringUtils"
)

// 将结构体转换为map，可以将任何结构体转换为map，并且能够处理嵌套结构体的情况
// 仅适用于gorm数据库对应字段，过滤掉gorm:"-"的字段，过滤掉0值
func StructToDbMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or struct pointer; got %v", v.Kind())
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		// 如果字段是一个结构体，递归转换
		if f.Kind() == reflect.Struct {
			embMap, err := StructToDbMap(f.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range embMap {
				out[k] = v
			}
			continue
		}
		// 忽略未导出字段
		if t.Field(i).PkgPath != "" {
			continue
		}
		// 忽略gorm:"-"字段
		if t.Field(i).Tag.Get("gorm") == "-" {
			continue
		}
		// 忽略0值字段
		if isZero(f.Interface()) {
			continue
		}
		out[stringUtils.ToUnderScoreCase(t.Field(i).Name)] = f.Interface()
	}

	return out, nil
}

// 判断值是否为零值
func isZero(i interface{}) bool {
	vi := reflect.ValueOf(i)
	switch vi.Kind() {
	case reflect.String:
		return vi.Len() == 0
	case reflect.Bool:
		return !vi.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vi.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return vi.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return vi.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return vi.IsNil()
	}
	return vi.Interface() == reflect.Zero(vi.Type()).Interface()
}
