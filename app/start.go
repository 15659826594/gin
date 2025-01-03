package app

import (
	"fmt"
	"gin"
	_ "gin/annotation"
	"gin/route"
)

type Application struct {
	Engine *gin.Engine
	Config *Config
}

func Run(config *Config) *Application {
	if config == nil {
		config = NewConfig()
	}
	if config.Debug {
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.DisableConsoleColor {
		gin.DisableConsoleColor()
	}

	appInst := &Application{
		Config: config,
	}

	return appInst
}

func (that *Application) SetEngine(engine *gin.Engine) *Application {
	that.Engine = engine
	return that
}

func (that *Application) Send(port int) error {
	if that.Engine == nil {
		that.Engine = gin.Default()
	}
	if that.Config.HTMLFolder == "" {
		that.Config.HTMLFolder = "application"
	}
	//加载html文件夹
	that.Engine.LoadHTMLFolder(that.Config.HTMLFolder)
	//自定义路由
	if that.Config.RouteRule != nil {
		that.Config.RouteRule(that.Engine)
	}
	//设置默认请求(无注解的action)
	if that.Config.Methods == nil {
		that.Config.Methods = []string{"GET", "POST"}
	}
	//通过控制器生成的路由
	that.builder(that.Engine, that.Config.Methods)

	if !gin.IsDebugging() {
		fmt.Printf("Listening and serving HTTP on :%d\n", port)
	}

	err := that.Engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}

// 构建路由
func (that *Application) builder(engine *gin.Engine, defaultMethod []string) {
	route.Builder(engine, defaultMethod)
}
