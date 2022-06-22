package cache

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var cacheObj *Cache

func TestMain(t *testing.M) {
	cacheObj = NewCache(WithMaxSize(5))

	code := t.Run()
	os.Exit(code)
}

func TestCache_Set(t *testing.T) {
	type args struct {
		key       string
		value     interface{}
		timestamp time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "测试存储key123",
			args: args{
				key:       "123",
				value:     123,
				timestamp: 5 * time.Second,
			},
		},
		{
			name: "测试存储key456",
			args: args{
				key:       "456",
				value:     456,
				timestamp: 5 * time.Second,
			},
		},
		{
			name: "测试存储key789",
			args: args{
				key:       "789",
				value:     789,
				timestamp: 5 * time.Second,
			},
		},
		{
			name: "测试存储key789改内容",
			args: args{
				key:       "789",
				value:     1000,
				timestamp: 15 * time.Second,
			},
		},
		{
			name: "测试存储key111",
			args: args{
				key:       "111",
				value:     111,
				timestamp: 5 * time.Second,
			},
		},
		{
			name: "测试存储key119",
			args: args{
				key:       "119",
				value:     119,
				timestamp: 5 * time.Second,
			},
		},
		{
			name: "测试存储key120",
			args: args{
				key:       "120",
				value:     120,
				timestamp: 5 * time.Second,
			},
		},
	}
	c := cacheObj
	for _, tt := range tests {
		if err := c.Set(tt.args.key, tt.args.value, tt.args.timestamp); (err != nil) != tt.wantErr {
			t.Errorf("Cache.Set() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
	tests1 := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
		{
			name: "测试获取",
			args: args{
				key: "119",
			},
			want:  119,
			want1: true,
		},
		{
			name: "测试重置过期时间",
			args: args{
				key:       "120",
				timestamp: 10 * time.Second,
			},
			want:  120,
			want1: true,
		},
		{
			name: "测试key已经被lru了",
			args: args{
				key: "123",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			got, got1 := c.Get(tt.args.key, tt.args.timestamp)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	tests3 := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{
			name: "测试容量",
			want: 5,
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			if got := c.Size(); got != tt.want {
				t.Errorf("Cache.Size() = %v, want %v", got, tt.want)
			}
		})
	}

	// 测试清空缓存
	c.Clear()
	v, has := c.Get("789", 0)
	if has {
		t.Logf("%v\n", v)
		t.FailNow()
	}
}

func TestCache_lru(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "lru 策略",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			c.lru()
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key       string
		timestamp time.Duration
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
		{
			name: "测试获取",
			args: args{
				key: "123",
			},
			want:  123,
			want1: false,
		},
		{
			name: "测试重置过期时间",
			args: args{
				key:       "123",
				timestamp: 10 * time.Second,
			},
			want:  123,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			got, got1 := c.Get(tt.args.key, tt.args.timestamp)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_Size(t *testing.T) {

	tests3 := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{
			name: "测试容量",
			want: 3,
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			if got := c.Size(); got != tt.want {
				t.Errorf("Cache.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Clear(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "清空缓存",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cacheObj
			c.Clear()
		})
	}
}
