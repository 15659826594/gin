package app

import (
	"gin"
	"html/template"
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
	Port                string //端口
	Debug               TriState
	TrustedProxies      []string
	Static              map[string]string
	StaticFile          map[string]string
	RouteRule           func(engine *gin.Engine)
	Methods             []string //默认添加的请求
	HTMLFolder          string   //html存放的目录
	DisableConsoleColor TriState //控制台颜色
	FuncMap             template.FuncMap
}

func NewConfig(config *Config) *Config {
	def := &Config{
		Port:           ":80",
		Debug:          True,
		TrustedProxies: []string{"127.0.0.1"}, // 设置 Gin 只信任本机的代理服务器
		Static: map[string]string{
			"/assets": "./public/assets",
		},
		StaticFile: map[string]string{
			"/favicon.ico": getFaviconIco("./public/assets/img/favicon.ico"),
		},
		Methods:             []string{"GET", "POST"},
		HTMLFolder:          "application",
		DisableConsoleColor: False,
		FuncMap: template.FuncMap{
			"htmlentities": htmlentities,
		},
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
