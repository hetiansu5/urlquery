package urlquery

import (
	"reflect"
	"strconv"
)

//translator from go basic structure to string

type valueEncode func(value reflect.Value) string

func boolEncode(value reflect.Value) string {
	if value.Bool() {
		return "1"
	} else {
		return "0"
	}
}

func intEncode(value reflect.Value) string {
	return strconv.FormatInt(value.Int(), 10)
}

func uintEncode(value reflect.Value) string {
	return strconv.FormatUint(value.Uint(), 10)
}

func floatEncode(value reflect.Value) string {
	return strconv.FormatFloat(value.Float(), 'f', -1, 64)
}

func stringEncode(value reflect.Value) string {
	return value.String()
}

func getEncodeFunc(kind reflect.Kind) valueEncode {
	switch kind {
	case reflect.Bool:
		return boolEncode
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncode
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncode
	case reflect.Float32, reflect.Float64:
		return floatEncode
	case reflect.String:
		return stringEncode
	default:
		return nil
	}
}
