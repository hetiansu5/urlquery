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

func unpackQueryKey(key string) (pre, suf string) {
	if len(key) > 3 && key[:3] == "%5B" {
		key = key[3:]
	}
	i := 0
	for i < len(key) {
		if i+3 <= len(key) {
			if key[i:i+3] == "%5D" {
				pre = key[0:i]
				suf = key[i+3:]
				return
			} else if key[i:i+3] == "%5B" {
				pre = key[0:i]
				suf = key[i:]
				return
			}
		}
		i++
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

//Is space character
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
