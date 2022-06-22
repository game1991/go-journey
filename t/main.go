package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	ip := "127.0.0.1:6060"

	for i := 0; i < (1 << 20); i++ {
		go timeFunc()
	}

	if err := http.ListenAndServe(ip, nil); err != nil {
		panic(err)
	}
}

func timeFunc() {
	time.AfterFunc(time.Hour, func() {
		fmt.Println("测试OOM")
	})
}
