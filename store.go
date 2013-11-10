package ujson

import (
	"errors"
	"fmt"
	"strconv"
)

type simpleStore struct{}

func (s simpleStore) NewObject() (interface{}, error) {
	return make(map[string]interface{}), nil
}

func (s simpleStore) NewArray() (interface{}, error) {
	a := make([]interface{}, 0)
	return &a, nil
}

func (s simpleStore) ObjectAddKey(mi interface{}, k interface{}, v interface{}) error {
	ks := k.(string)
	m := mi.(map[string]interface{})
	m[ks] = v
	return nil
}

func (s simpleStore) ArrayAddItem(ai interface{}, v interface{}) error {
	a := ai.(*[]interface{})
	*a = append(*a, v)
	return nil
}

func (s simpleStore) NewString(b []byte) (interface{}, error) {
	str, ok := unquote(b)
	if !ok {
		return nil, errors.New(fmt.Sprintf("failed to unquote string %s", b))
	}
	return str, nil
}

type numeric []byte

func (n numeric) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

func (n numeric) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

func (s simpleStore) NewNumeric(b []byte) (interface{}, error) {
	return numeric(b), nil
}

func (s simpleStore) NewTrue() (interface{}, error) {
	return true, nil
}

func (s simpleStore) NewFalse() (interface{}, error) {
	return false, nil
}

func (s simpleStore) NewNull() (interface{}, error) {
	return nil, nil
}
