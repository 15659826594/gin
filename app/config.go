package app

import (
	"gin"
	"gin/utils"
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
	ConfigFile          []string // 配置文件, 支持文件夹和文件
	RouteRule           func(*gin.Engine)
	NoRoute             func(*gin.Context)
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
			"/favicon.ico": utils.Ternary(utils.FileExists("favicon.ico"), "favicon.ico", "public/assets/img/favicon.ico").(string),
		},
		ConfigFile: []string{"internal/extra"},
		//RouteRule:           internal.Route,
		//NoRoute:             internal.NoRoute,
		Methods:             []string{"GET", "POST"},
		HTMLFolder:          "internal",
		DisableConsoleColor: False,
		FuncMap:             FuncMap,
	}

	if config == nil {
		return def
	}

	defValue := reflect.ValueOf(def).Elem()

	structValue := reflect.ValueOf(config).Elem()
	for i, lens := 0, structValue.NumField(); i < lens; i++ {
		field := structValue.Type().Field(i) // 获取字段类型
		value := structValue.Field(i)        // 获取字段值

		if !utils.Empty(value.Interface()) {
			defValue.FieldByName(field.Name).Set(value)
		}
	}

	//for _, files := range def.ConfigFile {
	//	if utils.IsDir(files) {
	//		//载入配置
	//		cfg.SearchFiles(files, func(file string, name string, args ...any) {
	//			_ = cfg.Load(file, name, args...)
	//		})
	//	} else {
	//		_ = cfg.Load(files, cfg.FileName(files))
	//	}
	//}
	//
	//if boolean, ok := cfg.Get("app_debug").(bool); ok && !boolean {
	//	// 如果是调试模式将version置为当前的时间戳可避免缓存
	//	cfg.Set("site.version", strconv.FormatInt(time.Now().Unix(), 10))
	//}

	return def
}
