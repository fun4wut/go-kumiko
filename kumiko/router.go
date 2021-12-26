package kumiko

import (
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFn
	trieRoot *trieNode
}

func newRouter() *router {
	return &router{
		handlers: map[string]HandlerFn{},
		trieRoot: &trieNode{
			pattern:   "",
			part:      "",
			children:  map[string]*trieNode{},
			wildNodes: []*trieNode{},
			isWild:    false,
		},
	}
}

func genRouterKey(method string, path string) string {
	return "/" + method + path
}

func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
}

func (r *router) addRoute(method string, path string, handler HandlerFn) {
	routerKey := genRouterKey(method, path)
	r.trieRoot.insert(splitPath(routerKey), 0)
	r.handlers[routerKey] = handler
}

func (r *router) getRoute(method string, path string) (*trieNode, map[string]string) {
	routerKey := genRouterKey(method, path)
	actualParts := splitPath(routerKey)
	pathParams := make(map[string]string)

	// 处理路径参数，把对应的参数填上去
	if node := r.trieRoot.find(actualParts, 0); node != nil {
		searchedParts := splitPath(node.pattern)
		for i, part := range searchedParts {
			switch part[0] {
			case ':':
				pathParams[part[1:]] = actualParts[i]
			case '*':
				pathParams[part[1:]] = strings.Join(actualParts[i:], "/")
			}
		}
		return node, pathParams
	}
	return nil, nil
}

func (r *router) handle(ctx *Ctx) {
	node, pathParam := r.getRoute(ctx.Method, ctx.Path)
	if node == nil {
		ctx.WriteText(http.StatusNotFound, "404 NOT FOUND")
		return
	}
	if fn, ok := r.handlers[node.pattern]; ok {
		ctx.pathParam = pathParam
		fn(ctx)
	} else {
		ctx.WriteText(http.StatusNotFound, "404 NOT FOUND")
	}
}
