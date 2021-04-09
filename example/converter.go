package main

import (
	"github.com/hetiansu5/urlquery"
	"fmt"
	"reflect"
	"strconv"
)

// test structure
type EncodeData struct {
	Id   int    `query:"id"`
	Name string `query:"name"`
}

func main() {
	data := EncodeData{
		Id:   2,
		Name: "Nick",
	}

	encoder := urlquery.NewEncoder()
	encoder.RegisterEncodeFunc(reflect.String, func(rv reflect.Value) string {
		return rv.String() + "Will"
	})

	//Marshal: from go structure to url query string
	bytes, err := encoder.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	//Unmarshal: from url query string to go structure
	parser := urlquery.NewParser()
	parser.RegisterDecodeFunc(reflect.Int, func(s string) (reflect.Value, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i + 10), nil
	})
	v := &EncodeData{}
	err = parser.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*v)
}
