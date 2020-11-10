package urlquery

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//parser from URL Query string to go structure

type parser struct {
	container  map[string]string
	err        error
	opts       options
	mutex      sync.Mutex
	urlEncoder UrlEncoder
}

func NewParser(opts ...Option) *parser {
	p := &parser{}
	for _, o := range opts {
		o.apply(&p.opts)
	}
	return p
}

func (p *parser) init(data []byte) (err error) {
	arr := bytes.Split(data, []byte(SymbolAnd))
	for _, value := range arr {
		ns := strings.SplitN(string(value), SymbolEqual, 2)
		if len(ns) > 1 {
			ns[0], err = p.urlEncoder.UnEscape(ns[0])
			if err != nil {
				return
			}

			ns[1], err = p.urlEncoder.UnEscape(ns[1])
			if err != nil {
				return
			}

			//If last two characters of key equal `[]`, repack it to `[{i++}]`
			l := len(ns[0])
			if l > 2 && ns[0][l-2:] == "[]" {
				//limit iteration to avoid attack of large or dead circle
				for i := 0; i < 1000000; i++ {
					tKey := ns[0][:l-2] + "[" + strconv.Itoa(i) + "]"
					if _, ok := p.container[tKey]; !ok {
						ns[0] = tKey
						break
					}
				}
			}

			p.container[ns[0]] = ns[1]
		}
	}
	return
}

//reset UrlEncoder
func (p *parser) resetUrlEncoder() {
	if p.opts.urlEncoder != nil {
		p.urlEncoder = p.opts.urlEncoder
	} else {
		p.urlEncoder = getUrlEncoder()
	}
}

//generate next parent node key
func (p *parser) genNextParentNode(parentNode, key string) string {
	return genNextParentNode(parentNode, key)
}

func (p *parser) parse(rv reflect.Value, parentNode string) {
	if p.err != nil {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr:
		//If Ptr is nil and can be set, Ptr should be initialized
		if rv.IsNil() {
			if rv.CanSet() {

				//lookup matched map data with prefix key
				matches := p.lookup(parentNode)
				// If none match keep nil
				if len(matches) == 0 {
					break
				}

				rv.Set(reflect.New(rv.Type().Elem()))
				p.parse(rv.Elem(), parentNode)
			}
		} else {
			p.parse(rv.Elem(), parentNode)
		}
	case reflect.Interface:
		p.parse(rv.Elem(), parentNode)
	case reflect.Map:
		if !rv.CanSet() {
			break
		}

		//limited condition of map key and value type
		//If not meet the condition, will return error
		if !isAccessMapKeyType(rv.Type().Key().Kind()) || !isAccessMapValueType(rv.Type().Elem().Kind()) {
			p.err = ErrInvalidMapKeyType{typ: rv.Type()}
			return
		}

		matches := p.lookup(parentNode)
		size := len(matches)

		if size == 0 {
			break
		}

		mapReflect := reflect.MakeMapWithSize(rv.Type(), size)
		for k, _ := range matches {
			reflectKey, err := p.decode(rv.Type().Key(), k)
			if err != nil {
				p.err = err
				return
			}

			value, ok := p.get(p.genNextParentNode(parentNode, k))
			if !ok {
				continue
			}

			reflectValue, err := p.decode(rv.Type().Elem(), value)
			if err != nil {
				p.err = err
				return
			}

			mapReflect.SetMapIndex(reflectKey, reflectValue)
		}
		rv.Set(mapReflect)
	case reflect.Array:
		for i := 0; i < rv.Cap(); i++ {
			p.parse(rv.Index(i), p.genNextParentNode(parentNode, strconv.Itoa(i)))
		}
	case reflect.Slice:
		if !rv.CanSet() {
			break
		}

		//lookup matched map data with prefix key
		matches := p.lookupForSlice(parentNode)
		if len(matches) == 0 {
			break
		}

		//get max cap of slice
		maxCap := 0
		for i, _ := range matches {
			if i+1 > maxCap {
				maxCap = i + 1
			}
		}

		//If slice is nil or cap of slice is less than max cap, slice should be reset correctly
		if rv.IsNil() || maxCap > rv.Cap() {
			rv.Set(reflect.MakeSlice(rv.Type(), maxCap, maxCap))
		}

		for i, _ := range matches {
			p.parse(rv.Index(i), p.genNextParentNode(parentNode, strconv.Itoa(i)))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			ft := rv.Type().Field(i)

			//specially handle anonymous fields
			if ft.Anonymous && rv.Field(i).Kind() == reflect.Struct {
				p.parse(rv.Field(i), parentNode)
				continue
			}

			tag := ft.Tag.Get("query")
			//all ignore
			if tag == "-" {
				continue
			}

			t := newTag(tag)
			name := t.getName()
			if name == "" {
				name = ft.Name
			}

			p.parse(rv.Field(i), p.genNextParentNode(parentNode, name))
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

	v, err := p.decode(rv.Type(), value)
	if err != nil {
		p.err = err
		return
	}

	rv.Set(v)
}

func (p *parser) decode(typ reflect.Type, value string) (v reflect.Value, err error) {
	decoder := getDecoder(typ.Kind())
	if decoder == nil {
		err = ErrUnhandledType{typ: typ}
		return
	}

	v, err = decoder.Decode(value)
	return
}

//lookup by prefix match from container variable
func (p *parser) lookup(prefix string) map[string]bool {
	data := map[string]bool{}
	for k, _ := range p.container {
		if strings.HasPrefix(k, prefix) {
			pre, _ := unpackQueryKey(k[len(prefix):])
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
		if err != nil {
			p.err = err
			break
		}
		data[i] = true
	}
	return data
}

//get value by key from container variable of map struct
func (p *parser) get(key string) (string, bool) {
	v, ok := p.container[key]
	return v, ok
}

/**
 * decode string to go structure
 * Thread safety
 */
func (p *parser) Unmarshal(data []byte, v interface{}) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	//for duplicate use
	p.container = map[string]string{}
	p.err = nil
	p.resetUrlEncoder()

	rv := reflect.ValueOf(v)
	reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidUnmarshalError{typ: reflect.TypeOf(v)}
	}

	err = p.init(data)
	if err != nil {
		return
	}

	p.parse(rv, "")

	//release resource
	p.container = nil
	return p.err
}

/**
 * decode string to go structure
 * Thread safety
 */
func Unmarshal(data []byte, v interface{}) error {
	p := NewParser()
	return p.Unmarshal(data, v)
}
