package lang

import (
	"gin/src/encoding"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// dictionary.Subs["zh"].Subs["cn"]
// dictionary.SearchSub(splitDictName("zh-cn"))
func init() {
	_, filename, _, _ := runtime.Caller(0)
	_ = Load(filepath.Clean(filename+"/../zh-cn.json"), "zh-cn")
}

var dictionary = NewDict()

/*I18n 国际化翻译
 * @param string 			   	$name 字符串
 * @param string    			$lang 作用域链式 zh-cn.index.user.vars <=
 * @param map[string]string    	$vars 自定义翻译(优先)
 * @return string
 */
func I18n(name string, lang string, args ...map[string]string) string {
	if lang == "" {
		lang = "zh-cn"
	}
	chains, err := dictChains(dictionary, splitDictName(lang))
	if err != nil {
		return ""
	}
	if len(args) > 0 {
		chains = append(chains, args[0])
	}
	for i := len(chains) - 1; i >= 0; i-- {
		if val, ok := chains[i][name]; ok {
			return val
		}
	}
	return name
}

/*Load
 * 载入文件 (".json", ".yml", ".yaml", ".xml")
 * @param  string $file 文件路径
 * @param  string $dictName 字典名称
 */
func Load(file string, dictName string) error {
	var dict *Dict
	var err error
	data, err := encoding.Parse(file)
	if err != nil {
		return err
	}
	dictName = strings.ToLower(dictName)
	if dictName == "" {
		dictName = "zh-cn"
	}
	dict, err = dictionary.SearchSub(splitDictName(dictName))
	if err != nil {
		if err.Error() == "dict not found" {
			dict, err = dictionary.CreateSub(splitDictName(dictName))
		} else {
			return err
		}
	}
	dict.mu.Lock()
	for k, v := range data {
		dict.Keys[k] = v.(string)
	}
	dict.mu.Unlock()
	return nil
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
