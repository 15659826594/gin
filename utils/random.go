package utils

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math"
)

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
		return Substr(StrShuffle(StrRepeat(pool, int(math.Ceil(float64(lens)/float64(len(pool)))))), 0, lens)
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
