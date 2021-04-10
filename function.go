package urlquery

import (
	"fmt"
	"math"
	"reflect"
)

var (
	//
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

// check if reflect.Kind of map's key is valid
func isAccessMapKeyType(kind reflect.Kind) bool {
	_, ok := accessMapTypes[kind]
	return ok
}

// check if reflect.Kind of map's value is valid
func isAccessMapValueType(kind reflect.Kind) bool {
	return isAccessMapKeyType(kind)
}

// students[0][id] -> students, [0][id]
// [students][0][id] -> students, [0][id]
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

// if key like `students[0]` , repack it to `students[]`
func repackArrayQueryKey(key string) string {
	l := len(key)
	if l > 0 && key[l-1] == ']' {
		for l--; l >= 0; l-- {
			if key[l] == '[' {
				return key[:l+1] + "]"
			}
		}
	}
	return key
}

// generate next parent node key
func genNextParentNode(parentNode, key string) string {
	if len(parentNode) > 0 {
		return parentNode + "[" + key + "]"
	}
	return key
}

// check if value is zero-value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !v.Index(i).IsZero() {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				return false
			}
		}
		return true
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(fmt.Sprintf("panic at isZeroValue %v", v.Kind()))
	}
}
