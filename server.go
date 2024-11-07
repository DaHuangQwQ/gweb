package gweb

import "net/http"

type HandleFunc func(*Context)

type Server interface {
	http.Handler

	Start(path string) error

	addRoute(method string, path string, handler HandleFunc, middleware ...Middleware)
}

var _ Server = &HttpServer{}

type HttpServer struct {
	router
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
	if !ok || info == nil {
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("404 page not found"))
		return
	}
	ctx.PathParams = info.pathParams
}
