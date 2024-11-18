package recover

import (
	"github.com/DaHuangQwQ/gweb/context"
	"github.com/DaHuangQwQ/gweb/middlewares"
	"github.com/DaHuangQwQ/gweb/types"
)

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	Log        func(ctx *context.Context)
}

func NewMiddlewareBuilder(statusCode int, data []byte, log func(ctx *context.Context)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		StatusCode: statusCode,
		Data:       data,
		Log:        log,
	}
}

func (m *MiddlewareBuilder) Build() middlewares.Middleware {
	return func(next types.HandleFunc) types.HandleFunc {
		return func(ctx *context.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.Data
					m.Log(ctx)
				}
			}()
			next(ctx)
		}
	}
}
