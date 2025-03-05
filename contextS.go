package gin

import (
	"gin/utils"
	"net/http"
	"net/url"
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

type RequestS struct {
	ctx            *Context
	method         string
	HandlerName    string
	modulename     string
	controllername string
	actionname     string
}

func (r *RequestS) Params() map[string]any {
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

func (r *RequestS) Param(key string) (value string) {
	if val, exist := r.Params()[key]; exist {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (r *RequestS) DefaultParam(key, defaultValue string) string {
	if val, exist := r.Params()[key]; exist {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (r *RequestS) Module(snake bool) string {
	s := r.modulename
	if s == "" {
		arr := strings.Split(r.HandlerName, "/")
		if len(arr) >= 3 {
			r.modulename = arr[2]
		}
		s = r.modulename
	}
	if snake {
		return utils.Camel2Snake(s)
	}
	return s
}

func (r *RequestS) Controller(snake bool) string {
	s := r.controllername
	if s == "" {
		arr := strings.Split(r.HandlerName, ".")
		if len(arr) >= 2 {
			r.controllername = arr[1]
		}
		s = r.controllername
	}
	if snake {
		return utils.Camel2Snake(s)
	}
	return s
}

func (r *RequestS) Action(snake bool) string {
	s := r.actionname
	if s == "" {
		arr := strings.Split(r.HandlerName, ".")
		if len(arr) >= 3 {
			r.actionname = arr[2]
		}
		s = r.actionname
	}
	if snake {
		return utils.Camel2Snake(s)
	}
	return s
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

	IController

	Auth IAuth

	RequestS
}

/*
SetHandlerName
设置action句柄
*/
func (c *Context) SetHandlerName(name string) {
	c.RequestS = RequestS{
		ctx:         c,
		method:      c.Request.Method,
		HandlerName: name,
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
	var msg, types string
	var data map[string]any
	var header map[string]string
	code := 1
	for i, arg := range args {
		switch i {
		case 0:
			msg = arg.(string)
		case 1:
			if arg == nil {
				continue
			} else if v, ok := arg.(H); ok {
				data = v
			} else if v, ok := arg.(map[string]any); ok {
				data = v
			}
		case 2:
			code = arg.(int)
		case 3:
			types = arg.(string)
		case 4:
			header = arg.(map[string]string)
		}
	}
	c.Result(data, code, msg, types, header)
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
	var msg, types string
	var code int
	var data map[string]any
	var header map[string]string
	for i, arg := range args {
		switch i {
		case 0:
			msg = arg.(string)
		case 1:
			if arg == nil {
				continue
			} else if v, ok := arg.(H); ok {
				data = v
			} else if v, ok := arg.(map[string]any); ok {
				data = v
			}
		case 2:
			code = arg.(int)
		case 3:
			types = arg.(string)
		case 4:
			header = arg.(map[string]string)
		}
	}
	c.Result(data, code, msg, types, header)
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
func (c *Context) Result(data any, code int, msg string, types string, header map[string]string) {
	if c.IController.IResult != nil {
		c.IController.Result(c, data, code, msg, types, header)
		return
	}

	var statuscode int
	result := H{
		"code": code,
		"msg":  msg,
		"time": time.Now().Unix(),
		"data": data,
	}
	// 如果未设置类型则使用默认类型判断
	if types == "" && c.IController.IResponseType != nil {
		types = c.IController.ResponseType(c)
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
	switch types {
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
