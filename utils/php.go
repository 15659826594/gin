package utils

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Comparable interface {
	int | byte | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | string
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

//////////// Directory/Filesystem Functions ////////////

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

// FileExists file_exists()
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
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
func IsDir(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	fm := fd.Mode()
	return fm.IsDir()
}

// Empty empty()
func Empty(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// ArrayIntersectKey map[string]any{"a":0,"b":1,"c":2} , map[string]any{"c":3} => map[string]any{"c":2}
func ArrayIntersectKey(target map[string]any, maps ...map[string]any) map[string]any {
	ret := map[string]any{}
	for s, a := range target {
		ret[s] = a
	}
	arrayIntersectKeyOne := func(t1 map[string]any, t2 map[string]any) map[string]any {
		r := make(map[string]any)
		for s, a := range t1 {
			r[s] = a
		}
		for s, _ := range r {
			if _, ok := t2[s]; !ok {
				delete(r, s)
			}
		}
		return r
	}
	for _, m := range maps {
		ret = arrayIntersectKeyOne(ret, m)
	}
	return ret
}

// ArrayFlip []string{"a","b","c"} => map[string]any{"a":0,"b":1,"c":2}
func ArrayFlip(slice []string) map[string]any {
	maps := make(map[string]any)
	for i, s := range slice {
		maps[s] = i
	}
	return maps
}

// Echo any类型转为字符串
func Echo(arg any) string {
	switch v := arg.(type) {
	case string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, complex64, complex128, uintptr:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	default:
		jsonData, err := json.Marshal(arg)
		if err != nil {
			return err.Error()
		}
		return string(jsonData)
	}
}

// ////////// Date/Time Functions ////////////

// Time time() 返回时间戳
func Time() int64 {
	return time.Now().Unix()
}

// Date 函数实现类似 PHP 的 date 方法
func Date(format string, args ...any) string {
	var timestamp int64
	for i, arg := range args {
		switch i {
		case 0:
			timestamp = arg.(int64)
		}
	}
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}
	t := time.Unix(timestamp, 0)
	layout := convertFormat(format)
	return t.Format(layout)
}

// convertFormat 将 PHP 的日期格式转换为 Go 的日期格式
func convertFormat(format string) string {
	replacements := map[string]string{
		"Y": "2006",
		"m": "01",
		"d": "02",
		"H": "15",
		"i": "04",
		"s": "05",
	}

	for phpFormat, goFormat := range replacements {
		format = strings.ReplaceAll(format, phpFormat, goFormat)
	}

	return format
}

// Checkdate checkdate()
// Validate a Gregorian date
// 判断是否是有效日期
func Checkdate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}

	return true
}

//////////// String Functions ////////////

// Ucfirst ucfirst()
// 首字母转为大写
func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// Lcfirst lcfirst()
// 首字母转为小写
func Lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// Ucwords ucwords()
// 每一个单词的首字母转为大写
func Ucwords(str string) string {
	isSeparator := func(r rune) bool {
		if r <= 0x7F {
			switch {
			case '0' <= r && r <= '9':
				return false
			case 'a' <= r && r <= 'z':
				return false
			case 'A' <= r && r <= 'Z':
				return false
			case r == '_':
				return false
			}
			return true
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}

		return unicode.IsSpace(r)
	}

	prev := ' '
	return strings.Map(func(r rune) rune {
		if isSeparator(prev) {
			prev = r
			return unicode.ToTitle(r)
		}
		prev = r
		return r
	}, str)
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

// JsonEncode 用于对变量进行 JSON 编码，该函数如果执行成功返回 JSON 数据
func JsonEncode(arg any) string {
	jsonData, err := json.Marshal(arg)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

// JsonDecode JSON 格式的字符串进行解码，并转换为 map[string]any
func JsonDecode(jsonStr string) map[string]any {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil
	}
	return result
}

// Ternary 三目运算符
// max := Ternary(a > b, a, b).(int)
func Ternary(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
