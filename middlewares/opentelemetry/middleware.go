package opentelemetry

import (
	"github.com/DaHuangQwQ/gweb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/DaHuangQwQ/gweb/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
	return &MiddlewareBuilder{Tracer: tracer}
}

func (m *MiddlewareBuilder) Build() gweb.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next gweb.HandleFunc) gweb.HandleFunc {
		return func(ctx *gweb.Context) {
			reqCtx := ctx.Req.Context()

			// 和客户端的 trace 结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))

			ctx.Req = ctx.Req.WithContext(reqCtx)

			next(ctx)

			span.SetName(ctx.MatchedRoute)
		}
	}
}
