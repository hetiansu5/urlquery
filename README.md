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
- Support to control whether ignoring Zero-value of struct member. It will cut down the length of encoded URL-Query string


### Quick Start
More to see [example](example/withoption.go)

```golang
package main

import (
	"github.com/hetiansu5/urlquery"
	"fmt"
)

type SimpleData struct {
	Id     int
	Name   string          `query:"name"`
}

func main() {
	data := SimpleData{
		Id:   2,
		Name: "test",
	}

	//Marshal: from go structure to http-query string
	bytes, err := urlquery.Marshal(data)
	fmt.Println(string(bytes), err)

	//Unmarshal: from http-query  string to go structure
	v := &SimpleData{}
	err = urlquery.Unmarshal(bytes, v)
	fmt.Println(*v, err)
}
```


### Attention
- For Map structure, Marshal supports map[Basic]Basic|Complex, but Unmarshal just supports map[Basic]Basic
- Default: ignoring Zero-value of struct member. You can enable it with `Option`
- Remember that: Byte is actually uint8, Rune is actually int32


### License
MIT