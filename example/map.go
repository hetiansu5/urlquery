package main

import (
	"github.com/hetiansu5/urlquery"
	"fmt"
)

func main() {
	runes := []rune("和")
	data := map[string]interface{}{
		"one": []int{1, 2, 3},
		"two": true,
		"three": map[uint8]float32{
			1: 1.222,
		},
		"four":  "test",
		"byte":  'f', //byte is actually uint8
		"float": 2.4343434,
		"rune":  runes[0], //rune is actually int32
	}

	bytes, err := urlquery.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	var m map[string]string
	s := "id=1&name=ab&arr%5B0%5D=6d"
	err = urlquery.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println(err)
		return
	}

	//remember: just id and name can be parsed
	//arr[0] is not match map[string]string
	fmt.Println(m)
}
