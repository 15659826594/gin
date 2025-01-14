package env

import (
	"fmt"
	Env "gin/src/encoding/env"
	"gin/utils"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var rootPath, _ = os.Getwd()

func init() {
	Setenv("THINK_VERSION", "1.0.0")
	Setenv("THINK_START_TIME", fmt.Sprintf("%d", time.Now().Unix()))
	Setenv("EXT", ".go")
	Setenv("DS", string(os.PathSeparator))
	Setenv("ROOT_PATH", rootPath+os.ExpandEnv("${DS}"))
	Setenv("THINK_PATH", os.ExpandEnv("${ROOT_PATH}think${DS}"))
	Setenv("LIB_PATH", os.ExpandEnv("${THINK_PATH}library${DS}"))
	Setenv("CORE_PATH", os.ExpandEnv("${LIB_PATH}think${DS}"))
	Setenv("TRAIT_PATH", os.ExpandEnv("${LIB_PATH}traits${DS}"))
	Setenv("APP_PATH", os.ExpandEnv("${ROOT_PATH}application${DS}"))
	Setenv("EXTEND_PATH", os.ExpandEnv("${ROOT_PATH}extend${DS}"))
	Setenv("VENDOR_PATH", os.ExpandEnv("${ROOT_PATH}vendor${DS}"))
	Setenv("RUNTIME_PATH", os.ExpandEnv("${ROOT_PATH}runtime${DS}"))
	Setenv("LOG_PATH", os.ExpandEnv("${RUNTIME_PATH}log${DS}"))
	Setenv("CACHE_PATH", os.ExpandEnv("${RUNTIME_PATH}cache${DS}"))
	Setenv("TEMP_PATH", os.ExpandEnv("${RUNTIME_PATH}temp${DS}"))
	Setenv("CONF_PATH", os.ExpandEnv("${APP_PATH}"))

	// 环境常量
	if runtime.GOOS == "windows" {
		Setenv("IS_WIN", "true")
	}

	// 加载环境变量配置文件
	if utils.IsFile(Getenv("ROOT_PATH") + ".env") {
		LoadEnv(Getenv("ROOT_PATH") + ".env")
	}
}

// LoadEnv 载入环境变量
func LoadEnv(filenames ...string) {
	for _, filename := range filenames {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		mapss, err := Env.UnmarshalBytes(bytes)
		if err != nil {
			continue
		}
		for k, v := range mapss {
			Setenv(k, v)
		}
	}
}

/*DefaultGetenv
 * 设置环境变量
 * @param string $key 健名
 * @return string $defaultValue 默认值
 */
func DefaultGetenv(key string, defaultValue string) string {
	key = filenamesOrDefault(key)
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

/*Getenv
 * 获取环境变量
 * @param string $key 健名 自动转为大写, . 转为 _
 * @return string
 */
func Getenv(key string) string {
	return os.Getenv(filenamesOrDefault(key))
}

/*Setenv
 * 设置环境变量
 * @param string $key 健名 自动转为大写, . 转为 _
 * @param string $value 健值
 * @return string
 */
func Setenv(key string, value string) {
	err := os.Setenv(filenamesOrDefault(key), value)
	if err != nil {
		log.Fatal(err)
	}
}

func Environ() []string {
	return os.Environ()
}

func Hostname() (name string, err error) {
	return os.Hostname()
}

func ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

func Clearenv() {
	os.Clearenv()
}

func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func filenamesOrDefault(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}
