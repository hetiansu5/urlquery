package main

import (
	"fmt"
	"github.com/hetiansu5/urlquery"
)

//Nested struct
type Point struct {
	X, Y int
}

// test structure
type Circle struct {
	Point //anonymous fields: specially handled
	R int
}

func main() {
	data := &Circle{R: 1}
	data.Point.X = 12
	data.Point.Y = 13
	//output &{{12 13} 1}
	fmt.Println(data)

	bytes, _ := urlquery.Marshal(data)
	//output X=12&Y=13&R=1
	fmt.Println(string(bytes))

	v := Circle{}
	urlquery.Unmarshal(bytes, &v)
	//output {{12 13} 1}
	fmt.Println(v)
}
