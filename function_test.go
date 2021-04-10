package urlquery

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Id int
}

func Test_unpackQueryKey_NotEnd(t *testing.T) {
	key := "hts"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "" {
		t.Error("unpack error")
	}
}

func Test_unpackQueryKey_RightSquareBracketEnd(t *testing.T) {
	key := "hts[0]"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "[0]" {
		t.Error("unpack error")
	}
}

func Test_unpackQueryKey_LeftSquareBracketEnd(t *testing.T) {
	key := "[hts][0]"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "[0]" {
		t.Error("unpack error")
	}
}

func Test_repackArrayQueryKey(t *testing.T) {
	key := "[hts][0]"
	target := repackArrayQueryKey(key)
	if target != "[hts][]" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_repackArrayQueryKey1(t *testing.T) {
	key := "hts]"
	target := repackArrayQueryKey(key)
	if target != "hts]" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_repackArrayQueryKey2(t *testing.T) {
	key := "[hts"
	target := repackArrayQueryKey(key)
	if target != "[hts" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_genNextParentNode(t *testing.T) {
	if genNextParentNode("", "test") != "test" {
		t.Error("failed to execute genNextParentNode")
	}

	if genNextParentNode("p", "test") != "p[test]" {
		t.Error("failed to execute genNextParentNode")
	}
}

func Test_isZeroValue_Bool(t *testing.T) {
	var a bool
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue is error for bool")
	}
}

func Test_isZeroValue_Complex(t *testing.T) {
	var a complex64
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue is error for complex64")
	}
}

func Test_isZeroValue_Float(t *testing.T) {
	var a float32
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue failed for float32")
	}
}

func Test_isZeroValue_Array(t *testing.T) {
	var a [3]int
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue failed for array")
	}

	a[1] = 0
	res = isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue failed for array 1")
	}

	a[2] = 2
	res = isZeroValue(reflect.ValueOf(a))
	if res != false {
		t.Error("isZeroValue failed for array 2")
	}
}

func Test_isZeroValue_Pointer(t *testing.T) {
	var a *int
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue failed for pointer")
	}
}

func Test_isZeroValue_Struct(t *testing.T) {
	a := testStruct{0}
	res := isZeroValue(reflect.ValueOf(a))
	if res != true {
		t.Error("isZeroValue failed for struct")
	}

	a.Id = 1
	res = isZeroValue(reflect.ValueOf(a))
	if res != false {
		t.Error("isZeroValue failed for struct 1")
	}
}

func Test_isZeroValue_UndefinedType(t *testing.T) {
	var a error
	defer func() {
		if err := recover(); err == nil {
			t.Error("dont not panic at undefined type")
		}
	}()
	isZeroValue(reflect.ValueOf(a))
}
