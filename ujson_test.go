package ujson

import (
	"encoding/json"
	"io/ioutil"
	// "log"
	"os"
	"testing"
)

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
	dec := NewDecoder(&jsDecStore{}, data)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// log.Printf("run %d", i)
		dec.Decode()
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
		json.Unmarshal(data, &m)
	}
	b.SetBytes(int64(len(data)))
}
