package gweb

import "net/http"

type HandleFunc func(*Context)

type Server interface {
	http.Handler

	Start(path string) error

	addRoute(method string, path string, handler HandleFunc, middleware ...Middleware)
}

var _ Server = &HttpServer{}

type HttpServerOption func(*HttpServer)

type HttpServer struct {
	Router
}

func NewHttpServer(opts ...HttpServerOption) *HttpServer {
	s := &HttpServer{
		Router: newRouter(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HttpServer) Start(path string) error {
	return http.ListenAndServe(path, h)
}

func (h *HttpServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	if info.n != nil {
		ctx.PathParams = info.pathParams
		ctx.MatchedRoute = info.n.route
	}

	var root HandleFunc = func(ctx *Context) {
		if !ok || info.n == nil || info.n.handler == nil {
			ctx.RespStatusCode = 404
			return
		}
		info.n.handler(ctx)
	}

	// 从后往前组装
	for i := len(info.mdls) - 1; i >= 0; i-- {
		root = info.mdls[i](root)
	}

	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 flashResp 是最后一个步骤
	//var m Middleware = func(next HandleFunc) HandleFunc {
	//	return func(ctx *Context) {
	//		next(ctx)
	//		h.flashResp(ctx)
	//	}
	//}
	//root = m(root)
	root(ctx)
}

func (h *HttpServer) Get(path string, hdl HandleFunc) {
	h.addRoute("GET", path, hdl)
}

func (h *HttpServer) Post(path string, hdl HandleFunc) {
	h.addRoute("POST", path, hdl)
}

func (h *HttpServer) Put(path string, hdl HandleFunc) {
	h.addRoute("PUT", path, hdl)
}

func (h *HttpServer) Delete(path string, hdl HandleFunc) {
	h.addRoute("DELETE", path, hdl)
}

func (h *HttpServer) Options(path string, hdl HandleFunc) {
	h.addRoute("OPTIONS", path, hdl)
}

func (h *HttpServer) Use(method string, path string, mdl ...Middleware) {
	h.addRoute(method, path, nil, mdl...)
}

func (h *HttpServer) UseAll(path string, middleware ...Middleware) {
	h.addRoute(http.MethodGet, path, nil, middleware...)
	h.addRoute(http.MethodPost, path, nil, middleware...)
	h.addRoute(http.MethodPut, path, nil, middleware...)
	h.addRoute(http.MethodDelete, path, nil, middleware...)
	h.addRoute(http.MethodOptions, path, nil, middleware...)
}
