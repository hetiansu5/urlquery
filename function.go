package urlquery

import (
	"reflect"
)

//common function library

var (
	accessMapTypes = map[reflect.Kind]bool{
		reflect.Bool:    true,
		reflect.Int:     true,
		reflect.Int8:    true,
		reflect.Int16:   true,
		reflect.Int32:   true,
		reflect.Int64:   true,
		reflect.Uint:    true,
		reflect.Uint8:   true,
		reflect.Uint16:  true,
		reflect.Uint32:  true,
		reflect.Uint64:  true,
		reflect.Uintptr: true,
		reflect.Float32: true,
		reflect.Float64: true,
		reflect.String:  true,
	}
)

func isAccessMapKeyType(kind reflect.Kind) bool {
	_, ok := accessMapTypes[kind]
	return ok
}

func isAccessMapValueType(kind reflect.Kind) bool {
	return isAccessMapKeyType(kind)
}

//students[0][id] -> students, [0][id]
//[students][0][id] -> students, [0][id]
func unpackQueryKey(key string) (pre, suf string) {
	if len(key) > 0 && key[0] == '[' {
		key = key[1:]
	}
	for i, v := range key {
		if v == ']' {
			pre = key[0:i]
			suf = key[i+1:]
			return
		} else if v == '[' {
			pre = key[0:i]
			suf = key[i:]
			return
		}
	}
	return key, ""
}

//Is Zero-value
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
