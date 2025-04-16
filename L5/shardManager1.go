package l5

import (
	"database/sql"
	"errors"
	"sync"
)

type Shard struct {
	Address string
	Number  int
}

type Manager struct {
	size int
	ss   *sync.Map
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
	}
}

func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}

func (m *Manager) ShardById(entityId int) (*Shard, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId / m.size
	if s, ok := m.ss.Load(n); ok {
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}

type Pool struct {
	sync.RWMutex
	cc map[string]*sql.DB
}

func NewPool() *Pool {
	return &Pool{
		cc: map[string]*sql.DB{},
	}
}

func (p *Pool) Connection(addr string) (*sql.DB, error) {
	p.RLock()
	if c, ok := p.cc[addr]; ok {
		defer p.RUnlock()
		return c, nil
	}
	p.RUnlock()
	p.Lock()
	defer p.Unlock()
	if c, ok := p.cc[addr]; ok {
		return c, nil
	}
	var err error
	p.cc[addr], err = sql.Open("postgres", addr)
	return p.cc[addr], err
}

//-----------------------------------

type User struct {
	UserId int
	Name   string
	Age    int
	Spouse int
}

func (u *User) connection() (*sql.DB, error) {
	s, err := m.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}
