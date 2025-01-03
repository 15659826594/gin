package app

import (
	"gin"
	"reflect"
)

type TriState int

const (
	Undefined TriState = iota
	False
	True
)

func (t TriState) Bool() bool {
	switch t {
	case Undefined:
		return false
	case False:
		return false
	case True:
		return true
	}
	return false
}

type Config struct {
	Port                int //端口
	Debug               TriState
	RouteRule           func(engine *gin.Engine)
	Methods             []string //默认添加的请求
	HTMLFolder          string   //html存放的目录
	DisableConsoleColor TriState //控制台颜色
	TrustedProxies      []string
}

func NewConfig(config *Config) *Config {
	def := &Config{
		Debug:               True,
		Port:                8080,
		TrustedProxies:      []string{"127.0.0.1"}, // 设置 Gin 只信任本机的代理服务器
		Methods:             []string{"GET", "POST"},
		HTMLFolder:          "application",
		DisableConsoleColor: False,
	}

	if config == nil {
		return def
	}

	defValue := reflect.ValueOf(def).Elem()

	structValue := reflect.ValueOf(config).Elem()
	for i, lens := 0, structValue.NumField(); i < lens; i++ {
		field := structValue.Type().Field(i) // 获取字段类型
		value := structValue.Field(i)        // 获取字段值

		if !gin.Empty(value.Interface()) {
			defValue.FieldByName(field.Name).Set(value)
		}
	}
	return def
}
