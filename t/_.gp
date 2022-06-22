package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "abc"
	//a := "a"
	//b := "b"
	//c := "c"

	str := strings.Split(s, "|")
	fmt.Println("str", str)
}
