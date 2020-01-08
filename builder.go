package query

import (
	"reflect"
	"net/url"
	"strconv"
	"sync"
	"fmt"
	"bytes"
)

type builder struct {
	output []byte
	mutex  sync.RWMutex
	err    error
}

func (q *builder) appendByte(b byte) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.output = append(q.output, b)
}

func (q *builder) appendBytes(bytes []byte) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.output = append(q.output, bytes...)
}

func (q *builder) appendString(s string) {
	q.appendBytes([]byte(s))
}

func (q *builder) GetBytes() []byte {
	return bytes.TrimRight(q.output, "&")
}

func (q *builder) buildQuery(rv reflect.Value, parentNode string) {
	if q.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			q.buildQuery(rv.Elem(), parentNode)
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			q.buildQuery(rv.MapIndex(key), genNextParentNode(parentNode, fmt.Sprint(key)))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			q.buildQuery(rv.Index(i), genNextParentNode(parentNode, strconv.Itoa(i)))
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

			q.buildQuery(rv.Field(i), genNextParentNode(parentNode, key))
		}
	default:
		q.appendKeyValue(parentNode, rv)
	}
}

func (q *builder) appendKeyValue(parentNode string, rv reflect.Value) {
	encoder := getEncoder(rv.Kind())
	if encoder == nil {
		s, _ := url.QueryUnescape(parentNode);
		q.err = ErrUnhandledType{key: s, t: rv.Type()}
		return
	}

	q.appendString(parentNode + "=" + url.QueryEscape(encoder.Encode(rv)) + "&")
	return
}

//HttpBuildQuery
func Marshal(data interface{}) ([]byte, error) {
	q := &builder{}
	rv := reflect.ValueOf(data)
	q.buildQuery(rv, "")
	if q.err != nil {
		return nil, q.err
	}
	return q.GetBytes(), nil
}
