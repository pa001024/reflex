package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
	"io"
)

// MD5 返回binary结果
func MD5(src string) (rst []byte) {
	h := md5.New()
	io.WriteString(h, src)
	rst = h.Sum(nil)
	return
}

// MD5 返回 hex结果
func MD5String(src string) (rst string) {
	rst = fmt.Sprintf("%x", MD5(src))
	return
}

// MD5 返回 hex结果大写
func MD5StringX(src string) (rst string) {
	rst = fmt.Sprintf("%X", MD5(src))
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
	dst := AES(byteKey, src, isEncode)
	rst = fmt.Sprintf("%x", dst)
	return
}

// AES 返回 hex结果大写
func AESStringX(byteKey, src []byte, isEncode bool) (rst string) {
	dst := AES(byteKey, src, isEncode)
	rst = fmt.Sprintf("%X", dst)
	return
}

// AES-CBC 返回binary结果
func AESCBC(byteKey, src []byte, isEncode bool) (rst []byte) {
	rst = make([]byte, len(src))
	b, err := aes.NewCipher(byteKey)
	if err != nil {
		return
	}
	iv := byteKey[:aes.BlockSize]
	var cc cipher.BlockMode
	var nsrc []byte
	if isEncode {
		cc = cipher.NewCBCEncrypter(b, iv)
		nsrc = PKCS5Padding(src, aes.BlockSize)
	} else {
		cc = cipher.NewCBCDecrypter(b, iv)
		nsrc = src
	}
	dst := make([]byte, len(nsrc))
	cc.CryptBlocks(dst, nsrc)
	if isEncode {
		rst = dst
	} else {
		rst = PKCS5UnPadding(dst)
	}
	return
}

// AES-CBC 返回 hex结果
func AESCBCString(byteKey, src []byte, isEncode bool) (rst string) {
	dst := AESCBC(byteKey, src, isEncode)
	rst = fmt.Sprintf("%x", dst)
	return
}

// AES-CBC 返回 hex结果大写
func AESCBCStringX(byteKey, src []byte, isEncode bool) (rst string) {
	dst := AESCBC(byteKey, src, isEncode)
	rst = fmt.Sprintf("%X", dst)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
