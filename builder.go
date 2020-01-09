package urlquery

import (
	"reflect"
	"strconv"
	"sync"
	"bytes"
)

//Translate from go structure data to a x-www-form-urlencoded form string

type builder struct {
	output []byte
	mutex  sync.RWMutex
	err    error
	opts   options
}

func NewBuilder(opts ...Option) *builder {
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
	if b.opts.urlEncoder != nil {
		return b.opts.urlEncoder.Escape(s)
	}
	return getUrlEncoder().Escape(s)
}

//generate next parent node key
func (b *builder) genNextParentNode(parentNode, key string) string {
	if len(parentNode) > 0 {
		return parentNode + b.urlEncode("["+key+"]")
	} else {
		return b.urlEncode(key)
	}
}

//unknown structure need to be detected and handled correctly
func (b *builder) buildQuery(rv reflect.Value, parentNode string, parentKind reflect.Kind) {
	if b.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			b.buildQuery(rv.Elem(), parentNode, rv.Kind())
		}
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
			tag := rt.Field(i).Tag.Get("query")
			key := rt.Field(i).Name

			if tag != "" {
				t := newTag(tag)
				if t.contains("outputIgnore", "ignore") {
					continue
				}
				//get the related name
				if t.getName() != "" {
					key = t.getName()
				}
			}

			b.buildQuery(rv.Field(i), b.genNextParentNode(parentNode, key), rv.Kind())
		}
	default:
		b.appendKeyValue(parentNode, rv, parentKind)
	}
}

//basic structure can be translated directly
func (b *builder) appendKeyValue(parentNode string, rv reflect.Value, parentKind reflect.Kind) {
	//when parent type is struct, empty value will be ignored by default. unless needEmptyValue is true.
	if parentKind == reflect.Struct {
		if !b.opts.needEmptyValue && isEmptyValue(rv) {
			return
		}
	}

	s, err := b.encode(rv)
	if err != nil {
		b.err = err
		return
	}

	b.appendString(parentNode + "=" + b.urlEncode(s) + "&")
}

func (b *builder) encode(rv reflect.Value) (s string, err error) {
	encoder := getEncoder(rv.Kind())
	if encoder == nil {
		err = ErrUnhandledType{typ: rv.Type()}
		return
	}

	s = encoder.Encode(rv)
	return
}

//encode go structure to string
func (b *builder) Marshal(data interface{}) ([]byte, error) {
	rv := reflect.ValueOf(data)
	b.buildQuery(rv, "", reflect.Interface)
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
