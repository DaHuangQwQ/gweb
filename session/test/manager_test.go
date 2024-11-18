package test

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/DaHuangQwQ/gweb/internal/context"
	"github.com/DaHuangQwQ/gweb/internal/types"
	"github.com/DaHuangQwQ/gweb/session"
	"github.com/DaHuangQwQ/gweb/session/cookie"
	"github.com/DaHuangQwQ/gweb/session/memory"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	s := gweb.NewHttpServer()

	m := session.Manager{
		SessCtxKey: "_sess",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessid",
			cookie.WithCookieOption(func(c *http.Cookie) {
				c.HttpOnly = true
			})),
	}

	s.Get("/login", func(ctx *context.Context) {
		// 登录校验
		id := uuid.New()
		sess, err := m.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		err = sess.Set(ctx.Req.Context(), "mykey", "some value")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})
	s.Get("/resource", func(ctx *context.Context) {
		sess, err := m.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		val, err := sess.Get(ctx.Req.Context(), "mykey")
		ctx.RespData = []byte(val)
	})

	s.Get("/logout", func(ctx *context.Context) {
		_ = m.RemoveSession(ctx)
	})

	s.UseAll("/*", func(next types.HandleFunc) types.HandleFunc {
		return func(ctx *context.Context) {
			if ctx.Req.URL.Path != "/login" {
				sess, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					return
				}
				ctx.UserValues["sess"] = sess
				_ = m.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	s.Start(":8081")
}
