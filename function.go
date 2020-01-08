package query

import (
	"net/url"
	"reflect"
)

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
		reflect.Float32: true,
		reflect.Float64: true,
		reflect.String:  true,
	}
)

func genNextParentNode(parentNode, key string) string {
	if len(parentNode) > 0 {
		return parentNode + url.QueryEscape("["+key+"]")
	} else {
		return url.QueryEscape(key)
	}
}

func unpackKey(key string) (pre, suf string) {
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

func isAccessMapKeyType(kind reflect.Kind) bool {
	_, ok := accessMapTypes[kind]
	return ok
}

func isAccessMapValueType(kind reflect.Kind) bool {
	return isAccessMapKeyType(kind)
}
