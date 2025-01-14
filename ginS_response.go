package gin

import (
	"fmt"
	"gin/config"
	"gin/utils"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (c *Context) Success(args ...any) {}
func (c *Context) Fetch(args ...any)   {}

const JsonUnescapedUnicode = 256

type Result struct {
	Code   int               `json:"code"`
	Msg    string            `json:"msg"`
	Time   int64             `json:"time"`
	Data   any               `json:"data"`
	Url    string            `json:"-"`
	Wait   int               `json:"-"`
	Type   string            `json:"-"`
	Header map[string]string `json:"-"`
}

type Response struct {
	// 原始数据
	data *Result
	// 当前的contentType
	contentType string
	// 字符集
	charset string
	//状态
	code int
	// 输出参数
	options map[string]string
	// header参数
	header  map[string]string
	content *Context
}

func Construct(context *Context, data *Result, code int, header map[string]string, options map[string]string) *Response {
	that := &Response{
		content:     context,
		contentType: "text/html",
		charset:     "utf-8",
		code:        200,
		header:      map[string]string{},
		options:     map[string]string{},
	}
	that.Data(data)
	if options != nil {
		for k, v := range options {
			that.options[k] = v
		}
	}
	that.ContentType(that.contentType, that.charset)
	for k, v := range header {
		that.header[k] = v
	}
	that.code = code
	return that
}

func Create(context *Context, data *Result, types string, code int, args ...map[string]string) *Response {
	header := make(map[string]string)
	options := make(map[string]string)
	for i, arg := range args {
		switch i {
		case 0:
			header = arg
		case 1:
			options = arg
		}
	}
	resp := Construct(context, data, code, header, options)
	switch types {
	case "json":
		resp.options["json_encode_param"] = fmt.Sprintf("%d", JsonUnescapedUnicode)
		resp.contentType = "application/json"
	case "jsonp":
		resp.options["var_jsonp_handler"] = "callback"
		resp.options["default_jsonp_handler"] = "jsonpReturn"
		resp.options["json_encode_param"] = fmt.Sprintf("%d", JsonUnescapedUnicode)
		resp.contentType = "application/javascript"
	case "xml":
		resp.options["root_node"] = "think"
		resp.options["root_attr"] = ""
		resp.options["item_node"] = "item"
		resp.options["item_key"] = "id"
		resp.options["encoding"] = "utf-8"
		resp.contentType = "text/xml"
	}
	return resp
}

// Render 处理数据
func (that *Response) Render() {
	fmt.Println(that.GetData())
}

// Options 输出的参数
func (that *Response) Options(options map[string]string) *Response {
	for k, v := range options {
		that.options[k] = v
	}
	return that
}

// Data 输出数据设置
func (that *Response) Data(data *Result) *Response {
	that.data = data
	return that
}

// Header 设置响应头
// aram string|array $name  参数名
// param string       $value 参数值
func (that *Response) Header(args ...any) *Response {
	if len(args) > 0 {
		if m, ok := args[0].(map[string]string); ok {
			for k, v := range m {
				that.header[k] = v
			}
		}
	} else if len(args) > 1 {
		name, ok := args[0].(string)
		value, ok1 := args[1].(string)
		if ok && ok1 {
			that.header[name] = value
		}
	}
	return that
}

// Code 发送HTTP状态
func (that *Response) Code(code int) *Response {
	that.code = code
	return that
}

// LastModified 设置最后修改日期
func (that *Response) LastModified(time int64) *Response {
	that.header["Last-Modified"] = fmt.Sprintf("%d", time)
	return that
}

// Expires HTTP 1.0，设置缓存的截止时间，在此之前，浏览器对缓存的数据不重新发请求。
// 它与Last-Modified/Etag结合使用，用来控制请求文件的有效时间，当请求数据在有效期内，浏览器从缓存获得数据。
func (that *Response) Expires(time int64) *Response {
	that.header["Expires"] = fmt.Sprintf("%d", time)
	return that
}

// ETag Etag由服务器端生成，客户端通过If-Match或者说If-None-Match这个条件判断请求来验证资源是否修改
// Etag 主要为了解决 Last-Modified 无法解决的一些问题
func (that *Response) ETag(eTag string) *Response {
	that.header["ETag"] = eTag
	return that
}

// CacheControl 页面缓存控制
func (that *Response) CacheControl(cache string) *Response {
	that.header["Cache-control"] = cache
	return that
}

// ContentType 页面输出类型
func (that *Response) ContentType(contentType string, args ...string) *Response {
	charset := "utf-8"
	if len(args) > 0 {
		charset = args[0]
	}
	that.header["Content-Type"] = fmt.Sprintf("%s; charset=%s", contentType, charset)
	return that
}

// GetHeader 获取头部信息
func (that *Response) GetHeader(name string) any {
	if name != "" {
		if v, ok := that.header[name]; ok {
			return v
		}
		return nil
	}
	return that.header
}

// GetData 获取原始数据
func (that *Response) GetData() *Result {
	return that.data
}

// GetCode 获取状态码
func (that *Response) GetCode() int {
	return that.code
}

/************************************/
/**********		view.go		 ********/
/************************************/

/*Assign
* 模板变量赋值
* @access protected
* @param  mixed $name  要显示的模板变量
* @param  mixed $value 变量的值
* @return $this
 */
func (c *Context) Assign(name string, value any) *Context {
	assignData := c.GetStringMap("__global__")
	if assignData == nil {
		assignData = make(map[string]interface{})
	}
	assignData[name] = value
	c.Set("__global__", assignData)
	return c
}

type View struct {
	engine *Context // 模板引擎实例
}

func (c *Context) View() *View {
	return &View{engine: c}
}

/*Success
 * 操作成功返回的数据
 * @param string $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为1
 * @param string $type   输出类型
 * @param array  $header 发送的 Header 信息
 */
func (v *View) Success(args ...any) {
	resp := &Result{
		Code:   1,
		Header: make(map[string]string),
	}
	for index, arg := range args {
		switch index {
		case 0:
			resp.Msg = arg.(string)
		case 1:
			resp.Data = arg
		case 2:
			resp.Code = arg.(int)
		case 3:
			resp.Type = arg.(string)
		case 4:
			resp.Header = arg.(map[string]string)
		}
	}
	v.Result(resp.Msg, resp.Data, resp.Code, resp.Type, resp.Header)
}

/*Error
 * 操作失败返回的数据
 * @param string $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为0
 * @param string $type   输出类型
 * @param array  $header 发送的 Header 信息
 */
func (v *View) Error(args ...any) {
	resp := &Result{
		Code:   0,
		Header: make(map[string]string),
	}
	for index, arg := range args {
		switch index {
		case 0:
			resp.Msg = arg.(string)
		case 1:
			resp.Data = arg
		case 2:
			resp.Code = arg.(int)
		case 3:
			resp.Type = arg.(string)
		case 4:
			resp.Header = arg.(map[string]string)
		}
	}
	v.Result(resp.Msg, resp.Data, resp.Code, resp.Type, resp.Header)
}

/*Result
 * 返回封装后的 API 数据到客户端
 * @param mixed  $msg    提示信息
 * @param mixed  $data   要返回的数据
 * @param int    $code   错误码，默认为0
 * @param string $type   输出类型，支持json/xml/jsonp
 * @param array  $header 发送的 Header 信息
 */
func (v *View) Result(msg string, data any, code int, types string, header map[string]string) {
	result := &Result{
		Code: code,
		Msg:  msg,
		Time: time.Now().Unix(),
		Data: data,
	}
	// 如果未设置类型则使用默认类型判断
	if types == "" {
		if value, ok := v.engine.Request.Header["Response-Type"]; ok {
			types = value[0]
		}
	}

	if statusCode, ok := header["statuscode"]; ok {
		code, _ = strconv.Atoi(statusCode)
		delete(header, "statuscode")
	} else {
		//未设置状态码,根据code值判断
		if code >= 1000 || code < 200 {
			code = 200
		}
	}
	resp := Create(v.engine, result, types, code).Header(header)
	resp.Render()
	Exit()
}

/*Fetch
* 解析和获取模板内容 用于输出
* @param string    $template 模板文件名或者内容
* @param array     $vars     模板输出变量
 */
func (v *View) Fetch(args ...any) {
	templ := defaultTemplName(2)
	resp := &Result{
		Code:   200,
		Url:    templ,
		Data:   make(map[string]any),
		Header: make(map[string]string),
	}
	for index, arg := range args {
		switch index {
		case 0:
			if val, ok := arg.(string); ok && val != "" {
				resp.Url = utils.URL(val, nil, templ)
			}
		case 1:
			if val, ok := arg.(map[string]any); ok {
				resp.Data = val
			}
		}
	}
	//如果没有后缀补上.html
	if !slices.Contains([]string{".html", ".tpl", ".tmpl"}, filepath.Ext(resp.Url)) {
		resp.Url += ".html"
	}
	obj := v.engine.GetStringMap("__global__")
	if obj == nil {
		obj = make(map[string]any)
	}
	for key, value := range resp.Data.(map[string]any) {
		obj[key] = value
	}
	v.engine.HTML(resp.Code, resp.Url, (H)(obj))
	Exit()
}

/************************************/
/**********		jump.go		 ********/
/************************************/

type Jump struct {
	context *Context
}

func (c *Context) Jump() *Jump {
	return &Jump{context: c}
}

/*Success
 * 操作成功跳转的快捷方法
 * @param mixed  $msg    提示信息
 * @param string $url    跳转的 URL 地址
 * @param mixed  $data   返回的数据
 * @param int    $wait   跳转等待时间
 * @param array  $header 发送的 Header 信息
 */
func (j *Jump) Success(args ...any) {
	Exit()
}

/*Error
 * 操作错误跳转的快捷方法
 * @param mixed  $msg    提示信息
 * @param string $url    跳转的 URL 地址
 * @param mixed  $data   返回的数据
 * @param int    $wait   跳转等待时间
 * @param array  $header 发送的 Header 信息
 */
func (j *Jump) Error(args ...any) {
	result := &Result{
		Code:   0,
		Wait:   3,
		Header: make(map[string]string),
	}
	types := j.context.getResponseType()
	for index, arg := range args {
		switch index {
		case 0:
			result.Msg = arg.(string)
		case 1:
			result.Url = arg.(string)
		case 2:
			result.Data = arg
		case 3:
			result.Wait = arg.(int)
		case 4:
			result.Header = arg.(map[string]string)
		}
	}
	if result.Url != "" {
		if !(strings.HasPrefix(result.Url, "/") || strings.HasPrefix(result.Url, ".")) {
			result.Url = "/" + result.Url
		}
	}

	//requset := RequsetS{context: j.context}

	response := Create(j.context, result, types, result.Code).Header(result.Header)
	response.Render()

	//if strings.ToLower(types) == "html" {
	//	j.context.HTML(http.StatusOK, config.Get("dispatch_error_tmpl").(string), gin.H{
	//		"lang": "zh-cn." + requset.Module(true),
	//		"code": 0,
	//		"msg":  result["msg"],
	//		"url":  result["url"],
	//		"wait": result["wait"],
	//	})
	//} else if strings.ToLower(types) == "json" {
	//	j.context.JSON(http.StatusOK, gin.H{
	//		"lang": "zh-cn." + requset.Module(true),
	//		"code": 0,
	//		"msg":  result["msg"],
	//		"time": time.Now().Unix(),
	//		"data": result["data"],
	//		"url":  result["url"],
	//		"wait": result["wait"],
	//	})
	//}
	Exit()
}

/*Result
 * 返回封装后的 API 数据到客户端
 * @param mixed  $data   要返回的数据
 * @param int    $code   返回的 code
 * @param mixed  $msg    提示信息
 * @param string $type   返回数据格式
 * @param array  $header 发送的 Header 信息
 */
func (j *Jump) Result(msg string, data any, code int, types string, header map[string]string) {
	result := &Result{
		Code: code,
		Msg:  msg,
		Time: time.Now().Unix(),
		Data: data,
	}
	if types == "" {
		types = j.context.getResponseType()
	}
	response := Create(j.context, result, types, result.Code).Header(header)
	response.Render()
	Exit()
}

/*Rredirect
 * URL 重定向
 * @param string    $url    跳转的 URL 表达式
 * @param array|int $params 其它 URL 参数
 * @param int       $code   http code
 * @param array     $with   隐式传参
 */
func (j *Jump) Rredirect() {

}

/*GetResponseType
 * 获取当前的 response 输出类型
 * @access protected
 * @return string
 */
func (j *Jump) GetResponseType() string {
	if j.context.IsAjax() {
		return config.Get("default_ajax_return").(string)
	} else {
		return config.Get("default_return_type").(string)
	}
}

// 根据方法名自动匹配模板
func defaultTemplName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	dir, file := filepath.Split(funcName)
	dirArr := strings.Split(filepath.ToSlash(dir), "/")
	fileArr := strings.Split(file, ".")
	fileArr[0] = "view"
	for i, lens := 1, len(fileArr); i < lens; i++ {
		fileArr[i] = utils.Camel2Snake(fileArr[i])
	}
	templ, _ := filepath.Rel("application", strings.Join([]string{strings.Join(dirArr[1:], "/"), strings.Join(fileArr, "/"), ".html"}, ""))
	return filepath.ToSlash(templ)
}
