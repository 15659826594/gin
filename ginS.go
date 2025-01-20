package gin

import (
	"errors"
	"fmt"
	"gin/src/html/template"
	"gin/utils"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
)

type IController interface {
	Initialize(*Context)
	Value() string
	NoNeedLogin() []string
	NoNeedRight() []string
	ResponseType() string
	BeforeAction() []HandlerFunc
	Exception() HandlerFunc
}

type IJump interface {
	Success(*Context, ...any)
	Error(*Context, ...any)
	Result(context *Context, data any, code int, msg string, types string, header map[string]string)
	Redirect(context *Context, url string, params map[string]string, code int, with map[string]string)
	GetResponseType(*Context) string
}

type IView interface {
	Assign(*Context, string, ...any)
	Fetch(*Context, ...any)
}

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
/**********    requset.go	 ********/
/************************************/

type RequsetS struct {
	context *Context
}

func (c *Context) Requests() *RequsetS {
	return &RequsetS{context: c}
}

func (r *RequsetS) IsGet() bool {
	return r.context.Request.Method == http.MethodGet
}

func (r *RequsetS) IsPost() bool {
	return r.context.Request.Method == http.MethodPost
}

func (r *RequsetS) IsPut() bool {
	return r.context.Request.Method == http.MethodPut
}

func (r *RequsetS) IsDelete() bool {
	return r.context.Request.Method == http.MethodDelete
}

func (r *RequsetS) IsHead() bool {
	return r.context.Request.Method == http.MethodHead
}

func (r *RequsetS) IsPatch() bool {
	return r.context.Request.Method == http.MethodPatch
}

func (r *RequsetS) IsOptions() bool {
	return r.context.Request.Method == http.MethodOptions
}

