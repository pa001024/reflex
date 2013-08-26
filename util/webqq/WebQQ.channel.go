package webqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/pa001024/MoeWorker/util"
)

const (
	CHANNEL_URL     = "http://d.web2.qq.com/channel/"
	CHANNEL_REFERER = "http://d.web2.qq.com/proxy.html?v=20110331002&callback=1&id=2"
)

// 通用channel接口(GET)
func (this *WebQQ) channel(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"clientid":   {this.ClientId},
		"psessionid": {this.SessionId},
		"t":          {fmt.Sprintf("%d", time.Now().UnixNano()/1e6)},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.getWithReferer(fmt.Sprintf("%s%s?%s", CAPTCHA_URL, api, val.Encode()), CHANNEL_REFERER)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 通用channel接口(POST)
func (this *WebQQ) postChannel(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"r": {util.ToJson(
			append(args, "clientid", this.ClientId, "psessionid", this.SessionId)...,
		)},
		"clientid":   {this.ClientId},
		"psessionid": {this.SessionId},
	}
	l := len(args) + 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.postFormWithReferer(CAPTCHA_URL+api, CHANNEL_REFERER, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 登录
func (this *WebQQ) login2() (v *ResultLogin2, err error) {
	data, err := this.postChannel(CHANNEL_URL+"login2",
		"status", "online",
		"ptwebqq", this.PtWebQQ,
		"passwd_sig", "",
	)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 登录 消息结构
 -------------

 {
 	"retcode": 0,
 	"result": {
 		"uin": 2735284921,
 		"cip": 3080236829,
 		"index": 1075,
 		"port": 47943,
 		"status": "online",
 		"vfwebqq": "209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4",
 		"psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133392e372e3136300000443100000163026e0400b92209a36d0000000a404b454773376a7457646d00000028209da35a9665546efac6e1032577fd75e8fcae3e2d7a264fc64fe598064245285ae63a270bc204f4",
 		"user_state": 0,
 		"f": 0
 	}
 }
*/
type ResultLogin2 struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result struct {
		Uin        Uin    `json:"uin"`
		VerifyCode string `json:"vfwebqq"`
		SessionId  string `json:"psessionid"`
		Status     string `json:"status"`
		// CIP        uint32 `json:"cip"` // 没用
		// Index     uint32 `json:"index"` // 没用
		// Port      uint32 `json:"port"`       // 没用
		// UserState uint32 `json:"user_state"` // 没用
		// F          uint32 `json:"f"` // 没用
	} `json:"result"`
}

/*
 获取在线好友
 ------------

 {"retcode":0,"result":[{"uin":3255435951,"status":"online","client_type":1}]}
*/
func (this *WebQQ) get_online_buddies2() (v *ResultOnlineBuddies, err error) {
	data, err := this.channel("get_online_buddies2")
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

type ResultOnlineBuddies struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin        Uin    `json:"uin"`
		Status     string `json:"status"`
		ClientType uint32 `json:"client_type"`
	} `json:"result"`
}

/*
 正在输入
 --------

 {"retcode":0,"result":"ok"}
*/
func (this *WebQQ) input_notify2(to_uin string) (v *Result, err error) {
	data, err := this.channel("input_notify2", "to_uin", to_uin)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 窗口抖动
 --------

 to_uin:22607026
 clientid:38898497
 psessionid:...
 t:1375922771786
 {"retcode":0,"result":"ok"}
*/
func (this *WebQQ) shake2(to_uin string) (v *Result, err error) {
	data, err := this.channel("shake2", "to_uin", to_uin)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 获取什么列表?
 -------------

 {"retcode":0,"result":[{"uin":3255435951,"type":0},{"uin":221664830,"type":1}]}
*/
func (this *WebQQ) get_recent_list2() (v *ResultRecentList, err error) {
	data, err := this.postChannel("get_recent_list2")
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

type ResultRecentList struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin  Uin    `json:"uin"`
		Type uint32 `json:"type"`
	} `json:"result"`
}

/*
 发送私聊消息
 ------------

 r = {
 	"to": 3255435951,
 	"face": 552,
 	"content": "[\"asd\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
 	"msg_id": 38610005,
 	"clientid": "10861648",
 	"psessionid": "..."
 }
 {"retcode":0,"result":"ok"}
*/
func (this *WebQQ) send_buddy_msg2(to Uin, content ContentM, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_buddy_msg2(to = %s , content = %v , msg_id = %v)", to, content, msg_id)
	data, err := this.postChannel("send_buddy_msg2",
		"to", to,
		"face", "552", // 这是啥?
		"content", content.Encode().EncodeString(),
		"msg_id", msg_id,
	)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 发送群消息
 ----------

 r = {
 	"group_uin": 221664830,
 	"content": "[\"msg\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
 	"msg_id": 38610004,
 	"clientid": "10861648",
 	"psessionid": "..."
 }
 {"retcode":0,"result":"ok"}
*/
func (this *WebQQ) send_qun_msg2(group_uin Uin, content ContentM, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_qun_msg2(group_uin = %s , content = %v , msg_id = %v)", group_uin, content, msg_id)
	data, err := this.postChannel("send_qun_msg2",
		"group_uin", group_uin,
		"content", content.Encode().EncodeString(),
		"msg_id", msg_id,
	)
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取消息
func (this *WebQQ) poll2() (v *ResultPoll, err error) {
	util.DEBUG.Log("poll2()")
	data, err := this.postChannel("poll2")
	if err != nil {
		err = json.Unmarshal(data, v)
	}
	return
}

/*
 poll2 result结构
 ----------------

 {"retcode":102,"errmsg":""} // 什么都没有
 {"retcode":116,"p":"39bd5c71be123aaf451073154d52bfef78b1adaa0d087601"} // 什么都没有
 {
 	"retcode": 0,
 	"result": [{
 		"poll_type": "message",
 		"value": {...}
 	},{
 		"poll_type": "group_message",
 		"value": {...}
 	},{
 		"poll_type": "buddies_status_change",
 		"value": {...}
 	}]
 }
*/
type ResultPoll struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result []struct {
		Type  string          `json:"poll_type"`
		Value json.RawMessage `json:"value"` // 转到WebQQ.event.go
	} `json:"result"`
}
