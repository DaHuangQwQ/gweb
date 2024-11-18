package gweb

import (
	"github.com/DaHuangQwQ/gweb/context"
	"github.com/DaHuangQwQ/gweb/middlewares/accesslog"
	"github.com/DaHuangQwQ/gweb/middlewares/errhdl"
	"github.com/DaHuangQwQ/gweb/middlewares/prometheus"
	grecover "github.com/DaHuangQwQ/gweb/middlewares/recover"
	"net/http"
)

const instrumentationName = "github.com/DaHuangQwQ/gweb"

func Default(opts ...HttpServerOption) *HttpServer {
	server := NewHttpServer(opts...)

	server.UseAll("/*", accesslog.NewMiddlewareBuilder().Build())

	server.UseAll("/*", errhdl.NewMiddlewareBuilder().Build())

	//tracer := otel.GetTracerProvider().Tracer(instrumentationName)

	server.UseAll("/*", prometheus.NewMiddlewareBuilder("_DaHuangQwQ_gweb", "web", "http_request", "speed").Build())

	server.UseAll("/*", grecover.NewMiddlewareBuilder(http.StatusBadRequest, []byte("bad req"), func(ctx *context.Context) {
		println(ctx.RespStatusCode)
	}).Build())

	return server
}
