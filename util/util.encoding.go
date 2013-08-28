package util

import (
	"strings"
)

const (
	_BASE62_ST = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// 新浪微博base62编码
func Base62(src int64) (rst string) {
	for {
		a := src % 62
		rst = string(_BASE62_ST[a]) + rst
		src = src / 62
		if src <= 0 {
			break
		}
	}
	return rst
}

// 新浪微博base62解码
func DeBase62(src string) (rst int64) {
	for _, v := range src {
		a := int64(strings.IndexRune(_BASE62_ST, v))
		if a < 0 {
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
