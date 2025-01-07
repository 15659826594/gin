package app

import (
	"html"
	"html/template"
	"net/url"
	"os"
	"strings"
	"time"
)

// 判断根目录下是否存在logo
func faviconIco(def string) string {
	if _, err := os.Stat("favicon.ico"); os.IsNotExist(err) {
		return def
	}
	return "favicon.ico"
}

// FuncMap 自定义模板函数
var FuncMap = template.FuncMap{
	"htmlentities": html.EscapeString,
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
	"date": Date,
	"time": func() int64 {
		return time.Now().Unix()
	},
	"__": func(s string) string { // i18n 国际化翻译
		return s
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

// Date 函数实现类似 PHP 的 date 方法
func Date(format string, timestamp int64) string {
	t := time.Unix(timestamp, 0)
	layout := convertFormat(format)
	return t.Format(layout)
}

// convertFormat 将 PHP 的日期格式转换为 Go 的日期格式
func convertFormat(format string) string {
	replacements := map[string]string{
		"Y": "2006",
		"m": "01",
		"d": "02",
		"H": "15",
		"i": "04",
		"s": "05",
	}

	for phpFormat, goFormat := range replacements {
		format = strings.ReplaceAll(format, phpFormat, goFormat)
	}

	return format
}
