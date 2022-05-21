package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

func IsContain(s interface{}, k interface{}) (int, bool) {
	switch reflect.TypeOf(s).Kind() {
	case reflect.Slice:
		sliceValue := reflect.ValueOf(s)
		sliceLen := sliceValue.Len()
		for i := 0; i < sliceLen; i++ {
			if sliceValue.Index(i).Interface() == k {
				return i, true
			}
		}
		return -1, false
	case reflect.Map:
		mapValue := reflect.ValueOf(s)
		for _, key := range mapValue.MapKeys() {
			if key.Interface() == k {
				return 0, true
			}
		}
		return 0, false
	default:
		return -1, false
	}
}

func RemoveSliceDuplicate(s interface{}) {
	sliceValue := reflect.ValueOf(s).Elem()
	if sliceValue.Len() == 0 {
		return
	}

	seen := make(map[interface{}]bool)
	i := 0
	for j := 1; j < sliceValue.Len(); j++ {
		if !seen[sliceValue.Index(j).Interface()] {
			seen[sliceValue.Index(j).Interface()] = true
			sliceValue.Index(i + 1).Set(reflect.ValueOf(sliceValue.Index(j).Interface()))
			i++
		}
	}
	sliceValue.SetLen(i)
}

func GetClientIP(ctx *gin.Context) string {
	if ip := ctx.Request.Header.Get("X-Envoy-External-Address"); ip != "" {
		return ip
	}

	return ctx.ClientIP()
}

func DerefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func DerefInt(s *int) int {
	if s == nil {
		return 0
	}

	return *s
}

func GetInputFromSource(v interface{}) map[string]interface{} {
	input, ok := v.(map[string]interface{})
	if ok {
		return input
	}

	return map[string]interface{}{}
}

func GetSubMap(v interface{}, keys ...string) map[string]interface{} {
	curData := v
	var curMap map[string]interface{}
	var ok bool
	for _, key := range keys {
		curMap, ok = curData.(map[string]interface{})
		if !ok {
			return map[string]interface{}{}
		}

		curData, ok = curMap[key]
		if !ok {
			return map[string]interface{}{}
		}
	}

	result, ok := curData.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}

	return result
}

func GetSubSliceMap(v interface{}, keys ...string) []map[string]interface{} {
	var innerMap map[string]interface{}
	var ok bool
	if len(keys) == 1 {
		innerMap, ok = v.(map[string]interface{})
		if !ok {
			return []map[string]interface{}{}
		}
	} else {
		innerMap = GetSubMap(v, keys[:len(keys)-1]...)
	}

	curData, ok := innerMap[keys[len(keys)-1]].([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}

	var result []map[string]interface{}
	for _, datum := range curData {
		if m, ok := datum.(map[string]interface{}); ok {
			result = append(result, m)
		} else {
			return []map[string]interface{}{}
		}
	}

	return result
}

func GetSubInteger(v interface{}, keys ...string) *int {
	curData := v
	for _, key := range keys {
		curMap, ok := curData.(map[string]interface{})
		if !ok {
			return nil
		}

		curData, ok = curMap[key]
		if !ok {
			return nil
		}
	}

	switch curData.(type) {
	case float64:
		result := int(curData.(float64))
		return &result
	case int:
		result := curData.(int)
		return &result
	default:
		return nil
	}
}

func CompareIntSlice(input, output []int, msgFormat string) (map[string]interface{}, bool) {
	outputMap := make(map[int]bool)
	for _, v := range output {
		outputMap[v] = true
	}

	msg := make(map[string]interface{})
	for i, v := range input {
		if !outputMap[v] {
			msg[fmt.Sprintf(msgFormat, i)] = ErrorInputFail
		}
	}

	return msg, len(msg) == 0
}

func GetOnlyScalar(v interface{}, keys ...string) map[string]interface{} {
	result := GetSubMap(v, keys...)

	scalarMap := make(map[string]interface{})
	for k, v := range result {
		switch v.(type) {
		case int, int8, int16, int32, int64, *int, *int8, *int16, *int32, *int64,
			float32, float64, *float32, *float64, string, *string:
			scalarMap[k] = v
		}
	}

	return scalarMap
}

func MergeMap(v1 map[string]interface{}, v2 map[string]interface{}) map[string]interface{} {
	if v2 == nil && v1 != nil {
		return v1
	}
	if v1 == nil && v2 != nil {
		return v2
	}
	if v1 == nil && v2 == nil {
		return map[string]interface{}{}
	}

	for key, value := range v2 {
		v1[key] = value
	}
	return v1
}

func GetIntPointer(v int) *int {
	return &v
}

func GetBoolPointer(v bool) *bool {
	return &v
}

func GetStringPointer(v string) *string {
	return &v
}
