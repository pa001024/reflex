package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"
	"runtime"
	"strings"
)

const (
	BASE62_ST = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// unsafe
func ToJson(args ...interface{}) string {
	l := len(args) - 1
	m := make(map[string]interface{})
	for i := 0; i < l; i += 2 {
		m[args[i].(string)] = args[i+1]
	}
	bin, _ := json.Marshal(m)
	if bin == nil {
		return ""
	}
	return string(bin)
}

// unsafe
func MustParseUrl(rawurl string) *url.URL {
	y, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return y
}

// 解析js hex值 形如 \x00
func DecodeJsHex(src string) string {
	src = strings.Replace(src, `\x`, "", -1)
	v, _ := hex.DecodeString(src)
	return string(v)
}

// 使用前缀和后缀移除可嵌套的区块
func RemoveBlock(src, perfix, suffix string) (rst string) {
	buf := new(bytes.Buffer)
	c := 0
	mi, im := len([]rune(perfix)), len([]rune(suffix))
	s := []rune(src)
	rl := len(s) - im
	p, r := []rune(perfix)[0], []rune(suffix)[0]
	i, o := 0, 0
	for ; i < rl; i++ {
		if s[i] == p {
			if string(s[i:i+mi]) == perfix {
				if c == 0 {
					buf.WriteString(string(s[o:i]))
				}
				c++
				o = i + mi
			}
		} else if s[i] == r {
			if string(s[i:i+im]) == suffix {
				c--
				o = i + im
			}
		}
	}
	if c == 0 {
		buf.WriteString(string(s[o:]))
	}
	return buf.String()
}

// 快速Hash
func HashCode(str string) (hash int) {
	for _, v := range str {
		hash = int(v) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}

// 检查值是否在数组中
func IsInArray(arr []string, sub string) bool {
	for _, v := range arr {
		if v == sub {
			return true
		}
	}
	return false
}

// 检查数组元素是否重复
func IsArrayDuplicate(arr []string, arr2 []string) bool {
	if len(arr)*len(arr2) > 1000 {
		return IsArrayDuplicateOpt(arr, arr2)
	}
	// O(n^2)
	for _, v := range arr {
		for _, v2 := range arr2 {
			if v == v2 {
				return true
			}
		}
	}
	return false
}

// 检查数组元素是否重复 - 性能优化版(空间换时间)
func IsArrayDuplicateOpt(arr []string, arr2 []string) bool {
	// O((n1 log n1) + n2) :: len(n2) > len(n1)
	if len(arr) > len(arr2) {
		arr, arr2 = arr2, arr
	}
	m := make(map[string]bool)
	for _, v := range arr {
		m[v] = true
	}
	for _, v := range arr2 {
		if m[v] {
			return true
		}
	}
	return false
}

// MD5 返回binary结果
func Md5(src string) (rst string) {
	h := md5.New()
	io.WriteString(h, src)
	rst = fmt.Sprintf("%s", h.Sum(nil))
	return
}

// MD5 返回 hex结果
func Md5String(src string) (rst string) {
	h := md5.New()
	io.WriteString(h, src)
	rst = fmt.Sprintf("%x", h.Sum(nil))
	return
}

// MD5 返回 hex结果大写
func Md5StringX(src string) (rst string) {
	h := md5.New()
	io.WriteString(h, src)
	rst = fmt.Sprintf("%X", h.Sum(nil))
	return
}

// 转换字符串到int64
func ToInt64(src string) (rst int64) {
	for _, v := range src {
		rst = rst*10 + int64(v-'0')
	}
	return
}

// 转换到字符串
func ToString(src interface{}) string {
	return fmt.Sprint(src)
}

// 新浪微博base62编码
func Base62(src int64) (rst string) {
	for {
		a := src % 62
		rst = string(BASE62_ST[a]) + rst
		src = src / 62
		if src <= 0 {
			break
		}
	}
	return rst
}

// 新浪微博base62解码
func DeBase62(src string) (rst int64) {
	for i, v := range src {
		a := int64(strings.IndexRune(BASE62_ST, v))
		if a < 0 {
			Log("Unknonw Rune:", v, "[", i, "]", "continued")
			continue
		}
		rst = rst*62 + a
	}
	return
}

// 新浪微博base62分组编码
func Base62Url(mid string) (url string) {
	const STEP = 7
	for i := len(mid) - STEP; i > -STEP; i -= STEP {
		if i < 0 {
			url = Base62(ToInt64(mid[0:i+STEP])) + url
		} else {
			url = Base62(ToInt64(mid[i:i+STEP])) + url
		}
	}
	return
}

// 新浪微博base62分组解码
func DeBase62Url(url string) (mid string) {
	const STEP = 4
	for i := len(url) - STEP; i > -STEP; i -= STEP {
		if i < 0 {
			mid = ToString(DeBase62(url[0:i+STEP])) + mid
		} else {
			mid = ToString(DeBase62(url[i:i+STEP])) + mid
		}
	}
	return
}

// 出错崩溃
func Try(e error) {
	if e != nil {
		fup, file, line, _ := runtime.Caller(1)
		fu := runtime.FuncForPC(fup)
		panic(fmt.Sprintf("%s%s\n   at %s(%s:%v)", string([]byte{0xb, 0xad, 0xca, 0xfe}), e.Error(), fu.Name(), path.Base(file), line))
	}
}

// 供给panic后恢复
func Catch() {
	if e := recover(); e != nil {
		es := fmt.Sprint(e)
		if es[:4] == string([]byte{0xb, 0xad, 0xca, 0xfe}) {
			fmt.Println(es[4:])
		} else {
			panic(e)
		}
	}
}
