package byte

import (
	"encoding/json"
	"fmt"
)

type Foo struct {
	Name string
	Fool bool
	Age  int
	// Bug  *Bug
}

type Bug struct {
	Msg string
}

func Mashal(b []byte) error {
	v := &Foo{}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	fmt.Printf("%#v\n", v)
	return nil
}
