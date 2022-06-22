package main

import (
	"fmt"
	"os"
	"sort"
)

// Order 排序规则
type Order int

// 排序方式
const (
	ASCII     Order = iota + 1 //ascii码正序
	ASCIIDesc                  // ascii倒序
)

func main() {
	mp := make(map[string]interface{})
	mp["aAx"] = 1
	mp["AAx"] = 2
	mp["Aax"] = 3
	mp["aax"] = 4
	mp["aaa"] = 5

	keys := make([]string, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}

	fmt.Println("排序前", keys)

	var order Order
LOOP:
	fmt.Fprintln(os.Stdout, "请选择排序方式：0、退出游戏；1、ASCII码；2、ASCII码倒序")
	fmt.Scanf("%d", &order)

	switch order {
	case 0:
		os.Exit(0)
	case ASCII:
		Asc(keys)
	case ASCIIDesc:
		Desc(keys)
	default:
		fmt.Fprintln(os.Stderr, "输入有误，请重新输入")
		goto LOOP
	}

	fmt.Println("排序后", keys)
}

// Desc ...
func Desc(keys []string) {
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
}

// Asc ...
func Asc(keys []string) {
	sort.Strings(keys)
}
