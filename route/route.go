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

var Router *Tree

// Register 注册路由规则
func Register(cStruct gin.IController) *Tree {
	if Router == nil {
		Router = NewTree()
	}
	_, filename, _, _ := runtime.Caller(1)
	module := Router.Module(filename)

	//astFile := parseFile(filename)

	controller := NewController(cStruct)

	//为方法绑定注解
	//if controller != nil {
	//	for _, action := range controller.Actions {
	//		comments, err := astFile.GetComments(controller.Name, action.Name)
	//		if err == nil {
	//			action.Annotations = comments
	//		}
	//	}
	//}

	module.Controllers = append(module.Controllers, controller)

	return Router
}

// Builder 构建路由
func Builder(engine *gin.Engine, defaultMethod []string) {
	if Router == nil {
		return
	}

	for _, version := range Router.Versions {
		// 当版本为application时忽略模块
		level1 := engine.Group(version.Path())
		for _, module := range version.Modules {
			level2 := level1.Group(module.Path())
			for _, controller := range module.Controllers {
				var level3 *gin.RouterGroup
				if strings.HasPrefix(controller.Path(), "/") {
					level3 = engine.Group(controller.Path())
				} else {
					level3 = level2.Group(controller.Path())
				}
				//异常捕获 | 控制器的方法挂在到上下文 | 前置操作(多个)
				handlersChain := append([]gin.HandlerFunc{controller.Exception(), func(c *gin.Context) {
					c.Set("__jump__", controller.IJump)
					c.Set("__view__", controller.IView)
				}}, controller.BeforeAction()...)

				level3.Use(handlersChain...)

				for _, action := range controller.Actions {
					//方法的位置
					handlerName := fmt.Sprintf("%s.%s.%s", module.AbsolutePath, controller.Name, action.Name)
					for _, relativePath := range action.Paths() {
						tmpRG := level3
						if strings.HasPrefix(relativePath, "/") {
							tmpRG = engine.Group("")
							tmpRG.Use(handlersChain...)
						}
						createURL(tmpRG, action.Methods(defaultMethod), relativePath, []gin.HandlerFunc{controller.Initialize, action.Handler}, handlerName)
					}
				}
			}
		}
	}
}

func createURL(group *gin.RouterGroup, httpMethods []string, url string, handler []gin.HandlerFunc, handlerName string) {
	//执行顺序 设置handlerName, 控制器Init方法, action方法
	handlers := append([]gin.HandlerFunc{func(c *gin.Context) {
		c.Set("__handler_name__", handlerName)
	}}, handler...)

	for _, method := range httpMethods {
		if slices.Contains(httpMethods, "Any") {
			group.Any(url, handlers...)
			break
		}
		group.Handle(method, url, handlers...)
	}
	if gin.IsDebugging() {
		httpMethod := strings.Join(httpMethods, " ")
		absolutePath := filepath.ToSlash(path.Clean(group.BasePath() + "/" + url))
		nuHandlers := len(group.Handlers) + len(handlers)
		fmt.Printf("[GIN-debug] %-10s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}
