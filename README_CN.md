### 简介
使用Go语言实现的application/x-www-form-urlencoded协议数据的转换器。

- 将x-www-form-urlencoded编码的字符串转换为Go数据结构
- 将Go语言数据结构转换为x-www-form-urlencoded编码的字符串


### Feature
- 支持丰富的Go数据结构互转：
    - 基础数据类型: 有符号整型(8、16、32、64) 无符号整形(8、16、32、64) 字符串 布尔值 浮点型(32、64) 字节 字面量
    - 复合数据类型: 数据 切片 Map 结构体
    - 嵌套数据体
- 支持自定义的Http query数据编码规则
