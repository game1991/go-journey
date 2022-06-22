package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestReq(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "http://localhost:9091/timeout", nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	client := http.DefaultClient
	client.Timeout = 30 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		t.Logf("请求失败[status]:%d[error]:%v\n", resp.StatusCode, string(bts))
		t.FailNow()
	}

	//json.NewDecoder(resp.Body).Decode()
	io.Copy(os.Stdout, resp.Body)
}
