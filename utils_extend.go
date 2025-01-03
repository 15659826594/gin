package gin

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"hash"
	"hash/adler32"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
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

// IsFile is_file()
func IsFile(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return !fd.IsDir()
}

// IsDir is_dir()
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

/*RandomAlnum
 * 生成数字和字母
 *
 * @param int $len 长度
 * @return string
 */
func RandomAlnum(len int) string {
	return ""
}

/*RandomAlpha
 * 生成数字和字母
 *
 * @param int $len 长度
 * @return string
 */
func RandomAlpha(len int) string {
	return ""
}

/*RandomNumeric
 * 生成指定长度的随机数字
 *
 * @param int $len 长度
 * @return string
 */
func RandomNumeric(len int) string {
	return RandomBuild("numeric", len)
}

/*RandomNozero
 * 生成指定长度的无0随机数字
 *
 * @param int $len 长度
 * @return string
 */
func RandomNozero(len int) string {
	return ""
}

/*RandomBuild 能用的随机数生成
 * @param string $type 类型 alpha/alnum/numeric/nozero/unique/md5/encrypt/sha1
 * @param int    $len  长度
 */
func RandomBuild(types string, lens int) string {
	switch types {
	case "alpha", "alnum", "numeric", "nozero":
		var pool string
		switch types {
		case "alpha":
			pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		case "alnum":
			pool = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		case "numeric":
			pool = "0123456789"
		case "nozero":
			pool = "123456789"
		}
		return Substr(StrShuffle(StrRepeat(pool, int(Ceil(float64(lens)/float64(len(pool)))))), 0, lens)
	case "unique", "md5":
		return "md5"
	case "encrypt", "sha1":
		return "encrypt"
	}

	return ""
}

// RandomUuid 获取全球唯一标识
func RandomUuid() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
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

type Comparable interface {
	int | byte | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | string
}

// Substr substr()
func Substr(str string, start uint, length int) string {
	if length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

// Strlen strlen()
func Strlen(str string) int {
	return len(str)
}

// StrRepeat str_repeat()
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// StrShuffle str_shuffle()
func StrShuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Ltrim ltrim()
func Ltrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, characterMask[0])
}

// MbStrlen mb_strlen()
func MbStrlen(str string) int {
	return utf8.RuneCountInString(str)
}

func SliceMerge[T Comparable](ss ...map[T]any) map[T]any {
	s := make(map[T]any)
	for _, v := range ss {
		for s2, a := range v {
			s[s2] = a
		}
	}
	return s
}

// ArrayRand array_rand()
func ArrayRand(elements []any) []any {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]any, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// Implode implode()
func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// InArray in_array()
// haystack supported _type: slice, array or map
func InArray(needle any, haystack any) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack _type muset be slice, array or map")
	}

	return false
}

func MtRand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r, _ := crand.Int(crand.Reader, big.NewInt(int64(max+1-min)))
	return int(r.Int64()) + min
}

// Ceil ceil()
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pathinfo pathinfo()
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// Usage:
// Pathinfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
func Pathinfo(path string, options int) map[string]string {
	if options == -1 {
		options = 1 | 2 | 4 | 8
	}
	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}
	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}
	if ((options & 4) == 4) || ((options & 8) == 8) {
		basename := ""
		if (options & 2) == 2 {
			basename = info["basename"]
		} else {
			basename = filepath.Base(path)
		}
		p := strings.LastIndex(basename, ".")
		filename, extension := "", ""
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}
		if (options & 4) == 4 {
			info["extension"] = extension
		}
		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}
	return info
}
