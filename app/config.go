package app

import (
	"gin"
	"reflect"
)

type Config struct {
	Port                int //端口
	TrustedProxies      []string
	Debug               bool
	RouteRule           func(engine *gin.Engine)
	Methods             []string //默认添加的请求
	HTMLFolder          string   //html存放的目录
	DisableConsoleColor bool     //控制台颜色
}

func NewConfig(config *Config) *Config {
	def := &Config{
		Debug:               true,
		Port:                8080,
		TrustedProxies:      []string{"127.0.0.1"},
		Methods:             []string{"GET", "POST"},
		HTMLFolder:          "application",
		DisableConsoleColor: false,
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
