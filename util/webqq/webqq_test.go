package webqq

import (
	"fmt"
	"github.com/pa001024/MoeWorker/util"
	"testing"
)

var (
	webqq_test = NewWebQQ(2735284921, util.Md5("12345678"))
)

// 测试密码算法
func TestGenPwd(t *testing.T) {
	p := webqq_test.genPwd("test")
	fmt.Println(p)
	if p != "EAF4F2B83C7CD5A59452145A2033CD9E" {
		t.Fail()
	}
}

// 测试Hash算法
func TestGenHash(t *testing.T) {
	webqq_test.PtWebQQ = "619c5c4ca4807d3b27aac8ab3a562d4165948290c4925686acaf73133c6ad727"
	h := webqq_test.genGetUserFriendsHash()
	fmt.Println(h)
	if h != "FCF5FBFA" {
		t.Fail()
	}
}
