package urlquery

import (
	"reflect"
	"strconv"
)

//translator from string to go basic structure

type valueDecode func(string) (reflect.Value, error)

func boolDecode(value string) (reflect.Value, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		err = ErrTranslated{err: err}
		return reflect.Value{}, err
	}
	return reflect.ValueOf(b), nil
}

func baseIntDecode(value string, bitSize int) (rv reflect.Value, err error) {
	v, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch bitSize {
	case 64:
		rv = reflect.ValueOf(v)
	case 32:
		rv = reflect.ValueOf(int32(v))
	case 16:
		rv = reflect.ValueOf(int16(v))
	case 8:
		rv = reflect.ValueOf(int8(v))
	case 0:
		rv = reflect.ValueOf(int(v))
	default:
		err = ErrUnsupportedBitSize{bitSize: bitSize}
	}
	return
}

func intDecode(bitSize int) valueDecode {
	return func(value string) (reflect.Value, error) {
		return baseIntDecode(value, bitSize)
	}
}

func baseUintDecode(value string, bitSize int) (rv reflect.Value, err error) {
	v, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch bitSize {
	case 64:
		rv = reflect.ValueOf(v)
	case 32:
		rv = reflect.ValueOf(uint32(v))
	case 16:
		rv = reflect.ValueOf(uint16(v))
	case 8:
		rv = reflect.ValueOf(uint8(v))
	case 0:
		rv = reflect.ValueOf(uint(v))
	default:
		err = ErrUnsupportedBitSize{bitSize: bitSize}
	}
	return
}

func uintDecode(bitSize int) valueDecode {
	return func(value string) (reflect.Value, error) {
		return baseUintDecode(value, bitSize)
	}
}

func uintPrtDecode(value string) (rv reflect.Value, err error) {
	v, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	return reflect.ValueOf(uintptr(v)), nil
}

func baseFloatDecode(value string, bitSize int) (rv reflect.Value, err error) {
	v, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch bitSize {
	case 64:
		rv = reflect.ValueOf(v)
	case 32:
		rv = reflect.ValueOf(float32(v))
	default:
		err = ErrUnsupportedBitSize{bitSize: bitSize}
	}
	return
}

func floatDecode(bitSize int) valueDecode {
	return func(value string) (reflect.Value, error) {
		return baseFloatDecode(value, bitSize)
	}
}

func stringDecode(value string) (reflect.Value, error) {
	return reflect.ValueOf(value), nil
}

func getDecodeFunc(kind reflect.Kind) valueDecode {
	switch kind {
	case reflect.Bool:
		return boolDecode
	case reflect.Int:
		return intDecode(0)
	case reflect.Int8:
		return intDecode(8)
	case reflect.Int16:
		return intDecode(16)
	case reflect.Int32:
		return intDecode(32)
	case reflect.Int64:
		return intDecode(64)
	case reflect.Uint:
		return uintDecode(0)
	case reflect.Uint8:
		return uintDecode(8)
	case reflect.Uint16:
		return uintDecode(16)
	case reflect.Uint32:
		return uintDecode(32)
	case reflect.Uint64:
		return uintDecode(64)
	case reflect.Uintptr:
		return uintPrtDecode
	case reflect.Float32:
		return floatDecode(32)
	case reflect.Float64:
		return floatDecode(64)
	case reflect.String:
		return stringDecode
	default:
		return nil
	}
}
