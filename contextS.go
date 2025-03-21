package gin

import (
	"fmt"
	"gin/utils"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IController struct {
	ISuccess
	IFail
	IResult
	IResponseType
	IAssign
	IFetch
}

type ISuccess interface {
	Success(*Context, ...any)
}

type IFail interface {
	Fail(*Context, ...any)
}

type IResult interface {
	Result(*Context, any, int, string, string, map[string]string)
}

type IResponseType interface {
	ResponseType(*Context) string
}

type IAssign interface {
	Assign(*Context, string, ...any)
}

type IFetch interface {
	Fetch(*Context, ...any)
}

type IAuth interface {
	GetUser() any
	Init(string) bool
	Register(string, string, string, string, map[string]any) bool
	Login(string, string) bool
	Logout() bool
	Changepwd(string, string, bool) bool
	Direct(int) bool
	Check(string, string) bool
	IsLogin() bool
	GetToken() string
	GetUserinfo() map[string]any
	GetRuleList() any
	GetRequestUri() string
	SetRequestUri(string)
	GetAllowFields() []string
	SetAllowFields([]string)
	Delete(int) bool
	GetEncryptPassword(string, string) string
	Match([]string) bool
	Keeptime(int64)
	SetError(error) IAuth
	GetError() string
}

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8
	fullPath string

	engine       *Engine
	params       *Params
	skippedNodes *[]skippedNode

	// This mutex protects Keys map.
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache caches the query result from c.Param.URL.Query().
	queryCache url.Values

	// formCache caches c.Param.PostForm, which contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite

	IController //来源于控制器

	Auth IAuth

	*RequestS
}

type RequestS struct {
	*Context
	mu             sync.RWMutex
	keys           map[string]any //通过assgin赋值
	HandlerName    string
	method         string
	versionname    string
	modulename     string
	controllername string
	actionname     string
}

func (r *RequestS) Arg(name string) string {
	var args = map[string]any{}
	for key, val := range r.Request.URL.Query() {
		if len(val) == 1 {
			args[key] = val[0]
		} else {
			args[key] = val
		}
	}
	err := r.Request.ParseMultipartForm(32 << 20)
	if err == nil {
		postForm := r.Request.PostForm
		for key, val := range postForm {
			if len(val) == 1 {
				args[key] = val[0]
			} else {
				args[key] = val
			}
		}
	}
	if _, exist := args[name]; !exist {
		return ""
	}
	if val, ok := args[name].(string); ok {
		return val
	} else if val, ok := args[name].([]string); ok {
		return strings.Join(val, ",")
	}
	return "is : " + reflect.ValueOf(args[name]).Kind().String()
}

func (r *RequestS) Args() map[string]any {
	var args = map[string]any{}
	for key, val := range r.Request.URL.Query() {
		if len(val) == 1 {
			args[key] = val[0]
		} else {
			args[key] = val
		}
	}
	err := r.Request.ParseMultipartForm(32 << 20)
	if err == nil {
		postForm := r.Request.PostForm
		for key, val := range postForm {
			if len(val) == 1 {
				args[key] = val[0]
			} else {
				args[key] = val
			}
		}
	}
	return args
}

func (r *RequestS) Version() string {
	if r.versionname == "internal" {
		return ""
	}
	return r.versionname
}

func (r *RequestS) Path() string {
	return strings.TrimPrefix(fmt.Sprintf("%s/%s/%s/%s", r.Version(), r.Module(), r.Controller(), r.Action()), "/")
}

func (r *RequestS) Module() string {
	return r.modulename
}

func (r *RequestS) Controller() string {
	return r.controllername
}

func (r *RequestS) Action() string {
	return r.actionname
}

func (r *RequestS) Langset(args ...string) string {
	language := r.Request.Header.Get("Accept-Language")
	if language == "" && len(args) > 0 {
		return args[0]
	}
	return language
}

func (r *RequestS) Url(url string, vars any, base any) string {
	var baseUrl string
	switch tmp := base.(type) {
	case string:
		baseUrl = tmp
	case bool:
		if tmp {
			baseUrl = r.Request.URL.String()
		}
	}
	return utils.URL(url, vars, baseUrl)
}

