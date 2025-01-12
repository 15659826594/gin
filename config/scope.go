package config

import "sync"

type IScope interface {
	Set(string, any)
	Get(string) (any, bool)
	GetStringMap(string) map[string]any
}

type Scope struct {
	mu   sync.RWMutex
	Keys map[string]any
}

func newScope() *Scope {
	return &Scope{Keys: make(map[string]any)}
}

func (s *Scope) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Keys == nil {
		s.Keys = make(map[string]any)
	}
	s.Keys[key] = value
}

func (s *Scope) Get(key string) (value any, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists = s.Keys[key]
	return
}

func (s *Scope) GetStringMap(key string) (sm map[string]any) {
	if val, ok := s.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

func findOrCreateScope(name string) *Scope {
	if val, ok := global[name]; ok {
		return val
	}
	global[name] = newScope()
	return global[name]
}
