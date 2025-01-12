package lang

import (
	"gin/src/encoding"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	_ = Load(filepath.Clean(filename + "/../zh-cn.json"))
}

var global = map[string]*Scope{
	"zh-cn": newScope(),
}

/*ParseGlob  批量载入
 * @param string 			   	$pattern 	匹配
 * @param map[string]string    	$prefix 	去掉前缀 application\/**\/lang => application
 * @param []string    			$allowFile 	检索文件类型
 * @demo
 * application\index\lang\zh-cn.json					=>	zh-cn.index
 * application\index\lang\zh-cn\ajax.json				=>	zh-cn.index.ajax
 * application\admin\lang\zh-cn\general\profile.json	=>	zh-cn.admin.general.profile
 */
func ParseGlob(pattern string, prefix string, allowFile []string) {
	filenames, _ := filepath.Glob(pattern)
	for _, filename := range filenames {
		_ = filepath.WalkDir(filename, func(path string, d fs.DirEntry, err error) error {
			ext := filepath.Ext(path)
			if slices.Contains(allowFile, ext) {
				name := strings.TrimSuffix(path, ext) //去掉文件后缀
				name, _ = filepath.Rel(prefix, name)  //去掉application
				arr := strings.Split(filepath.ToSlash(name), "/")
				if arr[1] == "lang" {
					arr = slices.Delete(arr, 1, 2) //去掉lang
				}
				arr[0], arr[1] = arr[1], arr[0]
				name = strings.Join(arr, ".")
				_ = Load(path, name)
			}
			return nil
		})
	}
}

/*Load
 * 载入文件 (".json", ".yml", ".yaml", ".xml")
 * @param  string $file 文件路径
 * @param  string $args 作用域
 */
func Load(file string, args ...any) error {
	data, err := encoding.Parse(file)
	if err != nil {
		return err
	}
	var scope IScope
	scope = global["zh-cn"]
	for i, arg := range args {
		switch i {
		case 0:
			if s, ok := arg.(IScope); ok {
				scope = s
			} else if s1, ok1 := arg.(string); ok1 {
				scope = findOrCreateScope(s1)
			}
		}
	}

	if s, ok := scope.(*Scope); ok {
		for k, v := range data {
			s.Keys[k] = v
		}
		return nil
	}

	for k, v := range data {
		Set(k, v, scope)
	}
	return nil
}

/*I18n 国际化翻译
 * @param string 			   	$name 字符串
 * @param map[string]string    	$vars 自定义翻译(优先)
 * @param string    			$lang 作用域链式 zh-cn.index.user.vars <=
 * @return string
 */
func I18n(name string, vars map[string]string, lang string) string {
	if val, ok := vars[name]; ok {
		return val
	}
	if lang == "" {
		lang = "zh-cn"
	}
	arr := strings.Split(lang, ".")
	chains := make([]map[string]string, 0, len(arr)+1)
	createScopeChains(arr, &chains)
	if vars != nil {
		chains = append(chains, vars)
	}
	for i := len(chains) - 1; i >= 0; i-- {
		if v, ok := chains[i][name]; ok {
			return v
		}
	}
	return name
}

/*Set
 * 设置配置参数 name 为数组则为批量设置
 * @access public
 * @param  string|array $name  配置参数名（支持二级配置 . 号分割）
 * @param  mixed        $value 配置值
 * @param  string       $range 作用域
 * @return mixed
 */
func Set(name string, value any, args ...any) {
	if name == "" {
		return
	}
	var scope IScope
	scope = global["zh-cn"]
	for i, arg := range args {
		switch i {
		case 0:
			if s, ok := arg.(IScope); ok {
				scope = s
			} else if s1, ok1 := arg.(string); ok1 {
				scope = findOrCreateScope(s1)
			}
		}
	}
	scope.Set(name, value)
}
