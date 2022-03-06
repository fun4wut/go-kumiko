package kotoha

import (
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 接口型函数 https://geektutu.com/post/7days-golang-q1.html
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。
type Group struct {
	name       string
	getter     Getter     // 拿不到缓存值，最终的获取值的callback
	mainCache  cache      // 本地缓存
	peerPicker PeerPicker // 获取正确的远程节点
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeerPicker(picker PeerPicker) {
	g.peerPicker = picker
}

func AddGroup(name string, cacheBytes int, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get
/** 1. 先读本地缓存
2. 再看远程缓存
3. 最后执行本地callback
*/
func (g *Group) Get(key string) (ByteView, error) {
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit!")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	if g.peerPicker != nil {
		if peer, ok := g.peerPicker.PickPeer(key); ok {
			if val, err := g.loadFromRemote(peer, key); err == nil {
				return val, nil
			}
		}
	}
	return g.loadFromLocal(key)
}

func (g *Group) loadFromRemote(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	log.Println("[Remote] hit!")
	return ByteView{b: bytes}, nil
}

func (g *Group) loadFromLocal(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.mainCache.add(key, value)
	log.Println("[Local] hit!")
	return value, nil
}
