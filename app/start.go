package app

import (
	"fmt"
	"gin"
	_ "gin/annotation"
	"gin/route"
)

type Application struct {
	*gin.Engine
	*Config
}

func Run(engine *gin.Engine, config *Config) error {
	config = NewConfig(config)
	app := &Application{
		engine,
		config,
	}
	if config.Debug.Bool() {
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.DisableConsoleColor.Bool() {
		gin.DisableConsoleColor()
	}

	if app.Engine == nil {
		app.Engine = gin.Default()
		engine = app.Engine
	}
	engine.SetFuncMap(config.FuncMap)
	// 在初始化时。即，在注册任何路由或路由器在套接字中侦听之前
	engine.LoadHTMLFolder(config.HTMLFolder)

	//静态文件服务
	for i, s := range config.Static {
		engine.Static(i, s)
	}
	for i, s := range config.StaticFile {
		engine.StaticFile(i, s)
	}

	engine.Use(gin.RecoveryExit())

	var err error

	err = engine.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		return err
	}

	if config.RouteRule != nil {
		config.RouteRule(app.Engine)
	}

	route.Builder(engine, config.Methods)

	if !gin.IsDebugging() {
		fmt.Printf("[GIN-%s] Listening and serving HTTP on %s\n", gin.Mode(), config.Port)
	}

	err = app.Run(config.Port)
	if err != nil {
		return err
	}

	return nil
}
