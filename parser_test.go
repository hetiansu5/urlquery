package urlquery

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type testParseChild struct {
	Description string `query:"desc"`
	Long        uint16 `query:",vip"`
	Height      int    `query:"-"`
}

type testParseInfo struct {
	Id       int
	Name     string           `query:"name"`
	Child    testParseChild   `query:"child"`
	ChildPtr *testParseChild  `query:"childPtr"`
	Children []testParseChild `query:"children"`
	Params   map[byte]int8
	status   bool
	UintPtr  uintptr
	Tags     []int16 `query:"tags"`
	Int64    int64
	Uint     uint
	Uint32   uint32
	Float32  float32
	Float64  float64
	Bool     bool
	Inter    interface{} `query:"inter"`
}

type errorQueryEncoder struct {
	times   int
	errorAt int
}

func (q *errorQueryEncoder) Escape(s string) string {
	return s
}
func (q *errorQueryEncoder) UnEscape(s string) (string, error) {
	q.times++
	if q.times >= q.errorAt {
		return "", errors.New("failed")
	}
	return s, nil
}

func TestParser_Unmarshal_DuplicateCall(t *testing.T) {
	parser := NewParser()

	d1 := "desc=bb&Long=200"
	v1 := &testParseChild{}
	_ = parser.Unmarshal([]byte(d1), v1)

	d2 := "desc=a&Long=100"
	v2 := &testParseChild{}
	err := parser.Unmarshal([]byte(d2), v2)
	if err != nil {
		t.Error(err)
	}
	if v2.Description != "a" || v2.Long != 100 {
		t.Error("failed to Unmarshal duplicate call")
	}
}

func TestParser_Unmarshal_NestedStructure(t *testing.T) {
	var data = "Id=1&name=test&child[desc]=c1&child[Long]=10&childPtr[Long]=2&childPtr[Description]=b" +
		"&children[0][desc]=d1&children[1][Long]=12&children[5][desc]=d5&children[5][Long]=50&desc=rtt" +
		"&Params[120]=1&Params[121]=2&status=1&UintPtr=300&tags[]=1&tags[]=2&Int64=64&Uint=22&Uint32=5&Float32=1.3" +
		"&Float64=5.64&Bool=0&inter=ss"
	data = encodeSquareBracket(data)
	v := &testParseInfo{}
	err := Unmarshal([]byte(data), v)

	if err != nil {
		t.Error(err)
	}

	if v.Id != 1 {
		t.Error("Id wrong")
	}

	if v.Name != "test" {
		t.Error("Name wrong")
	}

	if v.Child.Description != "c1" || v.Child.Long != 10 || v.Child.Height != 0 {
		t.Error("Child wrong")
	}

	if v.ChildPtr == nil || v.ChildPtr.Description != "" || v.ChildPtr.Long != 2 || v.ChildPtr.Height != 0 {
		t.Error("ChildPtr wrong")
	}

	if len(v.Children) != 6 {
		t.Error("Children's length is wrong")
	}

	if v.Children[0].Description != "d1" {
		t.Error("Children[0] wrong")
	}

	if v.Children[1].Description != "" || v.Children[1].Long != 12 {
		t.Error("Children[1] wrong")
	}

	if v.Children[2].Description != "" || v.Children[3].Description != "" || v.Children[4].Description != "" {
		t.Error("Children[2,3,4] wrong")
	}

	if v.Children[5].Description != "d5" || v.Children[5].Long != 50 || v.Children[5].Height != 0 {
		t.Error("Children[5] wrong")
	}

	if len(v.Params) != 2 || v.Params[120] != 1 || v.Params[121] != 2 {
		t.Error("Params wrong")
	}

	if v.status != false {
		t.Error("status wrong")
	}

	if v.UintPtr != uintptr(300) {
		t.Error("UintPtr wrong")
	}

	if len(v.Tags) != 2 {
		t.Error("Tags wrong")
	}
}

func TestParser_Unmarshal_Map(t *testing.T) {
	var m map[string]string
	data := "id=1&name=ab&arr[0]=6d"
	data = encodeSquareBracket(data)
	err := Unmarshal([]byte(data), &m)

	if err != nil {
		t.Error(err)
	}

	if len(m) != 2 {
		t.Error("length is wrong")
	}
	if v1, ok1 := m["id"]; v1 != "1" || !ok1 {
		t.Error("map[id] is wrong")
	}
	if v2, ok2 := m["name"]; v2 != "ab" || !ok2 {
		t.Error("map[iname] is wrong")
	}
	if _, ok3 := m["arr%5B0%5D"]; ok3 {
		t.Error("map[arr%5B0%5D] should not be exist")
	}
}

func TestParser_Unmarshal_Slice(t *testing.T) {
	var slice []int
	slice = make([]int, 0)
	data := "1=20&3=30"
	err := Unmarshal([]byte(data), &slice)

	if err != nil {
		t.Error(err)
	}

	if len(slice) != 4 {
		t.Error("failed to Unmarshal slice")
	}
}

func TestParser_Unmarshal_Array(t *testing.T) {
	var arr [5]int
	data := "1=20&3=30"
	err := Unmarshal([]byte(data), &arr)

	if err != nil {
		t.Error(err)
	}

	if arr[1] != 20 || arr[3] != 30 || arr[0] != 0 {
		t.Error("failed to Unmarshal array")
	}
}

func TestParser_Unmarshal_Array_Failed(t *testing.T) {
	var arr [5]int
	data := "1=20&3=s"
	err := Unmarshal([]byte(data), &arr)

	if err == nil {
		t.Error("dont return error")
	}
}

type testParserPoint struct {
	X, Y int
}

type testParserCircle struct {
	testParserPoint
	R int
}

