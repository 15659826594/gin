package gin

type Req struct {
	*Context
}

// Requests 当前的操作名
func (c *Context) Requests() *Req {
	return &Req{c}
}

// Action 当前的操作名
func (c *Req) Action() string {
	value, _ := c.Get("Action")
	return value.(string)
}

// Controller 当前的控制器名
func (c *Req) Controller() string {
	value, _ := c.Get("Controller")
	return value.(string)
}

// Module 获取模块名
func (c *Req) Module() string {
	value, _ := c.Get("Module")
	return value.(string)
}

// Server 获取server参数
func (c *Req) Server(name string, args ...string) string {
	value := c.GetHeader(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Request 获取Requests变量
func (c *Req) Request(name string, args ...string) string {
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
func (c *Req) Post(name string, args ...string) string {
	if len(args) > 0 {
		return c.DefaultPostForm(name, args[0])
	}
	return c.PostForm(name)
}

// Cookie 获取Cookie
func (c *Req) Cookie(name string) string {
	if cookie, err := c.Context.Request.Cookie("token"); err == nil {
		return cookie.Value
	}
	return ""
}
