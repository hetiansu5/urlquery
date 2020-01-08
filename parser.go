package query

import (
	"reflect"
	"strings"
	"strconv"
	"bytes"
)

//translator from a x-www-form-urlencoded form string to go structure

type parser struct {
	container map[string]string
	err       error
}

func NewParser(data []byte) *parser {
	p := &parser{
		container: map[string]string{},
	}
	return p
}

func (p *parser) init(data []byte) {
	arr := bytes.Split(data, []byte("&"))
	for _, value := range arr {
		ns := strings.SplitN(string(value), "=", 2)
		if len(ns) > 1 {
			p.container[ns[0]] = ns[1]
		}
	}
}

func (p *parser) parse(rv reflect.Value, parentNode string) {
	if p.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			if rv.CanSet() {
				rv.Set(reflect.New(rv.Type().Elem()))
				p.parse(rv.Elem(), parentNode)
			}
		} else {
			p.parse(rv.Elem(), parentNode)
		}
	case reflect.Map:
		if !rv.CanSet() {
			break
		}

		//limit condition of map struct
		//If not meet the condition, will return error
		if !isAccessMapKeyType(rv.Type().Key().Kind()) || !isAccessMapValueType(rv.Type().Elem().Kind()) {
			p.err = ErrUnhandledType{typ: rv.Type()}
			return
		}

		matches := p.lookup(parentNode)
		size := len(matches)

		if size == 0 {
			break
		}

		mapReflect := reflect.MakeMapWithSize(rv.Type(), size)
		for k, _ := range matches {
			reflectKey, err := p.decode(rv.Type().Key().Kind(), k)
			if err != nil {
				return
			}

			value, ok := p.get(genNextParentNode(parentNode, k))
			if !ok {
				continue
			}

			reflectValue, err := p.decode(rv.Type().Elem().Kind(), value)
			if err != nil {
				return
			}

			mapReflect.SetMapIndex(reflectKey, reflectValue)
		}
		rv.Set(mapReflect)
	case reflect.Array:
		for i := 0; i < rv.Cap(); i++ {
			p.parse(rv.Index(i), genNextParentNode(parentNode, strconv.Itoa(i)))
		}
	case reflect.Slice:
		if !rv.CanSet() {
			break
		}

		matches := p.lookupForSlice(parentNode)
		if len(matches) == 0 {
			break
		}

		maxCap := 0
		for i, _ := range matches {
			if i+1 > maxCap {
				maxCap = i + 1
			}
		}

		if rv.IsNil() || maxCap > rv.Cap() {
			rv.Set(reflect.MakeSlice(rv.Type(), maxCap, maxCap))
		}

		for i, _ := range matches {
			p.parse(rv.Index(i), genNextParentNode(parentNode, strconv.Itoa(i)))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			tag := rv.Type().Field(i).Tag.Get("query")
			key := rv.Type().Field(i).Name

			if tag != "" {
				t := newTag(tag)
				if t.hasFlag("inputIgnore") || t.hasFlag("ignore") {
					continue
				}
				if t.getName() != "" {
					key = t.getName()
				}
			}

			p.parse(rv.Field(i), genNextParentNode(parentNode, key))
		}
	default:
		p.parseValue(parentNode, rv)
	}
}

func (p *parser) parseValue(parentNode string, rv reflect.Value) {
	if !rv.CanSet() {
		return
	}

	value, ok := p.get(parentNode)
	if !ok {
		return
	}

	v, err := p.decode(rv.Kind(), value)
	if err != nil {
		return
	}

	rv.Set(v)
}

func (p *parser) decode(kind reflect.Kind, value string) (v reflect.Value, err error) {
	decoder := getDecoder(kind)
	v, err = decoder.Decode(value)
	if err != nil {
		p.err = err
	}
	return
}

//lookup by prefix match from container variable
func (p *parser) lookup(prefix string) map[string]bool {
	data := map[string]bool{}
	for k, _ := range p.container {
		if strings.HasPrefix(k, prefix) {
			pre, _ := unpackKey(k[len(prefix):])
			data[pre] = true
		}
	}
	return data
}

//lookup by prefix match from container variable
func (p *parser) lookupForSlice(prefix string) map[int]bool {
	tmp := p.lookup(prefix)
	data := map[int]bool{}
	for k, _ := range tmp {
		i, err := strconv.Atoi(k)
		if err == nil {
			data[i] = true
		}
	}
	return data
}

//get value by key from container variable of map struct
func (p *parser) get(key string) (string, bool) {
	v, ok := p.container[key]
	return v, ok
}

//decode string to go structure
func (p *parser) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidUnmarshalError{typ: reflect.TypeOf(v)}
	}

	p.init(data)
	p.parse(rv, "")
	return p.err
}

//decode string to go structure
func Unmarshal(data []byte, v interface{}) error {
	p := NewParser(data)
	return p.Unmarshal(data, v)
}
