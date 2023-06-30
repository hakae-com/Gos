package gos

import (
	"html/template"
	"net/http"
	"path"
	"strings"
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
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
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

func (m *RouterGroup) User(middlewares ...HandlerFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
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

func (m *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(m.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (m *RouterGroup) Static(relativePath string, root string) {
	handle := m.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	m.GET(urlPattern, handle)
}

func (m *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, m)
}

func (m *Engine) SetFuncMap(funcMap template.FuncMap) {
	m.funcMap = funcMap
}

func (m *Engine) LoadHTMLGlob(pattern string) {
	m.htmlTemplates = template.Must(template.New("").Funcs(m.funcMap).ParseGlob(pattern))
}

func (m *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range m.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = m
	m.router.handle(c)
}
