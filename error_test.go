package urlquery

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrUnhandledType_Error(t *testing.T) {
	err := ErrUnhandledType{typ: reflect.TypeOf("s")}
	if err.Error() != "failed to unhandled type(string)" {
		t.Error(err.Error())
	}
}

func TestErrTranslated_Error(t *testing.T) {
	err1 := errors.New("new")
	err := ErrTranslated{err: err1}
	if err.Error() != "failed to translate:new" {
		t.Error(err.Error())
	}
}

func TestErrUnsupportedBitSize_Error(t *testing.T) {
	err := ErrUnsupportedBitSize{bitSize: 32}
	if err.Error() != "failed to handle unsupported bitSize(32)" {
		t.Error(err.Error())
	}
}

func TestErrInvalidMapKeyType_Error(t *testing.T) {
	var f float64
	f = 3.14
	err := ErrInvalidMapKeyType{typ: reflect.TypeOf(f)}
	if err.Error() != "failed to handle map key type(float64)" {
		t.Error(err.Error())
	}
}

func TestErrInvalidUnmarshalError_Error(t *testing.T) {
	err := ErrInvalidUnmarshalError{}
	if err.Error() != "failed to unmarshal(non-pointer)" {
		t.Error(err.Error())
	}
}

func TestErrInvalidMapValueType_Error(t *testing.T) {
	i := uint(2)
	err := ErrInvalidMapValueType{typ: reflect.TypeOf(i)}
	if err.Error() != "failed to handle map value type(uint)" {
		t.Error(err.Error())
	}
}
