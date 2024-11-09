package recover

import (
	"github.com/DaHuangQwQ/gweb"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder(http.StatusBadRequest, []byte("bad req"), func(ctx *gweb.Context) {
		println(ctx.RespStatusCode)
	})
	server := gweb.NewHttpServer()
	server.UseAll("/*", builder.Build())

	server.Get("/*", func(ctx *gweb.Context) {
		panic("panic!")
	})
	server.Start(":8081")
}
