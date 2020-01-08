package query

import (
	"reflect"
	"strconv"
	"fmt"
)

//translator from go basic structure to string

var (
	boolE      = boolEncoder{}
	intE       = intEncoder{}
	unitE      = uintEncoder{}
	floatE     = floatEncoder{}
	stringE    = stringEncoder{}
	encoderMap = map[reflect.Kind]Encoder{
		reflect.Bool:    boolE,
		reflect.Int:     intE,
		reflect.Int8:    intE,
		reflect.Int16:   intE,
		reflect.Int32:   intE,
		reflect.Int64:   intE,
		reflect.Uint:    unitE,
		reflect.Uint8:   unitE,
		reflect.Uint16:  unitE,
		reflect.Uint32:  unitE,
		reflect.Uint64:  unitE,
		reflect.Uintptr: unitE,
		reflect.Float32: floatE,
		reflect.Float64: floatE,
		reflect.String:  stringE,
	}
)

type Encoder interface {
	Encode(value reflect.Value) string
}

type boolEncoder struct{}

func (e boolEncoder) Encode(value reflect.Value) string {
	if value.Bool() {
		return "1"
	} else {
		return "0"
	}
}

type intEncoder struct{}

func (e intEncoder) Encode(value reflect.Value) string {
	return strconv.FormatInt(value.Int(), 10)
}

type uintEncoder struct{}

func (e uintEncoder) Encode(value reflect.Value) string {
	return strconv.FormatUint(value.Uint(), 10)
}

type floatEncoder struct{}

func (e floatEncoder) Encode(value reflect.Value) string {
	return strconv.FormatFloat(value.Float(), 'f', -1, 64)
}

type stringEncoder struct{}

func (e stringEncoder) Encode(value reflect.Value) string {
	return value.String()
}

type commonEncoder struct{}

func (e commonEncoder) Encode(value reflect.Value) string {
	return fmt.Sprint(value.Interface())
}

func getEncoder(kind reflect.Kind) Encoder {
	if encoder, ok := encoderMap[kind]; ok {
		return encoder
	} else {
		return nil
	}
}