func (r *RequsetS) IsAjax() bool {
	return r.context.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Langset 当前的语言
func (r *RequsetS) Langset(withPath bool, args ...string) string {
	language := "zh-cn"
	for index, arg := range args {
		switch index {
		case 0:
			language = arg
		}
	}
	if withPath {
		_, module, controller, action, err := separateHandlerName(r.context.GetString("__handler_name__"))
		if err == nil {
			return utils.Camel2Snake(strings.Join([]string{language, module, controller, action}, "."))
		}
	}
	return language
}

// Module 获取当前的模块名
func (r *RequsetS) Module(toSnake bool) string {
	_, module, _, _, err := separateHandlerName(r.context.GetString("__handler_name__"))
	if err != nil {
		return ""
	}
	if toSnake {
		return utils.Camel2Snake(module)
	}
	return module
}

// Controller 获取当前的控制器名
func (r *RequsetS) Controller(toSnake bool) string {
	_, _, controller, _, err := separateHandlerName(r.context.GetString("__handler_name__"))
	if err != nil {
		return ""
	}
	if toSnake {
		return utils.Camel2Snake(controller)
	}
	return controller
}

// Action 获取当前的操作名
func (r *RequsetS) Action(toSnake bool) string {
	_, _, _, action, err := separateHandlerName(r.context.GetString("__handler_name__"))
	if err != nil {
		return ""
	}
	if toSnake {
		return utils.Camel2Snake(action)
	}
	return action
}

func (r *RequsetS) Token(args ...string) string {
	name := "__token__"
	types := "md5"
	for i, arg := range args {
		switch i {
		case 0:
			if arg != "" {
				name = arg
			}
		case 1:
			if arg != "" {
				types = arg
			}
		}
	}
	return name + types
}

func separateHandlerName(handlerName string) (version string, module string, controller string, action string, err error) {
	if handlerName == "" {
		err = errors.New("undefined handlerName")
		return
	}
	arr := strings.Split(handlerName, "/")
	lens := len(arr)
	version = arr[lens-3]
	module = arr[lens-2]
	filego := strings.Split(arr[lens-1], ".")
	lens = len(filego)
	controller = filego[lens-2]
	action = filego[lens-1]
	return
}

// Server 获取server参数
func (r *RequsetS) Server(name string, args ...string) string {
	value := r.context.GetHeader(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Request 获取Requests变量
func (r *RequsetS) Request(name string, args ...string) string {
	value := r.context.PostForm(name)
	if value != "" {
		return value
	}
	value = r.context.Query(name)
	if value != "" {
		return value
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Post 获取post参数
func (r *RequsetS) Post(name string, args ...string) string {
	if len(args) > 0 {
		return r.context.DefaultPostForm(name, args[0])
	}
	return r.context.PostForm(name)
}

// Cookie 获取Cookie
func (r *RequsetS) Cookie(name string) string {
	if cookie, err := r.context.Request.Cookie(name); err == nil {
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
func (r *RequsetS) Url(url string, vars any, base any) string {
	var baseUrl string
	switch tmp := base.(type) {
	case string:
		baseUrl = tmp
	case bool:
		if tmp {
			baseUrl = r.context.Request.URL.String()
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

func (r *Result) ToStringMap(all bool) map[string]any {
	kind := reflect.ValueOf(r).Kind()
	if !(kind == reflect.Struct || kind == reflect.Ptr) {
		return nil
	}
	typeOf := reflect.TypeOf(r)
	valueOf := reflect.ValueOf(r)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	lens := typeOf.NumField()
	kv := make(map[string]any, lens)
	for i := 0; i < lens; i++ {
		field := typeOf.Field(i)
		fieldValue := valueOf.Field(i)
		tag := field.Tag.Get("json")
		if all {
			kv[utils.Camel2Snake(field.Name)] = fieldValue.Interface()
		} else {
			if tag != "-" {
				kv[tag] = fieldValue.Interface()
			}
		}
	}
	return kv
}

// 从上下文中获取Jump接口
func getIJump(c *Context) IJump {
	var boolean bool
	value, boolean := c.Get("__jump__")
	if !boolean {
		return nil
	}
	jump, boolean := value.(IJump)
	if !boolean {
		return nil
	}
	return jump
}

// 从上下文中获取View接口
func getIView(c *Context) IView {
	var boolean bool
	value, boolean := c.Get("__view__")
	if !boolean {
		return nil
	}
	view, boolean := value.(IView)
	if !boolean {
		return nil
	}
	return view
}

// Jump接口

func (c *Context) Suces(args ...any) {
	jump := getIJump(c)
	if jump == nil {
		return
	}
	jump.Success(c, args...)
}

func (c *Context) Eror(args ...any) {
	jump := getIJump(c)
	if jump == nil {
		return
	}
	jump.Error(c, args...)
}

func (c *Context) Result(data any, code int, msg string, types string, header map[string]string) {
	jump := getIJump(c)
	if jump == nil {
		return
	}
	jump.Result(c, data, code, msg, types, header)
}

func (c *Context) Redirect301(url string, params map[string]string, code int, with map[string]string) {
	jump := getIJump(c)
	if jump == nil {
		return
	}
	jump.Redirect(c, url, params, code, with)
}

func (c *Context) GetResponseType() string {
	jump := getIJump(c)
	if jump == nil {
		return ""
	}
	return jump.GetResponseType(c)
}

// View接口

/*Fetch
 * 解析和获取模板内容 用于输出
 * @param string    $template 模板文件名或者内容
 * @param array     $vars     模板输出变量
 */
func (c *Context) Fetch(args ...any) {
	view := getIView(c)
	if view == nil {
		templ, vars := fetchFunc(c, args...)
		c.HTML(http.StatusOK, templ, (H)(vars))
		return
	}
	view.Fetch(c, args...)
}

func (c *Context) FetchFunc(args ...any) (string, map[string]any) {
	return fetchFunc(c, args...)
}

func fetchFunc(c *Context, args ...any) (string, map[string]any) {
	vars := map[string]any{}
	global := c.GetStringMap("__global__")
	for k, v := range global {
		vars[k] = v
	}
	temp := c.templete()
	for i, arg := range args {
		switch i {
		case 0:
			relativepath := arg.(string)
			if relativepath == "" {
				continue
			}
			if strings.HasPrefix(relativepath, "/") {
				temp = relativepath
			} else {
				if !strings.HasPrefix(relativepath, ".") {
					relativepath = "./" + relativepath
				}
				temp = filepath.Clean(temp + "/../" + relativepath)
			}
		case 1:
			if val, ok := arg.(map[string]any); ok {
				for k2, v2 := range val {
					vars[k2] = v2
				}
			}
		}
	}
	temp, _ = filepath.Rel("application/", temp)
	temp = filepath.ToSlash(temp)
	if filepath.Ext(temp) == "" {
		temp += ".html"
	}
	return temp, vars
}

/*Assign
 * 模板变量赋值
 * @param mixed $name  变量名
 * @param mixed $value 变量值
 */
func (c *Context) Assign(name string, args ...any) {
	view := getIView(c)
	if view == nil {
		assignFunc(c, name, args...)
		return
	}
	view.Assign(c, name, args...)
}

func AssignFunc(c *Context, name any, args ...any) {
	assignFunc(c, name, args...)
}

func assignFunc(c *Context, name any, args ...any) {
	global := c.GetStringMap("__global__")
	if global == nil {
		global = make(map[string]any)
	}
	switch n := name.(type) {
	case string:
		if len(args) > 0 {
			global[n] = args[0]
		}
	default:
		for k, v := range utils.Iterator(name) {
			global[fmt.Sprintf("%s", k)] = v
		}
	}
	c.Set("__global__", global)
}

func (c *Context) MergeAssign(vars map[string]any) H {
	return mergeAssign(c, vars)
}

func mergeAssign(c *Context, vars map[string]any) H {
	assignGlobal := c.GetStringMap("__global__")
	for k, v := range assignGlobal {
		if _, ok := vars[k]; !ok {
			vars[k] = v
		}
	}
	return H(vars)
}

func (c *Context) templete() string {
	handlerName := c.GetString("__handler_name__")
	var basepath string
	if handlerName != "" {
		fileArr := strings.Split(handlerName, "/")
		lens := len(fileArr)
		version := fileArr[lens-3]
		module := fileArr[lens-2]
		filego := strings.Split(fileArr[lens-1], ".")
		lens2 := len(filego)
		controller := filego[lens2-2]
		action := filego[lens2-1]
		basepath = utils.Camel2Snake(fmt.Sprintf("%s/%s/view/%s/%s.html", version, module, controller, action))
	}
	return basepath
}

/************************************/
/**********  	 auth.go 	 ********/
/************************************/

type IAuth interface {
	GetUser()
	Init(token string) bool
	Register(username string, password string, email string, mobile string, extend []string) bool
	Login(account string, password string) bool
	Logout() bool
	Changepwd(newpassword string, oldpassword string, ignoreoldpassword bool) bool
	Direct(userId int) bool
	Check(path string, args ...string) bool
	IsLogin() bool
	GetToken() string
	GetUserinfo() map[string]any
	GetRuleList()
	GetRequestUri() string
	SetRequestUri(uri string)
	GetAllowFields() []string
	SetAllowFields(fields []string)
	Delete(userId int) bool
	GetEncryptPassword(password string, salt string) string
	Match(arr []string) bool
	SetKeeptime(keeptime int64)
	Render()
	SetError(err string)
	GetError() string
}

func (c *Context) Auth() IAuth {
	if val, ok := c.Get("__auth__"); !ok {
		return nil
	} else {
		return val.(IAuth)
	}
}
