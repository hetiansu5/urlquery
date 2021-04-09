[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/hetiansu5/urlquery?status.svg)](https://godoc.org/github.com/hetiansu5/urlquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/hetiansu5/urlquery)](https://goreportcard.com/report/github.com/hetiansu5/urlquery)
[![License](https://img.shields.io/github/license/hetiansu5/urlquery)](LICENSE)
[![codecov](https://codecov.io/gh/hetiansu5/urlquery/branch/master/graph/badge.svg?token=N5NWVBHRTY)](https://codecov.io/gh/hetiansu5/urlquery)

## Introduction
A URL Query string Encoder and Parser based on go.

- Parse from URL Query string to go structure
- Encode from go structure to URL Query string

## Keywords
x-www-form-urlencoded Query Encoder URL-Query Http-Query go

## Feature
- Support full go structure Translation
    - Basic Structure: Int[8,16,32,64] Uint[8,16,32,64] String Bool Float[32,64] Byte Rune
    - Complex Structure: Array Slice Map Struct
    - Nested Struct with above Basic or Complex Structure
- Support top-level structure: Map, Slice or Array, not only Struct
- Support self-defined URL Query Encode rule [example](example/withoption.go)
- Support self-defined key name relation rule [example](example/simple.go)
- Support self-defined value encode and decode function [example](example/converter.go)
- Support to control whether ignoring Zero-value of struct member [example](example/withoption.go)


## Quick Start
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
	//output Id=2&name=http%3A%2F%2Flocalhost%2Ftest.php%3Fid%3D2&Child%5Bstatus%5D=1&p%5Bone%5D=1&Array%5B%5D=2&Array%5B%5D=3&Array%5B%5D=300
	fmt.Println(string(bytes))

	//Unmarshal: from url query  string to go structure
	v := &SimpleData{}
	urlquery.Unmarshal(bytes, v)
	//output {Id:2, Name:"http://localhost/test.php?id=2", Child: SimpleChild{Status:true}, Params:map[one:1], Array:[2, 3, 300]}
	fmt.Println(*v)
}
```


### Attention
- For Map structure, Marshal supports map[Basic]Basic|Complex, but Unmarshal just supports map[Basic]Basic
- Default: ignoring Zero-value of struct member. You can enable it with [Option](example/withoption.go)
- Remember that: Byte is actually uint8, Rune is actually int32


### License
MIT
