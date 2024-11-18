package errhdl

import (
	"github.com/DaHuangQwQ/gweb/context"
	"github.com/DaHuangQwQ/gweb/middlewares"
	"github.com/DaHuangQwQ/gweb/types"
)

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: make(map[int][]byte),
	}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) {
	m.resp[status] = data
}

func (m *MiddlewareBuilder) Build() middlewares.Middleware {
	return func(next types.HandleFunc) types.HandleFunc {
		return func(ctx *context.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if !ok {
				ctx.RespData = resp
			}
		}
	}
}
