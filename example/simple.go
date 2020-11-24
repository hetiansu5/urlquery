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
	Id         int
	Name       string          `query:"name"`
	Child      SimpleChild
	Params     map[string]int8 `query:"p"`
	SliceChild []SimpleChild   `query:"s"`
	Password   string          `query:"-"` //- means ignoring password
	Array      [3]uint16
}

func main() {
	data := SimpleData{
		Id:   2,
		Name: "http://localhost/test.php?id=2",
		Child: SimpleChild{
			Status: true,
		},
		Params: map[string]int8{
			"one": 1,
			"two": 2,
		},
		SliceChild: []SimpleChild{
			{Status: true},
			{Name: "honey"},
		},
		Password: "abc",
		Array:    [3]uint16{2, 3, 300},
	}

	fmt.Println(data)

	//Marshal: from go structure to url query string
	bytes, err := urlquery.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	//Unmarshal: from url query string to go structure
	v := &SimpleData{}
	err = urlquery.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*v)
}
