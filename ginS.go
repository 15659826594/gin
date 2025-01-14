package gin

import (
	"errors"
	"gin/config"
	"gin/src/html/template"
	"gin/utils"
	"strings"
)

/************************************/
/**********  	 gin.go 	 ********/
/************************************/

// LoadHTMLFolder loads HTML files identified folder
// and associates the result with HTML renderer.
func (engine *Engine) LoadHTMLFolder(path string, rename func(name string) string) {
	left := engine.delims.Left
	right := engine.delims.Right
	templ := template.Must(template.WrapT(template.New("").Delims(left, right).Funcs(engine.FuncMap)).ParseFolder(path, rename))

	if IsDebugging() {
		debugPrintLoadTemplate(templ.UnWrapT())
	}

	engine.SetHTMLTemplate(templ.UnWrapT())
}

/************************************/
/**********    context.go	 ********/
/************************************/

func (c *Context) IsAjax() bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (c *Context) getResponseType() string {
	var ret any
	if c.IsAjax() {
		ret = config.Get("default_ajax_return")
	} else {
		ret = config.Get("default_return_type")
	}
	if v, ok := ret.(string); ok {
		return v
	}
	return "html"
}

// Langset 当前的语言
func (c *Context) Langset(args ...string) string {
	language := "zh-cn"
	for index, arg := range args {
		switch index {
		case 0:
			language = arg
		}
	}
	return language + "." + strings.Join(strings.Split(utils.Camel2Snake(c.GetString("Request.URL")), "/"), ".")
}

// AcceptLang 当前网页语言
func (c *Context) AcceptLang() string {
	return c.Request.Header.Get("Accept-Language")
}

type RequsetS struct {
	context *Context
}

// Requests 当前的操作名
func (c *Context) Requests() *RequsetS {
	return &RequsetS{context: c}
}

// Action 当前的操作名
func (c *RequsetS) Action(toSnake bool) string {
	arr := strings.Split(c.context.GetString("Request.URL"), "/")
	if toSnake {
		return utils.Camel2Snake(arr[2])
	}
	return arr[2]
}

// Controller 当前的控制器名
func (c *RequsetS) Controller(toSnake bool) string {
	arr := strings.Split(c.context.GetString("Request.URL"), "/")
	if toSnake {
		return utils.Camel2Snake(arr[1])
	}
	return arr[1]
}

// Module 获取模块名
func (c *RequsetS) Module(toSnake bool) string {
	arr := strings.Split(c.context.GetString("Request.URL"), "/")
	if toSnake {
		return utils.Camel2Snake(arr[0])
	}
	return arr[0]
}

