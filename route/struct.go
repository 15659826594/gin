package route

import (
	"gin"
	"gin/env"
	"gin/utils"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	sep      = string(filepath.Separator)
	basePath = env.Getenv("ROOT_PATH")
	osGetwd  = filepath.Dir(filepath.Dir(basePath))
)

type Tree struct {
	Versions []*Version
}

func NewTree() *Tree {
	return &Tree{}
}

func (that *Tree) Module(filename string) *Module {
	relativePath, _ := filepath.Rel(basePath, filepath.FromSlash(filename))
	paths := strings.Split(relativePath, sep)
	if len(paths) < 4 {
		return nil
	}
	versionName := paths[0]
	moduleName := paths[1]
	version := that.GetVersion(versionName)
	return version.GetModule(moduleName, filename)
}

func (that *Tree) GetVersion(name string) *Version {
	for _, v := range that.Versions {
		if v.Name == name {
			return v
		}
	}
	version := &Version{
		Name:    name,
		Modules: []*Module{},
	}
	that.Versions = append(that.Versions, version)
	return version
}

type Version struct {
	Name    string
	Modules []*Module
}

func (that *Version) GetModule(name string, filename string) *Module {
	for _, v := range that.Modules {
		if v.Name == name {
			return v
		}
	}
	absolutePath := filepath.Dir(strings.TrimSuffix(filename, filepath.Ext(filename)))
	absolutePath, _ = filepath.Rel(osGetwd, absolutePath)
	module := &Module{
		Name:         name,
		Controllers:  []*Controller{},
		AbsolutePath: filepath.ToSlash(absolutePath),
	}
	that.Modules = append(that.Modules, module)
	return module
}

func (that *Version) Path() string {
	if that.Name == "application" {
		return ""
	}
	return utils.Camel2Snake(that.Name)
}

type Module struct {
	Name         string
	Controllers  []*Controller
	AbsolutePath string
}

func (that *Module) Path() string {
	return utils.Camel2Snake(that.Name)
}

type Controller struct {
	gin.IController
	gin.IJump
	gin.IView
	Name    string //控制器名
	Value   string //路由
	Actions []*Action
}

func NewController(obj any) *Controller {
	object, ok := obj.(gin.IController)
	if !ok {
		return nil
	}
	typeOf := reflect.TypeOf(obj)
	valueOf := reflect.ValueOf(obj)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	controller := &Controller{
		IController: object,
		Name:        typeOf.Name(),
		Value:       object.Value(),
		Actions:     []*Action{},
	}
	//后续挂在到上下文
	if jump, isJump := obj.(gin.IJump); isJump {
		controller.IJump = jump
	}
	if view, isView := obj.(gin.IView); isView {
		controller.IView = view
	}
	//提取结构体中的方法
	for i, lens := 0, valueOf.NumMethod(); i < lens; i++ {
		methodReflect := valueOf.Method(i)
		if handlerFunc, ok := methodReflect.Interface().(func(*gin.Context)); ok {
			controller.Actions = append(controller.Actions, &Action{
				Name:    valueOf.Type().Method(i).Name,
				Handler: handlerFunc,
			})
		} else if handler, ok := methodReflect.Interface().(func() (gin.HandlerFunc, string, string)); ok {
			handlerFunc, path, method := handler()
			controller.Actions = append(controller.Actions, &Action{
				Name:    valueOf.Type().Method(i).Name,
				Handler: handlerFunc,
				paths:   strings.Split(path, ","),
				methods: strings.Split(method, ","),
			})
		}
	}
	return controller
}

func (that *Controller) Path() string {
	if that.Value != "" {
		return utils.Camel2Snake(that.Value)
	}
	return utils.Camel2Snake(that.Name)
}

type Action struct {
	Name    string //方法名
	Handler gin.HandlerFunc
	paths   []string
	methods []string
}

func (that *Action) Paths() []string {
	if len(that.paths) > 0 {
		return that.paths
	}
	return []string{utils.Camel2Snake(that.Name)}
}

func (that *Action) Methods(defMethods []string) []string {
	if len(that.methods) > 0 {
		defMethods = that.methods
	}
	for i, s := range defMethods {
		if strings.ToUpper(s) == "ANY" {
			return []string{"Any"}
		}
		defMethods[i] = strings.ToUpper(s)
	}
	return defMethods
}
