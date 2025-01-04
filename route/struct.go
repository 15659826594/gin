package route

import (
	"errors"
	"gin"
	. "gin/annotation"
	"gin/env"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	sep      = string(filepath.Separator)
	basePath = env.Get("ROOT_PATH")
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
	return gin.Camel2Snake(that.Name)
}

type Module struct {
	Name        string
	Controllers []*Controller
}

func (that *Module) Path() string {
	return gin.Camel2Snake(that.Name)
}

type Controller struct {
	IController
	Name    string //控制器名
	Value   string //路由
	Actions []*Action
}

type IController interface {
	Initialize(*gin.Context)
	GetValue() string
	GetNoNeedLogin() []string
	GetNoNeedRight() []string
	BeforeAction() []gin.HandlerFunc
	Exception() gin.HandlerFunc
}

func NewController(obj any) *Controller {
	object, ok := obj.(IController)
	if !ok {
		return (*Controller)(nil)
	}
	typeOf := reflect.TypeOf(obj)
	valueOf := reflect.ValueOf(obj)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	controller := &Controller{
		Name:  typeOf.Name(),
		Value: object.GetValue(),
	}

	controller.IController = object

	controller.Actions = collectAction(valueOf)
	return controller
}

func (that *Controller) Path() string {
	return gin.Camel2Snake(that.Name)
}

func (that *Controller) IsNil() bool {
	return that.Name == ""
}

type Action struct {
	Name        string //方法名
	Handler     gin.HandlerFunc
	Annotations []Annotation
}

func (that *Action) Path() string {
	return gin.Camel2Snake(that.Name)
}

// Mapping (路由组, 初始化方法, 默认请求)
func (that *Action) Mapping(group *gin.RouterGroup, init gin.HandlerFunc, def []string) ([]string, string, error) {
	var ginHandlers []gin.HandlerFunc
	if init != nil {
		ginHandlers = append(ginHandlers, init)
	}
	ginHandlers = append(ginHandlers, that.Handler)
	for _, annotation := range that.Annotations {
		fn := Mapping.Get(annotation.Name)
		if fn != nil {
			return fn()(group, ginHandlers, annotation.Attributes, that.Path())
		}
	}
	i, s, err := RequestMapping(def)(group, ginHandlers, nil, that.Path())
	if err != nil {
		return nil, "", err
	}
	return i, s, errors.New("default")
}

// collectAction 提取结构体的方法
func collectAction(cValueOf reflect.Value) []*Action {
	var actions []*Action
	for i, lens := 0, cValueOf.NumMethod(); i < lens; i++ {
		methodReflect := cValueOf.Method(i)
		fn, ok := methodReflect.Interface().(func(*gin.Context))
		if ok {
			actions = append(actions, &Action{
				Name:    cValueOf.Type().Method(i).Name,
				Handler: fn,
			})
		}
	}
	return actions
}

type Annotation struct {
	Name        string
	Attributes  map[string]string
	Description []string
}
