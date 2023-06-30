package gos

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (m *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	m.router[key] = handler
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
	key := req.Method + "-" + req.URL.Path
	if handler, ok := m.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}
