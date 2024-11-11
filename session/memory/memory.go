package memory

import (
	"context"
	"errors"
	"github.com/DaHuangQwQ/gweb/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Store struct {
	mutex      sync.RWMutex
	c          *cache.Cache
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
func NewStore(expiration time.Duration) *Store {
	return &Store{
		c:          cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (m *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	m.c.Set(sess.ID(), sess, m.expiration)
	return sess, nil
}

func (m *Store) Refresh(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	sess, ok := m.c.Get(id)
	if !ok {
		return errors.New("session not found")
	}
	m.c.Set(sess.(*memorySession).ID(), sess, m.expiration)
	return nil
}

func (m *Store) Remove(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.c.Delete(id)
	return nil
}

func (m *Store) Get(ctx context.Context, id string) (session.Session, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	sess, ok := m.c.Get(id)
	if !ok {
		return nil, errors.New("session not found")
	}
	return sess.(*memorySession), nil
}

type memorySession struct {
	mutex sync.RWMutex
	id    string
	data  map[string]string
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("not found the key")
	}
	return val, nil
}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.id
}
