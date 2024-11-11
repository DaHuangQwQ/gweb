package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// Store 管理 Session
type Store interface {
	Generate(ctx context.Context, id string) (Session, error)

	Refresh(ctx context.Context, id string) error

	Remove(ctx context.Context, id string) error

	Get(ctx context.Context, id string) (Session, error)
}

type Propagator interface {
	// Inject 将 session id 注入到里面
	// Inject 必须幂等
	Inject(id string, writer http.ResponseWriter) error

	// Extract 将 session id 从 http.Request 中提取出来
	Extract(req *http.Request) (string, error)

	// Remove 将 session id 从 http.ResponseWriter 中删除
	Remove(writer http.ResponseWriter) error
}
