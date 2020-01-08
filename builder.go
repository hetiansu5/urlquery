package query

import (
	"reflect"
	"strconv"
	"sync"
	"fmt"
	"bytes"
)

//Translate from go structure data to a x-www-form-urlencoded form string

type builder struct {
	output []byte
	mutex  sync.RWMutex
	err    error
	opts   builderOptions
}

func NewBuilder(opts ...BuilderOption) *builder {
	b := &builder{}
	for _, o := range opts {
		o.apply(&b.opts)
	}
	return b
}

func (b *builder) appendBytes(bytes []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.output = append(b.output, bytes...)
}

func (b *builder) appendString(s string) {
	b.appendBytes([]byte(s))
}

//when finish, get the result string
func (b *builder) GetBytes() []byte {
	return bytes.TrimRight(b.output, "&")
}

//urlEncode
func (b *builder) urlEncode(s string) string {
	if b.opts.u != nil {
		return b.opts.u.Escape(s)
	}
	return GetUrlEncoder().Escape(s)
}

//unknown structure need to be detected and handled correctly
func (b *builder) buildQuery(rv reflect.Value, parentNode string) {
	if b.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			b.buildQuery(rv.Elem(), parentNode)
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			b.buildQuery(rv.MapIndex(key), genNextParentNode(parentNode, fmt.Sprint(key)))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			b.buildQuery(rv.Index(i), genNextParentNode(parentNode, strconv.Itoa(i)))
		}
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			tag := rt.Field(i).Tag.Get("query")
			key := rt.Field(i).Name

			if tag != "" {
				t := newTag(tag)
				if t.hasFlag("outputIgnore") || t.hasFlag("ignore") {
					continue
				}
				if t.getName() != "" {
					key = t.getName()
				}
			}

			b.buildQuery(rv.Field(i), genNextParentNode(parentNode, key))
		}
	default:
		b.appendKeyValue(parentNode, rv)
	}
}

//basic structure can be translated directly
func (b *builder) appendKeyValue(parentNode string, rv reflect.Value) {
	encoder := getEncoder(rv.Kind())
	if encoder == nil {
		b.err = ErrUnhandledType{typ: rv.Type()}
		return
	}

	b.appendString(parentNode + "=" + b.urlEncode(encoder.Encode(rv)) + "&")
	return
}

//encode go structure to string
func (b *builder) Marshal(data interface{}) ([]byte, error) {
	rv := reflect.ValueOf(data)
	b.buildQuery(rv, "")
	if b.err != nil {
		return nil, b.err
	}
	return b.GetBytes(), nil
}

//encode go structure to string
func Marshal(data interface{}) ([]byte, error) {
	b := NewBuilder()
	return b.Marshal(data)
}