func TestParser_Unmarshal_AnonymousFields(t *testing.T) {
	v := &testParserCircle{}
	data := "X=12&Y=13&R=1"
	err := Unmarshal([]byte(data), &v)

	if err != nil {
		t.Error(err)
	}

	if v.X != 12 || v.Y != 13 || v.R != 1 {
		t.Error("failed to Unmarshal anonymous fields")
	}
}

type testFormat struct {
	Id uint64
	B  rune `query:"b"`
}

func TestParser_Unmarshal_UnmatchedDataFormat(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &testFormat{}
	err := Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrTranslated); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

func TestParser_Unmarshal_UnhandledType(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &map[interface{}]string{}
	err := Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrInvalidMapKeyType); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

type TestUnhandled struct {
	Id     int
	Params map[string]testFormat
}

func TestParser_Unmarshal_UnhandledType2(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &TestUnhandled{}
	parser := NewParser(WithQueryEncoder(defaultQueryEncoder))
	err := parser.Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrInvalidMapKeyType); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

func TestParser_init(t *testing.T) {
	query := &errorQueryEncoder{errorAt: 1}
	parser := NewParser(WithQueryEncoder(query))
	parser.resetQueryEncoder()
	var data = "Id=1&b=a"
	err := parser.init([]byte(data))
	if err == nil || err.Error() != "failed" {
		t.Error("init error")
	}
}

func TestParser_Unmarshal_InitError(t *testing.T) {
	query := &errorQueryEncoder{errorAt: 2}
	parser := NewParser(WithQueryEncoder(query))
	v := &TestUnhandled{}
	var data = "Id=1&b=a"
	err := parser.Unmarshal([]byte(data), v)
	if err == nil || err.Error() != "failed" {
		t.Error("init error")
	}
}

func TestParser_Unmarshal_NonPointer(t *testing.T) {
	parser := NewParser()
	var data = "Id=1&b=a"
	v := TestUnhandled{}
	err := parser.Unmarshal([]byte(data), v)
	if _, ok := err.(ErrInvalidUnmarshalError); !ok {
		t.Error("unmatched error")
	}
}

func TestParser_Unmarshal_MapKey_DecodeError(t *testing.T) {
	parser := NewParser()
	parser.RegisterDecodeFunc(reflect.String, nil)
	var data = "Id=1&b=2"
	v := &map[string]int{}
	err := parser.Unmarshal([]byte(data), v)
	if _, ok := err.(ErrUnhandledType); !ok {
		t.Error("unmatched error")
	}
}

func TestParser_Unmarshal_MapValue_DecodeError(t *testing.T) {
	parser := NewParser()
	parser.RegisterDecodeFunc(reflect.Int, nil)
	var data = "Id=1&b=2"
	v := &map[string]int{}
	err := parser.Unmarshal([]byte(data), v)
	if _, ok := err.(ErrUnhandledType); !ok {
		t.Error("unmatched error")
	}
}

func TestParser_RegisterDecodeFunc(t *testing.T) {
	parser := NewParser()
	parser.RegisterDecodeFunc(reflect.String, func(s string) (reflect.Value, error) {
		return reflect.ValueOf("11"), nil
	})
	f := parser.getDecodeFunc(reflect.String)
	v, _ := f("bb")
	if v.String() != "11" {
		t.Error("failed to RegisterDecodeFunc")
	}
}

func TestParser_lookupForSlice(t *testing.T) {
	var data = "Tags[s]=1&Tags[]=2"
	data = encodeSquareBracket(data)
	v := &struct {
		Tags []int
	}{}
	err := Unmarshal([]byte(data), v)
	if _, ok := err.(*strconv.NumError); !ok {
		t.Error("dont failed for wrong slice data")
	}
}

func TestParser_SliceEmpty(t *testing.T) {
	var data = ""
	data = encodeSquareBracket(data)
	v := &struct {
		Tags []int
	}{}
	_ = Unmarshal([]byte(data), v)
	if len(v.Tags) != 0 {
		t.Error("not empty slice")
	}
}

func TestParser_decode_UnhandledType(t *testing.T) {
	parser := NewParser()
	parser.RegisterDecodeFunc(reflect.String, nil)
	_, err := parser.decode(reflect.TypeOf(""), "s")
	if _, ok := err.(ErrUnhandledType); !ok {
		t.Error("unmatched error")
	}
}

func TestParser_parseForMap_CanSet(t *testing.T) {
	var x = 3.4
	v := reflect.ValueOf(x)
	parser := NewParser()
	parser.parseForMap(v, "")
}

func TestParser_parseForSlice_CanSet(t *testing.T) {
	var x = 3.4
	v := reflect.ValueOf(x)
	parser := NewParser()
	parser.parseForSlice(v, "")
}

//mock multi-layer nested structure,
//BenchmarkUnmarshal-4   	  208219	     14873 ns/op
func BenchmarkUnmarshal(b *testing.B) {
	var data = "Id=1&name=test&child[desc]=c1&child[Long]=10&childPtr[Long]=2&childPtr[Description]=b" +
		"&children[0][desc]=d1&children[1][Long]=12&children[5][desc]=d5&children[5][Long]=50&desc=rtt" +
		"&Params[120]=1&Params[121]=2&status=1&UintPtr=300"
	data = encodeSquareBracket(data)

	for i := 0; i < b.N; i++ {
		v := &testParseInfo{}
		err := Unmarshal([]byte(data), v)
		if err != nil {
			b.Error(err)
		}
	}
}

func encodeSquareBracket(data string) string {
	data = strings.ReplaceAll(data, "[", "%5B")
	data = strings.ReplaceAll(data, "]", "%5D")
	return data
}
