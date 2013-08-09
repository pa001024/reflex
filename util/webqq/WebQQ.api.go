package webqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/pa001024/MoeWorker/util"
)

// 用户ID
type Uin uint64

func (u Uin) String() string { return fmt.Sprint(uint64(u)) }

// 群信息ID
type GCode uint64

func (u GCode) String() string { return fmt.Sprint(uint64(u)) }

// 真实QQ号
type Account uint64

func (u Account) String() string { return fmt.Sprint(uint64(u)) }

const (
	API_URL = "http://s.web2.qq.com/api/"
)

// 通用API接口(GET)
func (this *WebQQ) api(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"vfwebqq": {this.VerifyCode},
		"t":       {fmt.Sprint(time.Now().UnixNano() / 1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.GetWithReferer(API_URL + api + "?" + val.Encode())
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
			append(args, "vfwebqq", this.VerifyCode)...,
		)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.PostFormWithReferer(API_URL+api, val)
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

// 获取群成员信息
func (this *WebQQ) get_group_member_stat2(gcode GCode) (v *ResultGroupMemberStat, err error) {
	DEBUG.Logf("get_group_member_stat2(gcode = %v)", gcode)
	data, err := this.api("get_group_member_stat2", "gcode", gcode)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 群成员信息 result结构
 ---------------------

 {"retcode":0,"result":{"stats":[{"client_type":41,"uin":2735284921,"stat":10}],"gcode":738328699}}
*/
type ResultGroupMemberStat struct {
	Code   int `json:"retcode"`
	Result struct {
		Stats []Stat
	} `json:"result"`
}

/*
 获取自己的群名片
 ----------------

 {"retcode":10001} // 没有
*/
func (this *WebQQ) get_self_business_card2(gcode string) (v *ResultSelfBusinessCard, err error) {
	DEBUG.Logf("get_self_business_card2(gcode = %v)", gcode)
	data, err := this.api("get_self_business_card2", "gcode", gcode)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

type ResultSelfBusinessCard struct {
	Code   int `json:"retcode"`
	Result struct {
	} `json:"result"`
}

/*
 获取群信息
 ----------

 gcode:[738328699]
 retainKey:memo,gcode
*/
func (this *WebQQ) get_group_info(gcode, retainKey string) (v *ResultGroupInfo, err error) {
	DEBUG.Logf("get_group_info(gcode = %v, retainKey = %v)", gcode, retainKey)
	data, err := this.api("get_group_info", "gcode", gcode, "retainKey", retainKey)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 群信息 result结构
 -----------------

 {"retcode":0,"result":[{"memo":"","gcode":738328699}]}
*/
type ResultGroupInfo struct {
	Code   int `json:"retcode"`
	Result []struct {
		Memo  string `json:"memo"`
		GCode GCode  `json:"gcode"`
	} `json:"result"`
}

// 获取群详细信息
func (this *WebQQ) get_group_info_ext2(gcode string) (v *ResultGroupInfoExt, err error) {
	DEBUG.Logf("get_group_info_ext2(gcode = %v)", gcode)
	data, err := this.api("get_group_info_ext2", "gcode", gcode)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 群详细信息 result结构
 ---------------------

 {"retcode":0,"result":{"stats":[{...}],"minfo":[{...}],"ginfo":{...},"cards":[{...}],"vipinfo":[{...}]}}
*/
type ResultGroupInfoExt struct {
	Code   int `json:"retcode"`
	Result struct {
		Stats       []Stat       `json:"stats"`
		MemberInfos []MemberInfo `json:"minfo"`
		GroupInfo   GroupInfo    `json:"ginfo"`
		Cards       []Card       `json:"cards"`
		VipInfo     []VipInfo    `json:"vipinfo"`
	} `json:"result"`
}

// 客户端类型 (TODO:待补完)
const (
	ClientTypePC    = 1
	ClientTypeWebQQ = 41
)

/*
 群成员在线信息
 --------------

 {"client_type": 1,"uin": 3255435951,"stat": 10}
*/
type Stat struct {
	Uin        Uin    `json:"uin"`
	ClientType uint32 `json:"client_type"`
	Stat       uint32 `json:"stat"`
}

/*
 群成员信息
 ----------

 {"nick": "菊菊菊菊菊菊","province": "省","gender": "male","uin": 2735284921,"country": "国家","city": "城市"}
*/
type MemberInfo struct {
	Nick     string `json:"nick"`
	Province string `json:"province"`
	Gender   string `json:"gender"`
	Uin      Uin    `json:"uin"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

/*
 群详细信息
 ----------

 {
 	"face": 0,
 	"memo": "",
 	"class": 10021,
 	"fingermemo": "",
 	"code": 738328699,
 	"createtime": 1311148090,
 	"flag": 16777217,
 	"level": 3,
 	"name": "群名",
 	"gid": 221664830,
 	"owner": 3255435951,
 	"members": [{...}],
 	"option": 2
 }
*/
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

/*
 群成员
 ------

 0 = 普通成员
 84 = 群主 01010100

 {"muin": 3255435951,"mflag": 84}
*/
type Member struct {
	Uin  Uin    `json:"muin"`
	Flag uint32 `json:"mflag"`
}

/*
 群名片
 ------

 {"muin": 2735284921,"card": "花花花花"}
*/
type Card struct {
	Uin  Uin    `json:"muin"`
	Card string `json:"card"`
}

/*
 会员信息
 --------

 {"vip_level":5,"u":3255435951,"is_vip":1}
*/
type VipInfo struct {
	VipLevel string `json:"vip_level"`
	Uin      Uin    `json:"u"`
	IsVip    string `json:"is_vip"`
}

/*
 获取显示ID
 ----------

 type: 好友=1 群=4
*/
func (this *WebQQ) get_friend_uin2(tuin Uin, type_ string) (v *ResultFriendUin, err error) {
	DEBUG.Logf("get_friend_uin2(tuin = %v, type = %v)", tuin, type_)
	data, err := this.api("get_friend_uin2", "tuin", tuin, "type", type_)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取显示ID result结构
 ---------------------

 {"retcode":0,"result":{"uiuin":"","account":165640562,"uin":738328699}}
*/
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
	DEBUG.Logf("get_friend_uin2(tuin = %v)", tuin)
	data, err := this.api("get_friend_info2", "tuin", tuin)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取好友信息 result结构
 -----------------------

 {
	"retcode": 0,
	"result": {
		"uin": 2735284921,
		"face": 552,
		"allow": 0,
		"constel": 9,
		"blood": 0,
		"stat": 10,
		"vip_info": 0,
		"shengxiao": 4,
		"client_type": 41,
		"birthday": {...},
		"occupation": "",
		"phone": "",
		"college": "",
		"homepage": "",
		"country": "中国",
		"city": "舟山",
		"personal": "",
		"nick": "菊菊菊菊菊菊",
		"email": "",
		"province": "浙江",
		"gender": "male",
		"mobile": ""
	}
 }
*/
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

/*
 生日
 ----

 {"month": 10,"year": 1999,"day": 1}
*/
type Birthday struct {
	Month uint32 `json:"month"`
	Year  uint32 `json:"year"`
	Day   uint32 `json:"day"`
}

/*
 获取好友验证类型
 ----------------

 ?retainKey=allow
*/
func (this *WebQQ) get_allow_info2(tuin string) (v *ResultAllowInfo, err error) {
	data, err := this.api("get_allow_info2", "tuin", tuin, "retainKey", "allow")
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 好友验证类型 result结构
 -----------------------

 {"retcode":0,"result":{"allow":1}}
*/
type ResultAllowInfo struct {
	Code   int `json:"retcode"`
	Result struct {
		Allow uint32 `json:"allow"`
	} `json:"result"`
}

// 获取签名
func (this *WebQQ) get_single_long_nick2(tuin string) (v *ResultSingleLongNick, err error) {
	data, err := this.api("get_single_long_nick2", "tuin", tuin)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取签名 result结构
 -------------------

 {"retcode":100000} // 没有签名
 {"retcode":0,"result":[{"uin":3255435951,"lnick":"AFK"}]}
*/
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
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取等级 result结构
 -------------------

 {"retcode":100000} // 没有签名
 {"retcode":0,"result":{"level":39,"days":1728,"hours":21201,"remainDays":32,"tuin":3255435951}}
*/
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

// TODO: 获取不知什么玩意 先不写
func (this *WebQQ) get_discus_list() (v *ResultDiscusList, err error) {
	WARN.Log("calling an todo function get_discus_list()")
	// this.api("get_discus_list")
	return
}

/*
 上面这不知什么玩意的 result结构
 -------------------------------

 {"retcode":0,"result":{"dnamelist":[]}}
*/
type ResultDiscusList struct {
	Code   int `json:"retcode"`
	Result struct {
		Dnamelist []interface{} `json:"dnamelist"`
	} `json:"result"`
}

/*
 设置签名(POST)
 --------------

 r = {"nlk":"签名","vfwebqq":"..."}
 RESULT:
 {"retcode":0,"result":{"result":0}}
*/
func (this *WebQQ) set_long_nick2(nlk string) (v *Result, err error) {
	data, err := this.postApi("set_long_nick2", "nlk", nlk)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取好友列表(POST)
 ------------------

 r = {"h":"hello","hash":"ABF6A3FE","vfwebqq":"..."}
*/
func (this *WebQQ) get_user_friends2() (v *ResultUserFrientds, err error) {
	data, err := this.postApi("get_user_friends2",
		"h", "hello",
		"hash", this.GenHash(),
	)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取好友列表 result结构
 -----------------------

 {
 	"retcode": 0,
 	"result": {
 		"friends": [{...}],
 		"marknames": [],
 		"categories": [],
 		"vipinfo": [{
 			"vip_level": 5,
 			"u": 3255435951,
 			"is_vip": 1
 		}],
 		"info": [{
 			"face": 564,
 			"flag": 16777798,
 			"nick": "XpA.IDEA",
 			"uin": 3255435951
 		}]
 	}
 }
*/
type ResultUserFrientds struct {
	Code   int `json:"retcode"`
	Result struct {
		Friends    []Friend    `json:"friends"`
		Marknames  []Markname  `json:"marknames"`
		Categories []Categorie `json:"categories"`
		VipInfo    []VipInfo   `json:"vipinfo"`
		Info       []Info      `json:"info"`
	} `json:"result"`
}

/*
 好友
 ----

 {"flag":0,"uin":3255435951,"categories":0}
*/
type Friend struct {
}

/*
 备注
 ----

 未知
*/
type Markname struct {
}

/*
 分类
 ----

 未知
*/
type Categorie struct {
}

// 用户简单信息
type Info struct {
	Face string `json:"face"`
	Flag string `json:"flag"`
	Nick string `json:"nick"`
	Uin  string `json:"uin"`
}

/*
 修改群名片(POST)
 ----------------

 r = {"gcode":738328699,"markname":"测试","vfwebqq":"..."}
 {"retcode":6}
*/
func (this *WebQQ) update_group_info2(gcode GCode, markname string) (v *Result, err error) {
	data, err := this.postApi("update_group_info2", "gcode", gcode, "markname", markname)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 退出群(POST)
 ----------------

 r = {"gcode":738328699,"vfwebqq":"..."}
 {"retcode":10001} // 失败
*/
func (this *WebQQ) quit_group2(gcode GCode) (v *Result, err error) {
	data, err := this.postApi("quit_group2", "gcode", gcode)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取群列表(POST)
 ----------------

 r = {"vfwebqq":"..."}
*/
func (this *WebQQ) get_group_name_list_mask2() {
	this.postApi("get_group_name_list_mask2")
}

/*
 群列表 result结构
 -----------------

 {
 	"retcode": 0,
 	"result": {
 		"gmasklist": [],
 		"gnamelist": [{...}],
 		"gmarklist": []
 	}
 }
*/
type ResultGroupNameListMask struct {
	Code   int `json:"retcode"`
	Result struct {
		Gmasklist  []interface{} `json:"gmasklist"` // TODO
		GroupNames []GroupName   `json:"gnamelist"`
		GroupMarks []GroupMark   `json:"gmarklist"`
	} `json:"result"`
}

/*
 群名
 ----

 {"flag": 16777217,"name": "群名","gid": 221664830,"code": 738328699}
*/
type GroupName struct {
	Gid  Uin    `json:"gid"`
	Code GCode  `json:"code"`
	Flag uint32 `json:"flag"`
	Name string `json:"name"`
}

/*
 群备注
 ------

 未知
*/
type GroupMark struct {
}

/*
 查找好友(需要验证码)
 --------------------

 {"retcode":0,"result":{"face":147,"birthday":{"month":4,"year":1996,"day":9},"occupation":"待业/无业/失业","phone":"","allow":4,"college":"","constel":3,"blood":0,"stat":20,"homepage":"","country":"中国","city":"温州","uiuin":"","personal":"爱上你，不是因为你给了我需要的东西，而是因为你给了我从未有过的感觉。","nick":"小さくて暗い","shengxiao":1,"email":"","token":"e5c67a587252c899fac04fce508372b0e650b1622faa7230","province":"浙江","account":1159549778,"gender":"male","tuin":2265896498,"mobile":"-"}}
*/
func (this *WebQQ) search_qq_by_uin2(tuin, verifysession, code string) {
	this.api("search_qq_by_uin2", "tuin", tuin, "verifysession", verifysession, "code", code)
}

/*
 [算法] 获取好友列表的hash算法
 -----------------------------

 "h":"hello"
 1. 这是一个32位分组密码
 2. 取十进制的uin逐位与
 2. ptwebqq中的每4个char值分组 OR 入c , 溢出则从 0 开始计算
 3. 再 XOR 入d
 4. 得出最后32位hash d 以hex形式返回

 function(b, i) {
     for (var a = [], s = 0; s < b.length; s++)
         a[s] = b.charAt(s) - 0;
     for (var j = 0, d = -1, s = 0; s < a.length; s++) {
         j += a[s];
         j %= i.length;
         var c = 0;
         if (j + 4 > i.length)
             for (var l = 4 + j - i.length, x = 0; x < 4; x++)
                 c |= x < l ? (i.charCodeAt(j + x) & 255) << (3 - x) * 8 : (i.charCodeAt(x - l) & 255) << (3 - x) * 8;
         else
             for (x = 0; x < 4; x++)
                 c |= (i.charCodeAt(j + x) & 255) << (3 - x) * 8;
         d ^= c
     }
     a = [];
     a[0] = d >> 24 & 255;
     a[1] = d >> 16 & 255;
     a[2] = d >> 8 & 255;
     a[3] = d & 255;
     d = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"];
     s = "";
     for (j = 0; j < a.length; j++)
         s += d[a[j] >> 4 & 15], s += d[a[j] & 15];
     return s
 }
*/
func (this *WebQQ) GenHash() string {
	uin := []byte(this.Uin.String())
	for i, v := range uin {
		uin[i] = v - 48
	}
	j := uint32(0)
	pt := []byte(this.PtWebQQ)
	ptlen := uint32(len(pt))
	d := uint32(0xffffffff)
	for i, _ := range uin {
		j = (j + uint32(uin[i])) % ptlen
		c := uint32(0)
		if j+4 > ptlen {
			l := 4 + j - ptlen
			for x := uint32(0); x < 4; x++ {
				if x < l {
					c |= uint32(pt[j+x]) << ((3 - x) * 8)
				} else {
					c |= uint32(pt[x-l]) << ((3 - x) * 8)
				}
			}
		} else {
			for x := uint32(0); x < 4; x++ {
				c |= uint32(pt[j+x]) << ((3 - x) * 8)
			}
		}
		d ^= c
	}
	return fmt.Sprintf("%X", d)
}

/*
 [算法] 获取好友列表的hash算法2
 ------------------------------

 function(b, i) {
     for (var a = i + "password error", s = "", j = []; ; )
         if (s.length <= a.length) {
             if (s += b, s.length == a.length)
                 break
         } else {
             s =
             s.slice(0, a.length);
             break
         }
     for (var d = 0; d < s.length; d++)
         j[d] = s.charCodeAt(d) ^ a.charCodeAt(d);
     a = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"];
     s = "";
     for (d = 0; d < j.length; d++)
         s += a[j[d] >> 4 & 15], s += a[j[d] & 15];
     return s
 }
*/
func (this *WebQQ) GenHash2() (ret string) {
	a := this.PtWebQQ + "password error"
	uin := []byte(this.Uin.String())

	s := make([]byte, 0, len(a))
	for {
		if len(s) < len(a) {
			s = append(s, uin...)
			if len(s) == len(a) {
				break
			}
		} else {
			s = s[:len(a)]
			break
		}
	}

	j := make([]byte, len(s))
	for i, _ := range j {
		j[i] = s[i] ^ a[i]
	}

	const key = "0123456789ABCDEF"

	for i := 0; i < len(a); i++ {
		ret += string(key[j[i]>>4&15])
		ret += string(key[j[i]&15])
	}
	return
}
