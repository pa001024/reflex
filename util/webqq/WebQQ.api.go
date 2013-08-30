package webqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pa001024/MoeWorker/util"
)

const (
	_API_URL     = "http://s.web2.qq.com/api/"
	_API_REFERER = "http://s.web2.qq.com/proxy.html?v=20110412001&callback=1&id=3"
)

// 通用API接口(GET)
func (this *WebQQ) api(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.vfwebqq},
		"t":       {util.JsCurrentTime()},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.getWithReferer(_API_URL+api+"?"+val.Encode(), _API_REFERER)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用API接口(POST)
func (this *WebQQ) postApi(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"r": {util.ToJson(
			append(args, "vfwebqq", this.vfwebqq)...,
		)},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.postFormWithReferer(_API_URL+api, _API_REFERER, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用API接口(POST2)
func (this *WebQQ) rawPostApi(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.vfwebqq},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.getWithReferer(_API_URL+api+"?"+val.Encode(), _API_REFERER)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用 Result结构
type Result struct {
	Code int `json:"retcode"`
}

// 获取群成员在线信息
func (this *WebQQ) get_group_member_stat2(gcode GCode) (v *ResultGroupMemberStat, err error) {
	util.DEBUG.Logf("get_group_member_stat2(gcode = %v)", gcode)
	data, err := this.api("get_group_member_stat2", "gcode", gcode)
	if err == nil {
		v = &ResultGroupMemberStat{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 群成员在线信息 result结构
type ResultGroupMemberStat struct {
	Code   int `json:"retcode"`
	Result struct {
		Stats []MemberStat
	} `json:"result"`
}

// 获取自己的群名片
func (this *WebQQ) get_self_business_card2(gcode string) (v *ResultSelfBusinessCard, err error) {
	util.DEBUG.Logf("get_self_business_card2(gcode = %v)", gcode)
	data, err := this.api("get_self_business_card2", "gcode", gcode)
	if err == nil {
		v = &ResultSelfBusinessCard{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 自己群名片 result结构
type ResultSelfBusinessCard struct {
	Code   int `json:"retcode"`
	Result struct {
	} `json:"result"`
}

// 获取群信息
// gcode:[738328699] retainKey:memo,gcode
func (this *WebQQ) get_group_info(gcode, retainKey string) (v *ResultGroupInfo, err error) {
	util.DEBUG.Logf("get_group_info(gcode = %v, retainKey = %v)", gcode, retainKey)
	data, err := this.api("get_group_info", "gcode", gcode, "retainKey", retainKey)
	if err == nil {
		v = &ResultGroupInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 群信息 result结构
type ResultGroupInfo struct {
	Code   int `json:"retcode"`
	Result []struct {
		Memo  string `json:"memo"`
		GCode GCode  `json:"gcode"`
	} `json:"result"`
}

// 获取群详细信息
func (this *WebQQ) get_group_info_ext2(gcode string) (v *ResultGroupInfoExt, err error) {
	util.DEBUG.Logf("get_group_info_ext2(gcode = %v)", gcode)
	data, err := this.api("get_group_info_ext2", "gcode", gcode)
	if err == nil {
		v = &ResultGroupInfoExt{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 群详细信息 result结构
type ResultGroupInfoExt struct {
	Code   int `json:"retcode"`
	Result struct {
		Stats       []MemberStat `json:"stats"`
		MemberInfos []MemberInfo `json:"minfo"`
		GroupInfo   GroupInfo    `json:"ginfo"`
		Cards       []Card       `json:"cards"`
		VipInfo     []VipInfo    `json:"vipinfo"`
	} `json:"result"`
}

type ClientType uint32

// 客户端类型 (TODO:待补完)
const (
	ClientTypePC    ClientType = 1
	ClientTypeWebQQ            = 41
)

// 群成员在线信息
type MemberStat struct {
	Uin        Uin    `json:"uin"`
	ClientType uint32 `json:"client_type"`
	Stat       uint32 `json:"stat"`
}

// 群成员信息
type MemberInfo struct {
	Nick     string `json:"nick"`
	Province string `json:"province"`
	Gender   string `json:"gender"`
	Uin      Uin    `json:"uin"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

// 群详细信息
type GroupInfo struct {
	Face       uint32 `json:"face"`
	Class      uint32 `json:"class"`
	Memo       string `json:"memo"`
	Fingermemo string `json:"fingermemo"`
	Code       GCode  `json:"code"` // 群的信息ID gcode
	Createtime uint64 `json:"createtime"`
	Flag       uint32 `json:"flag"`
	Level      uint32 `json:"level"`
	Name       string `json:"name"`  // 群名
	Gid        Uin    `json:"gid"`   // 群的消息Uin group_uin
	Owner      Uin    `json:"owner"` // 群主Uin
	Members    string `json:"members"`
}

// 群成员 0 = 普通成员 84 = 群主 01010100
type Member struct {
	Uin  Uin    `json:"muin"`
	Flag uint32 `json:"mflag"`
}

// 群名片
type Card struct {
	Uin  Uin    `json:"muin"`
	Card string `json:"card"`
}

// 会员信息
type VipInfo struct {
	VipLevel string `json:"vip_level"`
	Uin      Uin    `json:"u"`
	IsVip    string `json:"is_vip"`
}

// 获取显示ID type: 好友=1 群=4
func (this *WebQQ) get_friend_uin2(tuin Uin, type_ string) (v *ResultFriendUin, err error) {
	util.DEBUG.Logf("get_friend_uin2(tuin = %v, type = %v)", tuin, type_)
	data, err := this.api("get_friend_uin2", "tuin", tuin, "type", type_)
	if err == nil {
		v = &ResultFriendUin{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取显示ID result结构
type ResultFriendUin struct {
	Code   int `json:"retcode"`
	Result struct {
		Uin     Uin     `json:"uin"`
		Account Account `json:"account"`
		Uiuin   string  `json:"uiuin"`
	} `json:"result"`
}

// 获取好友信息
func (this *WebQQ) get_friend_info2(tuin string) (v *ResultFriendInfo, err error) {
	util.DEBUG.Logf("get_friend_uin2(tuin = %v)", tuin)
	data, err := this.api("get_friend_info2", "tuin", tuin)
	if err == nil {
		v = &ResultFriendInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取好友信息 result结构
type ResultFriendInfo struct {
	Code   int        `json:"retcode"`
	Result FriendInfo `json:"result"`
}

// 用户信息
type FriendInfo struct {
	Uin        Uin      `json:"uin"`
	Face       uint32   `json:"face"`
	Allow      uint32   `json:"allow"`
	Constel    uint32   `json:"constel"`
	Blood      uint32   `json:"blood"`
	Stat       uint32   `json:"stat"`
	VipInfo    uint32   `json:"vip_info"`
	Shengxiao  uint32   `json:"shengxiao"`
	ClientType uint32   `json:"client_type"`
	Birthday   Birthday `json:"birthday"`
	Occupation string   `json:"occupation"`
	Phone      string   `json:"phone"`
	College    string   `json:"college"`
	Homepage   string   `json:"homepage"`
	Country    string   `json:"country"`
	City       string   `json:"city"`
	Personal   string   `json:"personal"`
	Nick       string   `json:"nick"`
	Email      string   `json:"email"`
	Province   string   `json:"province"`
	Gender     string   `json:"gender"`
	Mobile     string   `json:"mobile"`
}

// 生日
type Birthday struct {
	Month uint32 `json:"month"`
	Year  uint32 `json:"year"`
	Day   uint32 `json:"day"`
}

// 获取好友验证类型
func (this *WebQQ) get_allow_info2(tuin string) (v *ResultAllowInfo, err error) {
	data, err := this.api("get_allow_info2", "tuin", tuin, "retainKey", "allow")
	if err == nil {
		v = &ResultAllowInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 好友验证类型 result结构
type ResultAllowInfo struct {
	Code   int `json:"retcode"`
	Result struct {
		Allow uint32 `json:"allow"`
	} `json:"result"`
}

// 获取签名
func (this *WebQQ) get_single_long_nick2(tuin string) (v *ResultSingleLongNick, err error) {
	data, err := this.api("get_single_long_nick2", "tuin", tuin)
	if err == nil {
		v = &ResultSingleLongNick{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取签名 result结构
type ResultSingleLongNick struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin      Uin    `json:"uin"`
		LongNick string `json:"lnick"`
	} `json:"result"`
}

// 获取等级
func (this *WebQQ) get_qq_level2(tuin string) (v *ResultQQLevel, err error) {
	data, err := this.api("get_qq_level2", "tuin", tuin)
	if err == nil {
		v = &ResultQQLevel{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取等级 result结构
type ResultQQLevel struct {
	Code   int `json:"retcode"`
	Result struct {
		Uin        Uin    `json:"tuin"`
		Level      uint32 `json:"level"`
		Days       uint32 `json:"days"`
		Hours      uint32 `json:"hours"`
		RemainDays uint32 `json:"remainDays"`
	} `json:"result"`
}

// 获取讨论组列表
func (this *WebQQ) get_discus_list() (v *ResultDiscusList, err error) {
	data, err := this.api("get_discus_list")
	if err == nil {
		v = &ResultDiscusList{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 讨论组列表 result结构
type ResultDiscusList struct {
	Code   int `json:"retcode"`
	Result struct {
		Dnamelist []interface{} `json:"dnamelist"`
	} `json:"result"`
}

// 设置签名(POST)
func (this *WebQQ) set_long_nick2(nlk string) (v *Result, err error) {
	data, err := this.postApi("set_long_nick2", "nlk", nlk)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取好友列表(POST)
func (this *WebQQ) get_user_friends2() (v *ResultUserFrientds, err error) {
	data, err := this.postApi("get_user_friends2",
		"h", "hello",
		"hash", this.genGetUserFriendsHash(),
	)
	if err == nil {
		v = &ResultUserFrientds{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取好友列表 result结构
type ResultUserFrientds struct {
	Code   int `json:"retcode"`
	Result struct {
		Friends    []Friend           `json:"friends"`
		Marknames  []MarkName         `json:"marknames"`
		Categories []FrientdCategorie `json:"categories"`
		VipInfo    []VipInfo          `json:"vipinfo"`
		Info       []FrientdInfo      `json:"info"`
	} `json:"result"`
}

// 好友信息子结构
type Friend struct {
	Uin        Uin   `json:"uin"`
	Flag       int32 `json:"flag"`
	Categories int32 `json:"categories"`
}

// 备注
type MarkName struct {
	Uin      Uin    `json:"uin"`
	Markname string `json:"markname"`
}

// 好友分类子结构
type FrientdCategorie struct {
	Index uint32 `json:"index"`
	Sort  uint32 `json:"sort"`
	Name  string `json:"name"`
}

// 用户简单信息
type FrientdInfo struct {
	Uin  Uin    `json:"uin"`
	Face uint32 `json:"face"`
	Flag uint32 `json:"flag"`
	Nick string `json:"nick"`
}

// 修改群名片(POST)
func (this *WebQQ) update_group_info2(gcode GCode, markname string) (v *Result, err error) {
	data, err := this.postApi("update_group_info2", "gcode", gcode, "markname", markname)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 退出群(POST)
func (this *WebQQ) quit_group2(gcode GCode) (v *Result, err error) {
	data, err := this.postApi("quit_group2", "gcode", gcode)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取群列表(POST)
func (this *WebQQ) get_group_name_list_mask2() (v *ResultGroupNameListMask, err error) {
	data, err := this.postApi("get_group_name_list_mask2")
	if err == nil {
		v = &ResultGroupNameListMask{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 群列表 result结构
type ResultGroupNameListMask struct {
	Code   int `json:"retcode"`
	Result struct {
		GroupMasks []interface{} `json:"gmasklist"` // TODO
		GroupNames []GroupName   `json:"gnamelist"`
		GroupMarks []MarkName    `json:"gmarklist"`
	} `json:"result"`
}

// 群名
type GroupName struct {
	Gid  Uin    `json:"gid"`
	Code GCode  `json:"code"`
	Flag uint32 `json:"flag"`
	Name string `json:"name"`
}

// 查找好友(需要验证码)
func (this *WebQQ) search_qq_by_uin2(tuin, verifysession, code string) (v *Result, err error) {
	data, err := this.api("search_qq_by_uin2", "tuin", tuin, "verifysession", verifysession, "code", code)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// [POST] 批量获取VIP信息
func (this *WebQQ) batch_get_vipinfo(ul []Uin) (v *ResultBatchGetVipInfo, err error) {
	b, _ := json.Marshal(ul)
	data, err := this.rawPostApi("batch_get_vipinfo", "ul", string(b))
	if err == nil {
		v = &ResultBatchGetVipInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取VIP信息 result结构
type ResultBatchGetVipInfo struct {
	Code   int `json:"retcode"`
	Result struct {
		VipInfo []VipInfo `json:"vipinfo"`
	} `json:"result"`
}

// 获取陌生人信息
func (this *WebQQ) get_stranger_info2(tuin Uin) (v *ResultStrangerInfo, err error) {
	data, err := this.api("get_stranger_info2", "tuin", tuin, "gid", "0", "code", "")
	if err == nil {
		v = &ResultStrangerInfo{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取陌生人信息 result结构
type ResultStrangerInfo struct {
	Code   int        `json:"retcode"`
	Result FriendInfo `json:"result"`
}

// 添加好友(验证)
func (this *WebQQ) add_need_verify2(account, myallow, groupid, msg, token string) (v *Result, err error) {
	data, err := this.api("add_need_verify2",
		"account", account,
		"myallow", myallow,
		"groupid", groupid,
		"msg", msg,
		"token", token,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}
