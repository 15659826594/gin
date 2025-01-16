package route

import (
	"fmt"
	"gin"
	"gin/annotation"
	"path"
	"path/filepath"
	"runtime"
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

	astFile := parseFile(filename)

	controller := NewController(cStruct)

	//为方法绑定注解
	if controller != nil {
		for _, action := range controller.Actions {
			comments, err := astFile.GetComments(controller.Name, action.Name)
			if err == nil {
				action.Annotations = comments
			}
		}
	}

	module.Controllers = append(module.Controllers, controller)

	return Router
}

// Builder 构建路由
func Builder(engine *gin.Engine, defaultMethod []string) {
	if Router == nil {
		return
	}

	for _, version := range Router.Versions {
		// http://localhost:8080/api/user/index 当版本为application时忽略模块
		// http://localhost:8080/v2/api/user/index
		level1 := engine.Group(version.Path())
		for _, module := range version.Modules {
			level2 := level1.Group(module.Path())
			for _, controller := range module.Controllers {
				var level3 *gin.RouterGroup
				// http://localhost:8080/user/index  斜杠开头 , 忽略上级路由
				if strings.HasPrefix(controller.Path(), "/") {
					level3 = engine.Group(controller.Path())
				} else {
					level3 = level2.Group(controller.Path())
				}
				//异常捕获 | 前置操作(多个)
				handlersChain := append([]gin.HandlerFunc{controller.Exception()}, controller.BeforeAction()...)
				//将控制器的方法挂在到上下文
				handlersChain = append(handlersChain, func(c *gin.Context) {
					c.Set("__jump__", controller.IJump)
					c.Set("__view__", controller.IView)
				})
				level3.Use(handlersChain...)

				for _, action := range controller.Actions {
					//方法的位置
					handlerName := fmt.Sprintf("%s.%s.%s", module.AbsolutePath, controller.Name, action.Name)
					flag := false
					for _, anno := range action.Annotations {
						myfun := annotation.Get(anno.Name)
						if myfun != nil {
							tmpRG := level3
							httpMethods, uri := myfun(anno.Name, anno.Attributes)
							if uri == "" {
								uri = action.Path()
							}
							// http://localhost:8080/index  斜杠开头 , 忽略上级路由
							if strings.HasPrefix(uri, "/") {
								tmpRG = engine.Group("")
								tmpRG.Use(handlersChain...)
							}
							createURL(tmpRG, httpMethods, uri, []gin.HandlerFunc{controller.Initialize, action.Handler}, handlerName)
							flag = true
						}
					}
					//如果没有路由注解则自动生成方法
					if !flag {
						createURL(level3, defaultMethod, action.Path(), []gin.HandlerFunc{controller.Initialize, action.Handler}, handlerName)
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
		if method == "Any" {
			group.Any(url, handlers...)
			continue
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
