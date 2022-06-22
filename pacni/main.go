package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

func main() {
	str := " aa b c "

	fmt.Println(strings.TrimSpace(str))

	var b []string

	b = nil

	fmt.Println(len(b))

	fmt.Println("woshi main")
	time.Sleep(time.Second * 5)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
				debug.PrintStack()
			}
		}()

		panic("我是未知错误")
	}()
	select {
	case <-time.After(time.Second * 10):
		fmt.Println("我完成了")
	}
}
