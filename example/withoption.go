package main

import (
	"github.com/hetiansu5/urlquery"
	"fmt"
	"net/url"
)

type OptionChild struct {
	Status bool `query:"status"`
	Name   string
}

type OptionData struct {
	Id     int
	Name   string          `query:"name"`
	Child  OptionChild     `query:"c"`
	Params map[string]int8 `query:"p"`
	Slice  []OptionChild
}

type SelfUrlEncoder struct{}

func (u SelfUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

func (u SelfUrlEncoder) UnEscape(s string) (string, error) {
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
		urlquery.WithUrlEncoder(SelfUrlEncoder{}))
	bytes, err := builder.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	//Unmarshal: from http-query  string to go structure
	v := &OptionData{}
	parser := urlquery.NewParser(urlquery.WithUrlEncoder(SelfUrlEncoder{}))
	err = parser.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*v)
}
