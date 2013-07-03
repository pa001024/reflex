package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"
)

var (
	DEBUG = false
	IP    string
)

const (
	BASE62_ST = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	ip, _ := CheckIP()
	IP = ip
}
func HashCode(str string) (hash int) {
	for _, v := range str {
		hash = int(v) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}

func IsInArray(arr []string, sub string) bool {
	for _, v := range arr {
		if v == sub {
			return true
		}
	}
	return false
}
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
func IsArrayDuplicateOpt(arr []string, arr2 []string) bool {
	// O((n1 log n1) + n2)
	m := make(map[string]bool)
	for _, v := range arr {
		m[v] = true
	}
	for _, v := range arr2 {
		if m[v] == true {
			return true
		}
	}
	return false
}
func Md5String(src string) (rst string) {
	h := md5.New()
	io.WriteString(h, src)
	rst = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func ToInt64(src string) (rst int64) {
	for _, v := range src {
		rst = rst*10 + int64(v-'0')
	}
	return
}
func ToString(src interface{}) string {
	return fmt.Sprint(src)
}
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

func CheckIP() (ip string, err error) {
	res, err := http.Get("http://checkip.dyndns.com/")
	if err != nil {
		panic("CheckIP: App is Offine! Dead.")
	}
	bin, err := ioutil.ReadAll(res.Body)
	str := string(bin)
	if len(str) < 92 {
		panic("CheckIP: Bad Response!")
	}
	ip = str[76 : len(str)-14]
	return
}

func Log(args ...interface{}) {
	if DEBUG {
		DebugLog(args...)
	} else {
		TestLog(args...)
	}
}

func TestLog(args ...interface{}) {
	log.Println(args...)
}
func DebugLog(args ...interface{}) {
	fup, file, line, _ := runtime.Caller(2)
	fu := runtime.FuncForPC(fup)
	log.Println(args...)
	fmt.Println("        at", fu.Name()+"("+path.Base(file)+":"+ToString(line)+")")
}
