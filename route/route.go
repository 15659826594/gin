package route

import (
	"bufio"
	"errors"
	"fmt"
	"gin"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
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
						methodName, uri, err := action.Mapping(level3, controller.NeedAuth(action.Name), defaultMethod)
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

type AstFile struct {
	*ast.File
}

func parseFile(file string) AstFile {
	fileSet := token.NewFileSet()
	astParser, _ := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	return AstFile{astParser}
}

// GetComments 获取注释
func (that AstFile) GetComments(structName string, operationName string) ([]Annotation, error) {
	for _, decl := range that.Decls {
		// 检查是否是函数声明
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				recType := fn.Recv.List[0].Type
				// 结构体名称
				var sName string
				switch recType.(type) {
				case *ast.StarExpr: //指针方法
					sName = recType.(*ast.StarExpr).X.(*ast.Ident).Name
				case *ast.Ident: //结构体方法
					sName = recType.(*ast.Ident).Name
				}
				//查找指定的方法
				if sName != structName || fn.Name.Name != operationName {
					continue
				}
				if fn.Doc != nil {
					comments := that.ScanComments(fn.Doc.Text())
					if comments != nil {
						return comments, nil
					}
				}
				return nil, errors.New("not comments")
			}
		}
	}
	return nil, errors.New("not found")
}

// ScanComments 逐行扫描注释
func (that AstFile) ScanComments(text string) []Annotation {
	reader := strings.NewReader(text)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var annotations []string
	//逐行扫描
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "@") {
			annotations = append(annotations, line[1:])
		}
	}
	return that.ResolveAnnotation(annotations)
}

var bracketReg = regexp.MustCompile(`\((.*?)\)`)                      // 匹配括号内的内容
var annotationNameReg = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]+`) // 匹配方法名
var commaReg = regexp.MustCompile(`(?:[^,"]+|"[^"]*")+`)              // 逗号拆分参数, 排除双引号内

// ResolveAnnotation 解析注解
func (that AstFile) ResolveAnnotation(collection []string) []Annotation {
	if collection == nil {
		return nil
	}

	var annotations []Annotation
	for _, text := range collection {
		//注解的方法名
		annotationName := annotationNameReg.FindString(text)
		if annotationName != "" {
			text = strings.TrimPrefix(text, annotationName)
			annotation := Annotation{
				Name:        annotationName,
				Attributes:  make(map[string]string),
				Description: make([]string, 0),
			}
			bracket := bracketReg.FindString(text)
			//匹配到空格
			if bracket != "" {
				text = strings.Replace(text, bracket, "", 1)
				arguments := bracket[1 : len(bracket)-1]
				//逗号拆分参数, 排除双引号内
				args := commaReg.FindAllString(arguments, -1)
				// @Test(value="v",method="GET,POST")
				for _, arg := range args {
					mapping := strings.Split(arg, "=")
					//@Test("index")
					if len(mapping) == 1 {
						annotation.Attributes[""] = strings.TrimSpace(mapping[0])
					}
					//@Test(value="index")
					if len(mapping) == 2 {
						annotation.Attributes[strings.TrimSpace(mapping[0])] = strings.TrimSpace(mapping[1])
					}
				}
			}
			//解析类型
			for _, v := range strings.Split(text, " ") {
				if v != "" {
					annotation.Description = append(annotation.Description, v)
				}
			}
			annotations = append(annotations, annotation)
		}
	}
	return annotations
}
