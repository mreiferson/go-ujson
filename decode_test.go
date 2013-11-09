package ujson

import (
	"testing"
)

type jsDecStore struct {}

func (s *jsDecStore) NewObject() (interface{}, error) {
	return make(map[string]interface{}), nil
}

func (s *jsDecStore) NewArray() (interface{}, error) {
	a := make([]interface{}, 0)
	return &a, nil
}

func (s *jsDecStore) ObjectAddKey(o interface{}, k string, v interface{}) error {
	m := o.(map[string]interface{})
	m[k] = v
	return nil
}

func (s *jsDecStore) ArrayAddItem(o interface{}, v interface{}) error {
	a := o.(*[]interface{})
	*a = append(*a, v)
	return nil
}

func (s *jsDecStore) NewString(b []byte) (interface{}, error) {
	return b, nil
}

func (s *jsDecStore) NewNumeric(b []byte) (interface{}, error) {
	return b, nil
}

func (s *jsDecStore) NewTrue() (interface{}, error) {
	return true, nil
}

func (s *jsDecStore) NewFalse() (interface{}, error) {
	return false, nil
}

func (s *jsDecStore) NewNull() (interface{}, error) {
	return nil, nil
}

func TestDecode(t *testing.T) {
	testData := []byte(`{ "test": "hello world", "t2": ["a", 3, "c"], "asdf4": 0.14159, "sf": { "v": [4, 5], "z": "hw2" } }`)
	dec := NewDecoder(&jsDecStore{}, testData)
	dec.Decode()
}
