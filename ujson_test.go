package ujson

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	testData := []byte(`{
		"test": "hello world",
		"i64": 123456,
		"t2": ["a", 3, "c"],
		"asdf4": 0.14159,
		"sf": {
			"v": [4, 5],
			"z": "hw2"
		},
		"unicode": "M\u00fcNSTER"
	}`)
	obj, err := NewFromBytes(testData)
	if err != nil {
		t.Fatalf(err.Error())
	}
	testStr := obj.Get("test").String()
	if testStr != "hello world" {
		t.Fatalf(`key "test" (%s) should == "hello world"`, testStr)
	}
	testInt64 := obj.Get("i64").Int64()
	if testInt64 != 123456 {
		t.Fatalf(`key "i64" (%d) should == 123456`, testInt64)
	}
	testUnicodeStr := obj.Get("unicode").String()
	if testUnicodeStr != "MüNSTER" {
		t.Fatalf(`key "unicode" (%s) should == "MüNSTER"`, testUnicodeStr)
	}
}

func BenchmarkUjson(b *testing.B) {
	b.StopTimer()
	f, err := os.Open("testdata/small.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	dec := NewDecoder(simpleStore{}, data)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := dec.Decode()
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
	b.SetBytes(int64(len(data)))
}

func BenchmarkStdLib(b *testing.B) {
	b.StopTimer()
	f, err := os.Open("testdata/small.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		err := json.Unmarshal(data, &m)
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
	b.SetBytes(int64(len(data)))
}
