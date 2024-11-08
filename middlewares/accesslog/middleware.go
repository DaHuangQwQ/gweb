package accesslog

import (
	"encoding/json"
	"gweb"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(log string) {
			println(log)
		}}
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() gweb.Middleware {
	return func(next gweb.HandleFunc) gweb.HandleFunc {
		return func(ctx *gweb.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
