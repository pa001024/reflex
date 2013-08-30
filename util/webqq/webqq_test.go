package webqq

import (
	"encoding/json"
	"fmt"
	"github.com/pa001024/MoeWorker/util"
	"testing"
)

var (
	webqq_test = NewWebQQ(2735284921, GenEncryptPassword(2735284921, util.Md5("12345678")))
)

// 测试密码算法
func TestGenPwd(t *testing.T) {
	p := webqq_test.genPwd("test")
	fmt.Println(p)
	if p != "CA6345A361AE56B1AC3FD3457E3A851A" {
		t.Fail()
	}
}

// 测试Hash算法
func TestGenHash(t *testing.T) {
	webqq_test.ptwebQQ = "619c5c4ca4807d3b27aac8ab3a562d4165948290c4925686acaf73133c6ad727"
	h := webqq_test.genGetUserFriendsHash()
	fmt.Println(h)
	if h != "FCF5FBFA" {
		t.Fail()
	}
}

// 测试Hash算法
func TestDecodeContent(t *testing.T) {
	c := `{"msg_id":11026,"from_uin":2825512944,"to_uin":2735284921,"msg_id2":184082,"msg_type":43,"reply_ip":176756947,"group_code":3447839986,"send_uin":630101664,"seq":435,"time":1377778336,"info_seq":165640562,"content":[["font",{"size":10,"color":"000000","style":[0,0,0],"name":"\u5FAE\u8F6F\u96C5\u9ED1"}],["cface",{"name":"{DC854AD1-A1DB-0855-E267-2CC14D14345A}.gif","file_id":4086436741,"key":"j9v3yUNifMrByvWC","server":"112.90.78.75:0"}],"2",["cface",{"name":"{DB563BEA-1B58-76FB-EF3B-952C78E841F5}.gif","file_id":3746191035,"key":"j9v3yUNifMrByvWC","server":"112.90.78.75:0"}],"3 "]}`
	v := &EventGroupMessage{}
	json.Unmarshal([]byte(c), v)
	cm := v.Content.Decode()
	fmt.Printf("%s\n", cm.Message)
	if cm.Message != "[cface={DC854AD1-A1DB-0855-E267-2CC14D14345A}.gif]2[cface={DB563BEA-1B58-76FB-EF3B-952C78E841F5}.gif]3 " {
		t.Fail()
	}
}
