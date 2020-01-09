package urlquery

import (
	"reflect"
	"strconv"
)

//translator from string to go basic structure

var (
	decoderMap = map[reflect.Kind]valueDecoder{
		reflect.Bool:    boolDecoder{},
		reflect.Int:     intDecoder{},
		reflect.Int8:    intDecoder{8},
		reflect.Int16:   intDecoder{16},
		reflect.Int32:   intDecoder{32},
		reflect.Int64:   intDecoder{64},
		reflect.Uint:    uintDecoder{},
		reflect.Uint8:   uintDecoder{8},
		reflect.Uint16:  uintDecoder{16},
		reflect.Uint32:  uintDecoder{32},
		reflect.Uint64:  uintDecoder{64},
		reflect.Uintptr: uintptrDecoder{},
		reflect.Float32: floatDecoder{32},
		reflect.Float64: floatDecoder{64},
		reflect.String:  stringDecoder{},
	}
)

type valueDecoder interface {
	Decode(value string) (reflect.Value, error)
}

type boolDecoder struct{}

func (e boolDecoder) Decode(value string) (rv reflect.Value, err error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	rv = reflect.ValueOf(b)
	return
}

type intDecoder struct {
	bitSize int
}

func (e intDecoder) Decode(value string) (rv reflect.Value, err error) {
	v, err := strconv.ParseInt(value, 10, e.bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch e.bitSize {
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
		err = ErrUnsupportedBitSize{bitSize: e.bitSize}
	}
	return
}

type uintDecoder struct {
	bitSize int
}

func (e uintDecoder) Decode(value string) (rv reflect.Value, err error) {
	v, err := strconv.ParseUint(value, 10, e.bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch e.bitSize {
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
		err = ErrUnsupportedBitSize{bitSize: e.bitSize}
	}
	return
}

type uintptrDecoder struct {
	bitSize int
}

func (e uintptrDecoder) Decode(value string) (rv reflect.Value, err error) {
	v, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}

	return reflect.ValueOf(uintptr(v)), nil
}

type floatDecoder struct {
	bitSize int
}

func (e floatDecoder) Decode(value string) (rv reflect.Value, err error) {
	v, err := strconv.ParseFloat(value, e.bitSize)
	if err != nil {
		err = ErrTranslated{err: err}
		return
	}
	switch e.bitSize {
	case 64:
		rv = reflect.ValueOf(v)
	case 32:
		rv = reflect.ValueOf(float32(v))
	default:
		err = ErrUnsupportedBitSize{bitSize: e.bitSize}
	}
	return
}

type stringDecoder struct{}

func (e stringDecoder) Decode(value string) (rv reflect.Value, err error) {
	return reflect.ValueOf(value), nil
}

func getDecoder(kind reflect.Kind) valueDecoder {
	if decoder, ok := decoderMap[kind]; ok {
		return decoder
	} else {
		return nil
	}
}
