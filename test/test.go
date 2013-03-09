package main

import (
	"../../go-ujson"
	"log"
)

func main() {
	b := []byte(`{ "test": "hello world", "t2": ["a", 3, "c"], "asdf4": 0.14159, "sf": { "v": [4, 5], "z": "hw2" } }`)
	v, err := ujson.Unmarshal(b)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("%#v", v)
	s1, _ := v.Get("test").String()
	log.Printf("%s", s1)
	s2, _ := v.Get("sf").Get("z").String()
	log.Printf("%s", s2)
}
