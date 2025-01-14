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
	return version.GetModule(moduleName)
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

func (that *Version) GetModule(name string) *Module {
	for _, v := range that.Modules {
		if v.Name == name {
			return v
		}
	}
	module := &Module{
		Name:        name,
		Controllers: []*Controller{},
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
	Name        string
	Controllers []*Controller
}

func (that *Module) Path() string {
	return utils.Camel2Snake(that.Name)
}

type Controller struct {
	IController
	Name    string //控制器名
	Value   string //路由
	Actions []*Action
}

type IController interface {
	Initialize(*gin.Context)
	Value() string
	NoNeedLogin() []string
	NoNeedRight() []string
	ResponseType() string
	BeforeAction() []gin.HandlerFunc
	Exception() gin.HandlerFunc
}

func NewController(obj any) *Controller {
	object, ok := obj.(IController)
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
		object, typeOf.Name(), object.Value(), []*Action{},
	}
	//提取结构体中的方法
	for i, lens := 0, valueOf.NumMethod(); i < lens; i++ {
		methodReflect := valueOf.Method(i)
		if fn, isgc := methodReflect.Interface().(func(*gin.Context)); isgc {
			controller.Actions = append(controller.Actions, &Action{
				Name:    valueOf.Type().Method(i).Name,
				Handler: fn,
			})
		}
	}

	return controller
}

func (that *Controller) Path() string {
	return utils.Camel2Snake(that.Name)
}

type Action struct {
	Name        string //方法名
	Handler     gin.HandlerFunc
	Annotations []Annotation
}

func (that *Action) Path() string {
	return utils.Camel2Snake(that.Name)
}

type Annotation struct {
	Name        string
	Attributes  map[string]string
	Description []string
}
