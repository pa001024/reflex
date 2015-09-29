package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
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

// 返回当前秒数
func JsCurrentSecond() int64 {
	return time.Now().UnixNano() / 1e9
}

// 返回JS风格的当前时间字符串
func JsCurrentTime() string {
	return fmt.Sprint(time.Now().UnixNano() / 1e6)
}

// 解析JS风格的64位长整形
func DecodeJsUint64BE(src string) (uin uint64) {
	b := DecodeJsHex(src)
	return binary.BigEndian.Uint64(b)
}

// 返回JS风格的64位长整形 形如\x00\x00\x00\x00\x00\x00\x00\x01 (1)
func EncodeJsUint64BE(uin uint64) (src string) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, uin)
	return EncodeJsHex(buf.Bytes())
}

// 解析js hex值 形如 \x00
func DecodeJsHex(src string) []byte {
	src = strings.Replace(src, `\x`, "", -1)
	v, _ := hex.DecodeString(src)
	return v
}

// 编码js hex值 形如 \x00
func EncodeJsHex(src []byte) string {
	buf := &bytes.Buffer{}
	for _, v := range src {
		fmt.Fprintf(buf, "\\x%02x", v)
	}
	return buf.String()
}
