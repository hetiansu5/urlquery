### Introduce
an application/x-www-form-urlencoded agreement data translator by go.

- Translate from a x-www-form-urlencoded form string to go structure
- Translate from go structure data to a x-www-form-urlencoded form string


### Feature
- Support full go structure Translation
    - Basic Structure: Int(8、16、32、64) Uint(8、16、32、64) String Bool Float(32、64) Byte Rune
    - Complex Structure: Array Slice Map Struct
    - Nested Struct
- Support customized Http query encoding rule