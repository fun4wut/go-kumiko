package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type FnHash func([]byte) uint32

// Dict 一致性哈希算法 https://geektutu.com/post/geecache-day4.html
type Dict struct {
	fnHash   FnHash         // 哈希函数
	replicas int            // 虚拟节点的扩展因子
	keys     []int          // 哈希环
	hashDict map[int]string // vnode的hash -> 真实node的name
}

func New(replicas int, fnHash FnHash) *Dict {
	d := &Dict{
		fnHash:   fnHash,
		replicas: replicas,
		keys:     []int{},
		hashDict: make(map[int]string),
	}
	if d.fnHash == nil {
		d.fnHash = crc32.ChecksumIEEE
	}
	return d
}

// Add 往哈希环添加节点
func (d *Dict) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < d.replicas; i++ {
			// key{i}
			hashRes := (int)(d.fnHash([]byte(key + strconv.Itoa(i))))
			d.keys = append(d.keys, hashRes)
			d.hashDict[hashRes] = key
		}
	}
	sort.Ints(d.keys)
}

func (d *Dict) Remove(keys ...string) {
	for _, key := range keys {
		for i := 0; i < d.replicas; i++ {
			hashRes := (int)(d.fnHash([]byte(key + strconv.Itoa(i))))
			idx := sort.SearchInts(d.keys, hashRes)
			d.keys = append(d.keys[:idx], d.keys[idx+1:]...)
			delete(d.hashDict, hashRes)
		}
	}
}

func (d *Dict) Get(key string) string {
	if len(d.keys) == 0 {
		return ""
	}
	hashRes := d.fnHash([]byte(key))
	idx := sort.Search(len(d.keys), func(i int) bool {
		// 找到最靠近的一个节点
		return d.keys[i] >= (int)(hashRes)
	})
	// 如果search没找到，会返回n，取个模即可做到环的效果
	return d.hashDict[d.keys[idx%len(d.keys)]]
}
