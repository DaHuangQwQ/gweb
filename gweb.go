package gweb

import (
	"github.com/DaHuangQwQ/gweb/middlewares/accesslog"
	"github.com/DaHuangQwQ/gweb/middlewares/errhdl"
	"github.com/DaHuangQwQ/gweb/middlewares/opentelemetry"
	"github.com/DaHuangQwQ/gweb/middlewares/prometheus"
	grecover "github.com/DaHuangQwQ/gweb/middlewares/recover"
	"go.opentelemetry.io/otel"
	"net/http"
)

const instrumentationName = "github.com/DaHuangQwQ/gweb"

func Default(opts ...HttpServerOption) Server {
	server := NewHttpServer(opts...)

	server.UseAll("/*", accesslog.NewMiddlewareBuilder().Build())

	server.UseAll("/*", errhdl.NewMiddlewareBuilder().Build())

	tracer := otel.GetTracerProvider().Tracer(instrumentationName)

	server.UseAll("/*", opentelemetry.NewMiddlewareBuilder(tracer).Build())

	server.UseAll("/*", prometheus.NewMiddlewareBuilder("github.com/DaHuangQwQ/gweb", "web", "http_request", "speed").Build())

	server.UseAll("/*", grecover.NewMiddlewareBuilder(http.StatusBadRequest, []byte("bad req"), func(ctx *Context) {
		println(ctx.RespStatusCode)
	}).Build())

	return server
}
