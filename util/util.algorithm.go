package util

import (
	"bytes"
	"fmt"
	"math/rand"
)

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

// 生成N位随机数字
func Random(n int) (result string) {
	for _, v := range rand.Perm(n) {
		result += fmt.Sprint(v)
	}
	return
}

// 三元选择
func Sw(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
