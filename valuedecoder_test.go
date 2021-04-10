package urlquery

import (
	"reflect"
	"testing"
)

func Test_boolDecode_Error(t *testing.T) {
	_, err := boolDecode("d-22")
	if _, ok := err.(ErrTranslated); !ok {
		t.Error("unexpected error type")
	}
}

func Test_baseIntDecode_Error(t *testing.T) {
	_, err := baseIntDecode("d2", 8)
	if _, ok := err.(ErrTranslated); !ok {
		t.Error("unexpected error type")
	}

	_, err1 := baseIntDecode("2", 23)
	if _, ok := err1.(ErrUnsupportedBitSize); !ok {
		t.Error("unexpected error type")
	}
}

func Test_baseUintDecode_Error(t *testing.T) {
	_, err := baseUintDecode("d2", 8)
	if _, ok := err.(ErrTranslated); !ok {
		t.Error("unexpected error type")
	}

	_, err1 := baseUintDecode("2", 23)
	if _, ok := err1.(ErrUnsupportedBitSize); !ok {
		t.Error("unexpected error type")
	}
}

func Test_baseFloatDecode_Error(t *testing.T) {
	_, err := baseFloatDecode("d2.2", 32)
	if _, ok := err.(ErrTranslated); !ok {
		t.Error("unexpected error type")
	}

	_, err1 := baseFloatDecode("2.2", 63)
	if _, ok := err1.(ErrUnsupportedBitSize); !ok {
		t.Error("unexpected error type")
	}
}

func Test_baseUintPrtDecode_Error(t *testing.T) {
	_, err := uintPrtDecode("d2")
	if _, ok := err.(ErrTranslated); !ok {
		t.Error("unexpected error type")
	}
}

func Test_getDecodeFunc_Nil(t *testing.T) {
	v := getDecodeFunc(reflect.Chan)
	if v != nil {
		t.Error("getDecodeFunc of reflect.Chan should return nil")
	}
}
