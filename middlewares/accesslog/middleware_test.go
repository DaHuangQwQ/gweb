package accesslog

import (
	"gweb"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	b := &MiddlewareBuilder{}
	b.LogFunc(func(log string) {
		println(log)
	})
	s := gweb.NewHttpServer()
	s.Get("/", func(ctx *gweb.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *gweb.Context) {
		time.Sleep(time.Second)
		ctx.RespData = []byte("hello, user")
	})
	s.UseAll("/*", b.Build())
	s.Start(":8081")
}
