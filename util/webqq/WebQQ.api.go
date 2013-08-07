package webqq

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

const (
	API_URL = "http://s.web2.qq.com/api/"
)

// 通用API接口
func (this *WebQQ) GetAPI(api string, args ...string) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.VerifyCode},
		"t":       {fmt.Sprintf("%d", time.Now().UnixNano()/1e6)},
	}
	l = len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i], args[i+1])
	}
	res, err := this.GetWithReferer(API_URL + api + "?" + val.Encode())
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll()
	return
}

// 获取群成员信息
/*
{"retcode":0,"result":{"stats":[{"client_type":41,"uin":2735284921,"stat":10}],"gcode":738328699}}
*/
func (this *WebQQ) get_group_member_stat2(gcode string) {
	this.GetAPI("get_group_member_stat2", "gcode", gcode)
}

// 获取自己信息
/*
{"retcode":100003}
*/
func (this *WebQQ) get_self_business_card2(gcode string) {
	this.GetAPI("get_self_business_card2", "gcode", gcode)
}

// 获取自己信息
/*
{"retcode":0,"result":{"stats":[{"client_type":41,"uin":2735284921,"stat":10}],"minfo":[{"nick":"XpAhH","province":"","gender":"unknown","uin":3255435951,"country":"蒙古","city":"东方"},{"nick":"金馆长!","province":"","gender":"female","uin":873450563,"country":"蒙古","city":"东方"},{"nick":"菊菊菊菊菊菊","province":"浙江","gender":"male","uin":2735284921,"country":"中国","city":"舟山"}],"ginfo":{"face":0,"memo":"","class":10021,"fingermemo":"","code":738328699,"createtime":1311148090,"flag":16777217,"level":3,"name":".com","gid":221664830,"owner":3255435951,"members":[{"muin":3255435951,"mflag":68},{"muin":873450563,"mflag":0},{"muin":2735284921,"mflag":4}],"option":2},"vipinfo":[{"vip_level":5,"u":3255435951,"is_vip":1},{"vip_level":0,"u":873450563,"is_vip":0},{"vip_level":0,"u":2735284921,"is_vip":0}]}}
*/
func (this *WebQQ) get_group_info_ext2(gcode string) {
	this.GetAPI("get_group_info_ext2", "gcode", gcode)
}

// 获取自己信息
/*
{"retcode":0,"result":{"uiuin":"","account":165640562,"uin":738328699}}
*/
func (this *WebQQ) get_friend_uin2(tuin string, type_ string) {
	this.GetAPI("get_friend_uin2", "tuin", tuin, "type", type_)
}
