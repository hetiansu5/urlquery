### Introduce
A URL Query string Encoder and Parser by go.

- Parse from URL Query string to go structure
- Encode from go structure to URL Query string

### Keywords
x-www-form-urlencoded HTTP-Query URLEncode URL-Query go

### Feature
- Support full go structure Translation
    - Basic Structure: Int[8,16,32,64] Uint[8,16,32,64] String Bool Float[32,64] Byte Rune
    - Complex Structure: Array Slice Map Struct
    - Nested Struct with above Basic or Complex Structure
- Support top-level structure: Map, Slice or Array, not only Struct
- Support self-defined URL-Encode rule
- Support self-defined key name relation rule
- Support to control whether ignoring Zero-value of struct member


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
	Id         int
	Name       string          `query:"name"`
	Child      SimpleChild
	Params     map[string]int8 `query:"p"`
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
		},
		Array: [3]uint16{2, 3, 300},
	}

	//Marshal: from go structure to url query string
	bytes, _ := urlquery.Marshal(data)
	fmt.Println(string(bytes))

	//Unmarshal: from url query  string to go structure
	v := &SimpleData{}
	urlquery.Unmarshal(bytes, v)
	fmt.Println(*v)
}
```


### Attention
- For Map structure, Marshal supports map[Basic]Basic|Complex, but Unmarshal just supports map[Basic]Basic
- Default: ignoring Zero-value of struct member. You can enable it with `Option`
- Remember that: Byte is actually uint8, Rune is actually int32


### License
MIT