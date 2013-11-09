package ujson

import (
	"errors"
	"fmt"
	"log"
)

type JSON struct {
	root interface{}
}

func NewFromBytes(data []byte) (*JSON, error) {
	j := &JSON{}
	dec := NewDecoder(j, data)
	root, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	j.root = root
	return j, nil
}

func (j *JSON) NewObject() (interface{}, error) {
	return make(map[string]interface{}), nil
}

func (j *JSON) NewArray() (interface{}, error) {
	a := make([]interface{}, 0)
	return &a, nil
}

func (j *JSON) ObjectAddKey(mi interface{}, k interface{}, v interface{}) error {
	ks := k.(string)
	m := mi.(map[string]interface{})
	m[ks] = v
	return nil
}

func (j *JSON) ArrayAddItem(ai interface{}, v interface{}) error {
	a := ai.(*[]interface{})
	*a = append(*a, v)
	return nil
}

func (j *JSON) NewString(b []byte) (interface{}, error) {
	str, ok := unquote(b)
	if !ok {
		return nil, errors.New(fmt.Sprintf("failed to unquote string %s", b))
	}
	return str, nil
}

func (j *JSON) NewNumeric(b []byte) (interface{}, error) {
	return b, nil
}

func (j *JSON) NewTrue() (interface{}, error) {
	return true, nil
}

func (j *JSON) NewFalse() (interface{}, error) {
	return false, nil
}

func (j *JSON) NewNull() (interface{}, error) {
	return nil, nil
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *JSON) Get(key string) *JSON {
	m, err := j.MaybeMap()
	if err == nil {
		if val, ok := m[key]; ok {
			return &JSON{val}
		}
	}
	return &JSON{nil}
}

// Map guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").Map() {
//			fmt.Println(k, v)
//		}
func (j *JSON) Map(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Map() received too many arguments %d", len(args))
	}

	a, err := j.MaybeMap()
	if err == nil {
		return a
	}

	return def
}

// MaybeMap type asserts to `map`
func (j *JSON) MaybeMap() (map[string]interface{}, error) {
	if m, ok := (j.root).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// String guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").String(), js.Get("optional_param").String("my_default"))
func (j *JSON) String(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("String() received too many arguments %d", len(args))
	}

	s, err := j.MaybeString()
	if err == nil {
		return s
	}

	return def
}

// MaybeString type asserts to `string`
func (j *JSON) MaybeString() (string, error) {
	if s, ok := (j.root).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}
