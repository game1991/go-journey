package cache

import "time"

// Optional ...
type Optional func(*Option)

// Option ...
type Option struct {
	MaxSize           int           // 最大存放对象个数
	DefaultExpireTime time.Duration // 默认过期时间
}

func defaultOption() *Option {
	return &Option{
		MaxSize:           MaxSIZE,
		DefaultExpireTime: DEFAULTEXPIRETIME,
	}
}

// WithMaxSize ... 最大存放缓存对象个数
func WithMaxSize(maxSize int) Optional {
	return func(o *Option) {
		if maxSize != 0 {
			o.MaxSize = maxSize
		}
	}
}

// WithDefaultExpireTime ... 默认key过期时间
func WithDefaultExpireTime(defaultExpireTime time.Duration) Optional {
	return func(o *Option) {
		if defaultExpireTime != 0 {
			o.DefaultExpireTime = defaultExpireTime
		}
	}
}
