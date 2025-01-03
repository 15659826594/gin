package env

import (
	"fmt"
	"gin"
	"gin/lib/godotenv"
	"os"
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
		val = args[0]
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
	if gin.IsFile(defined("ROOT_PATH") + ".env.sample") {
		Load(os.Getenv("ROOT_PATH") + ".env.sample")
	}
}

func Load(path string) bool {
	err := godotenv.Load(path)
	if err != nil {
		return false
	}
	return true
}

func Get(name string, args ...string) string {
	result := defined(name)
	if result == "" && len(args) > 0 {
		return args[0]
	}
	return result
}

func Set(key string, val string) {
	defined(key, val)
}
