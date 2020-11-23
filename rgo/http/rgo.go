package http

import (
	"html/template"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type RgoHttp struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// New is the constructor
func New() *RgoHttp {
	rgoHttp := &RgoHttp{router: newRouter()}
	rgoHttp.RouterGroup = &RouterGroup{rgoHttp: rgoHttp}
	rgoHttp.groups = []*RouterGroup{rgoHttp.RouterGroup}
	return rgoHttp
}

// Default use Logger() & Recovery middlewares
func Default() *RgoHttp {
	rgoHttp := New()
	rgoHttp.Use(Logger(), Recovery())
	return rgoHttp
}

// for custom render function
func (this *RgoHttp) SetFuncMap(funcMap template.FuncMap) {
	this.funcMap = funcMap
}

func (this *RgoHttp) LoadHTMLGlob(pattern string) {
	this.htmlTemplates = template.Must(template.New("").Funcs(this.funcMap).ParseGlob(pattern))
}

// Run defines the method to start a http server
func (this *RgoHttp) Run(addr string) (err error) {
	return http.ListenAndServe(addr, this)
}

func (this *RgoHttp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range this.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.rgoHttp = this
	this.router.handle(c)
}
