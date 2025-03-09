package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/ripemd160"
	"hash"
	"hash/adler32"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"
)

/*URL
 * 生成 支持路由反射
 * @param string            				url 	路由地址
 * @param string|map[string]string          vars 	路由参数
 * @param string    						base 	例如:http://example.com/
 * @return string
 */
func URL(toUrl string, vars any, base string) string {
	var baseURL *url.URL
	toURL, _ := url.Parse(toUrl)
	if base != "" {
		baseURL, _ = url.Parse(base)
	}
	if baseURL != nil {
		toURL.Scheme = baseURL.Scheme
		toURL.Opaque = baseURL.Opaque
		toURL.User = baseURL.User
		toURL.Host = baseURL.Host
		if !strings.HasPrefix(toURL.Path, "/") {
			toURL.Path, _ = url.JoinPath(baseURL.Path, "../"+toURL.Path)
		}
	}
	switch tmp := vars.(type) {
	case string:
		if tmp != "" {
			search := toURL.Query()
			arr := strings.Split(tmp, "&")
			for _, v := range arr {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					search.Add(kv[0], kv[1])
				} else if len(kv) == 1 {
					search.Add(kv[0], "")
				}
			}
			toURL.RawQuery = search.Encode()
		}
	case map[string]string:
		if tmp != nil {
			search := toURL.Query()
			for k, v := range tmp {
				search.Set(k, v)
			}
			toURL.RawQuery = search.Encode()
		}
	}
	toURL.Path, _ = url.JoinPath(toURL.Path)
	return toURL.String()
}

// Camel2Snake 驼峰转下划线命名
func Camel2Snake(camel string) string {
	snake := regexp.MustCompile("([a-z0-9])([A-Z]+)").ReplaceAllString(camel, "${1}_${2}")
	snake = strings.ToLower(snake)
	return snake
}

// Snake2Camel 下划线转驼峰命名(pascal 大驼峰)
func Snake2Camel(s string, args ...bool) string {
	if s == "" {
		return s
	}
	pascal := false
	if len(args) > 0 {
		pascal = args[0]
	}
	re := regexp.MustCompile(`_([a-z])`)
	snake := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	if !pascal { // 小驼峰 (第一个单词小写)
		return strings.ToLower(snake)
	} else { // 大驼峰 (第一个单词大写)
		return strings.ToUpper(string(snake[0])) + snake[1:]
	}
}

func HashHmac(algo string, data string, key string) string {
	var m hash.Hash
	dataByte := []byte(data)
	keyByte := []byte(key)
	switch algo {
	case "sha256":
		m = hmac.New(sha256.New, keyByte)
	case "ripemd160":
		m = hmac.New(ripemd160.New, keyByte)
	default:
		panic("未定义加密类型")
	}
	_, _ = m.Write([]byte(dataByte))
	result := m.Sum(nil)
	return hex.EncodeToString(result)
}

func Hash(algo string, data string, binary bool) any {
	var m hash.Hash
	dataByte := []byte(data)
	switch algo {
	case "adler32":
		m = adler32.New()
	}
	_, _ = m.Write(dataByte)
	result := m.Sum(nil)
	if binary {
		return result
	}
	return hex.EncodeToString(result)
}

// Md5 字符串md5加密
func Md5(str string) string {
	hashs := md5.New()
	hashs.Write([]byte(str))
	return hex.EncodeToString(hashs.Sum(nil))
}

// Base64Encode base64编码
func Base64Encode(str string) string {
	input := []byte(str)
	return base64.StdEncoding.EncodeToString(input)
}

const (
	YEAR   = 31536000
	MONTH  = 2592000
	WEEK   = 604800
	DAY    = 86400
	HOUR   = 3600
	MINUTE = 60
)

/*DateUnixtime 获取一个基于时间偏移的Unix时间戳
 * @param string $type     时间类型，默认为day，可选minute,hour,day,week,month,quarter,year
 * @param int    $offset   时间偏移量 默认为0，正数表示当前type之后，负数表示当前type之前
 * @param string $position 时间的开始或结束，默认为begin，可选前(begin,start,first,front)，end
 * @param int    $year     基准年，默认为null，即以当前年为基准
 * @param int    $month    基准月，默认为null，即以当前月为基准
 * @param int    $day      基准天，默认为null，即以当前天为基准
 * @param int    $hour     基准小时，默认为null，即以当前年小时基准
 * @param int    $minute   基准分钟，默认为null，即以当前分钟为基准
 * @return int 处理后的Unix时间戳
 */
func DateUnixtime(params ...any) (times int64) {
	var date time.Time
	now := time.Now()
	lens := len(params)
	types := "day"
	offset := 0
	position := "begin"
	var year, month, day, hour, minute int
	if lens > 0 {
		types = params[0].(string)
	}
	if lens > 1 {
		offset = params[1].(int)
	}
	if lens > 2 {
		position = params[2].(string)
	}
	if lens > 3 {
		year = params[3].(int)
	} else {
		year = now.Year()
	}
	if lens > 4 {
		month = params[4].(int)
	} else {
		month = int(now.Month())
	}
	if lens > 5 {
		day = params[5].(int)
	} else {
		day = now.Day()
	}
	if lens > 6 {
		minute = params[6].(int)
	} else {
		minute = now.Minute()
	}
	flag := slices.Contains([]string{"begin", "start", "first", "front"}, position)
	timeMonth := time.Month(month)
	switch types {
	case "minute":
		if flag {
			date = time.Date(year, timeMonth, day, hour, minute+offset, 0, 0, time.Local)
		} else {
			date = time.Date(year, timeMonth, day, hour, minute+offset, 59, 0, time.Local)
		}
	case "hour":
		if flag {
			date = time.Date(year, timeMonth, day, hour+offset, 0, 0, 0, time.Local)
		} else {
			date = time.Date(year, timeMonth, day, hour+offset, 59, 59, 0, time.Local)
		}
	case "day":
		if flag {
			date = time.Date(year, timeMonth, day+offset, hour, 0, 0, 0, time.Local)
		} else {
			date = time.Date(year, timeMonth, day+offset, 23, 59, 59, 0, time.Local)
		}
	default:
		date = time.Date(year, timeMonth, day, hour, minute, 0, 0, time.Local)
	}
	return date.Unix()
}

// DateDaysInMonth 获取指定年月拥有的天数
func DateDaysInMonth(month int, year int) int {
	switch month {
	case 2:
		isLeap := (year%4 == 0 && year%100 != 0) || year%400 == 0
		if isLeap {
			return 29
		}
		return 28
	case 4, 6, 10, 11:
		return 30
	default:
		return 31
	}
}

// Ifor 三目运算符增强版
func Ifor(condition any, param any, optional ...any) any {
	var trueVal, falseVal any
	if len(optional) > 0 {
		trueVal = param
		falseVal = optional[0]
	} else {
		trueVal = condition
		falseVal = param
	}
	switch c := condition.(type) {
	case bool:
		if c {
			return trueVal
		} else {
			return falseVal
		}
	default:
		if !Empty(condition) {
			return trueVal
		} else {
			return falseVal
		}
	}
}

// Iterator 将any类型转为map[any]any
func Iterator(iter any) map[any]any {
	ret := map[any]any{}
	val := reflect.ValueOf(iter)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			ret[i] = val.Index(i).Interface()
		}
	case reflect.Map:
		for i, k := range val.MapKeys() {
			ret[i] = val.MapIndex(k).Interface()
		}
	default:
	}

	return ret
}
