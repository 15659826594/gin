package app

import (
	"fmt"
	"gin"
	"gin/route"
	"path/filepath"
)

type Application struct {
	*gin.Engine
	*Config
}

func create(engine *gin.Engine, config *Config) *Application {
	config = NewConfig(config)
	return &Application{
		engine,
		config,
	}
}

func (app *Application) Init() *Application {
	config := app.Config
	engine := app.Engine

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
	// LoadHTMLFolder包含模板内的define
	engine.LoadHTMLFolder(config.HTMLFolder, func(path string) string {
		path, _ = filepath.Rel("application", path)
		return filepath.ToSlash(path)
	})

	//静态文件服务
	for i, s := range config.Static {
		engine.Static(i, s)
	}
	for i, s := range config.StaticFile {
		engine.StaticFile(i, s)
	}

	_ = engine.SetTrustedProxies(config.TrustedProxies)

	//自定义路由
	if config.RouteRule != nil {
		config.RouteRule(app.Engine)
	}
	//404页面
	if config.NoRoute != nil {
		engine.NoRoute(config.NoRoute)
	}

	route.Builder(engine, config.Methods)

	return app
}

func New(engine *gin.Engine, config *Config) *Application {
	app := create(engine, config)
	return app.Init()
}

func (app *Application) Run() error {
	var err error
	if !gin.IsDebugging() {
		fmt.Println(fmt.Sprintf("[GIN-%s] Listening and serving HTTP on %s", gin.Mode(), app.Config.Port))
	}
	err = app.Engine.Run(app.Config.Port)
	if err != nil {
		return err
	}
	return nil
}

func Run(engine *gin.Engine, config *Config) error {
	return create(engine, config).Init().Run()
}
