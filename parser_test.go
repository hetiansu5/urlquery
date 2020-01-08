package query

import (
	"testing"
	"strings"
)

type ParseChild struct {
	Description string `query:"desc"`
	Long        uint16 `query:" vip"`
	Height      int    `query:" ignore"`
}

type ParseInfo struct {
	Id       int
	Name     string       `query:"name"`
	Child    ParseChild   `query:"child"`
	ChildPtr *ParseChild  `query:"childPtr"`
	Children []ParseChild `query:"children"`
	Params   map[string]int8
}

func TestUnmarshal_ComplexStruct(t *testing.T) {
	var data = "Id=1&name=test&child[desc]=c1&child[Long]=10&childPtr[Long]=2&childPtr[Description]=b" +
		"&children[0][desc]=d1&children[1][Long]=12&children[5][desc]=d5&children[5][Long]=50&desc=rtt&Params[a]=1&Params[b]=2"
	data = strings.ReplaceAll(data, "[", "%5B")
	data = strings.ReplaceAll(data, "]", "%5D")
	v := &ParseInfo{Params: map[string]int8{}}
	err := Unmarshal([]byte(data), v)

	if err != nil {
		t.Error(err)
	}

	if v.Id != 1 {
		t.Error("Id's data is wrong")
	}

	if v.Name != "test" {
		t.Error("Name's data is wrong")
	}

	if v.Child.Description != "c1" || v.Child.Long != 10 || v.Child.Height != 0 {
		t.Error("Child's data is wrong")
	}

	if v.ChildPtr == nil || v.ChildPtr.Description != "" || v.ChildPtr.Long != 2 || v.ChildPtr.Height != 0 {
		t.Error("ChildPtr's data is wrong")
	}

	if len(v.Children) != 6 {
		t.Error("Children's length is wrong")
	}

	if v.Children[0].Description != "d1" {
		t.Error("Children0's data is wrong")
	}

	if v.Children[1].Description != "" || v.Children[1].Long != 12 {
		t.Error("Children1's data is wrong")
	}

	if v.Children[2].Description != "" || v.Children[3].Description != "" || v.Children[4].Description != "" {
		t.Error("Children234's data is wrong")
	}

	if v.Children[5].Description != "d5" || v.Children[5].Long != 50 || v.Children[5].Height != 0 {
		t.Error("Children5's data is wrong")
	}

	if len(v.Params) != 2 || v.Params["a"] != 1 || v.Params["b"] != 2 {
		t.Error("Params's data is wrong")
	}
}

//TODO repeat Nesting
