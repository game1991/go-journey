package cache

import (
	"container/list"
	"sync"
	"time"
)

/*
lru淘汰策略就是:
如果当前key是被访问的key放到链表头部，
一直未被访问放到链表尾部，
如果key关联的过期时间大于当前时间过期，
则进行map删除该key，同时在链表中删除，
如果key对应的对象存放到链表的容量超过设置阈值，直接将链表尾部的数据丢弃
*/

/*
	目前key过期是被动过期方式，即当key被访问到时，判断是否过期，并且删除
	todo: 优化过期策略
*/

const (
	// MaxSIZE 最大存储容量
	MaxSIZE = 10000
	// DEFAULTEXPIRETIME 默认过期时间
	DEFAULTEXPIRETIME = 24 * time.Hour
)

// Entry 单个对象
type Entry struct {
	Key    string
	Data   interface{}
	Expire int64
}

// Cache ...
type Cache struct {
	maxEntry          int //链表存储对象最大容量
	defaultExpireTime time.Duration
	bucket            map[string]*list.Element
	lrulist           *list.List
	lock              *sync.RWMutex
	pool              *sync.Pool
}

// NewCache 生成缓存对象
func NewCache(opts ...Optional) *Cache {
	option := defaultOption()
	for _, opt := range opts {
		opt(option)
	}
	return &Cache{
		maxEntry:          option.MaxSize,
		defaultExpireTime: option.DefaultExpireTime,
		bucket:            make(map[string]*list.Element),
		lrulist:           list.New(),
		lock:              new(sync.RWMutex),
		pool: &sync.Pool{
			New: func() interface{} {
				return &Entry{}
			},
		},
	}
}

func (c *Cache) checkInit(key string) error {
	if c.bucket == nil || c.lrulist == nil || c.lock == nil || c.pool == nil {
		return ErrNotInit
	}

	if len(key) >= 1024 {
		return ErrTooLoogKey
	}

	return nil
}

// Set 设置kv和过期时间，lru算法
func (c *Cache) Set(key string, value interface{}, timestamp time.Duration) error {
	if err := c.checkInit(key); err != nil {
		return err
	}
	// 判断入参过期时间是否为0，如果小于等于0，则使用默认过期时间
	if timestamp <= 0 {
		timestamp = c.defaultExpireTime
	}
	//判断当前key是否存在，如果存在则进行更新
	if el := c.getListElement(key); el != nil {
		entry := c.getEntry(el)
		if entry != nil {
			c.lock.Lock()
			entry.Data = value
			c.resetExpireTime(entry, el, timestamp)
			c.lock.Unlock()
			return nil
		}
	}

	cacheElement := c.pool.Get()

	cacheEntry, ok := cacheElement.(*Entry)
	if !ok {
		cacheEntry = &Entry{}
	}
	cacheEntry.Key = key
	cacheEntry.Data = value
	cacheEntry.Expire = time.Now().Add(timestamp).Unix()
	c.lock.Lock()
	el := c.lrulist.PushFront(cacheEntry)
	c.bucket[key] = el
	c.lock.Unlock()

	// 淘汰策略，当超过容量的时候，进行链表尾部删除
	c.lru()
	return nil
}

// resetExpireTime 重置当前key的过期信息
func (c *Cache) resetExpireTime(entry *Entry, el *list.Element, timestamp time.Duration) {
	entry.Expire = time.Now().Add(timestamp).Unix()
	// 将当前数据放到链表前面
	c.lrulist.MoveToFront(el)
}

func (c *Cache) lru() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.lrulist.Len() > c.maxEntry {
		// 删除最后一个
		tmp := c.lrulist.Back()
		c.lrulist.Remove(tmp)
		item, ok := tmp.Value.(*Entry)
		if ok {
			// 从map中去除
			delete(c.bucket, item.Key)
			// 放入池中
			c.pool.Put(item)
		}

	}
}

// Get 获取key对应的value,如果key存在并且传参timestamp>0则认为重置key的时间
func (c *Cache) Get(key string, timestamp time.Duration) (interface{}, bool) {
	if err := c.checkInit(key); err != nil {
		return nil, false
	}
	el := c.getListElement(key)
	if el != nil {
		entry := c.getEntry(el)
		if entry != nil {
			// 如果查到对应的key，判断当前入参timestamp来决定是否重置key
			if timestamp > 0 {
				c.lock.Lock()
				c.resetExpireTime(entry, el, timestamp)
				value := entry.Data
				c.lock.Unlock()
				return value, true
			}
			// 过期情况
			if time.Now().Unix() > entry.Expire {
				// 过期进行删除key和链表数据移除
				c.lock.Lock()
				delete(c.bucket, key)
				c.lrulist.Remove(el)
				c.lock.Unlock()
				// 将entry对象放到池中
				c.pool.Put(entry)
				return nil, false
			}

			//没过期的情况下，由于key被访问，则将这个对象放到链表前面
			c.lock.Lock()
			c.lrulist.MoveToFront(el)
			value := entry.Data
			c.lock.Unlock()
			return value, true
		}
	}
	return nil, false
}

func (c *Cache) getEntry(el *list.Element) *Entry {
	if el != nil {
		c.lock.RLock()
		defer c.lock.RUnlock()
		if item, ok := el.Value.(*Entry); ok {
			return item
		}
	}
	return nil
}

func (c *Cache) getListElement(key string) *list.Element {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if el, ok := c.bucket[key]; ok {
		return el
	}
	return nil
}

// Delete 删除对应key
func (c *Cache) Delete(key string) error {
	if err := c.checkInit(key); err != nil {
		return err
	}
	c.lock.RLock()
	if el, ok := c.bucket[key]; ok {
		c.lock.RUnlock()

		c.lock.Lock()
		defer c.lock.Unlock()

		delete(c.bucket, key)
		c.lrulist.Remove(el)
		if item, ok := el.Value.(*Entry); ok {
			c.pool.Put(item)
		}
	}
	return nil
}

// Size 缓存区大小
func (c *Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.bucket == nil {
		return 0
	}
	return c.lrulist.Len()
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lrulist = list.New()
	c.bucket = make(map[string]*list.Element, c.maxEntry)
}
