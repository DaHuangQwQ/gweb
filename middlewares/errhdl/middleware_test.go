package errhdl

import (
	"github.com/DaHuangQwQ/gweb"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte("not found"))
	server := gweb.NewHttpServer()
	server.UseAll("/*", builder.Build())
	server.Start(":8081")
}
