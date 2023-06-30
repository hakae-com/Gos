package gos

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc
	index      int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (m *Context) Next() {
	m.index++
	s := len(m.handlers)
	for ; m.index < s; m.index++ {
		m.handlers[m.index](m)
	}
}

func (m *Context) Param(key string) string {
	value, _ := m.Params[key]
	return value
}

func (m *Context) PostForm(key string) string {
	return m.Req.FormValue(key)
}

func (m *Context) Query(key string) string {
	return m.Req.URL.Query().Get(key)
}

func (m *Context) Status(code int) {
	m.StatusCode = code
	m.Writer.WriteHeader(code)
}

func (m *Context) SetHeader(key string, value string) {
	m.Writer.Header().Set(key, value)
}

func (m *Context) String(code int, format string, values ...interface{}) {
	m.SetHeader("Content-Type", "text/plain")
	m.Status(code)
	_, err := m.Writer.Write([]byte(fmt.Sprintf(format, values)))
	if err != nil {
		return
	}
}

func (m *Context) JSON(code int, obj interface{}) {
	m.SetHeader("Content-Type", "application/json")
	m.Status(code)
	encoder := json.NewEncoder(m.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(m.Writer, err.Error(), 500)
	}
}

func (m *Context) Data(code int, data []byte) {
	m.Status(code)
	m.Writer.Write(data)
}

func (m *Context) HTML(code int, html string) {
	m.SetHeader("Content-Type", "text/html")
	m.Status(code)
	m.Writer.Write([]byte(html))
}
