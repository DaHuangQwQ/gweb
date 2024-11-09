package errhdl

import "github.com/DaHuangQwQ/gweb"

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

func (m *MiddlewareBuilder) Build() gweb.Middleware {
	return func(next gweb.HandleFunc) gweb.HandleFunc {
		return func(ctx *gweb.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if !ok {
				ctx.RespData = resp
			}
		}
	}
}
