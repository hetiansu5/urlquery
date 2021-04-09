package main

import (
	"fmt"
	"github.com/hetiansu5/urlquery"
	"net/url"
)

// A OptionChild is test structure
type OptionChild struct {
	Status bool `query:"status"`
	Name   string
}

// An OptionData is test structure
type OptionData struct {
	Id     int
	Name   string          `query:"name"`
	Child  OptionChild     `query:"c"`
	Params map[string]int8 `query:"p"`
	Slice  []OptionChild
}

// A SelfQueryEncoder is test structure
type SelfQueryEncoder struct{}

// Escape is test func
func (u SelfQueryEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

// UnEscape is test func
func (u SelfQueryEncoder) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}

func main() {
	data := OptionData{
		Id:   2,
		Name: "test",
		Child: OptionChild{
			Status: true,
		},
		Params: map[string]int8{
			"one": 1,
			"two": 2,
		},
		Slice: []OptionChild{
			{Status: true},
			{Name: "honey"},
		},
	}

	fmt.Println(data)

	//Marshal: from go structure to http-query string

	builder := urlquery.NewEncoder(urlquery.WithNeedEmptyValue(true),
		urlquery.WithQueryEncoder(SelfQueryEncoder{}))
	bytes, err := builder.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	//Unmarshal: from http-query  string to go structure
	v := &OptionData{}
	parser := urlquery.NewParser(urlquery.WithQueryEncoder(SelfQueryEncoder{}))
	err = parser.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*v)
}
