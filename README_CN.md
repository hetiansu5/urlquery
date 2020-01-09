### 简介
使用Go语言实现的application/x-www-form-urlencoded协议数据的转换器。

- 将x-www-form-urlencoded编码的字符串转换为Go数据结构
- 将Go语言数据结构转换为x-www-form-urlencoded编码的字符串

### 关键词
x-www-form-urlencoded HTTP-Query URLEncode URL-Query go

### 特性
- 支持丰富的Go数据结构互转：
    - 基础数据类型: 有符号整型[8,16,32,64] 无符号整形[8,16,32,64] 字符串 布尔值 浮点型[32,64] 字节 字面量
    - 复合数据类型: 数据 切片 Map 结构体
    - 嵌套结构体
- 支持自定义的URL-Encode编码规则，支持全局、局部设置方式，支持默认规则
- 支持自定义的键名映射规则（结构体Tag示例：`query:"name"`）
- 支持开启结构体零值忽略编码，减少编码后的URL-Query字符串长度


### 快速入门
更多查看[example目录](example/withoption.go)

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

### 注意事项
- 针对Map数据类型，Marshal可以支持map[基础数据类型]基础数据类型|复合数据类型，Unmarshal只能支持map[基础数据类型]基础数据类型
- 结构体零值忽略编码默认不开启，开启时需要关注自身业务逻辑的一致性问题
- 字节实际上是uint8，字面量是int32，所以编码后其实是整型，解码的时候也需要接收的是整型

### 许可
MIT
