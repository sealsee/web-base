package jsonUtils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sealsee/web-base/public/basemodel"
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
		tf := t.Field(i)
		if !f.CanInterface() {
			continue
		}
		// 如果是BaseTime类型，需要特殊处理。
		if marshaler, ok := f.Interface().(basemodel.BaseTime); ok {
			if !marshaler.IsZero() {
				text, err := marshaler.MarshalText()
				if err != nil {
					return nil, err
				}
				if string(text) != "" {
					out[stringUtils.ToUnderScoreCase(tf.Name)] = string(text)
				}
			}
		} else if f.Kind() == reflect.Struct { // 如果字段是一个结构体，递归转换
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
		if tf.PkgPath != "" {
			continue
		}
		// 忽略gorm:"-"字段
		gormTag := tf.Tag.Get("gorm")
		if gormTag == "-" {
			continue
		}
		// 忽略0值字段
		if isZero(f.Interface()) {
			continue
		}
		// 对应的列名称（tag中定义）
		realCol := ""
		gormTags := strings.Split(gormTag, ";")
		for _, v := range gormTags {
			if strings.Contains(v, "column:") {
				realCol = strings.Replace(v, "column:", "", -1)
			}
		}
		if realCol == "" {
			realCol = stringUtils.ToUnderScoreCase(tf.Name)
		}
		out[realCol] = f.Interface()
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
