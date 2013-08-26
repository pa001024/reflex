package webqq

import (
	"fmt"
	"github.com/pa001024/MoeWorker/util"
	"testing"
)

var (
	webqq_test = NewWebQQ("2735284921", util.Md5("3gEkpEOkUf"))
)

// 测试密码算法
func TestGenPwd(t *testing.T) {
	p := webqq_test.GenPwd(string([]byte{0x00, 0x00, 0x00, 0x00, 0xa3, 0x09, 0x22, 0xb9}), "awsz")
	fmt.Println(p)
	if p != "016A97B8E28587F7AFA2C66496948A65" {
		t.Fail()
	}
}

// 测试Hash算法
func TestGenHash(t *testing.T) {
	webqq_test.Uin = Uin(2735284921)
	webqq_test.PtWebQQ = "619c5c4ca4807d3b27aac8ab3a562d4165948290c4925686acaf73133c6ad727"
	h := webqq_test.GenHash()
	fmt.Println(h)
	if h != "FCF5FBFA" {
		t.Fail()
	}
}
