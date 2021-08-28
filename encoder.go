package urlquery

import (
	"bytes"
	"reflect"
	"strconv"
	"sync"
)

const (
	// SymbolEqual is key character of querystring
	SymbolEqual = "="
	// SymbolAnd is key character of querystring
	SymbolAnd = "&"
)

// A encoder from go structure data to URL Query string
type encoder struct {
	buffer        *bytes.Buffer
	err           error
	opts          options
	mutex         sync.Mutex
	queryEncoder  QueryEncoder
	encodeFuncMap map[reflect.Kind]valueEncode
}

// NewEncoder return new encoder object
// do some option initialization
func NewEncoder(opts ...Option) *encoder {
	b := &encoder{}
	for _, option := range opts {
		option(&b.opts)
	}
	b.encodeFuncMap = make(map[reflect.Kind]valueEncode)
	return b
}

// reset query encoder
func (b *encoder) resetQueryEncoder() {
	if b.opts.queryEncoder != nil {
		b.queryEncoder = b.opts.queryEncoder
	} else {
		b.queryEncoder = getQueryEncoder()
	}
}

// generate next parent node key
func (b *encoder) genNextParentNode(parentNode, key string) string {
	return genNextParentNode(parentNode, key)
}

// detect type of value via reflect, handle correctly
func (b *encoder) buildQuery(rv reflect.Value, parentNode string, parentKind reflect.Kind) {
	if b.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Map:
		b.buildQueryForMap(rv, parentNode)
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			b.buildQuery(rv.Index(i), b.genNextParentNode(parentNode, strconv.Itoa(i)), rv.Kind())
		}
	case reflect.Struct:
		b.buildQueryForStruct(rv, parentNode)
	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			b.buildQuery(rv.Elem(), parentNode, parentKind)
		}
	default:
		b.appendKeyValue(parentNode, rv, parentKind)
	}
}

// build query string for map value
func (b *encoder) buildQueryForMap(rv reflect.Value, parentNode string) {
	for _, key := range rv.MapKeys() {
		//If type of key is interface or ptr, check the pointed element of key
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
}

// build query string for struct value
func (b *encoder) buildQueryForStruct(rv reflect.Value, parentNode string) {
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
}

// basic structure can be translated directly
func (b *encoder) appendKeyValue(key string, rv reflect.Value, parentKind reflect.Kind) {
	//If parent type is struct and empty value will be ignored by default. unless needEmptyValue is true.
	if parentKind == reflect.Struct && !b.opts.needEmptyValue && isZeroValue(rv) {
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

	b.buffer.WriteString(b.queryEncoder.Escape(key) + SymbolEqual + b.queryEncoder.Escape(s) + SymbolAnd)
}

// encode a specified-type value to string
func (b *encoder) encode(rv reflect.Value) (s string, err error) {
	encodeFunc := b.getEncodeFunc(rv.Kind())
	if encodeFunc == nil {
		err = ErrUnhandledType{typ: rv.Type()}
		return
	}
	s = encodeFunc(rv)
	return
}

// get encode function for specified reflect kind
func (b *encoder) getEncodeFunc(kind reflect.Kind) valueEncode {
	if encodeFunc, ok := b.encodeFuncMap[kind]; ok {
		return encodeFunc
	}
	return getEncodeFunc(kind)
}

// register self-defined encode function for any reflect kind
func (b *encoder) RegisterEncodeFunc(kind reflect.Kind, encode valueEncode) {
	b.encodeFuncMap[kind] = encode
}

// Marshal do encoding go structure to string
// it is thread safety
func (b *encoder) Marshal(data interface{}) ([]byte, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	//for duplicate call
	b.buffer = new(bytes.Buffer)
	b.err = nil
	b.resetQueryEncoder()

	rv := reflect.ValueOf(data)
	b.buildQuery(rv, "", reflect.Interface)
	if b.err != nil {
		return nil, b.err
	}

	bs := b.buffer.Bytes()
	//release resource
	b.buffer = nil
	//do not forget to remove the last & character
	if len(bs) == 0 {
		return bs, nil
	}
	return bs[:len(bs)-1], nil
}

// Marshal do encoding go structure to string
// it is thread safety
func Marshal(data interface{}) ([]byte, error) {
	b := NewEncoder()
	return b.Marshal(data)
}
