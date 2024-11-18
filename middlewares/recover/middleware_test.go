package recover

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/DaHuangQwQ/gweb/context"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder(http.StatusBadRequest, []byte("bad req"), func(ctx *context.Context) {
		println(ctx.RespStatusCode)
	})
	server := gweb.NewHttpServer()
	server.UseAll("/*", builder.Build())

	server.Get("/*", func(ctx *context.Context) {
		panic("panic!")
	})
	server.Start(":8081")
}
