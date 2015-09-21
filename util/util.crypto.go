package util

import (
	"crypto/aes"
	"crypto/md5"
	"fmt"
	"io"
)

// MD5 返回binary结果
func Md5(src string) (rst []byte) {
	h := md5.New()
	io.WriteString(h, src)
	rst = h.Sum(nil)
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

// AES 返回binary结果
func AES(byteKey, src []byte, isEncode bool) (rst []byte) {
	rst = make([]byte, len(src))
	c, err := aes.NewCipher(byteKey)
	if err != nil {
		return
	}
	if isEncode {
		c.Encrypt(rst, src)
	} else {
		c.Decrypt(rst, src)
	}
	return
}

// AES 返回 hex结果
func AESString(byteKey, src []byte, isEncode bool) (rst string) {
	dst := make([]byte, len(src))
	c, err := aes.NewCipher(byteKey)
	if err != nil {
		return
	}
	if isEncode {
		c.Encrypt(dst, src)
	} else {
		c.Decrypt(dst, src)
	}
	rst = fmt.Sprintf("%x", dst)
	return
}

// AES 返回 hex结果大写
func AESStringX(byteKey, src []byte, isEncode bool) (rst string) {
	dst := make([]byte, len(src))
	c, err := aes.NewCipher(byteKey)
	if err != nil {
		return
	}
	if isEncode {
		c.Encrypt(dst, src)
	} else {
		c.Decrypt(dst, src)
	}
	rst = fmt.Sprintf("%X", dst)
	return
}
