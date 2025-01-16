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
func Register(cStruct IController) *Tree {
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
				//异常捕获 | 前置操作(多个)
				handlersChain := append([]gin.HandlerFunc{controller.Exception()}, controller.BeforeAction()...)
				level3.Use(handlersChain...)

				for _, action := range controller.Actions {
					//方法的位置
					handlerName := fmt.Sprintf("%s.%s.%s", module.AbsolutePath, controller.Name, action.Name)
					level3.Use(func(c *gin.Context) {
						c.Set("Request.URL", strings.Join([]string{module.Name, controller.Name, action.Name}, "/"))
					})
					flag := false
					for _, anno := range action.Annotations {
						myfun := annotation.Get(anno.Name)
						if myfun != nil {
							httpMethods, uri := myfun(anno.Name, anno.Attributes)
							if uri == "" {
								uri = action.Path()
							}
							if strings.HasPrefix(uri, "/") {
								level3 = engine.Group("")
								level3.Use(handlersChain...)
							}
							createURL(level3, httpMethods, uri, controller.Initialize, handlerName)
							flag = true
						}
					}
					//如果没有路由注解则自动生成方法
					if !flag {
						createURL(level3, defaultMethod, action.Path(), controller.Initialize, handlerName)
					}
				}
			}
		}
	}
}

func createURL(group *gin.RouterGroup, httpMethods []string, url string, handler gin.HandlerFunc, handlerName string) {
	for _, method := range httpMethods {
		if method == "Any" {
			group.Any(url, handler)
			continue
		}
		group.Handle(method, url, handler)
	}
	httpMethod := strings.Join(httpMethods, " ")
	absolutePath := filepath.ToSlash(path.Clean(group.BasePath() + "/" + url))
	nuHandlers := len(group.Handlers) + 1
	fmt.Printf("[GIN-debug] %-10s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
}
