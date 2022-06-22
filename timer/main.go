package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

const (
	// Layout 时间格式
	Layout = "2006/01/02 15:04:05.000"
)

func main() {
	ip := "127.0.0.1:6061"
	fmt.Printf("程序启动了%v,ip:%v\n", time.Now().Format(Layout), ip)
	tw := newTimeWheel()

	for i := 0; i < 100; i++ {
		time.Sleep(300 * time.Millisecond)
		go tw.Set(3 * time.Second)
	}

	if err := http.ListenAndServe("127.0.0.1:6061", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// TimerWheel 定时器控制器
type TimerWheel struct {
	tm     *time.Timer
	stop   chan struct{}
	expire int64
}

func newTimeWheel() *TimerWheel {
	return &TimerWheel{
		tm:   time.NewTimer(0),
		stop: make(chan struct{}),
	}
}

// Set ...
func (tw *TimerWheel) Set(expire time.Duration) {

	tw.expire = time.Now().Add(expire).Unix()
	tw.tm.Reset(expire)
	go tw.clean()
}

// Clean ...
func (tw *TimerWheel) Clean() {
}

func (tw *TimerWheel) clean() {
	select {
	case t := <-tw.tm.C:
		fmt.Println("清理过期key", time.Now().Format(Layout), t.Format(Layout))
	}
}
