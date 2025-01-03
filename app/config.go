package app

import "gin"

type Config struct {
	Engine              *gin.Engine
	Debug               bool
	RouteRule           func(engine *gin.Engine)
	Methods             []string
	HTMLFolder          string
	DisableConsoleColor bool
}

type ConfigOption func(config *Config)

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{
		Debug:               false,
		Methods:             []string{"Get", "POST"},
		HTMLFolder:          "application",
		DisableConsoleColor: true,
	}
	for _, option := range options {
		option(config)
	}
	return config
}

func WithEngine(c *gin.Engine) ConfigOption {
	return func(config *Config) {
		config.Engine = c
	}
}

func WithDebug(boolean bool) ConfigOption {
	return func(config *Config) {
		config.Debug = boolean
	}
}

func WithRouteRule(fn func(engine *gin.Engine)) ConfigOption {
	return func(config *Config) {
		config.RouteRule = fn
	}
}

// WithMethods ["GET","POST"]设置默认请求(无注解的action)
func WithMethods(methods []string) ConfigOption {
	return func(config *Config) {
		config.Methods = methods
	}
}

// WithDisableConsoleColor 禁止日志的颜色
func WithDisableConsoleColor(disableConsoleColor bool) ConfigOption {
	return func(config *Config) {
		config.DisableConsoleColor = disableConsoleColor
	}
}
