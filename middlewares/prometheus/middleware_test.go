package prometheus

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/DaHuangQwQ/gweb/internal/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	s := gweb.NewHttpServer()
	s.Get("/", func(ctx *context.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *context.Context) {
		time.Sleep(time.Second)
	})

	s.UseAll("/*", (&MiddlewareBuilder{
		Namespace: "github.com/DaHuangQwQ/gweb",
		Subsystem: "web",
		Name:      "http_request",
		Help:      "这是测试例子",
	}).Build())
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 一般来说，在实际中我们都会单独准备一个端口给这种监控
		http.ListenAndServe(":2112", nil)
	}()
	s.Start(":8081")
}
