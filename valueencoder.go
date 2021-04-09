package urlquery

import (
	"reflect"
	"strconv"
)

// A valueEncode is a converter from go basic structure to string
type valueEncode func(value reflect.Value) string

// converter from bool to string
func boolEncode(value reflect.Value) string {
	if value.Bool() {
		return "1"
	}
	return "0"
}

// converter from int(8-64) to string
func intEncode(value reflect.Value) string {
	return strconv.FormatInt(value.Int(), 10)
}

// converter from uint(8-64) to string
func uintEncode(value reflect.Value) string {
	return strconv.FormatUint(value.Uint(), 10)
}

// converter from float,double to string
func floatEncode(value reflect.Value) string {
	return strconv.FormatFloat(value.Float(), 'f', -1, 64)
}

// converter from string to string
func stringEncode(value reflect.Value) string {
	return value.String()
}

// get encode func for specified reflect kind
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
