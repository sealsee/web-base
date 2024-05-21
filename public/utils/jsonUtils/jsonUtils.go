package jsonUtils

import (
	"bytes"
	"encoding/json"

	"github.com/sealsee/web-base/public/utils/stringUtils"
)

// 结构体转换map
func StructToMap(structContent interface{}) map[string]interface{} {
	var structMap map[string]interface{}
	if marshalContent, err := json.Marshal(structContent); err != nil {
		panic(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&structMap); err != nil {
			panic(err)
		} else {
			for k, v := range structMap {
				structMap[k] = v
			}
		}
	}
	return structMap
}

// 结构体转换map，key使用下划线命名格式
func StructToMapWithKeyUnderLine(structContent interface{}) map[string]interface{} {
	var structMap map[string]interface{}
	if marshalContent, err := json.Marshal(structContent); err != nil {
		panic(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&structMap); err != nil {
			panic(err)
		} else {
			for k, v := range structMap {
				delete(structMap, k)
				underLineKey := stringUtils.ToUnderScoreCase(k)
				if v != "0001-01-01 00:00:00" { // 去除时间为空的字段
					structMap[underLineKey] = v
				}
			}
		}
	}
	return structMap
}
