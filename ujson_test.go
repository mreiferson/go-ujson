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
		"f64": 1234.56,
		"t2": ["a", 15, false, 3.14],
		"asdf4": 0.14159,
		"sf": {
			"v": [4, 5],
			"z": "hw2"
		},
		"unicode": "M\u00fcNSTER",
		"bool": true
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
	testFloat64 := obj.Get("f64").Float64()
	if testFloat64 != 1234.56 {
		t.Fatalf(`key "f64" (%f) should == 1234.56`, testFloat64)
	}
	testBool := obj.Get("bool").Bool()
	if !testBool {
		t.Fatalf(`key "bool" (%s) should == true`, testBool)
	}
	testArray := obj.Get("t2").Array()
	if len(testArray) != 4 {
		t.Fatalf(`length of key "t2" (%d) should == 4`, len(testArray))
	} else if testArray[0].String() != "a" {
		t.Fatalf(`first element of key "t2" (%s) should == "a"`, testArray[0])
	} else if testArray[1].Int64() != 15 {
		t.Fatalf(`second element of key "t2" (%d) should == 15`, testArray[1])
	} else if testArray[2].Bool() {
		t.Fatalf(`third element of key "t2" (%s) should == false`, testArray[2])
	} else if testArray[3].Float64() != 3.14 {
		t.Fatalf(`fourth element of key "t2" (%f) should == 3.14`, testArray[3])
	}

	fallback := []interface{}{"string", 1, true, 2.3}
	testArrayFallback := obj.Get("non-existant").Array(fallback)
	if len(testArrayFallback) != len(fallback) {
		t.Fatalf(`length of array (%d) should == %d`, len(testArrayFallback), len(fallback))
	} else if testArrayFallback[0].String() != "string" {
		t.Fatalf(`first element of array (%s) should == "string"`, testArrayFallback[0])
	} else if testArrayFallback[1].Int64() != 1 {
		t.Fatalf(`second element of array (%d) should == 1`, testArrayFallback[1])
	} else if !testArrayFallback[2].Bool() {
		t.Fatalf(`third element of array (%s) should == true`, testArrayFallback[2])
	} else if testArrayFallback[3].Float64() != 2.3 {
		t.Fatalf(`fourth element of array (%f) should == 2.3`, testArrayFallback[3])
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
