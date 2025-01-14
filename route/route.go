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
				level3 := level2.Group(controller.Path())
				//异常捕获 | 前置操作(多个)
				handlersChain := append([]gin.HandlerFunc{controller.Exception()}, controller.BeforeAction()...)
				handlersChain = append(handlersChain, controller.Initialize)
				level3.Use(handlersChain...)

				for _, action := range controller.Actions {
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
								createURL(level3, httpMethods, uri, controller.Initialize, "fastadmin/application.Route.func1")
							} else if strings.HasPrefix(uri, "/") {
								absoluteGroup := engine.Group("")
								absoluteGroup.Use(handlersChain...)
								createURL(absoluteGroup, httpMethods, uri, controller.Initialize, "fastadmin/application.Route.func1")
							}
							flag = true
						}
					}
					//如果没有路由注解则自动生成方法
					if !flag {
						createURL(level3, defaultMethod, action.Path(), controller.Initialize, "fastadmin/application.Route.func1")
					}
					//fmt.Printf("[GIN-debug] %-6s %-25s --> %s (%d handlers)\n", "GET", "url", "fastadmin/application.Route.func1", 4)
					//methodName, uri, err := action.Mapping(level3, controller.Initialize, defaultMethod)
					//if err == nil {
					//	callback([]string{version.Path(), module.Name, controller.Name, action.Name, uri, strings.Join(methodName, " "), fmt.Sprintf("%s/%s/%s/%s", version.Path(), module.Path(), controller.Path(), uri)})
					//} else if err.Error() == "default" {
					//	callback([]string{version.Path(), module.Name, controller.Name, action.Name, uri, fmt.Sprintf("def(%s)", strings.Join(methodName, " ")), fmt.Sprintf("%s/%s/%s/%s", version.Path(), module.Path(), controller.Path(), action.Path())})
					//} else if err.Error() == "invalid" {
					//	callback([]string{version.Path(), module.Name, controller.Name, action.Name, uri, "error", "invalid"})
					//}
				}
			}
		}
	}
}

func createURL(group *gin.RouterGroup, httpMethods []string, url string, handler gin.HandlerFunc, handlerName string) {
	for _, method := range httpMethods {
		if method == "Any" {
			group.Any(url).Use(handler)
			continue
		}
		group.Handle(method, url).Use(handler)
	}
	//httpMethod, absolutePath, handlerName string, nuHandlers int
	url = filepath.ToSlash(path.Clean(group.BasePath() + "/" + url))
	fmt.Printf("[GIN-debug]  %-25s --> %s (%-10s)\n", url, handlerName, strings.Join(httpMethods, " "))
}
