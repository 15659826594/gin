package env

import (
	"fmt"
	encodingEnv "gin/src/encoding/env"
	"gin/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const DS = string(os.PathSeparator)

var rootPath, _ = os.Getwd()

func defined(name string, args ...string) string {
	name = strings.ToUpper(strings.ReplaceAll(name, ".", "_"))
	var val string
	if len(args) > 0 {
		//统一分隔符
		val = filepath.ToSlash(args[0])
		err := os.Setenv(name, val)
		if err != nil {
			return val
		}
	}
	return os.Getenv(name)
}

func init() {
	defined("THINK_VERSION", "1.0.0")
	defined("THINK_START_TIME", fmt.Sprintf("%d", time.Now().Unix()))
	defined("EXT", ".go")
	defined("DS", DS)
	defined("THINK_PATH", rootPath+DS+"think"+DS)
	defined("LIB_PATH", defined("THINK_PATH")+"library"+DS)
	defined("CORE_PATH", defined("LIB_PATH")+"think"+DS)
	defined("TRAIT_PATH", defined("LIB_PATH")+"traits"+DS)
	defined("APP_PATH", rootPath+DS+"application"+DS)
	defined("ROOT_PATH", rootPath+DS)
	defined("EXTEND_PATH", rootPath+DS+"extend"+DS)
	defined("VENDOR_PATH", rootPath+DS+"vendor"+DS)
	defined("RUNTIME_PATH", rootPath+DS+"runtime"+DS)
	defined("LOG_PATH", defined("RUNTIME_PATH")+"log"+DS)
	defined("CACHE_PATH", defined("RUNTIME_PATH")+"cache"+DS)
	defined("TEMP_PATH", defined("RUNTIME_PATH")+"temp"+DS)
	defined("CONF_PATH", defined("APP_PATH"))
	defined("CONF_EXT", defined("EXT"))

	// 环境常量
	if runtime.GOOS == "windows" {
		defined("IS_WIN", "true")
	}

	// 加载环境变量配置文件
	if utils.IsFile(defined("ROOT_PATH") + ".env") {
		Load(os.Getenv("ROOT_PATH") + ".env")
	}
}

/*Load
 * 载入环境变量文件
 * @param string $filenames[] 文件名切片
 * @return bool
 */
func Load(filenames ...string) bool {
	err := encodingEnv.Load(filenames...)
	if err != nil {
		return false
	}
	return true
}

/*Get
 * 设置环境变量
 * @param string $name 健名
 * @param string $default 默认值可选
 * @return string
 */
func Get(name string, args ...string) string {
	result := defined(name)
	if result == "" && len(args) > 0 {
		return args[0]
	}
	return result
}

/*Set
 * 获取环境变量
 * @param string $key 健名
 * @param string $val 健值
 * @return string
 */
func Set(key string, val string) {
	defined(key, val)
}
