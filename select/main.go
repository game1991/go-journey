package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	//ch1 := make(chan int)

	go func() {
		var i = 1
		for {
			i++
			time.Sleep(time.Second * 1)
			ch <- i
			//time.Sleep(time.Second)
		}
	}()

	for i := 0; ; i++ {
		select {
		case x := <-ch:
			go pL(x)
			// case <-timeAfter():
			// 	println(time.Now().Unix())
			// }
		}
		fmt.Println(i)
	}
}

func pL(x int) {
	time.Sleep(time.Second * 15)
	println(x)
}

func timeAfter() <-chan time.Time {
	t := time.After(30 * time.Second)
	fmt.Println("通道2", time.Now())
	return t
}
