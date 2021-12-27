package kumiko

import (
	"net/http"
)

type HandlerFn func(ctx *Ctx)

type RouterGroup struct {
	prefix string
	top    *Engine
}

type Engine struct {
	router       *router
	*RouterGroup             // 全局对象本身也有路由组能力
	middlewares  []HandlerFn // 这里做个简化，中间件对所有路由都生效
}

func New() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{top: e}
	// 使用上错误恢复的中间件
	e.Use(Recovery())
	return e
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	top := g.top
	return &RouterGroup{
		prefix: g.prefix + prefix,
		top:    top,
	}
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

func (e *Engine) Use(fns ...HandlerFn) {
	e.middlewares = append(e.middlewares, fns...)
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 每次请求都构造一个上下文给handler使用
	ctx := newCtx(writer, request)
	ctx.handlers = e.middlewares
	e.router.handle(ctx)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
