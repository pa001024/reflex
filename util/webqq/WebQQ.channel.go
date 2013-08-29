package webqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pa001024/MoeWorker/util"
)

const (
	_CHANNEL_URL     = "http://d.web2.qq.com/channel/"
	_CHANNEL_REFERER = "http://d.web2.qq.com/proxy.html?v=20110331002&callback=1&id=2"
)

// 通用channel接口(GET)
func (this *WebQQ) channel(api string, args ...interface{}) (body []byte, err error) {
	val := url.Values{
		"clientid":   {this.clientid},
		"psessionid": {this.psessionid},
		"t":          {util.JsCurrentTime()},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.getWithReferer(fmt.Sprintf("%s%s?%s", _CHANNEL_URL, api, val.Encode()), _CHANNEL_REFERER)
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
			append(args, "clientid", this.clientid, "psessionid", this.psessionid)...,
		)},
		"clientid":   {this.clientid},
		"psessionid": {this.psessionid},
	}
	l := len(args) - 1
	for i := 0; i < l; i += 2 {
		val.Add(args[i].(string), fmt.Sprint(args[i+1]))
	}
	res, err := this.postFormWithReferer(_CHANNEL_URL+api, _CHANNEL_REFERER, val)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}

// 登录
func (this *WebQQ) login2() (v *ResultLogin2, err error) {
	data, err := this.postChannel("login2", "status", "online", "ptwebqq", this.ptwebQQ, "passwd_sig", "")
	if err == nil {
		v = &ResultLogin2{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 登录 返回值结构
type ResultLogin2 struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result struct {
		Uin        Uin    `json:"uin"`
		VerifyCode string `json:"vfwebqq"`
		SessionId  string `json:"psessionid"`
		Status     string `json:"status"`
	} `json:"result"`
}

// 获取消息
func (this *WebQQ) poll2() (v *ResultPoll, err error) {
	data, err := this.postChannel("poll2")
	if err == nil {
		v = &ResultPoll{}
		err = json.Unmarshal(data, v)
	}
	return
}

// poll2 result结构
type ResultPoll struct {
	Code   int    `json:"retcode"`
	Msg    string `json:"errmsg"`
	Result []struct {
		Type  string          `json:"poll_type"`
		Value json.RawMessage `json:"value"` // see WebQQ.event.go
	} `json:"result"`
}

// 获取在线好友
func (this *WebQQ) get_online_buddies2() (v *ResultOnlineBuddies, err error) {
	data, err := this.channel("get_online_buddies2")
	if err == nil {
		v = &ResultOnlineBuddies{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 在线好友result结构
type ResultOnlineBuddies struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin        Uin    `json:"uin"`
		Status     string `json:"status"`
		ClientType uint32 `json:"client_type"`
	} `json:"result"`
}

// 正在输入
func (this *WebQQ) input_notify2(to_uin string) (v *Result, err error) {
	data, err := this.channel("input_notify2", "to_uin", to_uin)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 窗口抖动
func (this *WebQQ) shake2(to_uin string) (v *Result, err error) {
	data, err := this.channel("shake2", "to_uin", to_uin)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 获取当前聊天列表
func (this *WebQQ) get_recent_list2() (v *ResultRecentList, err error) {
	data, err := this.postChannel("get_recent_list2")
	if err == nil {
		v = &ResultRecentList{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 当前聊天列表result结构 {"retcode":0,"result":[{"uin":3255435951,"type":0},{"uin":221664830,"type":1}]}
type ResultRecentList struct {
	Code   int `json:"retcode"`
	Result []struct {
		Uin  Uin    `json:"uin"`
		Type uint32 `json:"type"`
	} `json:"result"`
}

// 发送私聊消息
func (this *WebQQ) send_buddy_msg2(to Uin, content ContentM, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_buddy_msg2(to = %s , content = %v , msg_id = %v)", to, content, msg_id)
	data, err := this.postChannel("send_buddy_msg2",
		"to", to,
		"face", "552", // 这是啥?
		"content", content.Encode().EncodeString(),
		"msg_id", msg_id,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}

// 发送群消息
func (this *WebQQ) send_qun_msg2(group_uin Uin, content ContentM, msg_id uint32) (v *Result, err error) {
	util.DEBUG.Logf("send_qun_msg2(group_uin = %s , content = %v , msg_id = %v)", group_uin, content, msg_id)
	data, err := this.postChannel("send_qun_msg2",
		"group_uin", group_uin,
		"content", content.Encode().EncodeString(),
		"msg_id", msg_id,
	)
	if err == nil {
		v = &Result{}
		err = json.Unmarshal(data, v)
	}
	return
}
