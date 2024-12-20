package opentelemetry

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/DaHuangQwQ/gweb/context"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := &MiddlewareBuilder{
		Tracer: tracer,
	}

	server := gweb.NewHttpServer()
	server.UseAll("/*", builder.Build())

	server.Get("/*", func(ctx *context.Context) {
		_, span := tracer.Start(ctx.Req.Context(), "test")
		defer span.End()
	})

	server.Start(":8081")
}
