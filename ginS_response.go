package gin

import "fmt"

const JSON_UNESCAPED_UNICODE = 256

type Response struct {
	// 原始数据
	data any
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
	content any
}

func Construct(data any, code int, header map[string]string, options map[string]string) *Response {
	that := &Response{
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

func Create(data any, types string, code int, args ...map[string]string) *Response {
	header := make(map[string]string)
	options := make(map[string]string)
	if len(args) > 0 {
		header = args[0]
	}
	if len(args) > 1 {
		options = args[1]
	}
	resp := Construct(data, code, header, options)
	switch types {
	case "json":
		resp.options["json_encode_param"] = fmt.Sprintf("%d", JSON_UNESCAPED_UNICODE)
		resp.contentType = "application/json"
	case "jsonp":
		resp.options["var_jsonp_handler"] = "callback"
		resp.options["default_jsonp_handler"] = "jsonpReturn"
		resp.options["json_encode_param"] = fmt.Sprintf("%d", JSON_UNESCAPED_UNICODE)
		resp.contentType = "application/javascript"
	case "view":
		resp.contentType = "text/html"
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

// Send 发送数据到客户端
func (that *Response) Send(c *Context) {
	switch that.contentType {
	case "application/json":
		c.JSON(that.code, that.data)
	case "application/jsonp":
		c.JSONP(that.code, that.data)
	case "text/xml":
		c.XML(that.code, that.data)
	}
}

// Output 处理数据
func (that *Response) Output(data any) any {
	return data
}

// Options 输出的参数
func (that *Response) Options(options map[string]string) *Response {
	for k, v := range options {
		that.options[k] = v
	}
	return that
}

// Data 输出数据设置
func (that *Response) Data(data any) *Response {
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

// Content 设置页面输出内容
func (that *Response) Content(content any) *Response {
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
func (that *Response) GetData() any {
	return that.data
}

// GetContent 获取输出数据
func (that *Response) GetContent() any {
	if that.content == nil {
		content := that.Output(that.data)

		that.content, _ = content.(string)
	}
	return that.content
}

// GetCode 获取状态码
func (that *Response) GetCode() int {
	return that.code
}
