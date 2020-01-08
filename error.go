package query

import (
	"reflect"
	"strconv"
)

type ErrUnhandledType struct {
	key string
	t   reflect.Type
}

func (e ErrUnhandledType) Error() string {
	return "Type(" + e.t.String() + ") of key(" + e.key + ") is unhandled";
}

type ErrInvalidUnmarshalError struct {
	Type reflect.Type
}

func (e ErrInvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "query: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "query: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "query: Unmarshal(nil " + e.Type.String() + ")"
}

type ErrUnsupportedBitSize struct {
	bitSize int
}

func (e ErrUnsupportedBitSize) Error() string {
	return "bitSize(" + strconv.Itoa(e.bitSize) + ") is unsupported"
}
