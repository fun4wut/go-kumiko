package kotoha

import (
	"awesomeProject/kotoha/hash"
	"awesomeProject/kumiko"
	"log"
	"net/http"
	"sync"
)

const defaultBasePath = "/__kotoha"

// HttpPool http 服务端，处理http请求，如果本地处理不了，则再利用 httpGetter 转给其他远程服务端处理
type HttpPool struct {
	Self        string
	BasePath    string
	mu          sync.Mutex
	peers       *hash.Dict
	httpGetters map[string]*httpGetter // url -> getter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		Self:     self,
		BasePath: defaultBasePath,
	}
}

// AddPeer 添加远程节点
func (p *HttpPool) AddPeer(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = hash.New(3, nil)
	p.peers.Add(peers...)
	p.httpGetters = map[string]*httpGetter{}
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseUrl: peer + p.BasePath,
		}
	}
}

// PickPeer 给定一个key，然后选取一个节点
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 如果 peer == self，那就如密传如密了，直接降级到本地callback
	if peer := p.peers.Get(key); peer != "" && peer != p.Self {
		log.Printf("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HttpPool)(nil)

// HandleGet 处理过来的http请求
func HandleGet() kumiko.HandlerFn {
	return func(ctx *kumiko.Ctx) {
		groupName, _ := ctx.GetPathParam("groupname")
		key, _ := ctx.GetPathParam("key")
		group := GetGroup(groupName)
		v, err := group.Get(key)
		if err != nil {
			ctx.WriteText(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
		ctx.Writer.Write(v.ByteSlice())
	}
}
