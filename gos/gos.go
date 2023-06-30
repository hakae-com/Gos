package gos

import (
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (m *RouterGroup) Group(prefix string) *RouterGroup {
	engine := m.engine
	newGroup := &RouterGroup{
		prefix: m.prefix + prefix,
		parent: m,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (m *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := m.prefix + comp
	m.engine.router.addRoute(method, pattern, handler)
}

func (m *RouterGroup) GET(pattern string, handler HandlerFunc) {
	m.addRoute("GET", pattern, handler)
}

func (m *RouterGroup) POST(pattern string, handler HandlerFunc) {
	m.addRoute("POST", pattern, handler)
}

func (m *RouterGroup) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, m)
}

func (m *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	m.engine.router.handle(c)
}
