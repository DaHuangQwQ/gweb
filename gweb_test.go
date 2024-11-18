package gweb

import (
	"github.com/DaHuangQwQ/gweb/context"
	"testing"
)

func TestGweb(t *testing.T) {
	server := Default()

	server.Get("/*", func(ctx *context.Context) {
		_ = ctx.RespJSONOK("hello world")
		return
	})

	_ = server.Start(":8081")
}
