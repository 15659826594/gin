package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type IScope interface {
	Set(string, any)
	Get(string) (any, bool)
	GetStringMap(string) map[string]any
}

type Scope struct {
	Keys map[string]any
}

func newScope() *Scope {
	return &Scope{Keys: make(map[string]any)}
}

func (s *Scope) Set(key string, value any) {
	if s.Keys == nil {
		s.Keys = make(map[string]any)
	}
	s.Keys[key] = value
}

func (s *Scope) Get(key string) (value any, exists bool) {
	value, exists = s.Keys[key]
	return
}

func (s *Scope) GetStringMap(key string) (sm map[string]any) {
	if val, ok := s.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

var global = map[string]*Scope{
	"_sys_": newScope(),
}

func getOrCreateScope(name string) *Scope {
	if val, ok := global[name]; ok {
		return val
	}
	global[name] = newScope()
	return global[name]
}

/*Load
 * 载入文件
 * @param  string $file 文件路径
 * @param  string $name key
 * @param  string $args 作用域
 */
func Load(file string, name string, args ...any) error {
	var data map[string]any
	ext := filepath.Ext(file)
	switch ext {
	case ".json":
		tmp, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		err = json.Unmarshal(tmp, &data)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown file")
	}
	var scope IScope
	scope = global["_sys_"]
	for i, arg := range args {
		switch i {
		case 0:
			if s, ok := arg.(IScope); ok {
				scope = s
			} else if s1, ok1 := arg.(string); ok1 {
				scope = getOrCreateScope(s1)
			}
		}
	}
	Set(name, data, scope)
	return nil
}

/*SearchFiles
 * 搜索配置文件
 * @param  string $folder 文件夹
 */
func SearchFiles(folder string, callback func(string, string, ...any)) {
	_ = filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		ext := filepath.Ext(path)
		if ext == ".json" {
			name := strings.TrimSuffix(filepath.Base(path), ext)
			callback(path, name, nil)
		}
		return nil
	})
}

/*Get
 * 获取配置参数 为空则获取所有配置
 * @param  string $name 配置参数名（支持二级配置 . 号分割）
 * @param  string $range  作用域
 * @return mixed
 */
func Get(args ...any) any {
	var name string
	var scope IScope
	scope = global["_sys_"]
	for i, arg := range args {
		switch i {
		case 0:
			name = arg.(string)
		case 1:
			if s, ok := arg.(IScope); ok {
				scope = s
			} else if s1, ok1 := arg.(string); ok1 {
				scope = getOrCreateScope(s1)
			}
		}
	}
	if name == "" {
		return scope
	}
	deep := strings.Split(name, ".")
	if v, ok := scope.Get(deep[0]); ok {
		if len(deep) == 1 {
			return v
		} else {
			return rcteGet(v, deep[1:])
		}
	} else {
		return nil
	}
}

// 递归获取
func rcteGet(arr any, stack []string) any {
	if len(stack) == 1 {
		if v, ok := arr.(map[string]any); ok {
			return v[stack[0]]
		}
	} else {
		if v, ok := arr.(map[string]any); ok {
			return rcteGet(v[stack[0]], stack[1:])
		}
	}
	return nil
}

/*Set
 * 设置配置参数 name 为数组则为批量设置
 * @access public
 * @param  string|array $name  配置参数名（支持二级配置 . 号分割）
 * @param  mixed        $value 配置值
 * @param  string       $range 作用域
 * @return mixed
 */
func Set(name string, value any, args ...any) {
	if name == "" {
		return
	}
	var scope IScope
	scope = global["_sys_"]
	for i, arg := range args {
		switch i {
		case 0:
			if s, ok := arg.(IScope); ok {
				scope = s
			} else if s1, ok1 := arg.(string); ok1 {
				scope = getOrCreateScope(s1)
			}
		}
	}
	deep := strings.Split(name, ".")
	if len(deep) <= 1 {
		scope.Set(name, value)
	} else {
		l1 := scope.GetStringMap(deep[0])
		scope.Set(deep[0], rcteSet(l1, deep[1:], value))
	}
}

// 递归设置
func rcteSet(arr map[string]any, stack []string, value any) map[string]any {
	if arr == nil {
		arr = make(map[string]any)
	}
	key := stack[0]
	if len(stack) == 1 {
		arr[key] = value
	} else {
		var tmp map[string]any
		if v, ok := arr[key].(map[string]any); ok {
			tmp = v
		}
		arr[key] = rcteSet(tmp, stack[1:], value)
	}
	return arr
}
