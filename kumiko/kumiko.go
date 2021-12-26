package kumiko

import (
	"net/http"
)

type HandlerFn func(ctx *Ctx)

type Kumiko struct {
	*router
}

func New() *Kumiko {
	return &Kumiko{
		router: newRouter(),
	}
}

func (k *Kumiko) Get(path string, handler HandlerFn) {
	k.addRoute("GET", path, handler)
}

func (k *Kumiko) Post(path string, handler HandlerFn) {
	k.addRoute("POST", path, handler)
}

func (k *Kumiko) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 每次请求都构造一个上下文给handler使用
	ctx := newCtx(writer, request)
	k.handle(ctx)
}

func (k *Kumiko) Run(addr string) error {
	return http.ListenAndServe(addr, k)
}
