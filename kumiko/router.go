package kumiko

import "net/http"

type router struct {
	handlers map[string]HandlerFn
}

func newRouter() *router {
	return &router{handlers: map[string]HandlerFn{}}
}

func genRouterKey(method string, path string) string {
	return method + "^^" + path
}

func (r *router) addRoute(method string, path string, handler HandlerFn) {
	r.handlers[genRouterKey(method, path)] = handler
}

func (r *router) handle(ctx *Ctx) {
	key := genRouterKey(ctx.Method, ctx.Path)
	if fn, ok := r.handlers[key]; ok {
		fn(ctx)
	} else {
		ctx.WriteText(http.StatusNotFound, "404 NOT FOUND")
	}
}
