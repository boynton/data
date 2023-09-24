package data

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Foo struct {
	X int    `json:"x,omitempty"`
	Y string `json:"y,omitempty"`
}

func TestMap(test *testing.T) {
	m := NewMap[int]()
	m.Put("twentythree", 23)
	fmt.Println("m:", Pretty(m))
	m2 := NewMap[*Foo]()
	m2.Put("f1", &Foo{X: 23, Y: "hi"})
	m2.Put("f2", &Foo{X: 57, Y: "hello"})
	fmt.Println("m2:", Pretty(m2))
	s2 := Pretty(m2)
	var m2a Map[*Foo]
	err := json.Unmarshal([]byte(s2), &m2a)
	fmt.Println("m2a:", Pretty(m2a), err)
	for _, k := range m2a.Keys() {
		fmt.Println("k:", k)
		f := m2a.Get(k)
		fmt.Println(" ->", Pretty(f))
	}
	/*
		err := decode(jsonData, &m)
		if err != nil {
			test.Errorf("%v", err)
		} else {
			fmt.Println(Pretty(ts))
		}
	*/
}

func decode(jsonData string, expected interface{}) error {
	return json.Unmarshal([]byte(jsonData), &expected)
}
