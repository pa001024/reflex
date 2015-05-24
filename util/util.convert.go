package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// 转换uint64到BE bytes -> string
func Uint64BEString(v uint64) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return string(b)
}

// 转换BE bytes -> string到uint64
func BEStringUint64(src string) (v uint64) {
	return binary.BigEndian.Uint64([]byte(src))
}

// 判断字符串时都是纯数字
func IsNumber(src string) bool {
	for _, v := range src {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// 转换字符串到int
func ToInt(src string) (rst int) {
	for _, v := range src {
		rst = rst*10 + int(v-'0')
	}
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
func ToString(src ...interface{}) string {
	return fmt.Sprint(src...)
}

// 编码hex
func EncodeHexString(src []byte) string {
	return fmt.Sprintf("%x", src)
}

// 编码hex (大写)
func EncodeHexStringX(src []byte) string {
	return fmt.Sprintf("%X", src)
}

// 解析Host 到IP+端口
func ParseHost(host string) (ip net.IP, port int) {
	n := strings.IndexRune(host, ':')
	if n == -1 {
		ip = net.ParseIP(host)
		port = 80
	} else {
		ip = net.ParseIP(host[:n])
		port = ToInt(host[n+1:])
	}
	return
}
