package lang

import (
	"errors"
	"strings"
	"sync"
)

// Dict 字典结构体
type Dict struct {
	mu   sync.RWMutex
	Keys map[string]string
	Subs map[string]*Dict
}

func NewDict() *Dict {
	return &Dict{
		Keys: make(map[string]string),
		Subs: make(map[string]*Dict),
	}
}

// Search 当前字典中搜索Keys
func (d *Dict) Search(name string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if val, ok := d.Keys[name]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}

// Append 追加词典
func (d *Dict) Append(name string, value string) {
	d.mu.Lock()
	d.Keys[name] = value
	d.mu.Unlock()
}

// SearchSub 查找子词典
func (d *Dict) SearchSub(dictname any) (*Dict, error) {
	var dictName []string
	switch v := dictname.(type) {
	case string:
		dictName = splitDictName(v)
	case []string:
		dictName = v
	default:
		return nil, errors.New("arg error: Only supported string or []string")
	}
	lens := len(dictName)
	if lens == 0 {
		return nil, errors.New("dict path empty")
	}
	var dict *Dict
	if val, ok := d.Subs[dictName[0]]; !ok {
		return nil, errors.New("dict not found")
	} else {
		dict = val
	}
	if lens == 1 {
		return dict, nil
	} else {
		return dict.SearchSub(dictName[1:])
	}
}

// CreateSub 查找子词典
func (d *Dict) CreateSub(dictname any) (*Dict, error) {
	var dictName []string
	switch v := dictname.(type) {
	case string:
		dictName = splitDictName(v)
	case []string:
		dictName = v
	default:
		return nil, errors.New("arg error: Only supported string or []string")
	}
	lens := len(dictName)
	if lens == 0 {
		return d, nil
	}
	var dict *Dict
	if val, ok := d.Subs[dictName[0]]; !ok {
		d.Subs[dictName[0]] = NewDict()
		dict = d.Subs[dictName[0]]
	} else {
		dict = val
	}
	return dict.CreateSub(dictName[1:])
}

// 翻译链路
func dictChains(dict *Dict, dictName []string, args ...*[]map[string]string) ([]map[string]string, error) {
	var chains *[]map[string]string
	lens := len(dictName)
	if len(args) == 0 {
		tmp := make([]map[string]string, lens+1)
		chains = &tmp
	} else {
		chains = args[0]
		*chains = append(*chains, dict.Keys)
	}
	if lens == 0 {
		return *chains, nil
	} else {
		if v, ok := dict.Subs[dictName[0]]; ok {
			return dictChains(v, dictName[1:], chains)
		} else {
			return *chains, nil
		}
	}
}

// splitDictName zh-cn.index.user.index拆分成[zh,cn,index,user,index]
func splitDictName(dictName string) []string {
	var levels []string
	tmp := strings.Split(strings.ToLower(dictName), ".")
	for _, v := range tmp {
		levels = append(levels, strings.Split(v, "-")...)
	}
	return levels
}
