package webqq

import (
	"fmt"
	"math/rand"
)

// 生成随机数
func rand_r() string {
	return fmt.Sprint(rand.ExpFloat64())
}
