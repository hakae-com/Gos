package gos

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (m *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	m.router.addRoute(method, pattern, handler)
}

func (m *Engine) GET(pattern string, handler HandlerFunc) {
	m.addRoute("GET", pattern, handler)
}

func (m *Engine) POST(pattern string, handler HandlerFunc) {
	m.addRoute("POST", pattern, handler)
}

func (m *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, m)
}

func (m *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	m.router.handle(c)
}
