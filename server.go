package gweb

import (
	"github.com/DaHuangQwQ/gweb/context"
	"github.com/DaHuangQwQ/gweb/middlewares"
	"github.com/DaHuangQwQ/gweb/types"
	"log"
	"net/http"
)

type Server interface {
	http.Handler

	Start(path string) error

	addRoute(method string, path string, handler types.HandleFunc, middleware ...middlewares.Middleware)
}

var _ Server = &HttpServer{}

type HttpServerOption func(*HttpServer)

type HttpServer struct {
	Router

	log func(msg string, args ...any)
}

func NewHttpServer(opts ...HttpServerOption) *HttpServer {
	s := &HttpServer{
		Router: newRouter(),
		log:    log.Printf,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &context.Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HttpServer) Start(path string) error {
	return http.ListenAndServe(path, h)
}

func (h *HttpServer) serve(ctx *context.Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	if info.n != nil {
		ctx.PathParams = info.pathParams
		ctx.MatchedRoute = info.n.route
	}

	var root types.HandleFunc = func(ctx *context.Context) {
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

	// flashResp 是最后一个步骤
	var m middlewares.Middleware = func(next types.HandleFunc) types.HandleFunc {
		return func(ctx *context.Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (h *HttpServer) Get(path string, hdl types.HandleFunc) {
	h.addRoute("GET", path, hdl)
}

func (h *HttpServer) Post(path string, hdl types.HandleFunc) {
	h.addRoute("POST", path, hdl)
}

func (h *HttpServer) Put(path string, hdl types.HandleFunc) {
	h.addRoute("PUT", path, hdl)
}

func (h *HttpServer) Delete(path string, hdl types.HandleFunc) {
	h.addRoute("DELETE", path, hdl)
}

func (h *HttpServer) Options(path string, hdl types.HandleFunc) {
	h.addRoute("OPTIONS", path, hdl)
}

func (h *HttpServer) Use(method string, path string, mdl ...middlewares.Middleware) {
	h.addRoute(method, path, nil, mdl...)
}

func (h *HttpServer) UseAll(path string, middleware ...middlewares.Middleware) {
	h.addRoute(http.MethodGet, path, nil, middleware...)
	h.addRoute(http.MethodPost, path, nil, middleware...)
	h.addRoute(http.MethodPut, path, nil, middleware...)
	h.addRoute(http.MethodDelete, path, nil, middleware...)
	h.addRoute(http.MethodOptions, path, nil, middleware...)
}

func (h *HttpServer) flashResp(ctx *context.Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("flash resp write error:", err)
	}
}
