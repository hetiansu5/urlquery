package query

import (
	"testing"
)

type Child struct {
	Description string `query:"desc"`
	Height      int    `query:" ignore"`
}

type TestInfo struct {
	Id       int
	Name     string  `query:"name"`
	Child    Child   `query:"child"`
	ChildPtr *Child  `query:"childPtr"`
	Children []Child `query:"children"`

}

func TestMarshal(t *testing.T) {
	var (
		f32 = float32(1.2)
		f64 = float64(13.4343453535343242342)
		i8  = int8(3)
		i64 = int64(9999999 * 9999999)
		u64 = uint16(567)
	)
	params := map[string]interface{}{
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
			"child": TestInfo{
				Name: "child",
				Children: []Child{
					{Description: "d1", Height: 180},
					{Description: "d2"},
				},
				Child:    Child{Description: "c1"},
				ChildPtr: &Child{Description: "cptr"},
			},
		},
		"struct": TestInfo{
			Id:   222,
			Name: "test",
		},
	}

	_, err := Marshal(params)

	if err != nil {
		t.Error(err)
	}
}
