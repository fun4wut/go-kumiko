package kumiko

import (
	"net/http"
)

type HandlerFn func(ctx *Ctx)

type RouterGroup struct {
	middlewares []HandlerFn
	prefix      string
	parent      *RouterGroup
	top         *Kumiko
}

type Kumiko struct {
	router       *router
	*RouterGroup // 全局对象本身也有路由组能力
	groups       []*RouterGroup
}

func New() *Kumiko {
	k := &Kumiko{
		router: newRouter(),
	}
	k.RouterGroup = &RouterGroup{top: k}
	k.groups = []*RouterGroup{k.RouterGroup}
	return k
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	top := g.top
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		top:    top,
	}
	top.groups = append(top.groups, newGroup)
	return newGroup
}

// 在路由组上去定义路由方法，这样全局对象本身也可以享受
func (g *RouterGroup) addRoute(method, path string, handler HandlerFn) {
	g.top.router.addRoute(method, g.prefix+path, handler)
}

func (g *RouterGroup) Get(path string, handler HandlerFn) {
	g.addRoute("GET", path, handler)
}

func (g *RouterGroup) Post(path string, handler HandlerFn) {
	g.addRoute("POST", path, handler)
}

func (k *Kumiko) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 每次请求都构造一个上下文给handler使用
	ctx := newCtx(writer, request)
	k.router.handle(ctx)
}

func (k *Kumiko) Run(addr string) error {
	return http.ListenAndServe(addr, k)
}