// Server 获取server参数
func (c *RequsetS) Server(name string, args ...string) string {
	value := c.context.GetHeader(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Request 获取Requests变量
func (c *RequsetS) Request(name string, args ...string) string {
	value := c.context.PostForm(name)
	if value != "" {
		return value
	}
	value = c.context.Query(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Post 获取post参数
func (c *RequsetS) Post(name string, args ...string) string {
	if len(args) > 0 {
		return c.context.DefaultPostForm(name, args[0])
	}
	return c.context.PostForm(name)
}

// Cookie 获取Cookie
func (c *RequsetS) Cookie(name string) string {
	if cookie, err := c.context.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

/*Url 生成
 * @param string        $url 路由地址(可以是相对路径)
 * @param string|array  $vars 变量
 * @param bool|string   $base base路径
 * @return string
 */
func (c *RequsetS) Url(url string, vars any, base any) string {
	var baseUrl string
	switch tmp := base.(type) {
	case string:
		baseUrl = tmp
	case bool:
		if tmp {
			baseUrl = c.context.Request.URL.String()
		}
	}
	return utils.URL(url, vars, baseUrl)
}

// Exit 中断执行
func Exit() {
	panic(errors.New("exit"))
}

// ExceptionHandle 正常中断后续请求
func ExceptionHandle() HandlerFunc {
	return func(c *Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if err := rec.(error); err.Error() == "exit" {
				c.Abort()
			} else {
				panic(rec)
			}
		}()
		c.Next()
	}
}

/************************************/
/**********response_writer.go********/
/************************************/

//// Success 操作成功返回的数据
///**
// * @param string $msg    提示信息
// * @param mixed  $data   要返回的数据
// * @param int    $code   错误码，默认为1
// * @param string $type   输出类型
// * @param array  $header 发送的 Header 信息
// */
//func (c *Context) Success(args ...any) {
//	resp := new(Result)
//	for index, arg := range args {
//		switch index {
//		case 0:
//			resp.Msg = arg.(string)
//		case 1:
//			resp.Data = arg
//		case 2:
//			if arg != nil {
//				resp.Code = arg.(int)
//			} else {
//				resp.Code = 1
//			}
//		case 3:
//			resp.Type = arg.(string)
//		case 4:
//			resp.Header = arg.(map[string]string)
//		}
//	}
//	c.Result(resp.Msg, resp.Data, resp.Code, resp.Type, resp.Header)
//}
//
///*Fail 操作失败返回的数据
// * @param string $msg    提示信息
// * @param mixed  $data   要返回的数据
// * @param int    $code   错误码，默认为0
// * @param string $type   输出类型
// * @param array  $header 发送的 Header 信息
// */
//func (c *Context) Fail(args ...any) {
//	resp := new(Result)
//	for index, arg := range args {
//		switch index {
//		case 0:
//			resp.Msg = arg.(string)
//		case 1:
//			resp.Data = arg
//		case 2:
//			if arg != nil {
//				resp.Code = arg.(int)
//			} else {
//				resp.Code = 0
//			}
//		case 3:
//			resp.Type = arg.(string)
//		case 4:
//			resp.Header = arg.(map[string]string)
//		}
//	}
//	c.Result(resp.Msg, resp.Data, resp.Code, resp.Type, resp.Header)
//}
//
//// Result 返回封装后的 API 数据到客户端
//func (c *Context) Result(msg string, data any, code int, types string, header map[string]string) {
//	result := Result{
//		Code: code,
//		Msg:  msg,
//		Time: time.Now().Unix(),
//		Data: data,
//	}
//	// 如果未设置类型则使用默认类型判断
//	if types == "" {
//		if value, ok := c.Request.Header["Response-Type"]; ok {
//			types = value[0]
//		}
//	}
//
//	if statusCode, ok := header["statuscode"]; ok {
//		code, _ = strconv.Atoi(statusCode)
//		delete(header, "statuscode")
//	} else {
//		//未设置状态码,根据code值判断
//		if code >= 1000 || code < 200 {
//			code = 200
//		}
//	}
//	resp := Create(result, types, code).Header(header)
//	resp.Send(c)
//	Exit()
//}
//
///*Fetch
// * 解析和获取模板内容 用于输出
// * @param string    $template 模板文件名或者内容
// * @param array     $vars     模板输出变量
// */
//func (c *Context) Fetch(args ...any) {
//	templ := defaultTemplName(2)
//	temp := &(struct {
//		template string
//		vars     map[string]any
//		replace  map[string]string
//		config   map[string]string
//	}{
//		template: templ,
//		vars:     map[string]any{},
//		replace:  map[string]string{},
//		config:   map[string]string{},
//	})
//	for index, arg := range args {
//		switch index {
//		case 0:
//			if v, ok := arg.(string); ok {
//				if v == "" {
//					continue
//				} else if strings.HasPrefix(v, "/") {
//					temp.template = path.Clean(v)[1:]
//				} else if !strings.HasPrefix(v, ".") {
//					temp.template = filepath.Dir(temp.template) + "/" + v
//				} else if strings.HasPrefix(v, "./") {
//					temp.template = filepath.Dir(temp.template) + "/" + filepath.Base(v)
//				} else {
//					temp.template = filepath.Clean(temp.template + "/" + v)
//				}
//			}
//		case 1:
//			if v, ok := arg.(map[string]any); ok {
//				temp.vars = v
//			}
//		case 2:
//			if v, ok := arg.(map[string]string); ok {
//				temp.config = v
//			}
//		case 3:
//			if v, ok := arg.(map[string]string); ok {
//				temp.replace = v
//			}
//		}
//	}
//	//如果没有后缀补上.html
//	if !slices.Contains([]string{".html", ".tpl", ".tmpl"}, filepath.Ext(temp.template)) {
//		temp.template += ".html"
//	}
//	assignData := c.GetStringMap("GlobalAssign")
//	if assignData != nil {
//		for k, v := range assignData {
//			temp.vars[k] = v
//		}
//	}
//	temp.template = filepath.ToSlash(temp.template)
//	c.HTML(http.StatusOK, temp.template, (H)(temp.vars))
//	Exit()
//}
//
///*Assign
// * 模板变量赋值
// * @access protected
// * @param  mixed $name  要显示的模板变量
// * @param  mixed $value 变量的值
// * @return $this
// */
//func (c *Context) Assign(name string, value any) *Context {
//	assignData := c.GetStringMap("GlobalAssign")
//	if assignData == nil {
//		assignData = make(map[string]interface{})
//	}
//	assignData[name] = value
//	c.Set("GlobalAssign", assignData)
//	return c
//}
