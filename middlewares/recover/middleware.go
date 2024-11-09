package recover

import "github.com/DaHuangQwQ/gweb"

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	Log        func(ctx *gweb.Context)
}

func NewMiddlewareBuilder(statusCode int, data []byte, log func(ctx *gweb.Context)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		StatusCode: statusCode,
		Data:       data,
		Log:        log,
	}
}

func (m *MiddlewareBuilder) Build() gweb.Middleware {
	return func(next gweb.HandleFunc) gweb.HandleFunc {
		return func(ctx *gweb.Context) {
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
