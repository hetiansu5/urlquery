### Introduce
An application/x-www-form-urlencoded query translator by go.

- Translate from a x-www-form-urlencoded query string to go structure
- Translate from go structure to a x-www-form-urlencoded query string

### Keywords
x-www-form-urlencoded HTTP-Query URLEncode URL-Query go

### Feature
- Support full go structure Translation
    - Basic Structure: Int[8,16,32,64] Uint[8,16,32,64] String Bool Float[32,64] Byte Rune
    - Complex Structure: Array Slice Map Struct
    - Nested Struct with above Basic or Complex Structure
- Support self-defined URL-Encode rule
- Support self-defined key name relation rule
- Support enable ignoring encoding Zero-value from struct. It can cut down the length of encoded URL-Query string


### Quick Start
More to see [example](example/withoption.go)

```golang
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
```


### Attention
- For Map structure, Marshal supports map[Basic]Basic|Complex, but Unmarshal just supports map[Basic]Basic
- Default: disable ignoring encoding Zero-value of struct. When it is enabled, consistency should be evaluated in your project
- Remember that: Byte is actually uint8, Rune is actually int32


### License
MIT