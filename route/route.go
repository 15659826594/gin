package route

import (
	"fmt"
	"gin"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

type IException interface {
	Exception(*gin.Context)
}
type IInitialize interface {
	Initialize(*gin.Context)
}
type IBeforeAction interface {
	BeforeAction() []gin.HandlerFunc
}

var Router *Tree

// Register 注册路由规则
func Register(cStruct any) *Tree {
	if Router == nil {
		Router = NewTree()
	}
	_, filename, _, _ := runtime.Caller(1)
	module := Router.Module(filename)

	controller := NewController(cStruct)

	module.Controllers = append(module.Controllers, controller)

	return Router
}

// Builder 构建路由
func Builder(engine *gin.Engine, defaultMethod []string) {
	if Router == nil {
		return
	}

	engine.Use(Abort())

	for _, version := range Router.Versions {
		// 当版本为application时忽略模块
		verGroup := engine.Group(version.Path())
		for _, module := range version.Modules {
			modGroup := verGroup.Group(module.Path())
			for _, controller := range module.Controllers {
				var conGroup *gin.RouterGroup
				if strings.HasPrefix(controller.Path(), "/") {
					conGroup = engine.Group(controller.Path())
				} else {
					conGroup = modGroup.Group(controller.Path())
				}
				//中间件链
				var handlersChain []gin.HandlerFunc
				//异常捕获
				if initfunc, ok := controller.Raw.(IException); ok {
					handlersChain = append(handlersChain, initfunc.Exception)
				}
				//将控制器上的方法挂载到上下文
				handlersChain = append(handlersChain, mount(controller.Raw))
				//前置操作
				if initfunc, ok := controller.Raw.(IBeforeAction); ok {
					handlersChain = append(handlersChain, initfunc.BeforeAction()...)
				}

				for _, action := range controller.Actions {
					//方法的位置
					handlerName := fmt.Sprintf("%s.%s.%s", module.AbsolutePath, controller.Name, action.Name)
					for _, relativePath := range action.Paths() {
						var chains = make([]gin.HandlerFunc, len(handlersChain))
						copy(chains, handlersChain)
						tmpGroup := conGroup
						if strings.HasPrefix(relativePath, "/") {
							tmpGroup = engine.Group("")
						}
						chains = append([]gin.HandlerFunc{func(c *gin.Context) {
							c.SetContextS(handlerName)
						}}, chains...)
						//设置HandlerName -> 自定义异常处理 -> 控制器方法挂在到上下文 -> 控制器初始化方法 -> action方法
						if initfunc, ok := controller.Raw.(IInitialize); ok {
							chains = append(chains, initfunc.Initialize)
						}
						createURL(tmpGroup, action.Methods(defaultMethod), relativePath, append(chains, action.Handler), handlerName)
					}
				}
			}
		}
	}
}

func createURL(group *gin.RouterGroup, httpMethods []string, url string, handlers []gin.HandlerFunc, handlerName string) {
	for _, method := range httpMethods {
		if slices.Contains(httpMethods, "Any") {
			group.Any(url, handlers...)
			break
		}
		group.Handle(method, url, handlers...)
	}
	if gin.IsDebugging() {
		absolutePath := filepath.ToSlash(path.Clean(group.BasePath() + "/" + url))
		nuHandlers := len(group.Handlers) + len(handlers)
		fmt.Printf("[GIN-debug] %-10s %-25s --> %s (%d handlers)\n", strings.Join(httpMethods, " "), absolutePath, handlerName, nuHandlers)
	}
}

func mount(obj any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IController.ISuccess, _ = obj.(gin.ISuccess)
		c.IController.IFail, _ = obj.(gin.IFail)
		c.IController.IResult, _ = obj.(gin.IResult)
		c.IController.IResponseType, _ = obj.(gin.IResponseType)
		c.IController.IAssign, _ = obj.(gin.IAssign)
		c.IController.IFetch, _ = obj.(gin.IFetch)
	}
}

// Abort 中断, 不用每次都添加return
func Abort() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			rec := recover()
			if exc, ok := rec.(string); ok && exc == "Abort" {
				c.Abort()
			} else if rec != nil {
				panic(rec)
			}
		}()
		c.Next()
	}
}
