package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	api := &API{
		Name: "foo",
		Age:  18,
		Foo:  false,
		Msg:  "",
	}

	bts, err := json.Marshal(api)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("api序列化后", string(bts))

	rv := reflect.ValueOf(api)

	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Name
		// 找出tag对应的值
		tagValue, ok := StructTag(rt.Field(i).Tag).Lookup("json")
		if !ok {
			// 如果找不到对应的tag则跳过
			continue
		}

		split := strings.Split(tagValue, ",")

		for _, item := range split {
			if item == "omitempty" {

			}
		}

		value := rv.Field(i)
		if !value.CanInterface() {
			fmt.Println("当前对应的值类型不支持")
		}

		fmt.Printf("[key]:%v,[tagVaule]:%v,[v]:%v,[kind]:%v\n", key, tagValue, value.Interface(), value.Kind())

	}

	s := "a,c"
	index := strings.Index(s, ",")
	fmt.Printf("%#v,%#v\n", index, s[:index])
}

// API ...
type API struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Foo  bool   `json:"foo"`
	Msg  string `json:"msg,omitempty"`
}

// StructTag ...
type StructTag string

// Lookup 找出tag对应的value
func (tag StructTag) Lookup(key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
