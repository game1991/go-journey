package main

import (
	"container/list"
	"fmt"
)

func main() {

	l := list.New()

	l.PushFront(10)

	l.PushFront(9)

	l.PushFront(8)

	l.PushFront(7)

	if l.Len() > 3 {
		tmp := l.Back()
		l.Remove(tmp)

		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)
		}
	}

	l.Remove(l.Back())
}
