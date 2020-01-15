package urlquery

import (
	"reflect"
	"strconv"
	"bytes"
	"sync"
)

const (
	SymbolEqual = "="
	SymbolAnd   = "&"
)

//encoder from go structure data to URL Query string

type encoder struct {
	buffer *bytes.Buffer
	err    error
	opts   options
	mutex  sync.Mutex
}

func NewEncoder(opts ...Option) *encoder {
	b := &encoder{}
	for _, o := range opts {
		o.apply(&b.opts)
	}
	return b
}

//get UrlEncoder
func (b *encoder) getUrlEncoder() UrlEncoder {
	if b.opts.urlEncoder != nil {
		return b.opts.urlEncoder
	}
	return getUrlEncoder()
}

//generate next parent node key
func (b *encoder) genNextParentNode(parentNode, key string) string {
	return genNextParentNode(parentNode, key)
}

//unknown structure need to be detected and handled correctly
func (b *encoder) buildQuery(rv reflect.Value, parentNode string, parentKind reflect.Kind) {
	if b.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			//If type of key is interface or ptr, check the real element of key
			checkKey := key
			if key.Kind() == reflect.Interface || key.Kind() == reflect.Ptr {
				checkKey = checkKey.Elem()
			}

			//limited condition of map key type
			if !isAccessMapKeyType(checkKey.Kind()) {
				b.err = ErrInvalidMapKeyType{typ: checkKey.Type()}
				return
			}

			//encode key structure to string
			keyStr, err := b.encode(checkKey)
			if err != nil {
				b.err = err
				return
			}

			b.buildQuery(rv.MapIndex(key), b.genNextParentNode(parentNode, keyStr), rv.Kind())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			b.buildQuery(rv.Index(i), b.genNextParentNode(parentNode, strconv.Itoa(i)), rv.Kind())
		}
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			ft := rt.Field(i)
			//unexported
			if ft.PkgPath != "" && !ft.Anonymous {
				continue
			}

			//specially handle anonymous fields
			if ft.Anonymous && rv.Field(i).Kind() == reflect.Struct {
				b.buildQuery(rv.Field(i), parentNode, rv.Kind())
				continue
			}

			tag := ft.Tag.Get("query")
			//all ignore
			if tag == "-" {
				continue
			}

			t := newTag(tag)
			//get the related name
			name := t.getName()
			if name == "" {
				name = ft.Name
			}

			b.buildQuery(rv.Field(i), b.genNextParentNode(parentNode, name), rv.Kind())
		}
	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			b.buildQuery(rv.Elem(), parentNode, parentKind)
		}
	default:
		b.appendKeyValue(parentNode, rv, parentKind)
	}
}

//basic structure can be translated directly
func (b *encoder) appendKeyValue(key string, rv reflect.Value, parentKind reflect.Kind) {
	//If parent type is struct and empty value will be ignored by default. unless needEmptyValue is true.
	if parentKind == reflect.Struct && !b.opts.needEmptyValue && isEmptyValue(rv) {
		return
	}

	//If parent type is slice or array, then repack key. eg. students[0] -> students[]
	if parentKind == reflect.Slice || parentKind == reflect.Array {
		key = repackArrayQueryKey(key)
	}

	s, err := b.encode(rv)
	if err != nil {
		b.err = err
		return
	}

	b.buffer.WriteString(b.getUrlEncoder().Escape(key) + SymbolEqual + b.getUrlEncoder().Escape(s) + SymbolAnd)
}

func (b *encoder) encode(rv reflect.Value) (s string, err error) {
	encoder := getEncoder(rv.Kind())
	if encoder == nil {
		err = ErrUnhandledType{typ: rv.Type()}
		return
	}

	s = encoder.Encode(rv)
	return
}

/**
 * encode go structure to string
 * Thread safety
 */
func (b *encoder) Marshal(data interface{}) ([]byte, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	//for duplicate call
	b.buffer = new(bytes.Buffer)
	b.err = nil

	rv := reflect.ValueOf(data)
	b.buildQuery(rv, "", reflect.Interface)
	if b.err != nil {
		return nil, b.err
	}

	bs := b.buffer.Bytes()
	//release resource
	b.buffer = nil
	//do not forget to remove the last & character
	return bs[:len(bs)-1], nil
}

/**
 * encode go structure to string
 * Thread safety
 */
func Marshal(data interface{}) ([]byte, error) {
	b := NewEncoder()
	return b.Marshal(data)
}
