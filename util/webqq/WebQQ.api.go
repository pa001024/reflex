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
func (this *WebQQ) api(api string, args ...string) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.VerifyCode},
		"t":       {fmt.Sprint(time.Now().UnixNano() / 1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i], args[i+1])
	}
	res, err := this.GetWithReferer(API_URL + api + "?" + val.Encode())
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用API接口
func (this *WebQQ) postApi(api string, args ...string) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.VerifyCode},
		"t":       {fmt.Sprint(time.Now().UnixNano() / 1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i], args[i+1])
	}
	res, err := this.PostFormWithReferer(API_URL+api, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 获取群成员信息
/*
{"retcode":0,"result":{"stats":[{"client_type":41,"uin":2735284921,"stat":10}],"gcode":738328699}}
*/
func (this *WebQQ) get_group_member_stat2(gcode string) {
	this.api("get_group_member_stat2", "gcode", gcode)
}

// 获取自己的群名片
/*
{"retcode":100003}
*/
func (this *WebQQ) get_self_business_card2(gcode string) {
	this.api("get_self_business_card2", "gcode", gcode)
}

// 获取群信息
/*
gcode:[738328699]
retainKey:memo,gcode
{"retcode":0,"result":[{"memo":"","gcode":738328699}]}
*/
func (this *WebQQ) get_group_info(gcode, retainKey string) {
	this.api("get_group_info", "gcode", gcode, "retainKey", retainKey)
}

// 获取群详细信息
/*
{"retcode":0,"result":{"stats":[{"client_type":41,"uin":2735284921,"stat":10}],"minfo":[{"nick":"XpAhH","province":"","gender":"unknown","uin":3255435951,"country":"蒙古","city":"东方"},{"nick":"金馆长!","province":"","gender":"female","uin":873450563,"country":"蒙古","city":"东方"},{"nick":"菊菊菊菊菊菊","province":"浙江","gender":"male","uin":2735284921,"country":"中国","city":"舟山"}],"ginfo":{"face":0,"memo":"","class":10021,"fingermemo":"","code":738328699,"createtime":1311148090,"flag":16777217,"level":3,"name":".com","gid":221664830,"owner":3255435951,"members":[{"muin":3255435951,"mflag":68},{"muin":873450563,"mflag":0},{"muin":2735284921,"mflag":4}],"option":2},"vipinfo":[{"vip_level":5,"u":3255435951,"is_vip":1},{"vip_level":0,"u":873450563,"is_vip":0},{"vip_level":0,"u":2735284921,"is_vip":0}]}}
*/
func (this *WebQQ) get_group_info_ext2(gcode string) {
	this.api("get_group_info_ext2", "gcode", gcode)
}

// 获取显示ID
/*
type: 好友=1 群=4
{"retcode":0,"result":{"uiuin":"","account":165640562,"uin":738328699}}
*/
func (this *WebQQ) get_friend_uin2(tuin, type_ string) {
	this.api("get_friend_uin2", "tuin", tuin, "type", type_)
}

// 获取好友信息
/*
{"retcode":0,"result":{"face":552,"birthday":{"month":10,"year":1999,"day":1},"occupation":"","phone":"","allow":0,"college":"","uin":2735284921,"constel":9,"blood":0,"homepage":"","stat":10,"vip_info":0,"country":"中国","city":"舟山","personal":"","nick":"菊菊菊菊菊菊","shengxiao":4,"email":"","client_type":41,"province":"浙江","gender":"male","mobile":""}}
*/
func (this *WebQQ) get_friend_info2(tuin string) {
	this.api("get_friend_info2", "tuin", tuin)
}

// 获取好友验证类型
/*
?retainKey=allow
{"retcode":0,"result":{"allow":1}}
*/
func (this *WebQQ) get_allow_info2(tuin, type_ string) {
	this.api("get_allow_info2", "tuin", tuin, "retainKey", "allow")
}

// 获取签名
/*
{"retcode":100000} // 没有签名
{"retcode":0,"result":[{"uin":3255435951,"lnick":"AFK"}]}
*/
func (this *WebQQ) get_single_long_nick2(tuin string) {
	this.api("get_single_long_nick2", "tuin", tuin)
}

// 获取等级
/*
{"retcode":100000} // 没有签名
{"retcode":0,"result":{"level":39,"days":1728,"hours":21201,"remainDays":32,"tuin":3255435951}}
*/
func (this *WebQQ) get_qq_level2(tuin string) {
	this.api("get_qq_level2", "tuin", tuin)
}

// 设置签名(POST)
/*
POST:
r={"nlk":"qwwww","vfwebqq":"..."}
RESULT:
{"retcode":0,"result":{"result":0}}
*/
func (this *WebQQ) set_long_nick2(nlk string) {
	this.postApi("set_long_nick2", "nlk", nlk)
}

// 获取好友列表(POST)
/*
POST:
r={"h":"hello","hash":"ABF6A3FE","vfwebqq":"..."}
RESULT:
{"retcode":0,"result":{"friends":[],"marknames":[],"categories":[],"vipinfo":[],"info":[]}}
*/
func (this *WebQQ) get_user_friends2(tuin string) {
	this.postApi("get_user_friends2", "tuin", tuin)
}

// 获取群列表(POST)
/*
POST:
r={"vfwebqq":"..."}
RESULT:
{"retcode":0,"result":{"gmasklist":[],"gnamelist":[{"flag":16777217,"name":".com","gid":221664830,"code":738328699}],"gmarklist":[]}}
*/
func (this *WebQQ) get_group_name_list_mask2(tuin string) {
	this.postApi("get_group_name_list_mask2", "tuin", tuin)
}

// 查找好友(需要验证码)
/*
{"retcode":0,"result":{"face":147,"birthday":{"month":4,"year":1996,"day":9},"occupation":"待业/无业/失业","phone":"","allow":4,"college":"","constel":3,"blood":0,"stat":20,"homepage":"","country":"中国","city":"温州","uiuin":"","personal":"爱上你，不是因为你给了我需要的东西，而是因为你给了我从未有过的感觉。","nick":"小さくて暗い","shengxiao":1,"email":"","token":"e5c67a587252c899fac04fce508372b0e650b1622faa7230","province":"浙江","account":1159549778,"gender":"male","tuin":2265896498,"mobile":"-"}}
*/
func (this *WebQQ) search_qq_by_uin2(tuin, verifysession, code string) {
	this.api("search_qq_by_uin2", "tuin", tuin, "verifysession", verifysession, "code", code)
}
