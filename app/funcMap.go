package app

import (
	"gin"
	"gin/config"
	"gin/lang"
	"gin/utils"
	"html"
	"html/template"
	"strings"
)

// FuncMap 自定义模板函数
var FuncMap = template.FuncMap{
	//把一些字符转换为 HTML 实体
	"htmlentities": html.EscapeString,
	//把 HTML 实体转换为字符
	"html_entity_decode": func(str string) template.HTML {
		return template.HTML(str)
	},
	"date":        utils.Date,
	"time":        utils.Time,
	"echo":        utils.Echo,
	"json_encode": utils.JsonEncode,
	"json_decode": utils.JsonDecode,
	"ifor":        utils.Ifor,
	"ThinkConfig": config.Get,
	"bool": func(arg any) bool {
		return !utils.Empty(arg)
	},
	"i18n": func(str string, args ...gin.H) string {
		language := "zh-cn"
		for index, arg := range args {
			switch index {
			case 0:
				if v1, o1 := arg["lang"]; o1 {
					if v2, o2 := v1.(string); o2 {
						language = v2
					}
				}
			}
		}
		return lang.I18n(str, language)
	},
	"cdnurl": func(url string) string {
		return url
	},
	"url": func(toUrl string, args ...any) string {
		var base string
		var vars any
		for index, arg := range args {
			switch index {
			case 0:
				vars = arg
			case 1:
				if v, ok := arg.(string); ok {
					base = v
				}
			}
		}
		if !(strings.HasPrefix(toUrl, "/") || strings.HasPrefix(toUrl, ".")) {
			toUrl = "/" + toUrl
		}
		return utils.URL(toUrl, vars, base)
	},
	"default": func(arg any, def any) any {
		if utils.Empty(arg) {
			return def
		}
		return arg
	},
}
