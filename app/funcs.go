package app

import (
	"html"
	"os"
)

// 判断根目录下是否存在logo
func getFaviconIco(def string) string {
	if _, err := os.Stat("favicon.ico"); os.IsNotExist(err) {
		return def
	} else {
		return "favicon.ico"
	}
}

func htmlentities(s string) string {
	return html.EscapeString(s)
}