func (r *RequestS) IsAjax() bool {
	return r.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (r *RequestS) IsGet() bool {
	return r.method == http.MethodGet
}

func (r *RequestS) IsPost() bool {
	return r.method == http.MethodPost
}

func (r *RequestS) IsPut() bool {
	return r.method == http.MethodPut
}

func (r *RequestS) IsDelete() bool {
	return r.method == http.MethodDelete
}

func (r *RequestS) IsHead() bool {
	return r.method == http.MethodHead
}

func (r *RequestS) IsPatch() bool {
	return r.method == http.MethodPatch
}

func (r *RequestS) IsOptions() bool {
	return r.method == http.MethodOptions
}

// SetContextS 设置action句柄
func (c *Context) SetContextS(HandlerName string) {
	arr := strings.Split(filepath.ToSlash(utils.CaseSnake(HandlerName)), "/")
	l := len(arr)
	arr = append(arr[:l-1], strings.Split(arr[l-1], ".")[1:]...)
	l = len(arr)
	c.RequestS = &RequestS{
		Context:        c,
		keys:           map[string]any{},
		HandlerName:    HandlerName,
		method:         c.Request.Method,
		versionname:    arr[1],
		modulename:     arr[2],
		controllername: arr[l-2],
		actionname:     arr[l-1],
	}
}

/*Success
 * 操作成功返回的数据
 * @param string $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为1
 * @param string $type   输出类型
 * @param array  $header 发送的 Header 信息
 */
func (c *Context) Success(args ...any) {
	if c.IController.ISuccess != nil {
		c.IController.Success(c, args...)
		return
	}
	var msg, _type string
	var data map[string]any
	var header map[string]string
	code := 1
	for i, arg := range args {
		switch i {
		case 0:
			if val, ok := arg.(string); ok {
				msg = val
			}
		case 1:
			if arg == nil {
				continue
			} else if v, ok := arg.(H); ok {
				data = v
			} else if v, ok := arg.(map[string]any); ok {
				data = v
			}
		case 2:
			if val, ok := arg.(int); ok {
				code = val
			}
		case 3:
			if val, ok := arg.(string); ok {
				_type = val
			}
		case 4:
			if val, ok := arg.(map[string]string); ok {
				header = val
			}
		}
	}
	c.Result(data, code, msg, _type, header)
}

/*Fail
 * 操作失败返回的数据
 * @param string $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为0
 * @param string $type   输出类型
 * @param array  $header 发送的 Header 信息
 */
func (c *Context) Fail(args ...any) {
	if c.IController.IFail != nil {
		c.IController.Fail(c, args...)
		return
	}
	var msg, _type string
	var code int
	var data map[string]any
	var header map[string]string
	for i, arg := range args {
		switch i {
		case 0:
			if val, ok := arg.(string); ok {
				msg = val
			}
		case 1:
			if arg == nil {
				continue
			} else if v, ok := arg.(H); ok {
				data = v
			} else if v, ok := arg.(map[string]any); ok {
				data = v
			}
		case 2:
			if val, ok := arg.(int); ok {
				code = val
			}
		case 3:
			if val, ok := arg.(string); ok {
				_type = val
			}
		case 4:
			if val, ok := arg.(map[string]string); ok {
				header = val
			}
		}
	}
	c.Result(data, code, msg, _type, header)
}

/*Result
 * 返回封装后的 API 数据到客户端
 * @access protected
 * @param mixed  $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为0
 * @param string $type   输出类型，支持json/xml/jsonp
 * @param array  $header 发送的 Header 信息
 * @return void
 * @throws HttpResponseException
 */
func (c *Context) Result(data any, code int, msg string, _type string, header map[string]string) {
	if c.IController.IResult != nil {
		c.IController.Result(c, data, code, msg, _type, header)
		return
	}
	if header == nil {
		header = make(map[string]string)
	}

	var statuscode int
	result := H{
		"code": code,
		"msg":  msg,
		"time": strconv.FormatInt(time.Now().Unix(), 10),
		"data": data,
	}
	// 如果未设置类型则使用默认类型判断
	if _type == "" && c.IController.IResponseType != nil {
		_type = c.IController.ResponseType(c)
	}
	if sc, ok := header["statuscode"]; ok {
		statuscode, _ = strconv.Atoi(sc)
	} else {
		//未设置状态码,根据code值判断
		if code >= 1000 || code < 200 {
			statuscode = http.StatusOK
		} else {
			statuscode = code
		}
	}
	for key, val := range header {
		c.Header(key, val)
	}
	switch _type {
	case "json":
		c.JSON(statuscode, result)
	case "jsonp":
		c.JSONP(statuscode, result)
	case "xml":
		c.XML(statuscode, result)
	case "text":
		c.String(statuscode, msg)
	default:
		c.JSON(statuscode, result)
	}
	panic("Abort")
}

/*Assign
 * 模板变量赋值
 * @access protected
 * @param  mixed $name  要显示的模板变量
 * @param  mixed $value 变量的值
 * @return $this
 */
func (c *Context) Assign(name string, value any) *Context {
	c.RequestS.mu.Lock()
	defer c.RequestS.mu.Unlock()
	c.RequestS.keys[name] = value
	return c
}

/*Fetch
 * 解析和获取模板内容 用于输出
 * @param string    $template 模板文件名或者内容
 * @param array     $vars     模板输出变量
 * @param array     $replace 替换内容
 * @param array     $config     模板参数
 * @param bool      $renderContent     是否渲染内容
 * @return string
 * @throws Exception
 */
func (c *Context) Fetch(args ...any) {
	if c.IController.IFetch != nil {
		c.IController.Fetch(c, args...)
		return
	}
	var template = strings.TrimPrefix(fmt.Sprintf("%s/%s/view/%s/%s", c.Version(), c.Module(), c.Controller(), c.Action()), "/")
	var vars H
	//var replace = make(map[string]string)
	//var config = make(map[string]any)
	//var renderContent bool
	for i, arg := range args {
		switch i {
		case 0:
			if val, ok := arg.(string); ok {
				if val == "" {
					continue
				} else if strings.HasPrefix(val, "/") {
					template = val
					continue
				} else if !strings.HasPrefix(val, ".") {
					val = "./" + val
				}
				template = filepath.Clean(template + "/../" + val)
			}
		case 1:
			if arg == nil {
				continue
			} else if val, ok := arg.(H); ok {
				vars = val
			} else if val, ok := arg.(map[string]any); ok {
				vars = val
			}
			//case 2:
			//	if val, ok := arg.(map[string]string); ok {
			//		replace = val
			//	}
			//case 3:
			//	if val, ok := arg.(map[string]any); ok {
			//		config = val
			//	}
			//case 4:
			//	if val, ok := arg.(bool); ok {
			//		renderContent = val
			//	}
		}
	}
	if filepath.Ext(template) == "" {
		template += ".html"
	}
	c.RequestS.mu.Lock()
	c.RequestS.keys["Accept-Language"] = c.Langset()
	c.RequestS.keys["Requests-Url"] = c.Path()
	for s, a := range vars { //Assign参数合并
		c.RequestS.keys[s] = a
	}
	c.RequestS.mu.Unlock()
	c.HTML(http.StatusOK, strings.TrimPrefix(template, "/"), c.RequestS.keys)
	panic("Abort")
}
