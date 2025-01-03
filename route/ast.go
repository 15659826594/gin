package route

import (
	"bufio"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

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
