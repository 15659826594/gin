package annotation

import (
	"errors"
	"gin"
	"gin/utils"
	"strings"
)

var enum = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func RequestMapping(args ...any) Handler {
	methods := enum
	if len(args) > 0 {
		if v, ok := args[0].([]string); ok {
			methods = v
		}
	}
	return func(group *gin.RouterGroup, handler []gin.HandlerFunc, arguments map[string]string, defPath string) ([]string, string, error) {
		path := defPath
		if v, ok := arguments[""]; ok {
			path = strings.Trim(v, "\"")
		}
		if v, ok := arguments["value"]; ok {
			path = strings.Trim(v, "\"")
		}
		//转为驼峰命名法
		path = utils.Camel2Snake(path)

		if v, ok := arguments["method"]; ok {
			//去掉空格 , 移除左右" , 逗号拆分成切片
			list := strings.Split(strings.Trim(strings.ReplaceAll(v, " ", ""), "\""), ",")
			//获取有效的请求
			methods = intersectArray(list, enum)
		}
		for _, method := range methods {
			group.Handle(method, path, handler...)
		}
		if len(methods) != 0 {
			if len(methods) == len(enum) {
				return []string{"Any"}, path, nil
			}
			return methods, path, nil
		} else {
			return nil, path, errors.New("invalid")
		}
	}
}

// intersectArray 求两个切片的交集
func intersectArray(a []string, b []string) []string {
	var inter []string
	mp := make(map[string]bool)

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}
