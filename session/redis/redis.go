package redis

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/DaHuangQwQ/gweb/session"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:generate

var (
	//go:embed lua/cmd.lua
	lua string

	//go:embed lua/set.lua
	setLua string
)

var errSessionNotExist = errors.New("redis-session: session not found")

type StoreOption func(store *Store)

type Store struct {
	prefix     string
	client     redis.Cmdable
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		client:     client,
		prefix:     "session",
		expiration: time.Minute * 15,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	key := s.key(id)
	_, err := s.client.Eval(ctx, setLua, []string{key}, "_sess_id", id, s.expiration.Milliseconds()).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		key:    key,
		id:     id,
		client: s.client,
	}, nil
}

func (s *Store) key(id string) string {
	return fmt.Sprintf("%s_%s", s.prefix, id)
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := s.key(id)
	affected, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	if !affected {
		return errSessionNotExist
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	_, err := s.client.Del(ctx, s.key(id)).Result()
	return err
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := s.key(id)
	i, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, errSessionNotExist
	}
	return &Session{
		id:     id,
		key:    key,
		client: s.client,
	}, nil
}

type Session struct {
	key    string
	id     string
	client redis.Cmdable
}

// Set hset map[string]map[string]string
func (m *Session) Set(ctx context.Context, key string, val string) error {
	res, err := m.client.Eval(ctx, lua, []string{m.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotExist
	}
	return nil
}

func (m *Session) Get(ctx context.Context, key string) (string, error) {
	return m.client.HGet(ctx, m.key, key).Result()
}

func (m *Session) ID() string {
	return m.id
}
