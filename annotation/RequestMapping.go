package annotation

import (
	"path/filepath"
	"slices"
	"strings"
)

var enum = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func RequestMapping(httpMethod string, attr map[string]string) ([]string, string) {
	var httpMethods []string

	if httpMethod == "Param" {
		httpMethods = append(httpMethods, "Any")
	} else {
		httpMethods = append(httpMethods, strings.Split(httpMethod, ",")...)
	}

	var path string
	for k, v := range attr {
		switch k {
		case "", "path", "value":
			path = strings.Trim(strings.TrimSpace(v), "\"")
		case "method":
			if httpMethods[0] == "Any" {
				httpMethods = []string{}
				tmp := strings.Split(strings.Trim(v, "\""), ",")
				for _, val := range tmp {
					val = strings.TrimSpace(val)
					if slices.Contains(enum, val) {
						httpMethods = append(httpMethods, val)
					}
				}
			}
		}
	}
	return httpMethods, filepath.ToSlash(path)
}
