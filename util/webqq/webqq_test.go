package webqq

import (
	"fmt"
	"github.com/pa001024/MoeWorker/util"
	"testing"
)

var (
	webqq_test = NewWebQQ("youdontneedtoknow", util.Md5("3gEkpEOkUf"))
)

// 测试密码算法
func TestGenPwd(t *testing.T) {
	p := webqq_test.GenPwd(string([]byte{0x00, 0x00, 0x00, 0x00, 0xa3, 0x09, 0x22, 0xb9}), "ztwb")
	if p != "583D72C7D5B1D0064F1C8C86E5D4FE1E" {
		t.Fail()
	}
}
