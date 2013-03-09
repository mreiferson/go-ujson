package ujson

import (
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkUjson(b *testing.B) {
	b.StopTimer()
	f, err := os.Open("testdata.json")
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
		Unmarshal(data)
	}
	b.SetBytes(int64(len(data)))
}
