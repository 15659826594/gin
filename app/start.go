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

	var err error

	engine.LoadHTMLFolder(config.HTMLFolder)

	err = engine.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		return err
	}

	if config.RouteRule != nil {
		config.RouteRule(app.Engine)
	}

	route.Builder(engine, config.Methods)

	if !gin.IsDebugging() {
		fmt.Printf("[GIN-%s] Listening and serving HTTP on %s\n", gin.Mode(), fmt.Sprintf(":%d", config.Port))
	}

	err = app.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	return nil
}
