package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

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
