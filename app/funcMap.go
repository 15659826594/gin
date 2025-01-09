package app

import (
	"gin/config"
	"gin/utils"
	"html"
	"html/template"
	"net/url"
	"strings"
)

// FuncMap 自定义模板函数
var FuncMap = template.FuncMap{
	"htmlentities": html.EscapeString,
	"date":         utils.Date,
	"time":         utils.Time,
	"echo":         utils.Echo,
	"json_encode":  utils.JsonEncode,
	"json_decode":  utils.JsonDecode,
	"ifor":         utils.Ifor,
	"url": func(args ...any) string {
		var targetUrl, vars, currentURL string
		for index, arg := range args {
			switch index {
			case 0:
				if v, ok := arg.(string); ok {
					targetUrl = v
				}
			case 1:
				if v, ok := arg.(string); ok {
					vars = v
				}
			case 2:
				if v, ok := arg.(string); ok {
					currentURL = v
				}
			}
		}
		return Url(targetUrl, vars, currentURL).String()
	},
	"__": func(s string) string { // i18n 国际化翻译
		return s
	},
	"default": func(arg any, arg1 any) any {
		if utils.Empty(arg) {
			return arg1
		}
		return arg
	},
	"ThinkConfig": func(s string) any {
		return config.Get(s)
	},
}

/*Url 生成
 * @param string        $url 路由地址
 * @param string|array  $vars 变量
 * @param bool|string   $currentURL 当前路径
 * @return string
 */
func Url(targetUrl string, vars string, currentURL string) *url.URL {
	toURL, _ := url.Parse(targetUrl)

	path := toURL.Path
	if len(path) > 0 && (path[0] != '.' && path[0] != '/') {
		toURL.Path = "/" + path
	}

	if vars != "" {
		qu := toURL.Query().Encode()
		if qu != "" {
			vars = qu + "&" + vars
		}
		toURL.RawQuery = vars
	}

	if currentURL != "" {
		fromUrl, err := url.Parse(currentURL)
		if strings.HasPrefix(toURL.Path, "/") {
			fromUrl.Path = ""
		}
		if err == nil {
			newPath, err1 := url.JoinPath(fromUrl.Path, toURL.Path)
			if err1 == nil {
				toURL.Scheme = fromUrl.Scheme
				toURL.Host = fromUrl.Host
				toURL.Path = newPath
			}
		}
	}
	return toURL
}
