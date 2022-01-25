package lru

import "container/list"

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  int
	nBytes    int
	ll        *list.List
	cache     map[string]*list.Element
	onRemoved func(key string, value Value)
}

func New() *Cache {
	return &Cache{
		maxBytes:  1 << 10,
		nBytes:    0,
		ll:        list.New(),
		cache:     map[string]*list.Element{},
		onRemoved: nil,
	}
}

func (c *Cache) Get(key string) (val Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Front()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= len(kv.key) + kv.value.Len()
		if c.onRemoved != nil {
			c.onRemoved(kv.key, kv.value)
		}
	}
}

func (c *Cache) Set(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		c.nBytes += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		e := &entry{
			key:   key,
			value: value,
		}
		ele := c.ll.PushBack(e)
		c.cache[key] = ele
		c.nBytes += value.Len() + len(key)
	}
	for c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}
