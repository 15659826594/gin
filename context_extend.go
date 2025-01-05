package gin

type MyRequest struct {
	*Context
}

// GetRequest 当前的操作名
func (c *Context) GetRequest() *MyRequest {
	return &MyRequest{c}
}

// Action 当前的操作名
func (c *MyRequest) Action() string {
	value, _ := c.Get("Action")
	return value.(string)
}

// Controller 当前的控制器名
func (c *MyRequest) Controller() string {
	value, _ := c.Get("Controller")
	return value.(string)
}

// Module 获取模块名
func (c *MyRequest) Module() string {
	value, _ := c.Get("Module")
	return value.(string)
}

// Server 获取server参数
func (c *MyRequest) Server(name string, args ...string) string {
	value := c.GetHeader(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Request 获取request变量
func (c *MyRequest) Request(name string, args ...string) string {
	value := c.PostForm(name)
	if value != "" {
		return value
	}
	value = c.Query(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Post 获取post参数
func (c *MyRequest) Post(name string, args ...string) string {
	if len(args) > 0 {
		return c.DefaultPostForm(name, args[0])
	}
	return c.PostForm(name)
}
