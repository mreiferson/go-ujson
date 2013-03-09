package ujson

import (
	"io/ioutil"
	"log"
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
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("run %d", i)
		Unmarshal(data)
	}
	b.SetBytes(int64(len(data)))
}
