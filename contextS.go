package gin

import (
	"fmt"
	"gin/utils"
	"net/http"
	"net/url"
	"path/filepath"
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

	Auth IAuth

	IController

	*ContextS
}

type ContextS struct {
	ctx            *Context
	tmplDataMu     sync.RWMutex
	tmplData       map[string]any
	method         string
	HandlerName    string
	versionname    string
	modulename     string
	controllername string
	actionname     string
}

func (r *ContextS) Params() map[string]any {
	arr := make(map[string]any)
	query := r.ctx.Request.URL.Query()
	for k, v := range query {
		arr[k] = v
	}
	err := r.ctx.Request.ParseMultipartForm(32 << 20)
	if err == nil {
		postForm := r.ctx.Request.PostForm
		for k, v := range postForm {
			arr[k] = v
		}
	}
	return arr
}

func (r *ContextS) Param(key string) (value string) {
	if val, exist := r.Params()[key]; exist {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (r *ContextS) DefaultParam(key, defaultValue string) string {
	if val, exist := r.Params()[key]; exist {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (r *ContextS) Version() string {
	if r.versionname == "application" {
		return ""
	}
	return r.versionname
}

func (r *ContextS) Path() string {
	return strings.TrimPrefix(fmt.Sprintf("%s/%s/%s/%s", r.Version(), r.Module(), r.Controller(), r.Action()), "/")
}

func (r *ContextS) Module() string {
	return r.modulename
}

func (r *ContextS) Controller() string {
	return r.controllername
}

func (r *ContextS) Action() string {
	return r.actionname
}

func (r *ContextS) Langset(args ...string) string {
	language := r.ctx.Request.Header.Get("Accept-Language")
	if language == "" && len(args) > 0 {
		return args[0]
	}
	return language
}

func (r *ContextS) Url(url string, vars any, base any) string {
	var baseUrl string
	switch tmp := base.(type) {
	case string:
		baseUrl = tmp
	case bool:
		if tmp {
			baseUrl = r.ctx.Request.URL.String()
		}
	}
	return utils.URL(url, vars, baseUrl)
}

func (r *ContextS) IsAjax() bool {
	return r.ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (r *ContextS) IsGet() bool {
	return r.method == http.MethodGet
}

func (r *ContextS) IsPost() bool {
	return r.method == http.MethodPost
}

func (r *ContextS) IsPut() bool {
	return r.method == http.MethodPut
}

func (r *ContextS) IsDelete() bool {
	return r.method == http.MethodDelete
}

func (r *ContextS) IsHead() bool {
	return r.method == http.MethodHead
}

func (r *ContextS) IsPatch() bool {
	return r.method == http.MethodPatch
}

func (r *ContextS) IsOptions() bool {
	return r.method == http.MethodOptions
}

// SetContextS 设置action句柄
func (c *Context) SetContextS(HandlerName string) {
	arr := strings.Split(filepath.ToSlash(utils.Camel2Snake(HandlerName)), "/")
	l := len(arr)
	arr = append(arr[:l-1], strings.Split(arr[l-1], ".")[1:]...)
	l = len(arr)
	c.ContextS = &ContextS{
		ctx:            c,
		method:         c.Request.Method,
		tmplData:       map[string]any{},
		HandlerName:    HandlerName,
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
		"time": time.Now().Unix(),
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
	c.tmplDataMu.Lock()
	defer c.tmplDataMu.Unlock()
	c.tmplData[name] = value
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
	//fmt.Println(replace, config, renderContent)
	c.tmplData["Accept-Language"] = c.Langset()
	c.tmplData["Request-Url"] = c.Path()
	//Assign参数合并
	for s, a := range vars {
		c.tmplData[s] = a
	}
	c.HTML(http.StatusOK, strings.TrimPrefix(template, "/"), c.tmplData)
	panic("Abort")
}
