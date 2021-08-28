package urlquery

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type builderChild struct {
	Description string `query:"desc"`
	Long        uint16 `query:",vip"`
	Height      int    `query:"-"`
}

type builderInfo struct {
	Id       int
	Name     string         `query:"name"`
	Child    builderChild   `query:"child"`
	ChildPtr *builderChild  `query:"childPtr"`
	Children []builderChild `query:"children"`
	Params   map[string]rune
	status   bool
	UintPtr  uintptr
}

func TestEncoder_Marshal(t *testing.T) {
	data := getMockData()

	SetGlobalQueryEncoder(defaultQueryEncoder)
	_, err := Marshal(data)
	SetGlobalQueryEncoder(nil)

	if err != nil {
		t.Error(err)
	}
}

func TestEncoder_Marshal_Struct(t *testing.T) {
	data := getMockData2()

	encoder := NewEncoder(WithQueryEncoder(DefaultQueryEncoder{}), WithNeedEmptyValue(false))
	bytes, err := encoder.Marshal(data)

	if err != nil {
		t.Error(err)
		return
	}

	v := builderInfo{}
	err = Unmarshal(bytes, &v)
	if err != nil {
		t.Error(err)
		return
	}

	if v.Name != "child" || v.status != false || v.Child.Height != 0 || len(v.Children) != 5 {
		fmt.Println(v.Name, v.status, v.Child.Height, v.Children)
		t.Error("Marshal Unmarshal is not equal")
		return
	}
}

func TestEncoder_Marshal_NilPtr_Struct(t *testing.T) {
	data := getMockData3()

	bytes, err := Marshal(data)

	if err != nil {
		t.Error(err)
		return
	}

	v := builderInfo{}
	err = Unmarshal(bytes, &v)
	if err != nil {
		t.Error(err)
		return
	}

	if v.Name != "child3" || v.status != false || v.Child.Height != 0 || len(v.Children) != 5 {
		fmt.Println(v.Name, v.status, v.Child.Height, v.Children)
		t.Error("Marshal is not equal to Unmarshal")
		return
	}

	if v.ChildPtr != nil {
		t.Error("The child pointer should be nil not ", v.ChildPtr)
		return
	}
	if v.Params != nil {
		t.Error("The params map should be nil")
		return
	}
}

func TestEncoder_Marshal_Slice(t *testing.T) {
	data := []string{"a", "b"}
	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if string(bytes) != "0=a&1=b" {
		t.Error("failed to Marshal slice")
	}
}

func TestEncoder_Marshal_Array(t *testing.T) {
	data := [3]int32{10, 200, 50}
	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if string(bytes) != "0=10&1=200&2=50" {
		t.Error("failed to Marshal slice")
	}
}

type TestPoint struct {
	X, Y int
}

type TestCircle struct {
	TestPoint
	R int
}

func TestEncoder_Marshal_AnonymousFields(t *testing.T) {
	data := &TestCircle{R: 1}
	data.TestPoint.X = 12
	data.TestPoint.Y = 13

	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "X=12&Y=13&R=1" {
		t.Error("failed to Marshal anonymous fields")
	}
}

func TestEncoder_Marshal_DuplicateCall(t *testing.T) {
	d1 := builderChild{
		Description: "a",
		Long:        10,
	}

	encoder := NewEncoder()
	encoder.RegisterEncodeFunc(reflect.Int64, func(value reflect.Value) string {
		return strconv.FormatInt(value.Int(), 10)
	})
	_, _ = encoder.Marshal(d1)

	d2 := builderChild{
		Description: "bb",
		Long:        200,
	}

	bytes2, err := encoder.Marshal(d2)
	if err != nil {
		t.Error(err)
	}

	if string(bytes2) != "desc=bb&Long=200" {
		t.Error("failed to Marshal duplicate call")
	}
}

func TestEncoder_Marshal_ErrInvalidMapKeyType(t *testing.T) {
	encoder := NewEncoder()
	data := map[complex64]int{
		complex(1, 2): 23,
	}
	_, err := encoder.Marshal(data)
	if _, ok := err.(ErrInvalidMapKeyType); !ok {
		t.Error("unmatched error")
	}
}

func TestEncoder_buildQuery_ReturnError(t *testing.T) {
	encoder := NewEncoder()
	encoder.err = errors.New("return")
	encoder.buildQuery(reflect.ValueOf("s"), "", reflect.Int)
	if encoder.err == nil || encoder.err.Error() != "return" {
		t.Error("unmatched error")
	}
}

func TestEncoder_encodeError(t *testing.T) {
	encoder := NewEncoder()
	data := map[string]complex64{
		"d": complex(1, 2),
	}
	_, err := encoder.Marshal(data)
	if _, ok := err.(ErrUnhandledType); !ok {
		t.Error("unmatched error")
	}
}

func TestEncoder_encodeError2(t *testing.T) {
	encoder := NewEncoder()
	data := map[string]int{
		"d": 1,
	}
	encoder.RegisterEncodeFunc(reflect.String, nil)
	_, err := encoder.Marshal(data)
	if _, ok := err.(ErrUnhandledType); !ok {
		t.Error("unmatched error")
	}
}

func TestEncoder_RegisterEncodeFunc(t *testing.T) {
	encoder := NewEncoder()
	encoder.RegisterEncodeFunc(reflect.Int, nil)
	f := encoder.getEncodeFunc(reflect.Int)
	if f != nil {
		t.Error("failed to RegisterEncodeFunc")
	}
}

//BenchmarkMarshal-4     	  295726	     11902 ns/op
func BenchmarkMarshal(b *testing.B) {
	data := getMockData2()

	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func getMockData() map[string]interface{} {
	var (
		f32 = float32(1.2)
		i8  = int8(3)
		i64 = int64(9999999 * 9999999)
		u64 = uint16(567)
	)
	var f64 float64 = 13.4343453535343242342
	return map[string]interface{}{
		"id":     1,
		"fit":    true,
		"vip":    false,
		"desc":   "测试",
		"f32":    f32,
		"f64":    f64,
		"int8":   i8,
		"int64":  i64,
		"uint16": u64,
		"map": map[interface{}]interface{}{
			"caption": "test",
			5:         []int{11, 22},
			"child":   getMockData2(),
		},
		"struct": builderInfo{
			Id:   222,
			Name: "test",
		},
	}
}

func getMockData2() builderInfo {
	return builderInfo{
		Name: "child",
		Children: []builderChild{
			{Description: "d1", Height: 180},
			{Description: "d2", Long: 140},
			{Description: "d4"},
			{Description: "d5", Long: 1, Height: 20},
			{Description: "d6"},
		},
		Child:    builderChild{Description: "c1", Height: 20},
		ChildPtr: &builderChild{Description: "cptr", Long: 14, Height: 220},
		Params: map[string]rune{
			"abc":      111,
			"bbb":      222,
			"whoIsWho": 344340,
		},
		status:  true,
		UintPtr: uintptr(222),
	}
}

func getMockData3() builderInfo {
	return builderInfo{
		Name: "child3",
		Children: []builderChild{
			{Description: "d31", Height: 180},
			{Description: "d32", Long: 140},
			{Description: "d34"},
			{Description: "d35", Long: 1, Height: 20},
			{Description: "d36"},
		},
		Child:    builderChild{},
		ChildPtr: nil,
		Params:   nil,
		status:   true,
		UintPtr:  uintptr(2222),
	}
}

func TestEmptyStruct(t *testing.T) {
	data := &TestCircle{}
	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "" {
		t.Error("failed to Marshal anonymous fields")
	}
}
