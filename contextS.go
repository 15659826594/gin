package gin

import (
	"net/http"
	"net/url"
	"strconv"
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
	Register(string, string, string, string, map[string]string) bool
	Login(string, string) bool
	Logout()
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
	Match() bool
	Keeptime(int)
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

	handlerName string

	IController

	Auth IAuth
}

/*
SetHandlerName
设置action句柄
*/
func (c *Context) SetHandlerName(name string) {
	c.handlerName = name
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
	for i, val := range args {
		switch i {
		case 0:
			msg = val.(string)
		case 1:
			if v, ok := val.(H); ok {
				data = v
			} else {
				data = val.(map[string]any)
			}
		case 2:
			code = val.(int)
		case 3:
			types = val.(string)
		case 4:
			header = val.(map[string]string)
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
	for i, val := range args {
		switch i {
		case 0:
			msg = val.(string)
		case 1:
			if v, ok := val.(H); ok {
				data = v
			} else {
				data = val.(map[string]any)
			}
		case 2:
			code = val.(int)
		case 3:
			types = val.(string)
		case 4:
			header = val.(map[string]string)
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
