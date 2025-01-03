package http

import "gin"

type Request struct {
	*gin.Context
}

// Action 当前的操作名
func (that *Request) Action() string {
	value, _ := that.Get("Action")
	return value.(string)
}

// Controller 当前的控制器名
func (that *Request) Controller() string {
	value, _ := that.Get("Controller")
	return value.(string)
}

// Module 获取模块名
func (that *Request) Module() string {
	value, _ := that.Get("Module")
	return value.(string)
}

// Server 获取server参数
func (that *Request) Server(name string, args ...string) string {
	value := that.GetHeader(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Request 获取request变量
func (that *Request) Request(name string, args ...string) string {
	value := that.PostForm(name)
	if value != "" {
		return value
	}
	value = that.Query(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Post 获取post参数
func (that *Request) Post(name string, args ...string) string {
	if len(args) > 0 {
		return that.DefaultPostForm(name, args[0])
	}
	return that.PostForm(name)
}
