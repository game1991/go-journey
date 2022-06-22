package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// FOO ...
type FOO struct {
	Name    string
	Age     int
	UUID    string
	TimeInt int
	Bool    []bool
	Int     []int
	String  []string
	// Struct  []*FOO
	Inter []interface{}
	Map   map[string]interface{}
}

type foo struct {
	Name string
}

func main() {
	//JSONPrt()
	m := make(map[string]interface{})
	m["1"] = 1
	m["2"] = "2"
	m["3"] = true
	m["4"] = &foo{"4"}
	m["5"] = []map[string]interface{}{{"5": 1, "6": true}, {"11": 11, "10": "10"}}

	s := fmt.Sprintf("%v\n", NewFoo(m))
	os.Stdout.Write([]byte(s))
}

// JSONPrt ...
func JSONPrt() {
	foo := NewFoo(nil)

	bts, err := json.Marshal(foo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bts))
}

// NewFoo ...
func NewFoo(m map[string]interface{}) *FOO {
	return &FOO{
		Name:    "gg",
		Age:     22,
		UUID:    "12sasa12sdsad",
		TimeInt: int(time.Now().UTC().Unix()),
		Bool:    []bool{true, false},
		Map:     m,
	}
}
