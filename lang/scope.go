package lang

import (
	"strings"
	"sync"
)

type IScope interface {
	Set(string, any)
	Get(string) (any, bool)
	GetString(string) string
}

type Scope struct {
	mu     sync.RWMutex
	Keys   map[string]any
	Childs map[string]*Scope
}

func newScope() *Scope {
	return &Scope{
		Keys:   make(map[string]any),
		Childs: make(map[string]*Scope),
	}
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

func (s *Scope) GetString(key string) (s1 string) {
	if val, ok := s.Get(key); ok && val != nil {
		s1, _ = val.(string)
	}
	return
}

func findOrCreateScope(name string) *Scope {
	if name == "" {
		return nil
	}
	arr := strings.Split(name, ".")
	name = arr[0]
	if val, ok := global[name]; ok {
		if len(arr) > 1 {
			return findOrCreateScopeDeep(val, arr[1:])
		}
		return val
	}
	global[name] = newScope()
	return global[name]
}

func findOrCreateScopeDeep(scope *Scope, arr []string) *Scope {
	var s *Scope
	name := arr[0]
	if tmp, ok := scope.Childs[name]; ok {
		s = tmp
	} else {
		scope.Childs[name] = newScope()
		s = scope.Childs[name]
	}
	if len(arr) == 1 {
		return s
	}
	return findOrCreateScopeDeep(s, arr[1:])
}

func createScopeChains(lang []string, slicePtr *[]map[string]string) {
	var tmp *Scope
	for _, s := range lang {
		if tmp == nil {
			if _, ok := global[s]; !ok {
				return
			}
			tmp = global[s]
		} else {
			if _, ok := tmp.Childs[s]; !ok {
				return
			}
			tmp = tmp.Childs[s]
		}
		tmpSlice := map[string]string{}
		for k, v := range tmp.Keys {
			tmpSlice[k] = v.(string)
		}
		*slicePtr = append(*slicePtr, tmpSlice)
	}
}
