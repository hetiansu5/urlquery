package main

import (
	"fmt"
	"github.com/hetiansu5/urlquery"
)

type Point struct {
	X, Y int
}

type Circle struct {
	Point //anonymous fields: specially handled
	R int
}

func main() {
	data := &Circle{R: 1}
	data.Point.X = 12
	data.Point.Y = 13
	fmt.Println(data)

	bytes, _ := urlquery.Marshal(data)
	fmt.Println(string(bytes))

	v := Circle{}
	urlquery.Unmarshal(bytes, &v)
	fmt.Println(v)
}
