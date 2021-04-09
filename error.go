package urlquery

import (
	"reflect"
	"strconv"
)

// An ErrUnhandledType is a customized error
type ErrUnhandledType struct {
	typ reflect.Type
}

func (e ErrUnhandledType) Error() string {
	return "failed to unhandled type(" + e.typ.String() + ")"
}

// An ErrInvalidUnmarshalError is a customized error
type ErrInvalidUnmarshalError struct{}

func (e ErrInvalidUnmarshalError) Error() string {
	return "failed to unmarshal(non-pointer)"
}

// An ErrUnsupportedBitSize is a customized error
type ErrUnsupportedBitSize struct {
	bitSize int
}

func (e ErrUnsupportedBitSize) Error() string {
	return "failed to handle unsupported bitSize(" + strconv.Itoa(e.bitSize) + ")"
}

// An ErrTranslated is a customized error type
type ErrTranslated struct {
	err error
}

func (e ErrTranslated) Error() string {
	return "failed to translate:" + e.err.Error()
}

// An ErrInvalidMapKeyType is a customized error
type ErrInvalidMapKeyType struct {
	typ reflect.Type
}

func (e ErrInvalidMapKeyType) Error() string {
	return "failed to handle map key type(" + e.typ.String() + ")"
}

// An ErrInvalidMapValueType is a customized error
type ErrInvalidMapValueType struct {
	typ reflect.Type
}

func (e ErrInvalidMapValueType) Error() string {
	return "failed to handle map value type(" + e.typ.String() + ")"
}
