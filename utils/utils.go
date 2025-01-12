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
	"regexp"
	"strings"
	"time"
)

// Camel2Snake Camel to underline
func Camel2Snake(camel string) string {
	snake := regexp.MustCompile("([a-z0-9])([A-Z]+)").ReplaceAllString(camel, "${1}_${2}")
	snake = strings.ToLower(snake)
	return snake
}

// Snake2Camel Underline to camel
func Snake2Camel(s string) string {
	re := regexp.MustCompile(`_([a-z])`)
	snake := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	// 小驼峰 (第一个单词小写)
	if snake != "" {
		snake = strings.ToLower(snake)
	}
	return snake
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
	_, _ = m.Write([]byte(dataByte))
	result := m.Sum(nil)
	if binary {
		return result
	}
	return hex.EncodeToString(result)
}

// Md5 字符串md5加密
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
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
	flag := InArray(position, []string{"begin", "start", "first", "front"})
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
