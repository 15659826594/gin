package route

import (
	"fmt"
	"gin"
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

	if !controller.IsNil() {
		//为方法绑定注解
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

	gin.DebugPrintTable([]string{"version", "module", "controller", "action", "uri", "method", "url"}, func(callback func([]string)) {
		for _, version := range Router.Versions {
			level1 := engine.Group(version.Path())
			for _, module := range version.Modules {
				level2 := level1.Group(module.Path())
				for _, controller := range module.Controllers {
					level3 := level2.Group(controller.Path())
					//异常捕获 | 前置操作(多个)
					level3.Use(append([]gin.HandlerFunc{controller.Exception()}, controller.BeforeAction()...)...)
					for _, action := range controller.Actions {
						level3.Use(func(c *gin.Context) {
							c.Set("Module", module.Path())
							c.Set("Controller", controller.Path())
							c.Set("Action", action.Path())
						})
						methodName, uri, err := action.Mapping(level3, controller.Initialize, defaultMethod)
						if err == nil {
							callback([]string{version.Path(), module.Name, controller.Name, action.Name, uri, strings.Join(methodName, " "), fmt.Sprintf("%s/%s/%s/%s", version.Path(), module.Path(), controller.Path(), uri)})
						} else if err.Error() == "default" {
							callback([]string{version.Path(), module.Name, controller.Name, action.Name, uri, fmt.Sprintf("def(%s)", strings.Join(methodName, " ")), fmt.Sprintf("%s/%s/%s/%s", version.Path(), module.Path(), controller.Path(), action.Path())})
						} else if err.Error() == "invalid" {
							msg := fmt.Sprintf("%s: %s %s %s %s", err.Error(), version.Name, module.Name, controller.Name, action.Name)
							fmt.Printf("\033[1;31;40m%s\033[0m\n", msg)
						}
					}
				}
			}
		}
	})
}
