package main

import (
	"../../go-ujson"
	"log"
)

func main() {
	b := []byte(`{ "test": 1, "t2": ["a", 3, "c"], "asdf4": 0.14159 }`)
	v, err := ujson.Unmarshal(b)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("%#v", v)
}
