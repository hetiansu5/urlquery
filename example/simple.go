package main

import (
	"github.com/hetiansu5/urlquery"
	"fmt"
)

type SimpleChild struct {
	Status bool `query:"status"`
	Name   string
}

type SimpleData struct {
	Id     int
	Name   string          `query:"name"`
	Child  SimpleChild     `query:"c"`
	Params map[string]int8 `query:"p"`
	Slice  []SimpleChild
}

func main() {
	data := SimpleData{
		Id:   2,
		Name: "test",
		Child: SimpleChild{
			Status: true,
		},
		Params: map[string]int8{
			"one": 1,
			"two": 2,
		},
		Slice: []SimpleChild{
			{Status: true},
			{Name: "honey"},
		},
	}

	fmt.Println(data)

	//Marshal: from go structure to http-query string
	bytes, err := urlquery.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	//Unmarshal: from http-query  string to go structure
	v := &SimpleData{}
	err = urlquery.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*v)
}
